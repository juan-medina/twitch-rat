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

package imageButton

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/juan-medina/twitch-rat/internal/audio"
	"github.com/juan-medina/twitch-rat/internal/colors"
	"github.com/juan-medina/twitch-rat/internal/draw"
	"github.com/juan-medina/twitch-rat/internal/keys"
	"github.com/juan-medina/twitch-rat/internal/ui/button"
)

var (
	imageButtonEnabledColor  = colors.Yellow
	imageButtonHoverColor    = colors.Purple
	imageButtonPressedColor  = colors.Violet
	imageButtonDisabledColor = colors.DarkGray
)

type imageButton struct {
	id                  button.Id
	x, y                float64
	color               color.Color
	visible             bool
	sprite              draw.Sprite
	state               button.State
	timeToSendClick     int
	buttonClickCallback func(id button.Id)
	audioPlayer         audio.Player
}

func (b imageButton) GetId() button.Id {
	return b.id
}

func (b imageButton) Draw(screen *ebiten.Image) {
	if b.visible {
		b.sprite.SetColor(b.color)
		b.sprite.Draw(screen, b.x, b.y, false, false)
	}
}

func (b *imageButton) ChangeState(state button.State) {
	b.state = state
	switch state {
	case button.Hover:
		b.color = imageButtonHoverColor
	case button.Enabled:
		b.color = imageButtonEnabledColor
	case button.Pressed:
		b.color = imageButtonPressedColor
	case button.Disabled:
		b.color = imageButtonDisabledColor
	}
}

func (b *imageButton) SetVisible(visible bool) {
	b.visible = visible

}

func (b imageButton) hit(x, y float64) bool {
	w, h := b.sprite.Size()
	if x > b.x && x < b.x+w && y > b.y && y < b.y+h {
		return true
	}
	return false
}

func (b *imageButton) Move(x, y float64) {
	b.x = x
	b.y = y
}

func (b *imageButton) Update(elapsedTime int, mouseX, mouseY float64, leftJustPressed bool, leftPressed bool, keys keys.Keys) {
	if b.visible {
		if b.state != button.Disabled {
			if b.hit(mouseX, mouseY) {
				if leftJustPressed {
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

func (b *imageButton) OnButtonClickCallback(onButtonClick func(id button.Id)) {
	if onButtonClick != nil {
		b.buttonClickCallback = onButtonClick
	} else {
		b.buttonClickCallback = dummyButtonCallback
	}
}

func (b *imageButton) Click() {
	if b.state != button.Pressed {
		if b.timeToSendClick == 0 {
			b.ChangeState(button.Pressed)
			b.timeToSendClick = button.CLICK_SENT_DELAY
			b.audioPlayer.PlaySound(button.CLICK_SOUND)
		}
	}
}

func (b imageButton) Size() (width float64, height float64) {
	return b.sprite.Size()
}

func dummyButtonCallback(id button.Id) {}

func New(id button.Id, sprite draw.Sprite, audioPlayer audio.Player, state button.State) button.Button {
	audioPlayer.LoadSound(button.CLICK_SOUND)
	b := imageButton{
		id:                  id,
		visible:             false,
		sprite:              sprite,
		buttonClickCallback: dummyButtonCallback,
		audioPlayer:         audioPlayer,
	}
	b.ChangeState(state)
	return &b
}
