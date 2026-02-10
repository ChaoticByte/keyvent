package main

// Copyright (c) 2015 Marin (https://github.com/MarinX/keylogger)
// Copyright (c) 2026 Julian Müller (ChaoticByte)

import (
	"syscall"
	"unsafe"
)

// eventsize is size of structure of InputEvent
var eventsize = int(unsafe.Sizeof(InputEvent{}))

// InputEvent is the keyboard event structure itself
type InputEvent struct {
	Time  syscall.Timeval
	Code  uint16
	Value int32
}

// KeyString returns representation of pressed key as string
// eg enter, space, a, b, c...
func (i *InputEvent) KeyString() string {
	return keyCodeMap[i.Code]
}

// KeyPress is the value when we press the key on keyboard
func (i *InputEvent) KeyPress() bool {
	return i.Value == 1
}

// KeyRelease is the value when we release the key on keyboard
func (i *InputEvent) KeyRelease() bool {
	return i.Value == 0
}

// Time of the keypress in nanoseconds
func (i *InputEvent) TimeNano() int64 {
	return i.Time.Nano()
}

// KeyEvent is the keyboard event for up/down (press/release)
type KeyEvent int32

const (
	KeyPress   KeyEvent = 1
	KeyRelease KeyEvent = 0
)
