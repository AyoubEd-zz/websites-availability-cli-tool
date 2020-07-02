// package dashboard

// import (
// 	"fmt"
// 	"time"

// 	"github.com/ayoubed/datadog-home-project/statsagent"
// )

// func Run(urls []string, interval, span int) {
// 	go func() {
// 		ticker := time.NewTicker(time.Duration(interval) * time.Second)
// 		for {
// 			select {
// 			case t := <-ticker.C:
// 				fmt.Printf("[last %vs stats]---------------------%v-------------------\n", interval, t)
// 				res := statsagent.GetStats(urls, span)
// 				for k, v := range res {
// 					fmt.Println("> ", k)
// 					fmt.Printf("%+v\n", v)
// 				}
// 			}
// 		}
// 	}()
// }
package dashboard

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/ayoubed/datadog-home-project/statsagent"
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

	fmt.Println()
	fmt.Println()

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
						v.Title = fmt.Sprintf("Statistics for the last %vs(updated every %vs)", currentView.UpdateInterval, currentView.TimeFrame)
						v.FgColor = gocui.ColorCyan
						v.SelBgColor = gocui.ColorBlue
						v.SelFgColor = gocui.ColorBlack
						v.Wrap = true
						v.Clear()
						fmt.Fprintln(v, fmt.Sprintf("%-30v %25v %25v %+v\n", "Website", "Average Response Time", "Max Response Time", "Status Codes"))
						for k, value := range res {
							statusCodeSlice := make([]string, 0)
							for code, count := range value.StatusCodeCount {
								statusCodeSlice = append(statusCodeSlice, fmt.Sprintf("%v:%v", code, count))
							}
							statusCodeStr := fmt.Sprintf("(%v)", strings.Join(statusCodeSlice, " "))
							fmt.Fprintln(v, fmt.Sprintf("%-30v %25v %25v %-25v\n", k, value.AvgResponseTime, value.MaxResponseTime, statusCodeStr))
						}
						return nil
					})
				}
			}
		}(index, &views[index])
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func layout(g *gocui.Gui, views []View) func(*gocui.Gui) error {
	maxX, maxY := g.Size()
	return func(g *gocui.Gui) error {
		for index, view := range views {
			v, err := g.SetView(strconv.Itoa(view.TimeFrame), 0, index*(maxY/3), maxX, (index+1)*(maxY/3))
			if err != nil {
				if err != gocui.ErrUnknownView {
					log.Panic("Error setting views")
				}
				fmt.Fprintln(v, "One moment, we're getting the lastest statistics...")
			}
		}
		return nil
	}
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
