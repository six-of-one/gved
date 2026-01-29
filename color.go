package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/lusingander/colorpicker"
)

var (
	defaultColor = color.NRGBA{0xff, 0x00, 0xff, 0xff}
	Color2 = color.NRGBA{0x00, 0x00, 0x00, 0xff}
)

var wcolp fyne.Window
var tapd *tappableDisplayColor
/*
func main() {
	a := app.New()
	w := a.NewWindow("color picker sample")
	ww = w
	w.Resize(fyne.NewSize(400, 400))
	w.SetContent(colorCont(w))
	w.ShowAndRun()
}*/

func colorCont(w fyne.Window) fyne.CanvasObject {
	var bas_col color.Color
	bas_col = defaultColor

	tappableDisplayColor := newTappableDisplayColor(w)
	tappableDisplayColor.setColor(bas_col)

	simpleDisplayColor := newSimpleDisplayColor()
	picker := colorpicker.New(200, colorpicker.StyleHueCircle)
	picker.SetOnChanged(func(c color.Color) {
		bas_col = c
		simpleDisplayColor.setColor(bas_col)
		tappableDisplayColor.setColor(bas_col)
	})
	content := container.NewWithoutLayout(picker)
	button := widget.NewButton("Set:", func() {
		picker.SetColor(bas_col)
		dialog.ShowCustom("Select color", "OK", content, w)
	})
//	simpleDisplayColor.setColor(bas_col)

	tapd = tappableDisplayColor

	return container.New(
		layout.NewHBoxLayout(),
		layout.NewSpacer(),
		container.New(
			layout.NewVBoxLayout(),
			layout.NewSpacer(),
/*			button,
			container.New(
				layout.NewHBoxLayout(),
				layout.NewSpacer(),
				simpleDisplayColor.label,
				simpleDisplayColor.rect,
				layout.NewSpacer(),
			),
			layout.NewSpacer(), */
//			widget.NewLabel("Or tap rectangle"),
			container.New(
				layout.NewHBoxLayout(),
				layout.NewSpacer(),
				button,
				tappableDisplayColor.label,
				tappableDisplayColor.rect,
				layout.NewSpacer(),
			),
			layout.NewSpacer(),
		),
		layout.NewSpacer(),
	)
}

type simpleDisplayColor struct {
	label *widget.Label
	rect  *canvas.Rectangle
}

func newSimpleDisplayColor() *simpleDisplayColor {
	selectColorCode := widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true})
	selectColorRect := &canvas.Rectangle{}
	selectColorRect.SetMinSize(fyne.NewSize(30, 20))
	return &simpleDisplayColor{
		label: selectColorCode,
		rect:  selectColorRect,
	}
}

func (c *simpleDisplayColor) setColor(clr color.Color) {
		til := fmt.Sprintf("Select: %02x",clr)
		wcolp.SetTitle(til)
}

type tappableDisplayColor struct {
	label *widget.Label
	rect  colorpicker.PickerOpenWidget
}

func newTappableDisplayColor(w fyne.Window) *tappableDisplayColor {
	selectColorCode := widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true})
	selectColorRect := colorpicker.NewColorSelectModalRect(w, fyne.NewSize(30, 20), Color2)
	d := &tappableDisplayColor{
		label: selectColorCode,
		rect:  selectColorRect,
	}
	selectColorRect.SetOnChange(d.setColor)
	return d
}

func (c *tappableDisplayColor) setColor(clr color.Color) {
		til := fmt.Sprintf("Select: %02x",clr)
		wcolp.SetTitle(til)
	c.label.SetText(hexColorString(clr))
	c.rect.SetColor(clr)
	c.rect.Refresh()
}

func hexColorString(c color.Color) string {
	rgba, _ := c.(color.NRGBA)
	return fmt.Sprintf("#%.2X%.2X%.2X%.2X", rgba.R, rgba.G, rgba.B, rgba.A)
}
