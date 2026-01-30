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

func colorCont(wn fyne.Window) fyne.CanvasObject {
	var bas_col color.Color
	bas_col = HRGB{blotcol}

	blot_tap := newTappableDisplayColor(wn)
	blot_tap.setColor(bas_col)

	picker := colorpicker.New(200, colorpicker.StyleHueCircle)
	picker.SetOnChanged(func(c color.Color) {
		bas_col = c
		blot_tap.setColor(bas_col)
	})
	content := container.NewWithoutLayout(picker)
	blot_btn := widget.NewButton("blotter color", func() {
		picker.SetColor(bas_col)
		dialog.ShowCustom("Select color", "OK", content, wn)
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
		),
		layout.NewSpacer(),
	)
}

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
	newcl := fmt.Sprintf("Blotter Color: %02X",clr)
	statlin(cmdhin,newcl)
}

func hexColorString(c color.Color) string {
	rgba, _ := c.(color.NRGBA)
	return fmt.Sprintf("#%.2X%.2X%.2X%.2X", rgba.A, rgba.R, rgba.G, rgba.B)
}
