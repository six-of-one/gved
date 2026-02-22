package main

import (
	"fmt"
	"image/color"
	"image"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/lusingander/colorpicker"
	"github.com/disintegration/imaging"
)

// color fn() that were in render

type Color interface {
	IRGB() (irgb uint16)
}

type IRGB struct {
	irgb uint16
}

func (c IRGB) RGBA() (r, g, b, a uint32) {
	i := uint32(c.irgb&0xf000) >> 12
	r = uint32(c.irgb&0x0f00) >> 8 * i
	g = uint32(c.irgb&0x00f0) >> 4 * i
	b = uint32(c.irgb&0x000f) * i

	r = r << 8
	g = g << 8
	b = b << 8
	a = 0xffff

	return
}

// hex color triple, w/ possible alpha
// - and yes, you could just break down and insert color.RGBA{R: 205, G: 0, B: 205, A: 130}
// but i like this

type HColor interface {
	HRGB() (hrgb uint32)
}

type HRGB struct {
	hrgb uint32
}

func (c HRGB) RGBA() (r, g, b, a uint32) {
	a = uint32(c.hrgb&0xff000000) >> 24
	r = uint32(c.hrgb&0xff0000) >> 16
	g = uint32(c.hrgb&0x00ff00) >> 8
	b = uint32(c.hrgb&0x0000ff)

	r = r << 8
	g = g << 8
	b = b << 8
	if a == 0 { a = 0xff }	// an alpha of 0 seems to produce gray mush
	a = a << 8

	return
}

func irgb(c uint32) uint16 {

var ic uint16
	ic = 0

	a := (c & 0xff000000 >> 24) / 16
	r := (c & 0xff0000 >> 16) / 16
	g := (c & 0xff00 >> 8) / 16
	b := (c & 0xff) / 16

	ic = uint16(a * 0x1000 + r * 0x100 + g * 0x10 + b)
	return ic
}

// hue shift
// had to insert this from imaging source, somehow the github include doesnt... include it?? idk. i only work here

func AdjustHue(img image.Image, shift float64) *image.NRGBA {
	if math.Mod(shift, 360) == 0 {
		return imaging.Clone(img)
	}

	summand := shift / 360

	return imaging.AdjustFunc(img, func(c color.NRGBA) color.NRGBA {
		h, s, l := rgbToHSL(c.R, c.G, c.B)
		h += summand
		h = math.Mod(h, 1)
		//Adding 1 because Golang's Modulo function behaves differently to similar operators in most other languages.
		if h < 0 {
			h++
		}
		r, g, b := hslToRGB(h, s, l)
		return color.NRGBA{r, g, b, c.A}
	})
}

func hue(src *image.NRGBA, deg float64) *image.NRGBA {
	dst := AdjustHue(src, deg)
	return dst
}

// for maze output to se -- outputter is in pfrender
func ParseHexColor(s string) (c color.RGBA, err error) {
	c.A = 0xff
	switch len(s) {
	case 7:
		_, err = fmt.Sscanf(s, "#%02x%02x%02x", &c.R, &c.G, &c.B)
	case 4:
		_, err = fmt.Sscanf(s, "#%1x%1x%1x", &c.R, &c.G, &c.B)
		// Double the hex digits:
		c.R *= 17
		c.G *= 17
		c.B *= 17
	default:
		err = fmt.Errorf("invalid length, must be 7 or 4")

	}
	return
}

var (
	defaultColor = color.NRGBA{0xff, 0x00, 0xff, 0xff}
)

// vars edit colors

var lvl1col,lvl2col color.Color = HRGB{0xff0000ff},HRGB{0xffff0000}		// def level colors blue, red

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
	blot_hexent.Resize(fyne.Size{100, 36})
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
	col1_tap := newTappableDisplayColor(wn)
	col1_tap.setColor(lvl1col)

	col1_pick := colorpicker.New(200, colorpicker.StyleHueCircle)
	col1_pick.SetOnChanged(func(c color.Color) {
		lvl1col = c
		col1_tap.setColor(lvl1col)
	})
	col1_cont := container.NewWithoutLayout(col1_pick)
	col1_btn := widget.NewButton("level color 1 ", func() {
		col1_pick.SetColor(lvl1col)
		dialog.ShowCustom("Select color", "OK", col1_cont, wn)
	})

	col1_hexent := widget.NewEntry()
	col1_hexent.Resize(fyne.Size{100, 36})
	col1_hexent.SetText("FF00AAFF")
	col1_hexent.OnChanged = func(s string) {
		fmt.Sscanf(s,"%08x",&nc)
		lvl1col = HRGB{nc}
fmt.Printf("hex col: %v: %x - %s\n",lvl1col,nc,s)
		col1_tap.setColor(lvl1col)
		col1_tap.label.SetText(s)
	}
// level colors 2
	col2_tap := newTappableDisplayColor(wn)
	col2_tap.setColor(lvl2col)

	col2_pick := colorpicker.New(200, colorpicker.StyleHueCircle)
	col2_pick.SetOnChanged(func(c color.Color) {
		lvl2col = c
		col2_tap.setColor(lvl2col)
	})
	col2_cont := container.NewWithoutLayout(col2_pick)
	col2_btn := widget.NewButton("level color 2 ", func() {
		col2_pick.SetColor(lvl2col)
		dialog.ShowCustom("Select color", "OK", col2_cont, wn)
	})

	col2_hexent := widget.NewEntry()
	col2_hexent.Resize(fyne.Size{100, 36})
	col2_hexent.SetText("FF00AAFF")
	col2_hexent.OnChanged = func(s string) {
		fmt.Sscanf(s,"%08x",&nc)
		lvl2col = HRGB{nc}
fmt.Printf("hex col: %v: %x - %s\n",lvl2col,nc,s)
		col2_tap.setColor(lvl2col)
		col2_tap.label.SetText(s)
	}

// palette color testing

	cv1_label := widget.NewLabelWithStyle("0xFF00FF00", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true})
	cv1_rect := colorpicker.NewColorSelectModalRect(w, fyne.NewSize(30, 20), defaultColor)

	cv1_hexent := widget.NewEntry()
	cv1_hexent.Resize(fyne.Size{100, 36})
	cv1_hexent.SetText("FF0000FF")
	cv1_hexent.OnChanged = func(s string) {
		fmt.Sscanf(s,"%08x",&nc)
		hc := HRGB{nc}
		cv1_rect.SetColor(hc)
		ns := fmt.Sprintf("%04X",irgb(nc))
		cv1_label.SetText(ns)
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
				container.NewWithoutLayout(blot_hexent),
			),
			layout.NewSpacer(),
			container.New(
				layout.NewHBoxLayout(),
				col1_btn,
				col1_tap.label,
				col1_tap.rect,
				container.NewWithoutLayout(col1_hexent),
			),
			layout.NewSpacer(),
			container.New(
				layout.NewHBoxLayout(),
				col2_btn,
				col2_tap.label,
				col2_tap.rect,
				container.NewWithoutLayout(col2_hexent),
			),
			layout.NewSpacer(),
			container.New(
				layout.NewHBoxLayout(),
				cv1_label,
				cv1_rect,
				container.NewWithoutLayout(cv1_hexent),
			),
			layout.NewSpacer(),
		),
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
