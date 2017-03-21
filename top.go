package toplib

import (
	ui "github.com/gizak/termui"
	"github.com/gizak/termui/extra"
	"github.com/vektorlab/toplib/sample"
	"time"
)

// Item is a selectable interface with a unique ID
type Item interface {
	ID() string
}

// Section returns a renderable ui.Grid
type Section interface {
	Name() string
	Grid(Options) *ui.Grid
	Handlers(Options) map[string]func(ui.Event)
}

type Options struct {
	Recorder *Recorder
	Render   func()
}

// Top renders Sections which are periodically updated
type Top struct {
	Exit     chan bool
	Recorder *Recorder // Holds samples
	Sections []Section
	Tabpane  *extra.Tabpane
	Grid     *ui.Grid
	section  int
	Options  Options
}

func NewTop(sections []Section) *Top {
	top := &Top{
		Exit:     make(chan bool),
		Recorder: NewRecorder(),
		Sections: sections,
		Tabpane:  extra.NewTabpane(),
		Grid:     ui.NewGrid(),
	}
	top.Options = Options{
		Render: func() {
			render(top)
		},
		Recorder: top.Recorder,
	}
	return top
}

func handlers(top *Top) {
	ui.DefaultEvtStream.ResetHandlers()
	ui.Handle("/sys/kbd/q", func(ui.Event) {
		top.Exit <- true
	})
	ui.Handle("/sys/kbd/j", func(ui.Event) {
		top.Tabpane.SetActiveLeft()
		render(top)
	})
	ui.Handle("/sys/kbd/k", func(ui.Event) {
		top.Tabpane.SetActiveRight()
		render(top)
	})
	for path, fn := range top.Sections[top.section].Handlers(top.Options) {
		ui.Handle(path, fn)
	}
}

func render(top *Top) {
	handlers(top)
	tabs := []extra.Tab{}
	for _, section := range top.Sections {
		grid := section.Grid(top.Options)
		grid.Width = ui.TermWidth()
		grid.Align()
		tab := extra.NewTab(section.Name())
		tab.AddBlocks(grid)
		tabs = append(tabs, *tab)
	}
	top.Tabpane.SetTabs(tabs...)
	top.Tabpane.Width = ui.TermWidth()
	top.Tabpane.Align()
	ui.Clear()
	ui.Render(top.Tabpane)
}

func Run(top *Top, fn sample.SampleFunc) (err error) {
	if err = ui.Init(); err != nil {
		return err
	}
	defer ui.Close()
	tick := time.NewTicker(500 * time.Millisecond)
	go func() {
		for {
			select {
			case <-top.Exit:
				ui.StopLoop()
				break
			case <-tick.C:
				samples, err := fn()
				if err != nil {
					break
				}
				top.Recorder.Load(samples)
				handlers(top)
				render(top)
			}
		}
	}()
	render(top)
	ui.Loop()
	return err
}
