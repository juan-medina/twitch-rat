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

package button

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/juan-medina/twitch-rat/internal/audio"
	"github.com/juan-medina/twitch-rat/internal/colors"
	"github.com/juan-medina/twitch-rat/internal/ui/label"
)

var (
	enabledColor  = colors.DarkPurple
	hoverColor    = colors.Purple
	pressedColor  = colors.LightPurple
	disabledColor = colors.DarkGray

	buttonEnabledTextColor  = colors.Yellow
	buttonDisabledTextColor = colors.Gray
)

type Id int
type Button interface {
	GetId() Id
	Draw(screen *ebiten.Image)
	Update(mouseX, mouseY float64, leftPressed bool, elapsedTime int)
	OnButtonClick(onButtonClick func(id Id))
	SetVisible(visible bool)
	ChangeState(state State)
	Move(x, y float64)
}

type State int

const (
	Disabled State = iota
	Enabled
	Hover
	Pressed
)

type button struct {
	id                  Id
	x, y, w, h          float64
	color               color.Color
	visible             bool
	label               label.Label
	state               State
	timeToSendClick     int
	buttonClickCallback func(id Id)
	audioPlayer         audio.Player
}

func dummyButtonCallback(id Id) {}

const (
	CLICK_SENT_DELAY = 200
	CLICK_SOUND      = "embed/sounds/click.ogg"
)

func (b button) GetId() Id {
	return b.id
}

func (b button) Draw(screen *ebiten.Image) {
	if b.visible {
		vector.DrawFilledRect(screen, float32(b.x), float32(b.y), float32(b.w), float32(b.h), b.color, true)
		b.label.Draw(screen)
	}
}

func (b *button) ChangeState(state State) {
	b.state = state
	textColor := buttonEnabledTextColor
	switch state {
	case Hover:
		b.color = hoverColor
	case Enabled:
		b.color = enabledColor
	case Pressed:
		b.color = pressedColor
	case Disabled:
		b.color = disabledColor
		textColor = buttonDisabledTextColor
	}
	b.label.SetColor(textColor)
}

func (b *button) SetVisible(visible bool) {
	b.visible = visible
	b.label.SetVisible(visible)
}

func (b button) hit(x, y float64) bool {
	if x > b.x && x < b.x+b.w && y > b.y && y < b.y+b.h {
		return true
	}
	return false
}

func (b *button) Move(x, y float64) {
	b.x = x
	b.y = y

	dx, dy := b.label.Measure()
	tx := x + (b.w / 2) - (dx / 2)
	ty := y + (b.h / 2) - (dy / 2)

	b.label.Move(tx, ty)
}

func (b *button) Update(mouseX, mouseY float64, leftPressed bool, elapsedTime int) {
	if b.visible {
		if b.state != Disabled {
			if b.hit(mouseX, mouseY) {
				if leftPressed {
					if b.state != Pressed {
						if b.timeToSendClick == 0 {
							b.ChangeState(Pressed)
							b.timeToSendClick = CLICK_SENT_DELAY
							b.audioPlayer.PlaySound(CLICK_SOUND)
						}
					}
				} else {
					b.ChangeState(Hover)
				}
			} else {
				b.ChangeState(Enabled)
			}
		}
		if b.timeToSendClick != 0 {
			b.timeToSendClick -= elapsedTime
			if b.timeToSendClick <= 0 {
				b.timeToSendClick = 0
				b.click()
			}
		}
	}
}

func (b *button) OnButtonClick(onButtonClick func(id Id)) {
	if onButtonClick != nil {
		b.buttonClickCallback = onButtonClick
	} else {
		b.buttonClickCallback = dummyButtonCallback
	}
}

func (b *button) click() {
	b.buttonClickCallback(b.id)
}

func New(id Id, w, h float64, text string, face *text.GoTextFace, audioPlayer audio.Player, state State) Button {
	audioPlayer.LoadSound(CLICK_SOUND)
	b := button{
		id:                  id,
		w:                   w,
		h:                   h,
		label:               label.New(0, text, face, buttonEnabledTextColor),
		visible:             false,
		buttonClickCallback: dummyButtonCallback,
		audioPlayer:         audioPlayer,
	}
	b.ChangeState(state)
	return &b
}
