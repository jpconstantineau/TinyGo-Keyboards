package main

import (
	"machine"
	"time"
)

func main() {
	led1 := machine.LED1
	led2 := machine.LED2
	led3 := machine.LED3
	led4 := machine.LED4
	led1.Configure(machine.PinConfig{Mode: machine.PinOutput})
	led2.Configure(machine.PinConfig{Mode: machine.PinOutput})
	led3.Configure(machine.PinConfig{Mode: machine.PinOutput})
	led4.Configure(machine.PinConfig{Mode: machine.PinOutput})

	for {
		led1.Low()
		time.Sleep(500 * time.Millisecond)
		led1.High()
		time.Sleep(500 * time.Millisecond)
		led2.Low()
		time.Sleep(500 * time.Millisecond)
		led2.High()
		time.Sleep(500 * time.Millisecond)
		led3.Low()
		time.Sleep(500 * time.Millisecond)
		led3.High()
		time.Sleep(500 * time.Millisecond)
		led4.Low()
		time.Sleep(500 * time.Millisecond)
		led4.High()
		time.Sleep(500 * time.Millisecond)

	}
}
