package dashboard

import (
	"fmt"
	"os"

	"github.com/ayoubed/datadog-home-project/statsagent"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

// RunDashboard displays stats to the user
func main() {
	if err := ui.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
	}
	defer ui.Close()

	// should be used to feed data to the UI controller
	var statschan1, statschan2 chan statsagent.WebsiteStats

	var lastStats1, lastStats2 statsagent.WebsiteStats

	controller := newController()

	uiEvents := ui.PollEvents()
	for {
		select {
		case e := <-uiEvents:
			switch e.Type {
			case ui.KeyboardEvent:
				return
			case ui.ResizeEvent:
				controller.Resize()
			}
		case stats := <-statschan1:
			lastStats1 = stats
			for k, v := range lastStats1.Stats {
				fmt.Println(k, v)
			}
			fmt.Println()
			controller.Render("Me1", "Me6")
		case stats := <-statschan2:
			for k, v := range lastStats2.Stats {
				fmt.Println(k, v)
			}
			fmt.Println()
			fmt.Println()
			lastStats2 = stats
			controller.Render("Me1", "Me6")
		}
	}

}

func drawStats(height, width int, stats statsagent.WebsiteStats) {
	p0 := widgets.NewParagraph()
	if stats.ID == 1 {
		p0.Title = "1 Min Statistics(updated every 10 seconds)"
	} else {
		p0.Title = "1 Hour Statistics(updated every 1 minute)"
	}
	// lines := []string{}
	// for k, _ := range stats.Stats {
	// 	lines = append(lines, fmt.Sprintf("%v -----", k))
	// }
	// p0.Text = strings.Join(lines, "\n")
	p0.Text = fmt.Sprintf("%v", len(stats.Stats))
	p0.SetRect(0, (stats.ID-1)*height/3, width, stats.ID*height/3)

	ui.Render(p0)
}
