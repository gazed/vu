// Copyright © 2014 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package panel

// Control responds to user actions that are linked back to an application
// specific action handler. A Control is expected to be contained within a
// Panel or Section.
type Control interface {
	Widget // Control is a Widget.

	// AddReaction sets the trigger and related callback for when the trigger occurs.
	// The callback is called each time the trigger is pressed while the mouse
	// is over the control.
	AddReaction(trigger string, callback func()) Control
	RemReaction(trigger string) // Remove the indicated trigger.

	// Controls can use the widget image as a background as a hover indication.
	ImgOnHover() (show bool) // Get or,
	SetImgOnHover(show bool) // ...set the widget image behaviour.
}

// =============================================================================

// control implements Control.
type control struct {
	widget                      // control is a widget.
	reactions map[string]func() // Triggers and callback functions.
	hoverImg  bool              // Optional image behaviour.
}

// Control implementation.
func (c *control) ImgOnHover() (show bool) { return c.hoverImg }

// Control implementation.
func (c *control) SetImgOnHover(show bool) { c.hoverImg = show }

// AddReaction implements Control.
func (c *control) AddReaction(t string, cb func()) Control {
	c.reactions[t] = cb
	return c
}

// RemReaction implements Control.
func (c *control) RemReaction(t string) { delete(c.reactions, t) }

// Widget internal implementation.
// react checks if the callback is triggered. If the trigger input happened
// then call the available callback.
func (c *control) react(in *In) (overId uint) {
	if c.vis && c.Over(in.Mx, in.My) {
		overId = c.id

		// Trigger only on the first reaction that matches.
		// FUTURE: check if it makes any sense to trigger multiple?
		if len(c.reactions) > 0 {
			for trigger, callback := range c.reactions {
				if down, ok := in.Down[trigger]; ok && down == 1 {
					callback()
					return
				}
			}
		}
	}
	return
}

// =============================================================================

// FUTURE: Possibly add specific controls.
//    Popup dialog : A separate panel that can be shown/hidden.
//    Button       : Really just a control. Does it need any other behaviour.
//    Menu         : A popup list of controls, one column, no margins.
//    TabbedPane   : Harder... Overlapping panels. Tab image for each tab.
//    Others:
//       CheckBox
//       ComboBox
//       List
//       RadioButton
//       Slider
//       TextField ... depends on how text handling is done.
//       ScrollPane
