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

type ChargePumpSetting byte

const (
	ChargePumpDisabled ChargePumpSetting = 0x14
	ChargePumpEnabled  ChargePumpSetting = 0x10
)

type MemoryAddressingMode byte

const (
	HorizontalAddressing MemoryAddressingMode = 0
	VerticalAddressing   MemoryAddressingMode = 1
	PageAddressing       MemoryAddressingMode = 2
)

type SegmentRemapMode byte

const (
	Map0ToSeg0   SegmentRemapMode = 0
	Map127ToSeg0 SegmentRemapMode = 1
)

type ComOutputScanDirection byte

const (
	ScanAscending  ComOutputScanDirection = 0xc0
	ScanDescending ComOutputScanDirection = 0xc8
)

type VcomhDeselectLevel byte

const (
	VccTimesPoint65 VcomhDeselectLevel = 0x00
	VccTimesPoint77 VcomhDeselectLevel = 0x20
	VccTimesPoint83 VcomhDeselectLevel = 0x30
)

type ComPinConfig byte

const (
	SequentialComPinConfig  ComPinConfig = 0
	AlternativeComPinConfig ComPinConfig = 16
)

type LeftRightRemap byte

const (
	DisableComLeftRightRemap LeftRightRemap = 0
	EnableComLeftRightRemap  LeftRightRemap = 32
)

type Display interface {
	Reset() error

	Invert() error
	Uninvert() error
	TurnOn() error
	TurnOff() error
	SetChargePump(setting ChargePumpSetting) error
	ConfigureClock(clkDivRatio byte, oscFreqSetting byte) error
	ConfigureComPinsHardware(pinConfig ComPinConfig, leftRightRemap LeftRightRemap) error
	SetMultiplexRatio(ratio byte) error
	SetOffset(offset byte) error
	SetStartLine(startLine byte) error
	SetMemoryAddressingMode(mode MemoryAddressingMode) error
	SetSegmentRemap(mode SegmentRemapMode) error
	SetComOutputScanDirection(direction ComOutputScanDirection) error
	SetContrast(contrast byte) error
	SetPrechargePeriod(phase1Ticks byte, phase2Ticks byte) error
	SetVcomhDeselectLevel(level VcomhDeselectLevel) error
	ForceEntireDisplayOn() error
	StopForcingEntireDisplayOn() error
}

type display struct {
	spi      spi.WritableDevice
	dcPin    gpio.ValueSetter
	resetPin gpio.ValueSetter
}

func NewDisplay(spi spi.WritableDevice, dcPin gpio.ValueSetter, resetPin gpio.ValueSetter) Display {
	return &display{spi, dcPin, resetPin}
}

func (disp *display) Reset() error {
	err := disp.resetPin.SetValue(gpio.High)
	if err != nil {
		return err
	}

	time.Sleep(3 * time.Microsecond)

	err = disp.resetPin.SetValue(gpio.Low)
	if err != nil {
		return err
	}

	time.Sleep(3 * time.Microsecond)

	err = disp.resetPin.SetValue(gpio.High)
	return err
}

func (disp *display) sendCommand(data []byte) error {
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

func (disp *display) Invert() error {
	return disp.sendCommand([]byte{0xA7})
}

func (disp *display) Uninvert() error {
	return disp.sendCommand([]byte{0xA6})
}

func (disp *display) TurnOff() error {
	return disp.sendCommand([]byte{0xAE})
}

func (disp *display) TurnOn() error {
	return disp.sendCommand([]byte{0xAF})
}

func (disp *display) SetChargePump(setting ChargePumpSetting) error {
	return disp.sendCommand([]byte{0x8D, byte(setting)})
}

func (disp *display) ConfigureClock(clkDivRatio byte, oscFreqSetting byte) error {
	packedValue := (oscFreqSetting << 4) | clkDivRatio
	return disp.sendCommand([]byte{0x8D, packedValue})
}

func (disp *display) ForceEntireDisplayOn() error {
	return disp.sendCommand([]byte{0xA4})
}

func (disp *display) StopForcingEntireDisplayOn() error {
	return disp.sendCommand([]byte{0xA5})
}

func (disp *display) SetComOutputScanDirection(direction ComOutputScanDirection) error {
	return disp.sendCommand([]byte{byte(direction)})
}

func (disp *display) ConfigureComPinsHardware(pinConfig ComPinConfig, leftRightRemap LeftRightRemap) error {
	packedValue := 2 | byte(pinConfig) | byte(leftRightRemap)
	return disp.sendCommand([]byte{0xDA, packedValue})
}

func (disp *display) SetContrast(contrast byte) error {
	return disp.sendCommand([]byte{0x81, contrast})
}

func (disp *display) SetOffset(offset byte) error {
	return disp.sendCommand([]byte{0xD3, offset})
}

func (disp *display) SetStartLine(startLine byte) error {
	packedValue := 0x40 | startLine
	return disp.sendCommand([]byte{packedValue})
}

func (disp *display) SetMemoryAddressingMode(mode MemoryAddressingMode) error {
	return disp.sendCommand([]byte{0x20, byte(mode)})
}

func (disp *display) SetMultiplexRatio(ratio byte) error {
	return disp.sendCommand([]byte{0xA8, ratio})
}

func (disp *display) SetPrechargePeriod(phase1Ticks byte, phase2Ticks byte) error {
	packedValue := (phase2Ticks << 4) | phase1Ticks
	return disp.sendCommand([]byte{0xD9, packedValue})
}

func (disp *display) SetSegmentRemap(mode SegmentRemapMode) error {
	packedValue := 0xA0 | byte(mode)
	return disp.sendCommand([]byte{packedValue})
}

func (disp *display) SetVcomhDeselectLevel(level VcomhDeselectLevel) error {
	return disp.sendCommand([]byte{0xDB, byte(level)})
}
