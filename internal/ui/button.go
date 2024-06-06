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
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
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
	PLAY_BUTTON ButtonId = iota
	BACK_BUTTON
)

type button struct {
	id              ButtonId
	x, y, w, h      float64
	label           string
	color           color.Color
	visible         bool
	do              text.DrawOptions
	state           buttonState
	timeToSendClick int
	face            *text.GoTextFace
}

const (
	MAX_BUTTONS       = 10
	BUTTON_WIDTH      = 170
	BUTTON_HEIGHT     = 50
	BUTTON_GAP        = 10
	CLICK_SENT_DELAY  = 200
	BACK_BUTTON_WIDTH = 50
)

func (b button) draw(screen *ebiten.Image) {
	if b.visible {
		vector.DrawFilledRect(screen, float32(b.x), float32(b.y), float32(b.w), float32(b.h), b.color, true)
		text.Draw(screen, b.label, b.face, &b.do)
	}
}

func (b *button) changeState(state buttonState) {
	b.state = state
	textColor := buttonEnabledTextColor
	switch state {
	case buttonHover:
		b.color = hoverColor
	case buttonEnabled:
		b.color = enabledColor
	case buttonPressed:
		b.color = pressedColor
	case buttonDisabled:
		b.color = disabledColor
		textColor = buttonDisabledTextColor
	}

	b.do.ColorScale.Reset()
	b.do.ColorScale.ScaleWithColor(textColor)
}

func (b *button) setVisible(visible bool) {
	b.visible = visible
}

func (b button) hit(x, y float64) bool {
	if x > b.x && x < b.x+b.w && y > b.y && y < b.y+b.h {
		return true
	}
	return false
}

func (b *button) move(x, y float64) {
	b.x = x
	b.y = y

	dx, dy := text.Measure(b.label, b.face, 0)
	tx := x + (b.w / 2) - (dx / 2)
	ty := y + (b.h / 2) - (dy / 2)

	b.do.GeoM.Reset()
	b.do.GeoM.Translate(tx, ty)
}

func newButton(id ButtonId, w, h float64, label string, face *text.GoTextFace, state buttonState) button {
	do := text.DrawOptions{}
	b := button{
		id:      id,
		w:       w,
		h:       h,
		label:   label,
		visible: false,
		do:      do,
		face:    face,
	}
	b.changeState(state)
	return b
}
