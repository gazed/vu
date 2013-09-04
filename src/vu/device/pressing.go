// Copyright © 2013 Galvanized Logic Inc.
// Use is governed by a FreeBSD license found in the LICENSE file.

package device

// pressing keeps track of what keys and buttons the user is currently pressing.
// Each valid sequence of key or mouse presses is considered a user request
// for an action which may cause an application reaction.
//
// pressing is used by device to change the native layer raw events
// into key sequence strings.
type pressing struct {
	// each map index is a code representing a key, or mouse button
	// either by itself or in combination with modifier keys.
	allowed  map[int64]*pressed // Allowed keys, key sequences, or buttons
	pressing bool               // is the user still pressing something?
}

// pressed maps a key label to a pressed stage.
type pressed struct {
	label     string // Displayable key sequence.
	isPressed bool   // Track the key sequence up/down state.
}

// newPressedTracker initializes with the allowed key press combinations.
// This is currently:
//    - any non-mod-key by itself
//    - any one mod-key with a non-mod-key
func newPressedTracker(keys map[int]string, mods []int, modNames []string) *pressing {
	p := &pressing{}
	p.allowed = map[int64]*pressed{}

	// all regular keys
	for keyCode, keyLabel := range keys {
		actionCode := int64(keyCode)
		p.allowed[actionCode] = &pressed{keyLabel, false}

		// all regular keys with one mod key combinations.
		for cnt, mod1 := range mods {
			actionCode = int64(mod1)<<32 | int64(keyCode)
			actionLabel := modNames[cnt] + "-" + keyLabel
			p.allowed[actionCode] = &pressed{actionLabel, false}
		}
	}
	return p
}

// pressed is called on a regular basis with up or down key events.
// It keeps track of which keys and mouse buttons are currently pressed,
// in effect turning press events into pressed state.
func (p *pressing) pressed(pressEvent int, modMask int, isPressed bool) {
	actionCode := int64(modMask)<<32 | int64(pressEvent)
	if _, ok := p.allowed[actionCode]; ok {
		p.allowed[actionCode].isPressed = isPressed
		p.pressing = isPressed // key was pressed or released.

		// if the key is released then remove all action combinations related
		// to this action code. This is needed where the key or button is
		// lifted before the modifier keys.
		if !isPressed {
			for actionCode, action := range p.allowed {
				if action.isPressed {
					if actionCode&0xFFFF == int64(pressEvent) {
						p.allowed[actionCode].isPressed = false
					} else {
						p.pressing = true // user is still pushing something.
					}
				}
			}
		}
	}
	// ignore actions not in the map. eg: Ctrl-Alt-Cmd-Shift-*, etc.
}

// down returns the list of currently pressed key sequences where each
// key sequence is a request to do something.  Expected to be used
// by device to request actions for active key/mouse-button presses.
func (p *pressing) down() (down []string) {
	if p.pressing {
		for code, action := range p.allowed {
			if action.isPressed {
				down = append(down, p.allowed[code].label)
			}
		}
	}
	return
}
