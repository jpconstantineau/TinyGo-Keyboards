package main

import (
	"machine"
	"machine/usb/hid/keyboard"
	"time"
)

func blinkled() {
	// LED next to the reset button - this anonymous function will start as a concurrent go routine to make the LED blink...
	redled := machine.GPIO27
	redled.Configure(machine.PinConfig{Mode: machine.PinOutput})
	for {
		redled.High()
		time.Sleep(500 * time.Millisecond)
		redled.Low()
		time.Sleep(500 * time.Millisecond)
	}
}

func main() {
	// start keyboard first!
	kb := keyboard.Port()
	go blinkled()

	greenled := machine.GPIO28
	greenled.Configure(machine.PinConfig{Mode: machine.PinOutput})

	const columncount = 18
	const rowcount = 6
	const settlingtime = 75 // this is the electrical settling time when selecting a row.  In uSec
	const scanningtime = 10 // this is the time between scans. in mSec

	// Key Matrix Definition - Columns
	columns := [columncount]machine.Pin{machine.GPIO0, machine.GPIO1, machine.GPIO2, machine.GPIO3,
		machine.GPIO4, machine.GPIO5, machine.GPIO6, machine.GPIO7,
		machine.GPIO8, machine.GPIO9, machine.GPIO10, machine.GPIO11,
		machine.GPIO12, machine.GPIO13, machine.GPIO14, machine.GPIO15, machine.GPIO16, machine.GPIO17}

	// Key Matrix Definition - Rows
	rows := [rowcount]machine.Pin{machine.GPIO18, machine.GPIO19, machine.GPIO20, machine.GPIO21, machine.GPIO22, machine.GPIO26}

	for _, col := range columns {
		col.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	}

	for _, row := range rows {
		row.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	}

	// These are the trigger keys to select layers. - Order can be important
	specialkeys := []int{103}

	keys := []keyboard.Keycode{
		keyboard.KeyEsc, keyboard.ASCII00, keyboard.KeyF1, keyboard.KeyF2, keyboard.KeyF3, keyboard.KeyF4, keyboard.ASCII00, keyboard.KeyF5, keyboard.KeyF6, keyboard.KeyF7, keyboard.KeyF8, keyboard.KeyF9, keyboard.KeyF10, keyboard.KeyF11, keyboard.KeyF12, keyboard.KeyPrintscreen, keyboard.KeyScrollLock, keyboard.KeyPause,
		keyboard.KeyTilde, keyboard.Key1, keyboard.Key2, keyboard.Key3, keyboard.Key4, keyboard.Key5, keyboard.Key6, keyboard.Key7, keyboard.Key8, keyboard.Key9, keyboard.Key0, keyboard.KeyMinus, keyboard.KeyEqual, keyboard.ASCII00, keyboard.KeyBackspace, keyboard.KeyInsert, keyboard.KeyHome, keyboard.KeyPageUp,
		keyboard.KeyTab, keyboard.ASCII00, keyboard.KeyQ, keyboard.KeyW, keyboard.KeyE, keyboard.KeyR, keyboard.KeyT, keyboard.KeyY, keyboard.KeyU, keyboard.KeyI, keyboard.KeyO, keyboard.KeyP, keyboard.KeyLeftBrace, keyboard.KeyRightBrace, keyboard.KeyBackslash, keyboard.KeyDelete, keyboard.KeyEnd, keyboard.KeyPageDown,
		keyboard.KeyCapsLock, keyboard.ASCII00, keyboard.KeyA, keyboard.KeyS, keyboard.KeyD, keyboard.KeyF, keyboard.KeyG, keyboard.KeyH, keyboard.KeyJ, keyboard.KeyK, keyboard.KeyL, keyboard.KeySemicolon, keyboard.KeyQuote, keyboard.KeyEnter, keyboard.ASCII00, keyboard.ASCII00, keyboard.ASCII00, keyboard.ASCII00, keyboard.ASCII00,
		keyboard.KeyModifierLeftShift, keyboard.KeyZ, keyboard.KeyX, keyboard.KeyC, keyboard.KeyV, keyboard.KeyB, keyboard.KeyN, keyboard.KeyM, keyboard.KeyComma, keyboard.KeyPeriod, keyboard.KeySlash, keyboard.ASCII00, keyboard.KeyRightShift, keyboard.ASCII00, keyboard.ASCII00, keyboard.KeyUpArrow, keyboard.ASCII00,
		keyboard.KeyLeftCtrl, keyboard.KeyLeftGUI, keyboard.ASCII00, keyboard.KeyLeftAlt, keyboard.ASCII00, keyboard.ASCII00, keyboard.KeySpace, keyboard.ASCII00, keyboard.ASCII00, keyboard.ASCII00, keyboard.KeyRightAlt, keyboard.KeyRightGUI, keyboard.ASCII00, keyboard.ASCII00, keyboard.KeyRightCtrl, keyboard.KeyLeftArrow, keyboard.KeyDownArrow, keyboard.KeyRightArrow,

		keyboard.KeyEsc, keyboard.ASCII00, keyboard.KeyF1, keyboard.KeyF2, keyboard.KeyF3, keyboard.KeyF4, keyboard.ASCII00, keyboard.KeyF5, keyboard.KeyF6, keyboard.KeyF7, keyboard.KeyF8, keyboard.KeyF9, keyboard.KeyF10, keyboard.KeyF11, keyboard.KeyF12, keyboard.KeyPrintscreen, keyboard.KeyScrollLock, keyboard.KeyPause,
		keyboard.KeyTilde, keyboard.Key1, keyboard.Key2, keyboard.Key3, keyboard.Key4, keyboard.Key5, keyboard.Key6, keyboard.Key7, keyboard.Key8, keyboard.Key9, keyboard.Key0, keyboard.KeyMinus, keyboard.KeyEqual, keyboard.ASCII00, keyboard.KeyBackspace, keyboard.KeyInsert, keyboard.KeyHome, keyboard.KeyPageUp,
		keyboard.KeyTab, keyboard.ASCII00, keyboard.KeyQ, keyboard.KeyW, keyboard.KeyE, keyboard.KeyR, keyboard.KeyT, keyboard.KeyY, keyboard.KeyU, keyboard.KeyI, keyboard.KeyO, keyboard.KeyP, keyboard.KeyLeftBrace, keyboard.KeyRightBrace, keyboard.KeyBackslash, keyboard.KeyDelete, keyboard.KeyEnd, keyboard.KeyPageDown,
		keyboard.KeyCapsLock, keyboard.ASCII00, keyboard.KeyA, keyboard.KeyS, keyboard.KeyD, keyboard.KeyF, keyboard.KeyG, keyboard.KeyH, keyboard.KeyJ, keyboard.KeyK, keyboard.KeyL, keyboard.KeySemicolon, keyboard.KeyQuote, keyboard.KeyEnter, keyboard.ASCII00, keyboard.ASCII00, keyboard.ASCII00, keyboard.ASCII00, keyboard.ASCII00,
		keyboard.KeyModifierLeftShift, keyboard.KeyZ, keyboard.KeyX, keyboard.KeyC, keyboard.KeyV, keyboard.KeyB, keyboard.KeyN, keyboard.KeyM, keyboard.KeyComma, keyboard.KeyPeriod, keyboard.KeySlash, keyboard.ASCII00, keyboard.KeyRightShift, keyboard.ASCII00, keyboard.ASCII00, keyboard.KeyUpArrow, keyboard.ASCII00,
		keyboard.KeyLeftCtrl, keyboard.KeyLeftGUI, keyboard.ASCII00, keyboard.KeyLeftAlt, keyboard.ASCII00, keyboard.ASCII00, keyboard.KeySpace, keyboard.ASCII00, keyboard.ASCII00, keyboard.ASCII00, keyboard.KeyRightAlt, keyboard.KeyRightGUI, keyboard.ASCII00, keyboard.ASCII00, keyboard.KeyRightCtrl, keyboard.KeyLeftArrow, keyboard.KeyDownArrow, keyboard.KeyRightArrow,
	}

	layer := 0
	specialkeyid := -1
	//	specialkeypressed :=false
	var keystates [columncount * rowcount]bool // must be an array for the scan loop to work...

	// END OF SETUP
	for {
		// scan inputs
		for j, row := range rows {
			row.Configure(machine.PinConfig{Mode: machine.PinOutput})
			row.Low()
			time.Sleep(time.Microsecond * settlingtime)
			for i, key := range columns {
				keystates[j*columncount+i] = key.Get()
			}
			row.High()
			row.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
		}

		// process for layers
		layer = 0
		specialkeyid = -1
		//		specialkeypressed =false
		for i, triggerkey := range specialkeys {
			if !(keystates[triggerkey]) {
				layer = i + 1
				//				specialkeypressed =true
				specialkeyid = i
			}
		}

		// update green LED
		if (layer > 0) || (kb.CapsLockLed()) {
			greenled.High()

		} else {
			greenled.Low()
		}

		// send keypresses
		for i, state := range keystates {
			if i != specialkeyid {
				if state {
					kb.Up(keys[i+layer*rowcount*columncount])
				} else {
					kb.Down(keys[i+layer*rowcount*columncount])
				}
			}
		}

		// delay between scans
		time.Sleep(scanningtime * time.Millisecond)
	}
}
