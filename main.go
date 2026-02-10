package main

// Copyright (c) 2026 Julian Müller (ChaoticByte)

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
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

// global vars

var config Config
var xdg_runtime_dir string

var handler func(e InputEvent)

//

func PrintHelp() {
	arg1 := os.Args[0]
	fmt.Print(arg1[strings.LastIndex(arg1, "/")+1:] + " <command> [args...]\n\n")
	fmt.Print("Commands:\n\n")
	fmt.Print("  help\n")
	fmt.Print("    Print this help text\n\n")
	fmt.Print("  control <config>\n")
	fmt.Print("    Read the <config>-file and start listening for global hotkeys\n\n")
	fmt.Print("  dumpkeys\n")
	fmt.Print("    Print all keypresses to stdout\n\n")
	fmt.Printf("keyvent %s\n", Version)
}


func EncodeCmd(cmd uint32) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint32(buf, 4) // encode length
	binary.LittleEndian.PutUint32(buf[4:8], cmd) // encode cmd
	return buf
}

func HandleDumpKey(e InputEvent) {
	if _, found := keyCodeMap[e.Code]; found {
		fmt.Printf("\nKey: %s, Code: %d\n", e.KeyString(), e.Code)
	}
}

func HandleControl(e InputEvent) {
	if _, found := keyCodeMap[e.Code]; !found {
		return
	}
	// get command from key input
	keybindFound := true
	var cmd uint32
	switch e.Code {
	case config.KeyBinds.StartOrSplit:
		cmd = CmdStartSplit
	case config.KeyBinds.StopOrReset:
		cmd = CmdStopReset
	case config.KeyBinds.Cancel:
		cmd = CmdCancel
	case config.KeyBinds.Unsplit:
		cmd = CmdUnsplit
	case config.KeyBinds.SkipSplit:
		cmd = CmdSkip
	case config.KeyBinds.CloseLibreSplit:
		cmd = CmdExit
	default:
		keybindFound = false
	}
	if !keybindFound { return }
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

func main() {

	// parse cli
	if len(os.Args) < 2 {
		PrintHelp()
		os.Exit(1)
	}
	if os.Args[1] == "help" || os.Args[1] == "-h" || os.Args[1] == "--help" {
		PrintHelp()
		os.Exit(0)
	}
	if os.Args[1] == "control" && len(os.Args) < 3 {
		fmt.Print("Missing argument <config>\n\n")
		PrintHelp()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "control":
		// read config
		configData, err := os.ReadFile(os.Args[2])
		if err != nil { panic(err) }
		err = json.Unmarshal(configData, &config)
		if err != nil { panic(err) }
		// get xdg_runtime_dir
		xdg_runtime_dir = os.Getenv("XDG_RUNTIME_DIR")
		if xdg_runtime_dir == "" {
			xdg_runtime_dir = fmt.Sprintf("/run/user/%d", os.Getuid())
		}
		xdg_runtime_dir = strings.TrimRight(xdg_runtime_dir, "/")
		// set handler
		handler = HandleControl
	case "dumpkeys":
		handler = HandleDumpKey
	default:
		fmt.Print("Unknown command!\n\n")
		PrintHelp()
		os.Exit(1)
	}

	// setup keyboard event listeners
	kbs := FindAllKeyboardDevices()
	for _, kb := range kbs {
		kl, err := New(kb)
		if err != nil {
			panic(err)
		}
		c := kl.Read()
		defer kl.Close()
		defer close(c)
		go func () {
			for {
				e := <- c
				if e.KeyPress() {
					handler(e)
				}
			}
		}()
	}

	for { // wait
		fmt.Scanln()
	}

}
