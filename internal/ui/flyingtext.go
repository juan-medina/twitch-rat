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

package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/juan-medina/twitch-rat/internal/colors"
	"github.com/juan-medina/twitch-rat/internal/draw"
	"github.com/juan-medina/twitch-rat/internal/step"
)

type flyingText struct {
	text    string
	color   colors.CustomColor
	alpha   step.Value
	x, y    float64
	vy      float64
	visible bool
}

func (u *uiImpl) AddFlyingText(text string, color colors.CustomColor, x, y float64) {
	dx, _ := u.fontSmall.Measure(text, u.fontSmall.DefaultSize())
	px := x - (dx / 2)
	for i := range u.flyingTexts {
		if !u.flyingTexts[i].visible {
			u.flyingTexts[i].visible = true
			u.flyingTexts[i].text = text
			u.flyingTexts[i].color = color
			u.flyingTexts[i].x = px
			u.flyingTexts[i].y = y
			u.flyingTexts[i].vy = FLYING_TEXT_VY
			u.flyingTexts[i].alpha.Reset()
			return
		}
	}
	u.flyingTexts = append(u.flyingTexts, flyingText{
		visible: true,
		text:    text,
		color:   color,
		x:       px,
		y:       y,
		vy:      FLYING_TEXT_VY,
		alpha:   step.NewFromMiddleToPauseValue(255, 255, 0, FLYING_TIME_FULL, FLYING_TIME_TO_VANISH, FLYING_TIME_TO_FREE),
	})
}

func (f flyingText) draw(screen *ebiten.Image, font draw.Font) {
	if !f.visible {
		return
	}
	font.Draw(screen, f.text, f.x, f.y, font.DefaultSize(), f.color)
}

func (u *uiImpl) updateFlyingTexts(elapsedTime int) {
	for i := range u.flyingTexts {
		if u.flyingTexts[i].visible {
			u.flyingTexts[i].y -= (u.flyingTexts[i].vy * float64(elapsedTime))
			if u.flyingTexts[i].alpha.Update(elapsedTime) {
				newAlpha := uint8(u.flyingTexts[i].alpha.GetValue())
				u.flyingTexts[i].color = u.flyingTexts[i].color.NewWithAlpha(newAlpha)
				if u.flyingTexts[i].alpha.IsAtEnd() {
					u.flyingTexts[i].visible = false
				}
			}
		}
	}
}

func (u *uiImpl) drawFlyingText(screen *ebiten.Image) {
	for _, f := range u.flyingTexts {
		f.draw(screen, u.fontSmall)
	}
}
