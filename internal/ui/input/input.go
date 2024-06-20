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

package input

import (
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/juan-medina/twitch-rat/internal/audio"
	"github.com/juan-medina/twitch-rat/internal/colors"
	"github.com/juan-medina/twitch-rat/internal/draw"
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
	IsEditing() bool
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
	audioPlayer     audio.Player
	clipboard       Clipboard
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
	inputTextColor        = colors.Purple
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

		if i.IsEditing() {
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

func (i input) IsEditing() bool {
	return i.editing
}

func (i *input) SetVisible(visible bool) {
	i.visible = visible
	if i.IsEditing() {
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
			if !i.IsEditing() {
				i.edit()
			}
		} else {
			if i.IsEditing() {
				i.saveEdit()
			}
		}
	}

	if !i.IsEditing() {
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
		keys.SwallowKey(ebiten.KeyEnter)
		keys.SwallowKey(ebiten.KeyNumpadEnter)
	}

	if keys.IsDownNoRepeat(ebiten.KeyEscape) {
		playSound = true
		i.cancelEdit()
		keys.SwallowKey(ebiten.KeyEscape)
	}

	if keys.IsDown(ebiten.KeyControl) && keys.IsDownNoRepeat(ebiten.KeyV) {
		if i.IsEditing() {
			playSound = true
			keys.SwallowKey(ebiten.KeyControl)
			keys.SwallowKey(ebiten.KeyV)
			i.clipboard.Request(i.id)
		}
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

func (i *input) onClipboard(id Id, text string) {
	if i.visible && i.editing && i.id == id {
		text = strings.ReplaceAll(text, "\n", "")
		clipboardText := i.GetText() + text
		if len(clipboardText) > i.maxLength {
			clipboardText = clipboardText[:i.maxLength]
		}
		i.SetText(clipboardText)
	}
}

func New(id Id, w, h float64, maxLength int, initialText string, placeholder string, font draw.Font, lineHeight float64, audioPlayer audio.Player) Input {
	audioPlayer.LoadSound(MOUSE_CLICK_SOUND)
	audioPlayer.LoadSound(KEY_PRESS_SOUND)
	i := input{
		id:              id,
		w:               w,
		h:               h,
		maxLength:       maxLength,
		borderColor:     inputBorderColor,
		color:           inputColor,
		visible:         false,
		textLabel:       label.NewLabel(0, initialText, font, lineHeight, inputTextColor, nil),
		textPlaceHolder: label.NewLabel(0, placeholder, font, lineHeight, inputPlaceHolderColor, nil),
		caretLabel:      label.NewLabel(0, "|", font, lineHeight, caretColor, nil),
		caretAlpha:      step.NewPingPongValue(0, 1, 200, 100),
		audioPlayer:     audioPlayer,
	}

	i.clipboard = newClipboard()
	i.clipboard.OnReady(i.onClipboard)

	return &i
}
