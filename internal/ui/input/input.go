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

package input

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/juan-medina/twitch-rat/internal/audio"
	"github.com/juan-medina/twitch-rat/internal/colors"
	"github.com/juan-medina/twitch-rat/internal/keys"
	"github.com/juan-medina/twitch-rat/internal/step"
	"github.com/juan-medina/twitch-rat/internal/ui/label"
)

type Id int

type Input interface {
	GetId() Id
	Draw(screen *ebiten.Image)
	Update(mouseX, mouseY float64, leftPressed bool, keys keys.Keys, elapsedTime int)
	SetVisible(visible bool)
	Move(x, y float64)
	GetText() string
	SetText(text string)
}

type input struct {
	id         Id
	x, y, w, h float64
	maxLength  int

	borderColor     color.Color
	color           color.Color
	visible         bool
	textLabel       label.Label
	textPlaceHolder label.Label
	caretLabel      label.Label
	caretAlpha      step.Value
	editing         bool
	savedText       string
	face            *text.GoTextFace
	audioPlayer     audio.Player
}

const (
	INPUT_BORDER_SIZE = 5
	INPUT_LEFT_GAP    = 10
	INPUT_TOP_GAP     = 10
	MOUSE_CLICK_SOUND = "embed/sounds/mouse.ogg"
	KEY_PRESS_SOUND   = "embed/sounds/key.ogg"
)

var (
	inputBorderColor      = colors.DarkPurple
	inputColor            = colors.White
	inputTextColor        = colors.Black
	inputPlaceHolderColor = colors.Gray
	caretColor            = colors.Violet
)

func (i input) GetId() Id {
	return i.id
}

func (i input) Draw(screen *ebiten.Image) {
	if i.visible {
		vector.DrawFilledRect(screen, float32(i.x), float32(i.y), float32(i.w), float32(i.h), i.color, true)
		vector.StrokeRect(screen, float32(i.x), float32(i.y), float32(i.w), float32(i.h), INPUT_BORDER_SIZE, i.borderColor, true)
		i.textLabel.Draw(screen)

		if i.isEditing() {
			i.caretLabel.Draw(screen)
		} else {
			text := i.GetText()
			if text == "" {
				i.textPlaceHolder.Draw(screen)
			}
		}
	}
}

func (i *input) edit() {
	i.audioPlayer.PlaySound(MOUSE_CLICK_SOUND)
	i.editing = true
	i.savedText = i.textLabel.GetText()
	i.updateCaretPosition()
}

func (i *input) cancelEdit() {
	i.editing = false
	i.SetText(i.savedText)
}

func (i *input) saveEdit() {
	i.editing = false
	i.savedText = i.textLabel.GetText()
}

func (i input) isEditing() bool {
	return i.editing
}

func (i *input) SetVisible(visible bool) {
	i.visible = visible
	if i.isEditing() {
		i.saveEdit()
	}
	i.caretLabel.SetVisible(visible)
	i.textLabel.SetVisible(visible)
	i.textPlaceHolder.SetVisible(visible)
}

func (i *input) Update(mouseX, mouseY float64, leftPressed bool, keys keys.Keys, elapsedTime int) {

	if !i.visible {
		return
	}

	hit := i.hit(mouseX, mouseY)
	if hit {
		ebiten.SetCursorShape(ebiten.CursorShapeText)
	}

	if leftPressed {
		if hit {
			if !i.isEditing() {
				i.edit()
			}
		} else {
			if i.isEditing() {
				i.saveEdit()
			}
		}
	}

	if !i.isEditing() {
		return
	}

	playSound := false
	inputChar := keys.LastInputChar()
	if inputChar != 0 {
		playSound = true
		i.addLetter(inputChar)
	}

	if keys.IsDownNoRepeat(ebiten.KeyBackspace) {
		playSound = true
		i.removeLetter()
	}

	if keys.IsDownNoRepeat(ebiten.KeyEnter) || keys.IsDownNoRepeat(ebiten.KeyNumpadEnter) {
		playSound = true
		i.saveEdit()
	}

	if keys.IsDownNoRepeat(ebiten.KeyEscape) {
		playSound = true
		i.cancelEdit()
	}

	if playSound {
		i.audioPlayer.PlaySound(KEY_PRESS_SOUND)
	}
	i.updateCaretPosition()
	i.caretAlpha.Update(elapsedTime)

	if i.caretAlpha.Update(elapsedTime) {
		i.caretLabel.SetColor(caretColor)
		i.caretLabel.SetAlpha(i.caretAlpha.GetValue())
	}
}

func (i *input) addLetter(letter rune) {
	text := i.textLabel.GetText()
	if len(text) < i.maxLength {
		text += string(letter)
		i.SetText(text)
	}
}

func (i *input) removeLetter() {
	runes := []rune(i.textLabel.GetText())
	if len(runes) > 0 {
		runes = runes[:len(runes)-1]
		i.SetText(string(runes))
	}
}

func (i input) GetText() string {
	return i.textLabel.GetText()
}

func (i *input) SetText(text string) {
	i.textLabel.SetText(text)
	i.updateCaretPosition()
}

func (i *input) updateCaretPosition() {
	advance := 0.0
	if i.textLabel.GetText() != "" {
		advance, _ = i.textLabel.Measure()
	}
	i.caretLabel.Move(i.x+INPUT_LEFT_GAP+advance, i.y+INPUT_TOP_GAP)
}

func (i input) hit(x, y float64) bool {
	if x > i.x && x < i.x+i.w && y > i.y && y < i.y+i.h {
		return true
	}
	return false
}

func (i *input) Move(x, y float64) {
	i.x = x
	i.y = y
	i.updateCaretPosition()

	i.textLabel.Move(i.x+INPUT_LEFT_GAP, i.y+INPUT_TOP_GAP)
	i.textPlaceHolder.Move(i.x+INPUT_LEFT_GAP, i.y+INPUT_TOP_GAP)
}

func New(id Id, w, h float64, maxLength int, initialText string, placeholder string, face *text.GoTextFace, audioPlayer audio.Player) Input {
	audioPlayer.LoadSound(MOUSE_CLICK_SOUND)
	audioPlayer.LoadSound(KEY_PRESS_SOUND)
	return &input{
		id:              id,
		w:               w,
		h:               h,
		maxLength:       maxLength,
		borderColor:     inputBorderColor,
		color:           inputColor,
		visible:         false,
		textLabel:       label.New(0, initialText, face, inputTextColor),
		textPlaceHolder: label.New(0, placeholder, face, inputPlaceHolderColor),
		caretLabel:      label.New(0, "|", face, caretColor),
		caretAlpha:      step.NewPingPongValue(0, 1, 200, 100),
		face:            face,
		audioPlayer:     audioPlayer,
	}
}
