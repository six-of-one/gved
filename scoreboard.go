package main

import (
//	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/data/binding"
)

// score board (seen in splash screen rotate) and gameplay info

// display scores while in game, and some status

var tsb *fyne.Container				// title small brand
var sbtl *fyne.Container			// title upd for score board
var scorec *fyne.Container			// scoreboard contains
var scors *fyne.Container			// image for scores

func dlg_scboard() {

	wc := a.NewWindow("High Score!")
	wc.Resize(fyne.NewSize(270, 600))
	tsb = container.NewStack()
	sbtl = container.NewStack()
	scorec = container.NewStack()
	scors = container.NewStack()
	txtB := binding.NewString()
	txtWid := widget.NewEntryWithData(txtB)

	osb := container.New(
		layout.NewVBoxLayout(),
		tsb,
		txtWid,
	)
	gif_lodr("splash/sanc_tsb.gif", tsb, sbtl, "")
	sbtl.Resize(fyne.NewSize(270, 120))
	sbtl.Refresh()
	wc.SetContent(osb)
	wc.Show()
}