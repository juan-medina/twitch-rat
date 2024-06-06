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
	"bytes"
	"embed"
	"io/fs"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/juan-medina/twitch-rat/internal/keys"
)

type UI interface {
	Init(fileSystem embed.FS, keys keys.Keys)
	Update(elapsedTime int)
	Draw(screen *ebiten.Image)

	SetStatusMessage(message string)

	EnableButton(id ButtonId)
	DisableButton(id ButtonId)
	SetButtonVisible(id ButtonId, visible bool)
	OnButtonClick(callback func(id ButtonId))

	GetInputText(id InputId) string
	SetInputText(id InputId, text string)
	SetInputVisible(id InputId, visible bool)

	OnLayoutChange(width, height int)
}

type uiImpl struct {
	screenWidth   int
	screenHeight  int
	fileSystem    embed.FS
	faceSource    *text.GoTextFaceSource
	normalFace    *text.GoTextFace
	lastMessage   string
	lastMessageDO text.DrawOptions
	buttons       []button
	inputs        []input
	keys          keys.Keys
}

const (
	MAX_BUTTONS       = 10
	BUTTON_WIDTH      = 170
	BUTTON_HEIGHT     = 50
	BUTTON_GAP        = 10
	BACK_BUTTON_WIDTH = 50
)

func (u *uiImpl) Init(fileSystem embed.FS, keys keys.Keys) {
	u.fileSystem = fileSystem
	u.keys = keys

	fontBytes, err := fs.ReadFile(u.fileSystem, "embed/fonts/default.ttf")
	if err != nil {
		log.Fatal(err)
	}

	s, err := text.NewGoTextFaceSource(bytes.NewReader(fontBytes))
	if err != nil {
		log.Fatal(err)
	}

	u.faceSource = s

	u.normalFace = &text.GoTextFace{
		Source: u.faceSource,
		Size:   24,
	}

	u.buttons = make([]button, 0, MAX_BUTTONS)
	u.inputs = make([]input, 0, MAX_INPUTS)

	u.addInput(INPUT_CHANNEL, INPUT_WIDTH, INPUT_HEIGHT, 24, "", "Twitch Channel Name")
	u.addButton(PLAY_BUTTON, BUTTON_WIDTH, BUTTON_HEIGHT, "Play!", buttonEnabled)
	u.addButton(BACK_BUTTON, BACK_BUTTON_WIDTH, BUTTON_HEIGHT, "X", buttonEnabled)
}

func (u *uiImpl) Draw(screen *ebiten.Image) {
	text.Draw(screen, u.lastMessage, u.normalFace, &u.lastMessageDO)
	for _, b := range u.buttons {
		b.draw(screen)

	}

	for _, i := range u.inputs {
		i.draw(screen)
	}
}

func (u *uiImpl) Update(elapsedTime int) {
	xe, ye := ebiten.CursorPosition()
	x, y := float64(xe), float64(ye)
	pressed := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)

	u.updateButtons(x, y, pressed, elapsedTime)
	u.updateInputs(x, y, pressed, elapsedTime)
}

func (u *uiImpl) SetStatusMessage(message string) {
	u.lastMessage = message
}

func (u *uiImpl) OnLayoutChange(width, height int) {
	u.screenWidth = width
	u.screenHeight = height

	cx := float64(u.screenWidth / 2)
	cy := float64(u.screenHeight / 2)

	px := cx - INPUT_WIDTH - (BUTTON_GAP / 2)
	py := cy - (BUTTON_HEIGHT / 2)
	u.moveInput(INPUT_CHANNEL, px, py)

	px = cx + (BUTTON_GAP / 2)
	u.moveButton(PLAY_BUTTON, px, py)

	px = float64(u.screenWidth) - BACK_BUTTON_WIDTH
	py = 0
	u.moveButton(BACK_BUTTON, px, py)

	u.lastMessageDO.GeoM.Reset()
	gapX := u.normalFace.Size
	gapY := u.normalFace.Size * 1.5
	u.lastMessageDO.GeoM.Translate(gapX, float64(u.screenHeight)-gapY)
}

func (u *uiImpl) getButton(id ButtonId) *button {
	for i, b := range u.buttons {
		if b.id == id {
			return &u.buttons[i]
		}
	}
	return nil
}

func (u *uiImpl) changeButtonState(id ButtonId, state buttonState) {
	if b := u.getButton(id); b != nil {
		b.changeState(state)
	}
}

func (u *uiImpl) EnableButton(id ButtonId) {
	u.changeButtonState(id, buttonEnabled)
}

func (u *uiImpl) DisableButton(id ButtonId) {
	u.changeButtonState(id, buttonDisabled)
}

func (u *uiImpl) moveButton(id ButtonId, x, y float64) {
	if b := u.getButton(id); b != nil {
		b.move(x, y)
	}
}

func (u *uiImpl) addButton(id ButtonId, w, h float64, label string, state buttonState) {
	u.buttons = append(u.buttons, newButton(id, w, h, label, u.normalFace, state))
}

func (u *uiImpl) SetButtonVisible(id ButtonId, visible bool) {
	if b := u.getButton(id); b != nil {
		b.SetVisible(visible)
	}
}

func (u *uiImpl) OnButtonClick(callback func(id ButtonId)) {
	for i := range u.buttons {
		u.buttons[i].OnButtonClick(callback)
	}
}

func (u *uiImpl) updateButtons(mouseX, mouseY float64, leftPressed bool, elapsedTime int) {
	for i := range u.buttons {
		u.buttons[i].Update(mouseX, mouseY, leftPressed, elapsedTime)
	}
}

func (u *uiImpl) getInput(id InputId) *input {
	for i := range u.inputs {
		if u.inputs[i].id == id {
			return &u.inputs[i]
		}
	}
	return nil
}

func (ui *uiImpl) moveInput(id InputId, x, y float64) {
	if i := ui.getInput(id); i != nil {
		i.move(x, y)
	}
}

func (ui uiImpl) GetInputText(id InputId) string {
	if i := ui.getInput(id); i != nil {
		return i.GetText()
	}
	return ""
}

func (ui *uiImpl) SetInputText(id InputId, text string) {
	if i := ui.getInput(id); i != nil {
		i.SetText(text)
	}
}

func (ui *uiImpl) SetInputVisible(id InputId, visible bool) {
	if i := ui.getInput(id); i != nil {
		i.setVisible(visible)
	}
}

func (u *uiImpl) addInput(id InputId, w, h float64, maxLength int, initialText string, placeHolder string) {
	u.inputs = append(u.inputs, newInput(id, w, h, maxLength, initialText, u.normalFace, placeHolder))
}

func (u *uiImpl) updateInputs(mouseX, mouseY float64, leftPressed bool, elapsedTime int) {
	ebiten.SetCursorShape(ebiten.CursorShapeDefault)
	for i := range u.inputs {
		u.inputs[i].Update(mouseX, mouseY, leftPressed, u.keys, elapsedTime)
	}
}

func New() UI {
	return &uiImpl{}
}
