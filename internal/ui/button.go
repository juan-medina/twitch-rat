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
	CONNECT_BUTTON ButtonId = iota
	DISCONNECT_BUTTON
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
	MAX_BUTTONS      = 10
	BUTTON_WIDTH     = 170
	BUTTON_HEIGHT    = 50
	BUTTON_GAP       = 10
	CLICK_SENT_DELAY = 200
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
	if x > b.x && x < b.x+BUTTON_WIDTH && y > b.y && y < b.y+BUTTON_HEIGHT {
		return true
	}
	return false
}

func (u *uiImpl) EnableButton(id ButtonId) {
	u.changeButtonState(id, buttonEnabled)
}

func (u *uiImpl) DisableButton(id ButtonId) {
	u.changeButtonState(id, buttonDisabled)
}

func (u *uiImpl) changeButtonState(id ButtonId, state buttonState) {
	for i, b := range u.buttons {
		if b.id == id {
			u.buttons[i].changeState(state)
			return
		}
	}
}

func (u *uiImpl) OnButtonClick(callback func(id ButtonId)) {
	if callback != nil {
		u.onButtonClick = callback
	} else {
		u.onButtonClick = dummyCallback
	}
}

func (u *uiImpl) updateButtons(mouseX, mouseY float64, leftPressed bool, elapsedTime int) {
	for i, b := range u.buttons {
		if b.visible {
			if b.state != buttonDisabled {
				if b.hit(mouseX, mouseY) {
					if leftPressed {
						if b.state != buttonPressed {
							if b.timeToSendClick == 0 {
								u.changeButtonState(b.id, buttonPressed)
								u.buttons[i].timeToSendClick = CLICK_SENT_DELAY
							}
						}
					} else {
						u.changeButtonState(b.id, buttonHover)
					}
				} else {
					u.changeButtonState(b.id, buttonEnabled)
				}
			}
			if b.timeToSendClick != 0 {
				u.buttons[i].timeToSendClick -= elapsedTime
				if b.timeToSendClick <= 0 {
					u.buttons[i].timeToSendClick = 0
					u.onButtonClick(b.id)
				}
			}
		}
	}
}

func (u *uiImpl) addButton(id ButtonId, x, y, w, h float64, label string, state buttonState) {
	u.buttons = append(u.buttons, newButton(id, x, y, w, h, label, u.normalFace, state))
}

func (u *uiImpl) SetButtonVisible(id ButtonId, visible bool) {
	for i, b := range u.buttons {
		if b.id == id {
			u.buttons[i].setVisible(visible)
			return
		}
	}
}

func newButton(id ButtonId, x, y, w, h float64, label string, face *text.GoTextFace, state buttonState) button {
	do := text.DrawOptions{}
	dx, dy := text.Measure(label, face, 0)
	tx := x + (w / 2) - (dx / 2)
	ty := y + (h / 2) - (dy / 2)

	do.GeoM.Reset()
	do.GeoM.Translate(tx, ty)

	b := button{
		id:      id,
		x:       x,
		y:       y,
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
