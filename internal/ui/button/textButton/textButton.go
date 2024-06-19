/*
 *  Copyright (c) 2024 Juan Medina
 *
 *   All rights reserved. This software and related documentation are proprietary to Juan Medina.
 *
 *   This source code is for internal use only and may not be copied, modified, or distributed
 *   without the express written permission of Juan Medina. Any use of this software for any
 *   purpose other than its intended use is strictly prohibited and may result in severe civil
 *   and criminal penalties.
 *
 *   THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED,
 *   INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR
 *   PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE
 *   FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR
 *   OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
 *   DEALINGS IN THE SOFTWARE.
 */

package textButton

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/juan-medina/twitch-rat/internal/audio"
	"github.com/juan-medina/twitch-rat/internal/colors"
	"github.com/juan-medina/twitch-rat/internal/draw"
	"github.com/juan-medina/twitch-rat/internal/ui/button"
	"github.com/juan-medina/twitch-rat/internal/ui/label"
)

var (
	textButtonEnabledColor  = colors.DarkPurple
	textButtonHoverColor    = colors.Purple
	textButtonPressedColor  = colors.Violet
	textButtonDisabledColor = colors.DarkGray

	textButtonEnabledTextColor  = colors.LightYellow
	textButtonDisabledTextColor = colors.Gray
)

type textButton struct {
	id                  button.Id
	x, y, w, h          float64
	color               color.Color
	visible             bool
	label               label.Label
	state               button.State
	timeToSendClick     int
	buttonClickCallback func(id button.Id)
	audioPlayer         audio.Player
}

func (b textButton) GetId() button.Id {
	return b.id
}

func (b textButton) Draw(screen *ebiten.Image) {
	if b.visible {
		vector.DrawFilledRect(screen, float32(b.x), float32(b.y), float32(b.w), float32(b.h), b.color, true)
		b.label.Draw(screen)
	}
}

func (b *textButton) ChangeState(state button.State) {
	b.state = state
	textColor := textButtonEnabledTextColor
	switch state {
	case button.Hover:
		b.color = textButtonHoverColor
	case button.Enabled:
		b.color = textButtonEnabledColor
	case button.Pressed:
		b.color = textButtonPressedColor
	case button.Disabled:
		b.color = textButtonDisabledColor
		textColor = textButtonDisabledTextColor
	}
	b.label.SetColor(textColor)
}

func (b *textButton) SetVisible(visible bool) {
	b.visible = visible
	b.label.SetVisible(visible)
}

func (b textButton) hit(x, y float64) bool {
	if x > b.x && x < b.x+b.w && y > b.y && y < b.y+b.h {
		return true
	}
	return false
}

func (b *textButton) Move(x, y float64) {
	b.x = x
	b.y = y

	dx, dy := b.label.Measure()
	tx := x + (b.w / 2) - (dx / 2)
	ty := y + (b.h / 2) - (dy / 2)

	b.label.Move(tx, ty)
}

func (b *textButton) Update(mouseX, mouseY float64, leftPressed bool, elapsedTime int) {
	if b.visible {
		if b.state != button.Disabled {
			if b.hit(mouseX, mouseY) {
				if leftPressed {
					b.Click()
				} else {
					ebiten.SetCursorShape(ebiten.CursorShapePointer)
					b.ChangeState(button.Hover)
				}
			} else {
				b.ChangeState(button.Enabled)
			}
		}
		if b.timeToSendClick != 0 {
			b.timeToSendClick -= elapsedTime
			if b.timeToSendClick <= 0 {
				b.timeToSendClick = 0
				b.buttonClickCallback(b.id)
			}
		}
	}
}

func (b *textButton) OnButtonClickCallback(onButtonClick func(id button.Id)) {
	if onButtonClick != nil {
		b.buttonClickCallback = onButtonClick
	} else {
		b.buttonClickCallback = dummyButtonCallback
	}
}

func (b *textButton) Click() {
	if b.state != button.Pressed {
		if b.timeToSendClick == 0 {
			b.ChangeState(button.Pressed)
			b.timeToSendClick = button.CLICK_SENT_DELAY
			b.audioPlayer.PlaySound(button.CLICK_SOUND)
		}
	}
}

func dummyButtonCallback(id button.Id) {}

func New(id button.Id, w, h float64, text string, font draw.Font, lineHeight float64, audioPlayer audio.Player, state button.State) button.Button {
	audioPlayer.LoadSound(button.CLICK_SOUND)
	b := textButton{
		id:                  id,
		w:                   w,
		h:                   h,
		label:               label.NewLabel(0, text, font, lineHeight, textButtonEnabledTextColor),
		visible:             false,
		buttonClickCallback: dummyButtonCallback,
		audioPlayer:         audioPlayer,
	}
	b.ChangeState(state)
	return &b
}
