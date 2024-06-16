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

package label

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/juan-medina/twitch-rat/internal/draw"
)

type labelBMPF struct {
	id         Id
	text       string
	visible    bool
	lineHeight float64
	font       draw.Font
	color      color.Color
	x, y       float64
}

func (l labelBMPF) GetId() Id {
	return l.id
}

func (l *labelBMPF) Draw(screen *ebiten.Image) {
	if l.visible {
		l.font.Draw(screen, l.text, l.x, l.y, l.lineHeight, l.color)
	}
}

func (l *labelBMPF) Move(x float64, y float64) {
	l.x = x
	l.y = y
}

func (l *labelBMPF) SetText(text string) {
	l.text = text
}

func (l labelBMPF) GetText() string {
	return l.text
}

func (l labelBMPF) Measure() (float64, float64) {
	return l.font.Measure(l.text, l.lineHeight)
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
