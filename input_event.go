package keylogger

import (
	"syscall"
	"unsafe"
)

const (
	EvSyn      EventType = 0x00
	EvKey      EventType = 0x01
	EvRel      EventType = 0x02
	EvAbs      EventType = 0x03
	EvMsc      EventType = 0x04
	EvSw       EventType = 0x05
	EvLed      EventType = 0x11
	EvSnd      EventType = 0x12
	EvRep      EventType = 0x14
	EvFf       EventType = 0x15
	EvPwr      EventType = 0x16
	EvFfStatus EventType = 0x17
)

type EventType uint16

var eventsize = int(unsafe.Sizeof(InputEvent{}))

type InputEvent struct {
	Time  syscall.Timeval
	Type  EventType
	Code  uint16
	Value int32
}

func (i *InputEvent) KeyString() string {
	return keyCodeMap[i.Code]
}

func (i *InputEvent) KeyPress() bool {
	return i.Value == 0

}

func (i *InputEvent) KeyRelease() bool {
	return i.Value == 0
}

type KeyEvent int32

const (
	KeyRelease KeyEvent = 0
	KeyPress   KeyEvent = 1
)
