package dashboard

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/ayoubed/datadog-home-project/statsagent"
	"github.com/fatih/color"
	"github.com/jroimartin/gocui"
	"golang.org/x/sync/errgroup"
)

// View represents one entity on the Gui
// views display stats for a user-defined timeframe
// they are updated following a user-defined interval
type View struct {
	UpdateInterval int   `json:"updateInterval"`
	TimeFrame      int64 `json:"timeFrame"`
}

// Run displays the statistics, and alerts in our terminal
func Run(ctx context.Context, urls []string, views []View, alertc chan string, done context.CancelFunc) error {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return fmt.Errorf("error creating GUI: %v", err)
	}
	defer g.Close()

	// set the layout of the GUI
	g.SetManagerFunc(layout(g, views))

	// launch goroutines to continuously update our views
	errg, gctx := errgroup.WithContext(ctx)

	for _, view := range views {
		view := view
		errg.Go(func() error {
			return updateView(gctx, view, g, urls)
		})
	}

	errg.Go(func() error {
		return monitorAlertChan(gctx, g, alertc)
	})

	// Set key bindings for the GUI
	if err := initKeybindings(g, done); err != nil {
		return err
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		return fmt.Errorf("error while executing the dashboards main loop: %v", err)
	}

	if err := errg.Wait(); err != nil {
		return fmt.Errorf("dashboard process error: %v", err)
	}

	return nil
}

func initKeybindings(g *gocui.Gui, done context.CancelFunc) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			done()
			return quit(g, v)
		}); err != nil {
		return fmt.Errorf("error while trying to close our GUI: %v", err)
	}
	if err := g.SetKeybinding("stdin", gocui.KeyArrowUp, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			scrollView(v, -1)
			return nil
		}); err != nil {

		log.Panicln(err)
	}
	if err := g.SetKeybinding("stdin", gocui.KeyArrowDown, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			scrollView(v, 1)
			return nil
		}); err != nil {

		log.Panicln(err)
	}
	return nil
}

func updateView(ctx context.Context, currentView View, g *gocui.Gui, urls []string) error {

	ticker := time.NewTicker(time.Duration(currentView.UpdateInterval) * time.Second)
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return nil
		case t := <-ticker.C:
			// Grab the latest stats over the given timeframe
			res, err := statsagent.GetStats(urls, t, currentView.TimeFrame)
			if err != nil {
				return fmt.Errorf("error while getting stats to update view: %v", err)
			}

			// update the GUI with the latest stats
			g.Update(func(g *gocui.Gui) error {
				v, err := g.View(strconv.Itoa(int(currentView.TimeFrame)))
				if err != nil {
					return fmt.Errorf("error getting view in update function: %v", err)
				}
				v.Clear()

				// pretty print the stats to our view
				header := color.New(color.FgYellow, color.Bold)
				header.Fprintln(v, fmt.Sprintf("%-30v %21v %21v %21v %21v %21v %25v\n", "Website", "Availability", "Avg Response Time", "Max Response Time", "Avg TTFB", "Max TTFB", "Status Codes"))

				for _, url := range urls {
					value := res[url]
					statusCodeSlice := make([]string, 0)
					for code, count := range value.StatusCodeCount {
						statusCodeSlice = append(statusCodeSlice, fmt.Sprintf("%v:%v", code, count))
					}
					statusCodeStr := fmt.Sprintf("[%v]", strings.Join(statusCodeSlice, " "))
					fmt.Fprintln(v, fmt.Sprintf("%-30v %20.2f%% %21v %21v %21v %21v %25v", url, 100*value.Availability, value.AvgResponseTime, value.MaxResponseTime, value.AvgTimeToFirstByte, value.MaxTimeToFirstByte, statusCodeStr))
				}
				return nil
			})
		}
	}
}

func monitorAlertChan(ctx context.Context, g *gocui.Gui, alertc chan string) error {
	var alerts []string
	for {
		select {
		case alertMessage := <-alertc:
			alerts = append(alerts, alertMessage)
			if err := updateAlertView(g, alerts); err != nil {
				return fmt.Errorf("error while updating the alert view: %v", err)
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func updateAlertView(g *gocui.Gui, alerts []string) error {

	g.Update(func(g *gocui.Gui) error {
		v, err := g.View("alerts")
		if err != nil {
			return err
		}
		v.Clear()

		for i := len(alerts) - 1; i >= 0; i-- {
			fmt.Fprintln(v, alerts[i])
		}
		return nil
	})
	return nil
}

func layout(g *gocui.Gui, views []View) func(*gocui.Gui) error {
	maxX, maxY := g.Size()
	return func(g *gocui.Gui) error {
		// Set stats views
		numViews := len(views) + 1 // number of views, plus the alert channel
		for index, view := range views {
			v, err := g.SetView(strconv.Itoa(int(view.TimeFrame)), 0, index*(maxY/numViews), maxX, (index+1)*(maxY/numViews))
			v.FgColor = gocui.ColorCyan
			if err != nil {
				if err != gocui.ErrUnknownView {
					log.Panic("Error setting views")
				}

				loadingMessage := color.New(color.FgMagenta)
				loadingMessage.Fprintln(v, fmt.Sprintf("\n\n%v One moment, we're waiting for statistics for the last %vs...", "⌛ ", view.TimeFrame))
			}
			v.Title = fmt.Sprintf(" Statistics for the last %vs (updated every %vs) ", view.TimeFrame, view.UpdateInterval)
			v.Wrap = true
		}

		// Set alerts view
		v, err := g.SetView("alerts", 0, (numViews-1)*(maxY/numViews), maxX, maxY)
		v.FgColor = gocui.ColorCyan
		if err != nil {
			if err != gocui.ErrUnknownView {
				log.Panic("Error setting views")
			}
		}
		v.Title = fmt.Sprintf(" Alerts ")
		v.Wrap = true
		return nil
	}
}

func scrollView(v *gocui.View, dy int) error {
	if v != nil {
		v.Autoscroll = false
		ox, oy := v.Origin()
		if err := v.SetOrigin(ox, oy+dy); err != nil {
			return err
		}
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
