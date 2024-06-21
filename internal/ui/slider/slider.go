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

package slider

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/juan-medina/twitch-rat/internal/colors"
	"github.com/juan-medina/twitch-rat/internal/draw"
	"github.com/juan-medina/twitch-rat/internal/keys"
	"github.com/juan-medina/twitch-rat/internal/ui/label"
)

type Id int

type Slider interface {
	GetId() Id
	SetVisible(visible bool)
	Update(elapsedTime int, mouseX, mouseY float64, leftJustPressed bool, leftPressed bool, keys keys.Keys)
	Draw(screen *ebiten.Image)
	Move(x, y float64)
	SetValue(value float64)
	OnValueChangeCallback(onValueChange func(id Id, value float64))
}

var (
	sliderColor        = colors.DarkPurple
	markerNormalColor  = colors.DarkYellow
	markerHoverColor   = colors.Yellow
	markerPressedColor = colors.LightYellow
)

type status int

const (
	released status = iota
	pressed
)

type sliderImpl struct {
	id                  Id
	visible             bool
	color               color.Color
	x, y, w, h          float64
	valueLabel          label.Label
	value               float64
	markerWidth         float64
	markerHeight        float64
	markerX             float64
	markerY             float64
	markerColor         color.Color
	status              status
	valueChangeCallback func(id Id, value float64)
}

func dummyValueChangeCallback(id Id, value float64) {}

func (s *sliderImpl) GetId() Id {
	return s.id
}

func (s *sliderImpl) SetVisible(visible bool) {
	s.visible = visible
	s.valueLabel.SetVisible(visible)
}

func (s *sliderImpl) Update(elapsedTime int, mouseX, mouseY float64, leftJustPressed bool, leftPressed bool, keys keys.Keys) {
	if !s.visible {
		return
	}

	if s.isHit(mouseX, mouseY) {
		ebiten.SetCursorShape(ebiten.CursorShapePointer)
	}
	currentValue := s.value
	hit := s.isMakerHit(mouseX, mouseY)

	switch s.status {
	case released:
		if leftJustPressed {
			if hit {
				s.markerColor = markerPressedColor
				s.status = pressed
			} else {
				s.markerColor = markerNormalColor
				if s.isHit(mouseX, mouseY) {
					s.status = pressed
					s.markerColor = markerPressedColor
				}
			}
		} else {
			if hit {
				s.markerColor = markerHoverColor
			} else {
				s.markerColor = markerNormalColor
			}
		}
	case pressed:
		if !leftJustPressed && !leftPressed {
			if hit {
				s.markerColor = markerHoverColor
			} else {
				s.markerColor = markerNormalColor
			}
			s.status = released
		} else {
			s.markerColor = markerPressedColor
		}
	}
	if s.status == pressed {
		value := float64(mouseX-s.x) / float64(s.w)
		if value < 0 {
			value = 0
		}
		if value > 1 {
			value = 1
		}
		s.value = value
		if currentValue != value {
			s.valueChangeCallback(s.id, value)
			s.SetValue(value)
		}
	}

}

func (s sliderImpl) isMakerHit(x, y float64) bool {
	return x >= s.markerX && x <= s.markerX+s.markerWidth && y >= s.markerY && y <= s.markerY+s.markerHeight
}

func (s sliderImpl) isHit(x, y float64) bool {
	return x >= s.x && x <= s.x+s.w && y >= s.y && y <= s.y+s.h
}

func (s sliderImpl) Draw(screen *ebiten.Image) {
	if !s.visible {
		return
	}
	vector.DrawFilledRect(screen, float32(s.x), float32(s.y), float32(s.w), float32(s.h), s.color, true)
	vector.DrawFilledRect(screen, float32(s.markerX), float32(s.markerY), float32(s.markerWidth), float32(s.markerHeight), s.markerColor, true)
	s.valueLabel.Draw(screen)
}

func (s *sliderImpl) Move(x, y float64) {
	s.x = x
	s.y = y

	_, lh := s.valueLabel.Measure()

	cx := s.x + s.w
	cy := s.y + ((s.h - lh) / 2)
	s.valueLabel.Move(cx, cy)

	s.markerY = s.y - ((s.markerHeight - s.h) / 2)
	s.SetValue(s.value)
}

func (s *sliderImpl) SetValue(value float64) {
	s.value = value
	s.valueLabel.SetText(fmt.Sprintf("  %d%%", int(value*100)))
	s.markerX = s.x + s.w*value - (s.markerWidth / 2)
}

func (s *sliderImpl) OnValueChangeCallback(onValueChange func(id Id, value float64)) {
	if onValueChange == nil {
		s.valueChangeCallback = dummyValueChangeCallback
	} else {
		s.valueChangeCallback = onValueChange
	}

}

func New(id Id, w, h float64, font draw.Font, lineHeight float64, labelColor color.Color) Slider {
	markerWidth := w * 0.05
	markerHeight := h + (h * 0.40)

	return &sliderImpl{
		id:                  id,
		w:                   w,
		h:                   h,
		color:               sliderColor,
		valueLabel:          label.NewLabel(0, " 100%", font, lineHeight, labelColor, nil),
		markerWidth:         markerWidth,
		markerHeight:        markerHeight,
		markerColor:         markerNormalColor,
		status:              released,
		valueChangeCallback: dummyValueChangeCallback,
	}
}
