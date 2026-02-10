package main

// Copyright (c) 2015 Marin (https://github.com/MarinX/keylogger)
// Copyright (c) 2026 Julian Müller (ChaoticByte)

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"strings"
	"syscall"
)

// system paths
const SysDeviceNamePath = "/sys/class/input/event%d/device/name"
const DevInputEventPath = "/dev/input/event%d"

// use lowercase names for devices, as we turn the device input name to lower case
var restrictedDevices = devices{"mouse"}
var allowedDevices = devices{"keyboard", "logitech mx keys"}

// KeyListener wrapper around file descriptior
type KeyListener struct {
	fd *os.File
}

type devices []string

func (d *devices) hasDevice(str string) bool {
	for _, device := range *d {
		if strings.Contains(str, device) {
			return true
		}
	}

	return false
}

// New creates a new keylogger for a device path
func New(devPath string) (*KeyListener, error) {
	k := &KeyListener{}
	fd, err := os.OpenFile(devPath, os.O_RDWR, os.ModeCharDevice)
	if err != nil {
		if os.IsPermission(err) && !k.IsRoot() {
			return nil, errors.New("permission denied. run with root permission or use a user with access to " + devPath)
		}
		return nil, err
	}
	k.fd = fd
	return k, nil
}

// Like FindKeyboardDevice, but finds all devices which contain keyword 'keyboard'
// Returns an array of file paths which contain keyboard events
func FindAllKeyboardDevices() []string {
	valid := make([]string, 0)

	for i := range 255 {
		buff, err := os.ReadFile(fmt.Sprintf(SysDeviceNamePath, i))

		// prevent from checking non-existant files
		if os.IsNotExist(err) {
			break
		}
		if err != nil {
			continue
		}

		deviceName := strings.ToLower(string(buff))

		if restrictedDevices.hasDevice(deviceName) {
			continue
		} else if allowedDevices.hasDevice(deviceName) {
			valid = append(valid, fmt.Sprintf(DevInputEventPath, i))
		}
	}
	return valid
}

// IsRoot checks if the process is run with root permission
func (k *KeyListener) IsRoot() bool {
	return syscall.Getuid() == 0 && syscall.Geteuid() == 0
}

// Read from file descriptor
// Blocking call, returns channel
// Make sure to close channel when finish
func (k *KeyListener) Read() chan InputEvent {
	event := make(chan InputEvent)
	go func(event chan InputEvent) {
		for {
			e, err := k.read()
			if err != nil {
				close(event)
				break
			}

			if e != nil {
				event <- *e
			}
		}
	}(event)
	return event
}

// read from file description and parse binary into go struct
func (k *KeyListener) read() (*InputEvent, error) {
	buffer := make([]byte, eventsize)
	n, err := k.fd.Read(buffer)
	if err != nil {
		return nil, err
	}
	// no input, dont send error
	if n <= 0 {
		return nil, nil
	}
	return k.eventFromBuffer(buffer)
}

// eventFromBuffer parser bytes into InputEvent struct
func (k *KeyListener) eventFromBuffer(buffer []byte) (*InputEvent, error) {
	event := &InputEvent{}
	err := binary.Read(bytes.NewBuffer(buffer), binary.LittleEndian, event)
	return event, err
}

// Close file descriptor
func (k *KeyListener) Close() error {
	if k.fd == nil {
		return nil
	}
	return k.fd.Close()
}
