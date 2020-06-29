package dashboard

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type controller struct {
	Grid  *ui.Grid
	Text1 *widgets.Paragraph
	Text2 *widgets.Paragraph
}

func (c *controller) Resize() {
	c.resize()
	ui.Render(c.Grid)
}

func (c *controller) resize() {
	w, h := ui.TerminalDimensions()
	c.Grid.SetRect(0, 0, w, h)
}

func (c *controller) Render(str1, str2 string) {
	c.Text1.Text = str1
	c.Text2.Text = str2

	c.resize()
	ui.Render(c.Grid)
}

func (c *controller) initUI() {
	c.resize()
	c.Text1.Title = "Stats for the last 1 minute(updated every 10s)"
	c.Text2.Title = "Stats for the last 1 hour(updated every 1min)"

	c.Grid.Set(ui.NewRow(.4, c.Text1), ui.NewRow(.4, c.Text2))
}

func newController() *controller {
	controller := &controller{
		Grid:  ui.NewGrid(),
		Text1: widgets.NewParagraph(),
		Text2: widgets.NewParagraph(),
	}

	return controller
}
