
# (WIP) keyVent - LibreSplit Global Hotkeys

> [!WARN]  
> This is a work in progress!

> [!NOTE]  
> Only single-key hotkeys are supported.

This program uses udev to enable display server independent global hotkeys for [LibreSplit](https://libresplit.org/) on Linux.

It detects keyboard input devices in `/dev/input/event*` and reads key presses from there.
The program then controls LiveSplit through it's libresplit-ctl unix domain socket.

> [!NOTE]  
> This software is heavily based on https://github.com/MarinX/keylogger


## Compatibility

This program was tested on

- Debian 13 / Wayland


## Config

Example config:

```json
{
    "keybinds": {
        "start_or_split": 82,
        "stop_or_reset": 83,
        "unsplit": 75,
        "skip_split": 77,
        "cancel": 72,
        "close_libresplit": 0
    }
}
```

- Start/Split: Keypad 0
- Stop/Reset: Keypad Dot
- Undo Split: Keypad 4
- Skip Split: Keypad 6
- Cancel Run: Keypad 8
- Close Libresplit: disabled

The keybinds must be valid keycodes. You can run `./keyvent dumpkeys` to print all incoming key presses to the terminal, or have a look at [keymap.go](./keymap.go)   
Use `0` to disable the key bind.


## Usage

...
