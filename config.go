package main

// Copyright (c) 2026 Julian Müller (ChaoticByte)

type KeyBindsConfig struct {
	StartOrSplit uint16 `json:"start_or_split"`
	StopOrReset uint16 `json:"stop_or_reset"`
	Cancel uint16 `json:"cancel"`
	Unsplit uint16 `json:"unsplit"`
	SkipSplit uint16 `json:"skip_split"`
	CloseLibreSplit uint16 `json:"close_libresplit"`
}

type Config struct {
	KeyBinds KeyBindsConfig `json:"keybinds"`
}
