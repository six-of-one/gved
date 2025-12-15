package main

import "C"

import (
	"fmt"
//	"regexp"
//	"strconv"
//	"strings"
	"os"
//	"bufio"
//	"time"

	"fyne.io/fyne/v2"
//    "fyne.io/fyne/v2/app"
//	"fyne.io/fyne/v2/dialog"
//    "fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
)

var w fyne.Window
var a fyne.App

func typedRune(r rune) {

//	fmt.Printf("in keys event - %x\n",r)
	if r == 'q' { os.Exit(0) }

if deskCanvas, ok := w.Canvas().(desktop.Canvas); ok {
        deskCanvas.SetOnKeyDown(func(key *fyne.KeyEvent) {
            fmt.Printf("Desktop key down: %h\n", key.Name)
        })
        deskCanvas.SetOnKeyUp(func(key *fyne.KeyEvent) {
//            fmt.Printf("Desktop key up: %v\n", key)
			if key.Name == "Escape" { os.Exit(0) }
       })
    }
}