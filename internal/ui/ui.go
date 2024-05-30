/*
 * Copyright (c) 2023 Juan Antonio Medina Iglesias
 *
 *  Permission is hereby granted, free of charge, to any person obtaining a copy
 *  of this software and associated documentation files (the "Software"), to deal
 *  in the Software without restriction, including without limitation the rights
 *  to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *  copies of the Software, and to permit persons to whom the Software is
 *  furnished to do so, subject to the following conditions:
 *
 *  The above copyright notice and this permission notice shall be included in
 *  all copies or substantial portions of the Software.
 *
 *  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *  IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *  FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *  AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *  LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 *  THE SOFTWARE.
 */

package ui

import (
	"bytes"
	"embed"
	"image/color"
	"io/fs"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	yellow      = color.RGBA64{0xFFFF, 0xFFFF, 0x0000, 0xFFFF}
	lightPurple = color.RGBA64{0xFFFF, 0x0000, 0xFFFF, 0xFFFF}
	darkPurple  = color.RGBA64{0x6666, 0x0000, 0x6666, 0xFFFF}
	purple      = color.RGBA64{0x8888, 0x0000, 0x8888, 0x0000}
	gray        = color.RGBA64{0x6666, 0x6666, 0x6666, 0x0000}
	darkGray    = color.RGBA64{0x3333, 0x3333, 0x3333, 0x0000}

	enabledColor  = darkPurple
	hoverColor    = purple
	pressedColor  = lightPurple
	disabledColor = darkGray

	buttonEnabledTextColor  = yellow
	buttonDisabledTextColor = gray
)

type buttonState int

const (
	buttonDisabled buttonState = iota
	buttonEnabled
	buttonHover
	buttonPressed
)

type ButtonId int

const (
	CONNECT_BUTTON ButtonId = iota
)

type button struct {
	id         ButtonId
	x, y, w, h float64
	label      string
	color      color.Color
	visible    bool
	do         text.DrawOptions
	state      buttonState
}

const (
	MAX_BUTTONS   = 10
	BUTTON_WIDTH  = 150
	BUTTON_HEIGHT = 50
)

type UI interface {
	Init(fileSystem embed.FS, screenWidth, screenHeight int)
	Update()
	Draw(screen *ebiten.Image)
	SetStatusMessage(message string)
	OnButtonClick(callback func(id ButtonId))
	DisableButton(id ButtonId)
}

type ui struct {
	screenWidth   int
	screenHeight  int
	fileSystem    embed.FS
	faceSource    *text.GoTextFaceSource
	normalFace    *text.GoTextFace
	lastMessage   string
	lastMessageDO text.DrawOptions
	buttons       []button
	onButtonClick func(id ButtonId)
}

// init implements UI.
func (u *ui) Init(fileSystem embed.FS, width int, height int) {
	u.fileSystem = fileSystem
	u.screenWidth = width
	u.screenHeight = height

	fontBytes, err := fs.ReadFile(u.fileSystem, "embed/fonts/default.ttf")
	if err != nil {
		log.Fatal(err)
	}

	s, err := text.NewGoTextFaceSource(bytes.NewReader(fontBytes))
	if err != nil {
		log.Fatal(err)
	}

	u.faceSource = s

	u.normalFace = &text.GoTextFace{
		Source: u.faceSource,
		Size:   24,
	}

	u.lastMessageDO.GeoM.Reset()
	gapX := u.normalFace.Size
	gapY := u.normalFace.Size * 1.5
	u.lastMessageDO.GeoM.Translate(gapX, float64(u.screenHeight)-gapY)
	u.lastMessage = "Loading..."

	u.buttons = make([]button, 0, MAX_BUTTONS)

	bx := float64((u.screenWidth / 2) - (BUTTON_WIDTH / 2))
	by := float64((u.screenHeight / 2) - (BUTTON_HEIGHT / 2))
	u.addButton(CONNECT_BUTTON, bx, by, BUTTON_WIDTH, BUTTON_HEIGHT, "CONNECT")
}

// Draw implements UI.
func (u *ui) Draw(screen *ebiten.Image) {
	text.Draw(screen, u.lastMessage, u.normalFace, &u.lastMessageDO)
	for _, b := range u.buttons {
		if b.visible {
			vector.DrawFilledRect(screen, float32(b.x), float32(b.y), float32(b.w), float32(b.h), b.color, false)
			text.Draw(screen, b.label, u.normalFace, &b.do)
		}
	}
}

// Update implements UI.
func (u *ui) Update() {
	x, y := ebiten.CursorPosition()
	pressed := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	for _, b := range u.buttons {
		if b.visible {
			if b.state != buttonDisabled {
				if u.buttonHit(b, float64(x), float64(y)) {
					if pressed {
						if b.state != buttonPressed {
							u.changeButtonState(b.id, buttonDisabled)
							u.onButtonClick(b.id)
						}
					} else {
						u.changeButtonState(b.id, buttonHover)
					}
				} else {
					u.changeButtonState(b.id, buttonEnabled)
				}
			}
		}
	}
}

func (u *ui) SetStatusMessage(message string) {
	u.lastMessage = message
}

func (u *ui) addButton(id ButtonId, x, y, w, h float64, label string) {
	do := text.DrawOptions{}
	dx, dy := text.Measure(label, u.normalFace, 0)
	tx := x + (w / 2) - (dx / 2)
	ty := y + (h / 2) - (dy / 2)

	do.GeoM.Reset()
	do.GeoM.Translate(tx, ty)

	u.buttons = append(u.buttons, button{
		id:      id,
		x:       x,
		y:       y,
		w:       w,
		h:       h,
		label:   label,
		visible: true,
		do:      do,
	})

	u.changeButtonState(id, buttonEnabled)
}

func (u ui) buttonHit(b button, x, y float64) bool {
	if x > b.x && x < b.x+BUTTON_WIDTH && y > b.y && y < b.y+BUTTON_HEIGHT {
		return true
	}
	return false
}

func (u *ui) OnButtonClick(callback func(id ButtonId)) {
	u.onButtonClick = callback
}

// DisableButton implements UI.
func (u *ui) DisableButton(id ButtonId) {
	u.changeButtonState(id, buttonDisabled)
}

func (u *ui) changeButtonState(id ButtonId, state buttonState) {
	for i, b := range u.buttons {
		if b.id == id {
			u.buttons[i].state = state
			textColor := buttonEnabledTextColor
			switch state {
			case buttonHover:
				u.buttons[i].color = hoverColor
			case buttonEnabled:
				u.buttons[i].color = enabledColor
			case buttonPressed:
				u.buttons[i].color = pressedColor
			case buttonDisabled:
				u.buttons[i].color = disabledColor
				textColor = buttonDisabledTextColor
			}

			u.buttons[i].do.ColorScale.Reset()
			r := float32(textColor.R / 255.0)
			g := float32(textColor.G / 255.0)
			b := float32(textColor.B / 255.0)
			a := float32(textColor.A / 255.0)
			u.buttons[i].do.ColorScale.Scale(r, g, b, a)

			return
		}
	}

}

func New() UI {
	return &ui{
		onButtonClick: func(id ButtonId) {},
	}
}
