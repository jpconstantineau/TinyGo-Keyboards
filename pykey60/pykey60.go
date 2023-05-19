package main

import (
	"image/color"
	"machine"
	"time"

	"machine/usb/hid/keyboard"

	"tinygo.org/x/drivers/buzzer"
	"tinygo.org/x/drivers/ws2812"
)

type note struct {
	tone     float64
	duration float64
}

var leds [61]color.RGBA

func blinkled() {
	// LED next to the reset button - this anonymous function will start as a concurrent go routine to make the LED blink...
	led := machine.GPIO24
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	for {
		led.High()
		time.Sleep(500 * time.Millisecond)
		led.Low()
		time.Sleep(500 * time.Millisecond)
	}
}

func main() {
	// start keyboard first!
	kb := keyboard.Port()
	go blinkled()

	// Buzzer inside the keyboard
	bzrPin := machine.GPIO21
	bzrPin.Configure(machine.PinConfig{Mode: machine.PinOutput})

	bzr := buzzer.New(bzrPin)

	song := []note{
		{buzzer.E5, buzzer.Eighth},
		{buzzer.A4, buzzer.Eighth},
		{buzzer.E4, buzzer.Eighth},
	}

	for _, val := range song {
		bzr.Tone(val.tone, val.duration)
	}

	// RGB LEDs under each key
	var neo = machine.GPIO29
	neo.Configure(machine.PinConfig{Mode: machine.PinOutput})

	ws := ws2812.New(neo)
	rg := false

	for i := range leds {
		rg = !rg
		if rg {
			// Alpha channel is not supported by WS2812 so we leave it out
			leds[i] = color.RGBA{R: 0xff, G: 0x00, B: 0x00}
		} else {
			leds[i] = color.RGBA{R: 0x00, G: 0xff, B: 0x00}
		}
	}
	ws.WriteColors(leds[:])

	const columncount = 14
	const rowcount = 5
	const settlingtime = 75 // this is the electrical settling time when selecting a row.  In uSec
	const scanningtime = 10 // this is the time between scans. in mSec

	// Key Matrix Definition - Columns
	columns := [columncount]machine.Pin{machine.GPIO0, machine.GPIO1, machine.GPIO2, machine.GPIO3,
		machine.GPIO4, machine.GPIO5, machine.GPIO6, machine.GPIO7,
		machine.GPIO8, machine.GPIO9, machine.GPIO10, machine.GPIO11,
		machine.GPIO12, machine.GPIO13}

	// Key Matrix Definition - Rows
	rows := [rowcount]machine.Pin{machine.GPIO14, machine.GPIO15, machine.GPIO16, machine.GPIO17, machine.GPIO18}

	for _, col := range columns {
		col.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	}

	for _, row := range rows {
		row.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	}

	// These are the trigger keys to select layers. - Order can be important
	specialkeys := []int{67, 28}

	keys := []keyboard.Keycode{
		keyboard.KeyEsc, keyboard.Key1, keyboard.Key2, keyboard.Key3, keyboard.Key4, keyboard.Key5, keyboard.Key6, keyboard.Key7, keyboard.Key8, keyboard.Key9, keyboard.Key0, keyboard.KeyMinus, keyboard.KeyEqual, keyboard.KeyBackspace,
		keyboard.KeyTab, keyboard.KeyQ, keyboard.KeyW, keyboard.KeyE, keyboard.KeyR, keyboard.KeyT, keyboard.KeyY, keyboard.KeyU, keyboard.KeyI, keyboard.KeyO, keyboard.KeyP, keyboard.KeyLeftBrace, keyboard.KeyRightBrace, keyboard.KeyBackslash,
		keyboard.ASCII00, keyboard.KeyA, keyboard.KeyS, keyboard.KeyD, keyboard.KeyF, keyboard.KeyG, keyboard.KeyH, keyboard.KeyJ, keyboard.KeyK, keyboard.KeyL, keyboard.KeySemicolon, keyboard.KeyQuote, keyboard.ASCII00, keyboard.KeyEnter,
		keyboard.KeyModifierLeftShift, keyboard.KeyZ, keyboard.KeyX, keyboard.KeyC, keyboard.KeyV, keyboard.KeyB, keyboard.KeyN, keyboard.KeyM, keyboard.KeyComma, keyboard.KeyPeriod, keyboard.KeySlash, keyboard.ASCII00, keyboard.ASCII00, keyboard.KeyRightShift,
		keyboard.KeyLeftCtrl, keyboard.KeyLeftGUI, keyboard.KeyLeftAlt, keyboard.ASCII00, keyboard.ASCII00, keyboard.KeySpace, keyboard.ASCII00, keyboard.ASCII00, keyboard.ASCII00, keyboard.KeyRightAlt, keyboard.KeyRightGUI, keyboard.ASCII00, keyboard.ASCII00, keyboard.KeyRightCtrl,

		keyboard.KeyEsc, keyboard.KeyF1, keyboard.KeyF2, keyboard.KeyF3, keyboard.KeyF4, keyboard.KeyF5, keyboard.KeyF6, keyboard.KeyF7, keyboard.KeyF8, keyboard.KeyF9, keyboard.KeyF10, keyboard.KeyF11, keyboard.KeyF12, keyboard.KeyDelete,
		keyboard.KeyTab, keyboard.KeyQ, keyboard.KeyW, keyboard.KeyE, keyboard.KeyR, keyboard.KeyT, keyboard.KeyY, keyboard.KeyU, keyboard.KeyI, keyboard.KeyO, keyboard.KeyP, keyboard.KeyLeftBrace, keyboard.KeyRightBrace, keyboard.KeyBackslash,
		keyboard.KeyCapsLock, keyboard.KeyA, keyboard.KeyS, keyboard.KeyD, keyboard.KeyF, keyboard.KeyG, keyboard.KeyH, keyboard.KeyJ, keyboard.KeyK, keyboard.KeyL, keyboard.KeySemicolon, keyboard.KeyQuote, keyboard.ASCII00, keyboard.KeyEnter,
		keyboard.KeyModifierLeftShift, keyboard.KeyZ, keyboard.KeyX, keyboard.KeyC, keyboard.KeyV, keyboard.KeyB, keyboard.KeyN, keyboard.KeyM, keyboard.KeyComma, keyboard.KeyPeriod, keyboard.KeySlash, keyboard.ASCII00, keyboard.ASCII00, keyboard.KeyRightShift,
		keyboard.KeyLeftCtrl, keyboard.KeyLeftGUI, keyboard.KeyLeftAlt, keyboard.ASCII00, keyboard.ASCII00, keyboard.KeySpace, keyboard.ASCII00, keyboard.ASCII00, keyboard.ASCII00, keyboard.KeyRightAlt, keyboard.KeyRightGUI, keyboard.ASCII00, keyboard.ASCII00, keyboard.KeyRightCtrl,

		keyboard.KeyTilde, keyboard.KeyF1, keyboard.KeyF2, keyboard.KeyF3, keyboard.KeyF4, keyboard.KeyF5, keyboard.KeyF6, keyboard.KeyF7, keyboard.KeyF8, keyboard.KeyF9, keyboard.KeyF10, keyboard.KeyF11, keyboard.KeyF12, keyboard.KeyBackspace,
		keyboard.KeyTab, keyboard.KeyQ, keyboard.KeyW, keyboard.KeyE, keyboard.KeyR, keyboard.KeyT, keyboard.KeyY, keyboard.KeyU, keyboard.KeyI, keyboard.KeyO, keyboard.KeyP, keyboard.KeyLeftBrace, keyboard.KeyRightBrace, keyboard.KeyBackslash,
		keyboard.ASCII00, keyboard.KeyA, keyboard.KeyS, keyboard.KeyD, keyboard.KeyF, keyboard.KeyG, keyboard.KeyH, keyboard.KeyJ, keyboard.KeyK, keyboard.KeyL, keyboard.KeySemicolon, keyboard.KeyQuote, keyboard.ASCII00, keyboard.KeyEnter,
		keyboard.KeyModifierLeftShift, keyboard.KeyZ, keyboard.KeyX, keyboard.KeyC, keyboard.KeyV, keyboard.KeyB, keyboard.KeyN, keyboard.KeyM, keyboard.KeyComma, keyboard.KeyPeriod, keyboard.KeySlash, keyboard.ASCII00, keyboard.ASCII00, keyboard.KeyUpArrow,
		keyboard.KeyLeftCtrl, keyboard.KeyLeftGUI, keyboard.KeyLeftAlt, keyboard.ASCII00, keyboard.ASCII00, keyboard.KeySpace, keyboard.ASCII00, keyboard.ASCII00, keyboard.ASCII00, keyboard.KeyRightAlt, keyboard.KeyLeftArrow, keyboard.KeyDownArrow, keyboard.ASCII00, keyboard.KeyRightArrow,
	}

	layer := 0
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
		for i, triggerkey := range specialkeys {
			if !(keystates[triggerkey]) {
				layer = i + 1
			}
		}

		// send keypresses
		for i, state := range keystates {
			if state {
				kb.Up(keys[i+layer*rowcount*columncount])
			} else {
				kb.Down(keys[i+layer*rowcount*columncount])
			}
		}

		// delay between scans
		time.Sleep(scanningtime * time.Millisecond)
	}
}
