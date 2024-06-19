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

package label

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/juan-medina/twitch-rat/internal/draw"
)

type labelBMPF struct {
	id               Id
	text             string
	visible          bool
	lineHeight       float64
	font             draw.Font
	color            color.Color
	x, y             float64
	hasBackground    bool
	backgroundColor  color.Color
	bgX, bgY         float64
	bgW, bgH         float64
	expandBackground float64
}

func (l labelBMPF) GetId() Id {
	return l.id
}

func (l *labelBMPF) Draw(screen *ebiten.Image) {
	if l.visible {
		if l.hasBackground {
			vector.DrawFilledRect(screen, float32(l.bgX), float32(l.bgY), float32(l.bgW), float32(l.bgH), l.backgroundColor, false)
		}
		l.font.Draw(screen, l.text, l.x, l.y, l.lineHeight, l.color)
	}
}

func (l *labelBMPF) Move(x float64, y float64) {
	l.x = x
	l.y = y
	l.calculateBackground()
}

func (l *labelBMPF) SetText(text string) {
	l.text = text
	l.calculateBackground()
}

func (l labelBMPF) GetText() string {
	return l.text
}

func (l labelBMPF) Measure() (width float64, height float64) {
	width, height = l.font.Measure(l.text, l.lineHeight)
	width += l.expandBackground * 2
	height += l.expandBackground * 2
	return
}

func (l *labelBMPF) SetColor(color color.Color) {
	l.color = color
}
func (l labelBMPF) GetColor() color.Color {
	return l.color
}

func (l *labelBMPF) SetAlpha(alpha float32) {
	r, g, b, _ := l.color.RGBA()
	l.color = color.RGBA{uint8(r), uint8(g), uint8(b), uint8(alpha * 255)}
}

func (l *labelBMPF) SetVisible(visible bool) {
	l.visible = visible
}

func (l *labelBMPF) SetBackgroundColor(color color.Color, expand float64) {
	l.hasBackground = color != nil
	l.backgroundColor = color
	l.expandBackground = expand
	l.calculateBackground()
}

func (l *labelBMPF) calculateBackground() {
	if l.hasBackground {
		w, h := l.font.Measure(l.text, l.lineHeight)
		l.bgX = l.x - l.expandBackground
		l.bgY = l.y - l.expandBackground
		l.bgW = w + l.expandBackground*2
		l.bgH = h + l.expandBackground*2
	}
}

func NewLabel(id Id, text string, font draw.Font, lineHeight float64, color color.Color) Label {
	return &labelBMPF{
		id:         id,
		text:       text,
		visible:    false,
		lineHeight: lineHeight,
		font:       font,
		color:      color,
	}
}
