package main

import (
	"machine"
	"machine/usb/hid/keyboard"
	"time"
)

func main() {
	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})

	buttons := []machine.Pin{machine.GPIO3, machine.GPIO4, machine.GPIO5, machine.GPIO6,
		machine.GPIO7, machine.GPIO8, machine.GPIO9, machine.GPIO10,
		machine.GPIO15, machine.GPIO16, machine.GPIO17, machine.GPIO18,
		machine.GPIO19, machine.GPIO20, machine.GPIO21, machine.GPIO22}

	keys := []keyboard.Keycode{keyboard.KeyA, keyboard.KeyB, keyboard.KeyC, keyboard.KeyD,
		keyboard.KeyE, keyboard.KeyF, keyboard.KeyG, keyboard.KeyH,
		keyboard.KeyI, keyboard.KeyJ, keyboard.KeyK, keyboard.KeyL,
		keyboard.KeyM, keyboard.KeyN, keyboard.KeyO, keyboard.KeyP}

	for _, button := range buttons {
		button.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	}

	kb := keyboard.Port()

	for {
		for i, button := range buttons {
			if button.Get() {
				kb.Up(keys[i])
			} else {
				kb.Down(keys[i])
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
}
