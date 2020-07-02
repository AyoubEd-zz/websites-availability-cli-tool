package dashboard

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/ayoubed/datadog-home-project/statsagent"
	"github.com/fatih/color"
	"github.com/jroimartin/gocui"
)

type View struct {
	UpdateInterval int `json:"updateInterval"`
	TimeFrame      int `json:"timeFrame"`
}

// Run displays the statistics in terminal
func Run(urls []string, views []View) {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layout(g, views))

	updateViews(views, g, urls)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
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

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func updateViews(views []View, g *gocui.Gui, urls []string) {

	for index := range views {
		go func(currentIndex int, currentView *View) {
			ticker := time.NewTicker(time.Duration(currentView.UpdateInterval) * time.Second)
			for {
				select {
				case <-ticker.C:
					res := statsagent.GetStats(urls, currentView.TimeFrame)
					g.Update(func(g *gocui.Gui) error {
						v, err := g.View(strconv.Itoa(currentView.TimeFrame))
						if err != nil {
							return err
						}
						v.Clear()

						blue := color.New(color.FgYellow)
						blue.Fprintln(v, fmt.Sprintf("%-30v %20v %21v %21v %21v %21v %25v\n", "Website", "Availability", "Average Response Time", "Max Response Time", "Avg TTFB", "Max TTFB", "Status Codes"))

						for _, url := range urls {
							value := res[url]
							statusCodeSlice := make([]string, 0)
							for code, count := range value.StatusCodeCount {
								statusCodeSlice = append(statusCodeSlice, fmt.Sprintf("%v:%v", code, count))
							}
							statusCodeStr := fmt.Sprintf("[%v]", strings.Join(statusCodeSlice, " "))
							fmt.Fprintln(v, fmt.Sprintf("%-30v %20v %21v %21v %21v %21v %25v", url, value.Availability, value.AvgResponseTime, value.MaxResponseTime, value.AvgTimeToFirstByte, value.MaxTimeToFirstByte, statusCodeStr))
						}
						return nil
					})
				}
			}
		}(index, &views[index])
	}
}

func layout(g *gocui.Gui, views []View) func(*gocui.Gui) error {
	maxX, maxY := g.Size()
	return func(g *gocui.Gui) error {
		for index, view := range views {
			v, err := g.SetView(strconv.Itoa(view.TimeFrame), 0, index*(maxY/3), maxX, (index+1)*(maxY/3))
			v.FgColor = gocui.ColorCyan
			if err != nil {
				if err != gocui.ErrUnknownView {
					log.Panic("Error setting views")
				}
				fmt.Fprintln(v, fmt.Sprintf("One moment, we're getting statistics for the last %vs...", view.TimeFrame))
			}
			v.Title = fmt.Sprintf(" Statistics for the last %vs (updated every %vs) ", view.TimeFrame, view.UpdateInterval)
			v.Wrap = true
		}
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
