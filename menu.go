package main

import (
	"image"
	"fmt"
	"math"
	"os"
	"io/ioutil"
	"time"
//		"image/color"

	"fyne.io/fyne/v2"
    "fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
    "fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/container"
)

// menu system isolated from keyboard & control now

