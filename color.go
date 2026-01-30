package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/lusingander/colorpicker"
)

var (
	defaultColor = color.NRGBA{0xff, 0x00, 0xff, 0xff}
)

var wcolp fyne.Window

// vars edit colors

var lvl1col uint32
var lvl2col uint32

func colorCont(wn fyne.Window) fyne.CanvasObject {
	var bas_col color.Color
	bas_col = HRGB{blotcol}

	blot_tap := newTappableDisplayColor(wn)
	blot_tap.setColor(bas_col)

	blot_pick := colorpicker.New(200, colorpicker.StyleHueCircle)
	blot_pick.SetOnChanged(func(c color.Color) {
		bas_col = c
		blot_tap.setColor(bas_col)
	})
	blot_cont := container.NewWithoutLayout(blot_pick)
	blot_btn := widget.NewButton("blotter color", func() {
		blot_pick.SetColor(bas_col)
		dialog.ShowCustom("Select color", "OK", blot_cont, wn)
	})

	blot_hexent := widget.NewEntry()
	blot_hexent.SetText("FF00AAFF")
	var nc uint32
	blot_hexent.OnChanged = func(s string) {
		fmt.Sscanf(s,"%08x",&nc)
		bas_col = HRGB{nc}
fmt.Printf("hex col: %v: %x - %s\n",bas_col,nc,s)
		blot_tap.setColor(bas_col)
		blot_tap.label.SetText(s)
	}
// level colors 1
	bas_col = HRGB{lvl1col}
	col1_tap := newTappableDisplayColor(wn)
	col1_tap.setColor(bas_col)

	col1_pick := colorpicker.New(200, colorpicker.StyleHueCircle)
	col1_pick.SetOnChanged(func(c color.Color) {
		bas_col = c
		col1_tap.setColor(bas_col)
	})
	col1_cont := container.NewWithoutLayout(col1_pick)
	col1_btn := widget.NewButton("level color 1 ", func() {
		col1_pick.SetColor(bas_col)
		dialog.ShowCustom("Select color", "OK", col1_cont, wn)
	})

	col1_hexent := widget.NewEntry()
	col1_hexent.SetText("FF00AAFF")
	col1_hexent.OnChanged = func(s string) {
		fmt.Sscanf(s,"%08x",&nc)
		bas_col = HRGB{nc}
fmt.Printf("hex col: %v: %x - %s\n",bas_col,nc,s)
		col1_tap.setColor(bas_col)
		col1_tap.label.SetText(s)
	}
// level colors 2
	bas_col = HRGB{lvl2col}
	col2_tap := newTappableDisplayColor(wn)
	col2_tap.setColor(bas_col)

	col2_pick := colorpicker.New(200, colorpicker.StyleHueCircle)
	col2_pick.SetOnChanged(func(c color.Color) {
		bas_col = c
		col2_tap.setColor(bas_col)
	})
	col2_cont := container.NewWithoutLayout(col2_pick)
	col2_btn := widget.NewButton("level color 2 ", func() {
		col2_pick.SetColor(bas_col)
		dialog.ShowCustom("Select color", "OK", col2_cont, wn)
	})

	col2_hexent := widget.NewEntry()
	col2_hexent.SetText("FF00AAFF")
	col2_hexent.OnChanged = func(s string) {
		fmt.Sscanf(s,"%08x",&nc)
		bas_col = HRGB{nc}
fmt.Printf("hex col: %v: %x - %s\n",bas_col,nc,s)
		col2_tap.setColor(bas_col)
		col2_tap.label.SetText(s)
	}

	return container.New(
		layout.NewVBoxLayout(),
//		layout.NewSpacer(),
		container.New(
			layout.NewVBoxLayout(),
			layout.NewSpacer(),
			container.New(
				layout.NewHBoxLayout(),
				blot_btn,
				blot_tap.label,
				blot_tap.rect,
				blot_hexent,
			),
			layout.NewSpacer(),
			container.New(
				layout.NewHBoxLayout(),
				col1_btn,
				col1_tap.label,
				col1_tap.rect,
				col1_hexent,
			),
			layout.NewSpacer(),
			container.New(
				layout.NewHBoxLayout(),
				col2_btn,
				col2_tap.label,
				col2_tap.rect,
				col2_hexent,
			),
			layout.NewSpacer(),		),
		layout.NewSpacer(),
	)
}

// generics for ops

type tappableDisplayColor struct {
	label *widget.Label
	rect  colorpicker.PickerOpenWidget
}

func newTappableDisplayColor(w fyne.Window) *tappableDisplayColor {
	selectColorCode := widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true})
	selectColorRect := colorpicker.NewColorSelectModalRect(w, fyne.NewSize(30, 20), defaultColor)
	d := &tappableDisplayColor{
		label:  selectColorCode,
		rect:  selectColorRect,
	}
	selectColorRect.SetOnChange(d.setColor)
	return d
}

func (c *tappableDisplayColor) setColor(clr color.Color) {
		til := fmt.Sprintf("Select: %02X",clr)
		wcolp.SetTitle(til)
	c.label.SetText(hexColorString(clr))
	c.rect.SetColor(clr)
	c.rect.Refresh()
	ecolor = clr
	newcl := fmt.Sprintf("Color: %02X",clr)
	statlin(cmdhin,newcl)
}

func hexColorString(c color.Color) string {
	rgba, _ := c.(color.NRGBA)
	return fmt.Sprintf("#%.2X%.2X%.2X%.2X", rgba.A, rgba.R, rgba.G, rgba.B)
}
