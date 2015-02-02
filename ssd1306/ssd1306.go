// Package ssd1306 is a driver for SSD1306 monochrome OLED display modules.
//
// It works in terms of GPIO and SPI implementations provided elsewhere,
// implementing the go-gpio and go-spi interfaces. To use this driver on
// an embedded Linux system you may be able to use the GPIO and SPI
// implementations in go-linuxgpio and go-linuxspi, as long as your kernel
// has drivers that expose the hardware GPIO and SPI pins to userspace.
package ssd1306

import (
	"fmt"
	"github.com/apparentlymart/go-gpio/gpio"
	"github.com/apparentlymart/go-spi/spi"
	"time"
)

type Display interface {
	Reset() error

	InvertDisplay() error
	NormalDisplay() error
	DisplayOn() error
	DisplayOff() error
}

type display struct {
	spi spi.WritableDevice
	dcPin gpio.ValueSetter
	resetPin gpio.ValueSetter
}

func NewDisplay(spi spi.WritableDevice, dcPin gpio.ValueSetter, resetPin gpio.ValueSetter) (Display) {
	return &display{spi, dcPin, resetPin};
}

func (disp display) Reset() error {
	err := disp.resetPin.SetValue(gpio.High)
	if (err != nil) {
		return err
	}

	time.Sleep(3 * time.Microsecond)

	err = disp.resetPin.SetValue(gpio.Low)
	if (err != nil) {
		return err
	}

	time.Sleep(3 * time.Microsecond)

	err = disp.resetPin.SetValue(gpio.High)
	return err
}

func (disp display) sendCommand(data []byte) error {
	disp.dcPin.SetValue(gpio.Low)
	n, err := disp.spi.Write(data)
	if err != nil {
		return err
	}
	if n != len(data) {
		return fmt.Errorf("Short write")
	}
	//disp.dcPin.SetValue(gpio.High)
	return nil
}

func (disp display) InvertDisplay() error {
	return disp.sendCommand([]byte{0xA7})
}

func (disp display) NormalDisplay() error {
	return disp.sendCommand([]byte{0xA6})
}

func (disp display) DisplayOff() error {
	return disp.sendCommand([]byte{0xAE})
}

func (disp display) DisplayOn() error {
	return disp.sendCommand([]byte{0xAF})
}
