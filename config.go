package main

// Copyright (c) 2026 Julian Müller (ChaoticByte)

import (
	"encoding/json"
	"os"
)

type KeyBindsConfig struct {
	StartOrSplit uint16 `json:"start_or_split"`
	StopOrReset uint16 `json:"stop_or_reset"`
	Cancel uint16 `json:"cancel"`
	Unsplit uint16 `json:"unsplit"`
	SkipSplit uint16 `json:"skip_split"`
	CloseLibreSplit uint16 `json:"close_libresplit"`
}

var config struct {
	KeyBinds KeyBindsConfig `json:"keybinds"`
}

func ReadConfig() {
	// read config
	configData, err := os.ReadFile(os.Args[2])
	if err != nil { panic(err) }
	err = json.Unmarshal(configData, &config)
	if err != nil { panic(err) }
}
