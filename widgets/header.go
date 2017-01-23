package widgets

import (
	"fmt"
	"time"

	ui "github.com/gizak/termui"
)

type CTopHeader struct {
	Time  *ui.Par
	Count *ui.Par
}

func NewCTopHeader() *CTopHeader {
	return &CTopHeader{
		Time:  headerPar(timeStr()),
		Count: headerPar("-"),
	}
}

func (c *CTopHeader) Row() *ui.Row {
	c.Time.Text = timeStr()
	return ui.NewRow(
		ui.NewCol(2, 0, c.Time),
		ui.NewCol(2, 0, c.Count),
	)
}

func (c *CTopHeader) SetCount(val int) {
	c.Count.Text = fmt.Sprintf("%d containers", val)
}

func timeStr() string {
	return time.Now().Local().Format("15:04:05 MST")
}

func headerPar(s string) *ui.Par {
	p := ui.NewPar(fmt.Sprintf(" %s", s))
	p.Border = false
	p.Height = 1
	p.Width = 20
	p.TextFgColor = ui.ColorDefault
	p.TextBgColor = ui.ColorWhite
	p.Bg = ui.ColorWhite
	return p
}
