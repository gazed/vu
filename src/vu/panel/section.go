// Copyright © 2014 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package panel

// Section groups and positions a child panels and controls. A section can have
// a background image that is drawn behind the child controls and section.
// Sections also have a layout that calculates the child widgets sizes and
// locations.
type Section interface {
	Widget             // Section is a Widget
	Widgets() []Widget // Get the list of widgets in this panel.

	// Add* new child widgets. They child widgets size and location are
	// expected to be controlled by a layout.
	AddSection() Section // Add a child Section.
	AddControl() Control // Add a child Control.

	// Layouts concern is the positioning of a panels child widgets.
	SetLayout(l Layout) // Nil layouts ignored.
}

// =============================================================================

// section implements Section.
type section struct {
	widget           // A panel is a widget and conforms to iWidget.
	ww, wh  int      // Each section knows the overall window size.
	widgets []Widget // Child widgets.
	layout  Layout   // Layout manager.
}

// newSection creates and initializes a panel.
func newSection(ww, wh int) *section {
	s := &section{widget: newWidget(), ww: ww, wh: wh}
	s.layout = &GridLayout{Columns: 1, Margin: 0}
	s.widgets = []Widget{}
	return s
}

// Implements Section.
func (s *section) Widgets() []Widget { return s.widgets }

// Implements Section.
func (s *section) SetLayout(l Layout) {
	if l != nil {
		s.layout = l
		s.align()
	}
}

// Implements Section.
func (s *section) AddSection() Section {
	sub := newSection(s.ww, s.wh)
	s.widgets = append(s.widgets, sub)
	s.align()
	return sub
}

// Implements Section.
func (s *section) AddControl() Control {
	c := &control{widget: newWidget()}
	c.reactions = map[string]func(){}
	c.vis = true
	s.widgets = append(s.widgets, c)
	s.align()
	return c
}

// align ensures that all panels are aligned with the top panels width and
// height information and that all child widgets are aligned according to
// each panels layout.
func (s *section) align() {
	s.layout.Align(s, s.ww, s.wh) // align this panels child widgets.
	for _, w := range s.widgets { // align the child panels.
		if sub, ok := w.(*section); ok {
			sub.ww, sub.wh = s.ww, s.wh // propogate window size.
			sub.align()
		}
	}
}

// react filters the input down the hiearchy to the widgets that
// are under the mouse.
func (s *section) react(in *In) uint {
	if s.Over(in.Mx, in.My) {
		for _, w := range s.widgets {
			if id := w.react(in); id > 0 {
				return id
			}
		}
	}
	return 0
}
