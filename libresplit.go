package main

// Copyright (c) 2026 Julian Müller (ChaoticByte)

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
)

const LibreSplitSocket = "libresplit.sock"

const (
	CmdStartSplit uint32 = 0
	CmdStopReset uint32 = 1
	CmdCancel uint32 = 2
	CmdUnsplit uint32 = 3
	CmdSkip uint32 = 4
	CmdExit uint32 = 5 // close LibreSplit
)

func EncodeCmd(cmd uint32) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint32(buf, 4) // encode length
	binary.LittleEndian.PutUint32(buf[4:8], cmd) // encode cmd
	return buf
}

func SendCommand(cmd uint32) {
	// try opening socket
	var err error
	sock, err := net.Dial("unix", xdg_runtime_dir + "/" + LibreSplitSocket)
	if err != nil {
		fmt.Println("Could not connect to libresplit: ", err)
		return
	}
	// write to socket
	msg := EncodeCmd(cmd)
	n, err := sock.Write(msg)
	// fmt.Println(msg, n, err)
	sock.Close()
	if err != nil || n != len(msg) {
		if err == nil {
			err = errors.New("wrong amout of bytes was written")
		}
		fmt.Println("Could not communicate with libresplit socket: ", err)
	}
}
