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
	"image/color"
	"io/fs"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/juan-medina/twitch-rat/internal/colors"
	"github.com/juan-medina/twitch-rat/internal/keys"
	"github.com/juan-medina/twitch-rat/internal/ui/button"
	"github.com/juan-medina/twitch-rat/internal/ui/input"
	"github.com/juan-medina/twitch-rat/internal/ui/label"
)

type UI interface {
	Init(fileSystem embed.FS, keys keys.Keys)
	Update(elapsedTime int)
	Draw(screen *ebiten.Image)

	SetStatusMessage(message string, color color.Color)

	EnableButton(id button.Id)
	DisableButton(id button.Id)
	SetButtonVisible(id button.Id, visible bool)
	OnButtonClick(callback func(id button.Id))

	GetInputText(id input.Id) string
	SetInputText(id input.Id, text string)
	SetInputVisible(id input.Id, visible bool)

	OnLayoutChange(width, height float64)

	SetLabelText(id label.Id, text string)
	SetLabelColor(id label.Id, color color.Color)
	SetLabelVisible(id label.Id, visible bool)
}

var (
	normalLabelColor = colors.White
)

type uiImpl struct {
	screenWidth  float64
	screenHeight float64
	fileSystem   embed.FS
	faceSource   *text.GoTextFaceSource
	normalFace   *text.GoTextFace
	bigFace      *text.GoTextFace
	buttons      []button.Button
	inputs       []input.Input
	labels       []label.Label
	keys         keys.Keys
}

const (
	MAX_BUTTONS       = 10
	BUTTON_WIDTH      = 170.0
	BUTTON_HEIGHT     = 50.0
	BUTTON_GAP        = 10.0
	BACK_BUTTON_WIDTH = 50.0
	MAX_INPUTS        = 1
	INPUT_WIDTH       = BUTTON_WIDTH*2.0 + BUTTON_GAP*2.0
	INPUT_HEIGHT      = 50
	MAX_LABELS        = 10
)

const (
	PLAY_BUTTON button.Id = iota
	BACK_BUTTON
)

const (
	INPUT_CHANNEL input.Id = iota
)

const (
	LABEL_LAST_MESSAGE label.Id = iota
	LABEL_VERSION
	LABEL_TITLE
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

	u.bigFace = &text.GoTextFace{
		Source: u.faceSource,
		Size:   40,
	}

	u.AddLabel(LABEL_LAST_MESSAGE, "", u.normalFace, normalLabelColor)
	u.AddLabel(LABEL_TITLE, "Twitch Rat", u.bigFace, normalLabelColor)

	versionStr := "v0.0.0"
	if data, err := fileSystem.ReadFile("embed/version.txt"); err == nil {
		versionStr = "v" + string(data)
	}
	u.AddLabel(LABEL_VERSION, versionStr, u.normalFace, normalLabelColor)

	u.buttons = make([]button.Button, 0, MAX_BUTTONS)
	u.inputs = make([]input.Input, 0, MAX_INPUTS)

	u.addInput(INPUT_CHANNEL, INPUT_WIDTH, INPUT_HEIGHT, 24, "", "Twitch Channel")
	u.addButton(PLAY_BUTTON, BUTTON_WIDTH, BUTTON_HEIGHT, "Play!", button.Enabled)
	u.addButton(BACK_BUTTON, BACK_BUTTON_WIDTH, BUTTON_HEIGHT, "X", button.Enabled)

	u.SetLabelVisible(LABEL_LAST_MESSAGE, true)
	u.SetLabelVisible(LABEL_VERSION, true)
}

func (u *uiImpl) Draw(screen *ebiten.Image) {
	for _, l := range u.labels {
		l.Draw(screen)
	}

	for _, b := range u.buttons {
		b.Draw(screen)

	}

	for _, i := range u.inputs {
		i.Draw(screen)
	}
}

func (u *uiImpl) Update(elapsedTime int) {
	xe, ye := ebiten.CursorPosition()
	x, y := float64(xe), float64(ye)
	pressed := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)

	u.updateButtons(x, y, pressed, elapsedTime)
	u.updateInputs(x, y, pressed, elapsedTime)
}

func (u *uiImpl) SetStatusMessage(message string, color color.Color) {
	u.SetLabelColor(LABEL_LAST_MESSAGE, color)
	u.SetLabelText(LABEL_LAST_MESSAGE, message)
}

