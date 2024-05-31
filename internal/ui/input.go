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

type InputId int

const (
	INPUT_CHANNEL InputId = iota
)

type input struct {
	id          InputId
	x, y, w, h  float64
	text        string
	placeHolder string
	borderColor color.Color
	color       color.Color
	visible     bool
	do          text.DrawOptions
	face        *text.GoTextFace
}

const (
	MAX_INPUTS        = 1
	INPUT_BORDER_SIZE = 5
	INPUT_LEFT_GAP    = 10
	INPUT_TOP_GAP     = 10
	INPUT_WIDTH       = BUTTON_WIDTH*2 + BUTTON_GAP*2
	INPUT_HEIGHT      = 50
)

var (
	inputBorderColor = darkPurple
	inputColor       = white
	inputTextColor   = black
	inputEmptyColor  = gray
)

func (i input) draw(screen *ebiten.Image) {
	if i.visible {
		vector.DrawFilledRect(screen, float32(i.x), float32(i.y), float32(i.w), float32(i.h), i.color, true)
		vector.StrokeRect(screen, float32(i.x), float32(i.y), float32(i.w), float32(i.h), INPUT_BORDER_SIZE, i.borderColor, true)
		if i.text != "" {
			text.Draw(screen, i.text, i.face, &i.do)
		} else {
			text.Draw(screen, i.placeHolder, i.face, &i.do)
		}

	}
}

func (u *uiImpl) addInput(id InputId, x, y, w, h float64, initialText string, placeHolder string) {
	u.inputs = append(u.inputs, newInput(id, x, y, w, h, initialText, u.normalFace, placeHolder))
}

func newInput(id InputId, x, y, w, h float64, initialText string, face *text.GoTextFace, placeHolder string) input {
	do := text.DrawOptions{}
	do.GeoM.Reset()
	do.GeoM.Translate(x+INPUT_LEFT_GAP, y+INPUT_TOP_GAP)

	do.ColorScale.Reset()
	if initialText == "" {
		do.ColorScale.ScaleWithColor(inputEmptyColor)
	} else {
		do.ColorScale.ScaleWithColor(inputTextColor)
	}

	return input{
		id:          id,
		x:           x,
		y:           y,
		w:           w,
		h:           h,
		text:        initialText,
		placeHolder: placeHolder,
		borderColor: inputBorderColor,
		color:       inputColor,
		visible:     true,
		do:          do,
		face:        face,
	}
}
