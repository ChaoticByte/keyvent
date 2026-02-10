package main

// Copyright (c) 2026 Julian Müller (ChaoticByte)

import (
	"fmt"
	"os"
	"strings"
)

// handlers

var handler func(e InputEvent)

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
	// send command to libresplit
	SendCommand(cmd)
}

// cli

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
		ReadConfig()
		DetectXdgRuntimeDir()
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
				if e.KeyPress() { handler(e) }
			}
		}()
	}

	for { // wait
		fmt.Scanln()
	}

}
