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
	"github.com/juan-medina/twitch-rat/internal/keys"
	"github.com/juan-medina/twitch-rat/internal/step"
)

type InputId int

const (
	INPUT_CHANNEL InputId = iota
)

type input struct {
	id          InputId
	x, y, w, h  float64
	maxLength   int
	text        string
	placeHolder string
	borderColor color.Color
	color       color.Color
	visible     bool
	textDo      text.DrawOptions
	caretDo     text.DrawOptions
	caretAlpha  step.LoopValue
	focus       bool
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
		if i.focus {
			if i.text != "" {
				text.Draw(screen, i.text, i.face, &i.textDo)
			}
			text.Draw(screen, "|", i.face, &i.caretDo)
		} else {
			if i.text != "" {
				text.Draw(screen, i.text, i.face, &i.textDo)
			} else {
				text.Draw(screen, i.placeHolder, i.face, &i.textDo)
			}
		}

	}
}

func (u *uiImpl) addInput(id InputId, x, y, w, h float64, maxLength int, initialText string, placeHolder string) {
	u.inputs = append(u.inputs, newInput(id, x, y, w, h, maxLength, initialText, u.normalFace, placeHolder))
}

func (u *uiImpl) updateInputs(mouseX, mouseY float64, leftPressed bool, elapsedTime int) {
	for i, _ := range u.inputs {
		u.inputs[i].Update(mouseX, mouseY, leftPressed, u.keys, elapsedTime)
	}
}

func (i *input) Update(mouseX, mouseY float64, leftPressed bool, keys keys.Keys, elapsedTime int) {

	if leftPressed {
		if i.hit(mouseX, mouseY) {
			i.focus = true
		} else {
			i.focus = false
		}
	}

	if !i.focus {
		return
	}
	if key := keys.LastRepeatedKey(); key != ebiten.KeyMax {
		if key >= ebiten.Key0 && key <= ebiten.Key9 {
			i.addLetter(rune(key - ebiten.Key0 + '0'))
		} else if key >= ebiten.KeyNumpad0 && key <= ebiten.KeyNumpad9 {
			i.addLetter(rune(key - ebiten.KeyNumpad0 + '0'))
		} else if key >= ebiten.KeyA && key <= ebiten.KeyZ {
			i.addLetter(rune(key - ebiten.KeyA + 'a'))
		} else if key == ebiten.KeyMinus {
			if keys.IsKeyDown(ebiten.KeyShift) {
				i.addLetter('_')
			} else {
				i.addLetter('-')
			}
		} else if key == ebiten.KeyNumpadSubtract {
			i.addLetter('-')
		} else if key == ebiten.KeyShift {
			if keys.IsKeyDown(ebiten.KeyMinus) {
				i.addLetter('_')
			}
		} else if key == ebiten.KeyBackspace {
			i.removeLetter()
		}
	}

	if keys.IsKeyDown(ebiten.KeyEnter) || keys.IsKeyDown(ebiten.KeyNumpadEnter) {
		i.focus = false
	}

	if i.caretAlpha.Update(elapsedTime) {
		i.caretDo.ColorScale.Reset()
		i.caretDo.ColorScale.ScaleWithColor(inputTextColor)
		i.caretDo.ColorScale.ScaleAlpha(i.caretAlpha.GetValue())
	}
}

func (i *input) addLetter(letter rune) {
	if len(i.text) < i.maxLength {
		i.text += string(letter)
		i.updateTextColor()
		i.updateCaretPosition()
	}
}

func (i *input) removeLetter() {
	if len(i.text) > 0 {
		i.text = i.text[:len(i.text)-1]
		i.updateTextColor()
		i.updateCaretPosition()
	}
}

func (i *input) updateTextColor() {
	i.textDo.ColorScale.Reset()
	if i.text == "" {
		i.textDo.ColorScale.ScaleWithColor(inputEmptyColor)
	} else {
		i.textDo.ColorScale.ScaleWithColor(inputTextColor)
	}
}

func (i input) GetText() string {
	return i.text
}

func (i *input) SetText(text string) {
	i.text = text
	i.updateTextColor()
	i.updateCaretPosition()
}

func (i *input) updateCaretPosition() {
	i.caretDo.GeoM.Reset()
	i.caretDo.GeoM.Translate(i.x+INPUT_LEFT_GAP, i.y+INPUT_TOP_GAP)
	if i.text != "" {
		move, _ := text.Measure(i.text, i.face, 0)
		i.caretDo.GeoM.Translate(float64(move), 0)
	}
}

func (i input) hit(x, y float64) bool {
	if x > i.x && x < i.x+i.w && y > i.y && y < i.y+i.h {
		return true
	}
	return false
}

func (ui uiImpl) GetInputText(id InputId) string {
	for _, i := range ui.inputs {
		if i.id == id {
			return i.GetText()
		}
	}
	return ""
}

func (ui *uiImpl) SetInputText(id InputId, text string) {
	for iid, i := range ui.inputs {
		if i.id == id {
			ui.inputs[iid].SetText(text)
			return
		}
	}
}

func newInput(id InputId, x, y, w, h float64, maxLength int, initialText string, face *text.GoTextFace, placeHolder string) input {
	textDo := text.DrawOptions{}
	textDo.GeoM.Reset()
	textDo.GeoM.Translate(x+INPUT_LEFT_GAP, y+INPUT_TOP_GAP)

	textDo.ColorScale.Reset()
	if initialText == "" {
		textDo.ColorScale.ScaleWithColor(inputEmptyColor)
	} else {
		textDo.ColorScale.ScaleWithColor(inputTextColor)
	}
	caretDo := text.DrawOptions{}
	caretDo.GeoM.Reset()
	caretDo.GeoM.Translate(x+INPUT_LEFT_GAP, y+INPUT_TOP_GAP)
	caretDo.ColorScale.Reset()
	caretDo.ColorScale.ScaleWithColor(inputTextColor)
	return input{
		id:          id,
		x:           x,
		y:           y,
		w:           w,
		h:           h,
		maxLength:   maxLength,
		text:        initialText,
		placeHolder: placeHolder,
		borderColor: inputBorderColor,
		color:       inputColor,
		visible:     true,
		textDo:      textDo,
		caretDo:     caretDo,
		caretAlpha:  step.NewPingPongValue(0, 1, 200, 100),
		face:        face,
	}
}