func (u *uiImpl) OnLayoutChange(width, height float64) {
	u.screenWidth = width
	u.screenHeight = height

	cx := u.screenWidth / 2
	cy := u.screenHeight / 2

	px := cx - (INPUT_WIDTH / 2)
	py := cy - (INPUT_HEIGHT / 2) - BUTTON_GAP
	u.moveInput(INPUT_CHANNEL, px, py)

	px = cx - (BUTTON_WIDTH / 2)
	py = py + INPUT_HEIGHT + BUTTON_GAP
	u.moveButton(PLAY_BUTTON, px, py)

	px = u.screenWidth - BACK_BUTTON_WIDTH
	py = 0
	u.moveButton(BACK_BUTTON, px, py)

	gapX := u.normalFace.Size
	gapY := u.normalFace.Size * 1.5
	u.getLabel(LABEL_LAST_MESSAGE).Move(gapX, u.screenHeight-gapY)

	versionLabel := u.getLabel(LABEL_VERSION)
	cx, cy = versionLabel.Measure()
	versionLabel.Move(u.screenWidth-cx, u.screenHeight-cy)

	titleLabel := u.getLabel(LABEL_TITLE)
	cx, cy = titleLabel.Measure()
	titleLabel.Move(u.screenWidth/2-cx/2, u.screenHeight/2-cy-INPUT_HEIGHT-gapY)
}

func (u *uiImpl) getButton(id button.Id) button.Button {
	for i, b := range u.buttons {
		if b.GetId() == id {
			return u.buttons[i]
		}
	}
	return nil
}

func (u *uiImpl) changeButtonState(id button.Id, state button.State) {
	if b := u.getButton(id); b != nil {
		b.ChangeState(state)
	}
}

func (u *uiImpl) EnableButton(id button.Id) {
	u.changeButtonState(id, button.Enabled)
}

func (u *uiImpl) DisableButton(id button.Id) {
	u.changeButtonState(id, button.Disabled)
}

func (u *uiImpl) moveButton(id button.Id, x, y float64) {
	if b := u.getButton(id); b != nil {
		b.Move(x, y)
	}
}

func (u *uiImpl) addButton(id button.Id, w, h float64, label string, state button.State) {
	u.buttons = append(u.buttons, button.New(id, w, h, label, u.normalFace, state))
}

func (u *uiImpl) SetButtonVisible(id button.Id, visible bool) {
	if b := u.getButton(id); b != nil {
		b.SetVisible(visible)
	}
}

func (u *uiImpl) OnButtonClick(callback func(id button.Id)) {
	for i := range u.buttons {
		u.buttons[i].OnButtonClick(callback)
	}
}

func (u *uiImpl) updateButtons(mouseX, mouseY float64, leftPressed bool, elapsedTime int) {
	for i := range u.buttons {
		u.buttons[i].Update(mouseX, mouseY, leftPressed, elapsedTime)
	}
}

func (u *uiImpl) getInput(id input.Id) input.Input {
	for i := range u.inputs {
		if u.inputs[i].GetId() == id {
			return u.inputs[i]
		}
	}
	return nil
}

func (ui *uiImpl) moveInput(id input.Id, x, y float64) {
	if i := ui.getInput(id); i != nil {
		i.Move(x, y)
	}
}

func (ui uiImpl) GetInputText(id input.Id) string {
	if i := ui.getInput(id); i != nil {
		return i.GetText()
	}
	return ""
}

func (ui *uiImpl) SetInputText(id input.Id, text string) {
	if i := ui.getInput(id); i != nil {
		i.SetText(text)
	}
}

func (ui *uiImpl) SetInputVisible(id input.Id, visible bool) {
	if i := ui.getInput(id); i != nil {
		i.SetVisible(visible)
	}
}

func (u *uiImpl) addInput(id input.Id, w, h float64, maxLength int, initialText string, placeHolder string) {
	u.inputs = append(u.inputs, input.New(id, w, h, maxLength, initialText, placeHolder, u.normalFace))
}

func (u *uiImpl) updateInputs(mouseX, mouseY float64, leftPressed bool, elapsedTime int) {
	ebiten.SetCursorShape(ebiten.CursorShapeDefault)
	for i := range u.inputs {
		u.inputs[i].Update(mouseX, mouseY, leftPressed, u.keys, elapsedTime)
	}
}

func (u *uiImpl) AddLabel(id label.Id, text string, face *text.GoTextFace, color color.Color) {
	u.labels = append(u.labels, label.New(id, text, face, color))
}

func (u *uiImpl) getLabel(id label.Id) label.Label {
	for i := range u.labels {
		if u.labels[i].GetId() == id {
			return u.labels[i]
		}
	}
	return nil
}

func (u *uiImpl) SetLabelText(id label.Id, text string) {
	if l := u.getLabel(id); l != nil {
		l.SetText(text)
	}
}
func (u *uiImpl) SetLabelColor(id label.Id, color color.Color) {
	if l := u.getLabel(id); l != nil {
		l.SetColor(color)
	}
}

func (ui *uiImpl) SetLabelVisible(id label.Id, visible bool) {
	if l := ui.getLabel(id); l != nil {
		l.SetVisible(visible)
	}
}

func New() UI {
	return &uiImpl{}
}
