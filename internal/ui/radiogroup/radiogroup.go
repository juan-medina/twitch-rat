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

package radiogroup

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/juan-medina/twitch-rat/internal/audio"
	"github.com/juan-medina/twitch-rat/internal/colors"
	"github.com/juan-medina/twitch-rat/internal/draw"
	"github.com/juan-medina/twitch-rat/internal/keys"
)

var (
	radioGroupEnabledColor  = colors.PalePurple
	radioGroupHoverColor    = colors.Purple
	radioGroupPressedColor  = colors.DarkGray
	radioGroupSelectedColor = colors.Violet

	radioGroupTextColor         = colors.DarkGold
	radioGroupSelectedTextColor = colors.LightYellow
	radioGroupBackground        = colors.DarkPurple
)

const (
	BORDER_SIZE      = 4
	CLICK_SENT_DELAY = 200
	CLICK_SOUND      = "embed/sounds/click.ogg"
)

type Id int

type RadioGroup interface {
	Draw(screen *ebiten.Image)
	Update(elapsedTime int, mouseX float64, mouseY float64, leftJustPressed bool, leftPressed bool, keys keys.Keys)
	Move(x, y float64)
	GetId() Id
	SetVisible(visible bool)
	Measure() (float64, float64)
	SetSelected(index int)
	GetSelected() int
	OnChange(callback func(Id, int))
}

type radioState int

const (
	Enabled radioState = iota
	Hover
	Pressed
	Selected
)

type radio struct {
	string
	state      radioState
	x, y, w, h float64
	tx, ty     float64
}

type radioGroupImpl struct {
	id             Id
	x, y, w, h     float64
	font           draw.Font
	visible        bool
	radios         []radio
	selected       int
	audioPlayer    audio.Player
	changeCallback func(Id, int)
}

func (r *radioGroupImpl) Draw(screen *ebiten.Image) {
	if !r.visible {
		return
	}
	vector.DrawFilledRect(screen, float32(r.x), float32(r.y), float32(r.w), float32(r.h), radioGroupBackground, false)
	for i, rad := range r.radios {
		var color colors.CustomColor
		switch rad.state {
		case Enabled:
			color = radioGroupEnabledColor
		case Hover:
			color = radioGroupHoverColor
		case Pressed:
			color = radioGroupPressedColor
		case Selected:
			color = radioGroupSelectedColor
		}
		vector.DrawFilledRect(screen, float32(rad.x+r.x), float32(rad.y+r.y), float32(rad.w), float32(rad.h), color, false)

		if i == r.selected {
			color = radioGroupSelectedTextColor
		} else {
			color = radioGroupTextColor
		}
		r.font.Draw(screen, rad.string, rad.tx+r.x, rad.ty+r.y, r.font.DefaultSize(), color)
	}
}

func (r *radioGroupImpl) Update(elapsedTime int, mouseX float64, mouseY float64, leftJustPressed bool, leftPressed bool, keys keys.Keys) {
	if !r.visible {
		return
	}

	for i, rad := range r.radios {
		hit := mouseX > rad.x+r.x && mouseX < rad.x+r.x+rad.w && mouseY > rad.y+r.y && mouseY < rad.y+r.y+rad.h
		selected := i == r.selected
		if hit {
			ebiten.SetCursorShape(ebiten.CursorShapePointer)
			if leftJustPressed && rad.state != Pressed && !selected {
				r.audioPlayer.PlaySound(CLICK_SOUND)
				r.SetSelected(i)
				if r.changeCallback != nil {
					r.changeCallback(r.id, i)
				}
				return
			} else {
				r.radios[i].state = Hover
			}
		} else {
			if selected {
				r.radios[i].state = Selected
			} else {
				r.radios[i].state = Enabled
			}
		}
	}
}

func (r *radioGroupImpl) Move(x float64, y float64) {
	r.x = x
	r.y = y
}

func (r *radioGroupImpl) GetId() Id {
	return r.id
}

func (r *radioGroupImpl) SetVisible(visible bool) {
	r.visible = visible
}

func (r *radioGroupImpl) SetSelected(index int) {
	if index < 0 || index >= len(r.radios) {
		return
	}
	for i := range r.radios {
		if i == index {
			r.radios[i].state = Selected
		} else {
			r.radios[i].state = Enabled
		}
	}
	r.selected = index
}

func (r *radioGroupImpl) GetSelected() int {
	return r.selected
}

func (r radioGroupImpl) Measure() (float64, float64) {
	return r.w, r.h
}

func (r *radioGroupImpl) OnChange(callback func(Id, int)) {
	r.changeCallback = callback
}

func New(id Id, w, h float64, font draw.Font, audioPlayer audio.Player, options ...string) RadioGroup {
	totalOptions := len(options)
	radios := make([]radio, totalOptions)

	radioHeight := h - 2*BORDER_SIZE
	radioWidth := (w - float64((1+totalOptions)*BORDER_SIZE)) / float64(totalOptions)
	startX := float64(BORDER_SIZE)

	for i, option := range options {
		tw, th := font.Measure(option, font.DefaultSize())
		tx := startX + (radioWidth / 2) - (tw / 2)
		ty := BORDER_SIZE + (radioHeight / 2) - (th / 2)
		radios[i] = radio{
			string: option,
			state:  Enabled,
			x:      startX,
			y:      BORDER_SIZE,
			w:      radioWidth,
			h:      radioHeight,
			tx:     tx,
			ty:     ty,
		}
		startX += radioWidth + float64(BORDER_SIZE)
	}

	r := radioGroupImpl{
		id:          id,
		w:           w,
		h:           h,
		font:        font,
		radios:      radios,
		audioPlayer: audioPlayer,
	}
	r.SetSelected(0)
	return &r
}
