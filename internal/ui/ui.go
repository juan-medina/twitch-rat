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
	"github.com/juan-medina/twitch-rat/internal/audio"
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
	ClickButton(id button.Id)
	SetButtonClickCallback(callback func(id button.Id))

	GetInputText(id input.Id) string
	SetInputText(id input.Id, text string)
	SetInputVisible(id input.Id, visible bool)
	IsInputEditing(id input.Id) bool

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
	smallFace    *text.GoTextFace
	normalFace   *text.GoTextFace
	bigFace      *text.GoTextFace
	buttons      []button.Button
	inputs       []input.Input
	labels       []label.Label
	keys         keys.Keys
	audioPlayer  audio.Player
	licenseText  string
	aboutText    string
}

const (
	MAX_BUTTONS                  = 10
	MENU_START                   = 250.0
	BUTTON_WIDTH                 = 180.0
	BUTTON_HEIGHT                = 50.0
	SMALL_BUTTON_WIDTH           = 85.0
	SMALL_BUTTON_HEIGHT          = 25.0
	BUTTON_GAP                   = 20.0
	BACK_BUTTON_WIDTH            = 35.0
	BACK_BUTTON_HEIGHT           = 35.0
	TITLE_TO_ELEMENTS_SEPARATION = 50.0
	MAX_INPUTS                   = 1
	INPUT_WIDTH                  = BUTTON_WIDTH*2.0 + BUTTON_GAP*2.0
	INPUT_HEIGHT                 = 50
	MAX_LABELS                   = 10
)

const (
	PLAY_BUTTON button.Id = iota
	BACK_BUTTON
	OPTIONS_BUTTON
	ABOUT_BUTTON
	SUBMENU_ABOUT_BACK_BUTTON
	SUBMENU_OPTION_BACK_BUTTON
	ACCEPT_LICENSE_BUTTON
)

const (
	INPUT_CHANNEL input.Id = iota
)

const (
	LABEL_LAST_MESSAGE label.Id = iota
	LABEL_VERSION
	LABEL_TITLE
	LABEL_LICENSE
	LABEL_ABOUT_MESSAGE
)

func (u *uiImpl) Init(fileSystem embed.FS, keys keys.Keys) {
	u.fileSystem = fileSystem
	u.keys = keys

	data, _ := fileSystem.ReadFile("embed/text/license.txt")
	u.licenseText = string(data)

	data, _ = fileSystem.ReadFile("embed/text/about.txt")
	u.aboutText = string(data)

	fontBytes, err := fs.ReadFile(u.fileSystem, "embed/fonts/default.ttf")
	if err != nil {
		log.Fatal(err)
	}

	s, err := text.NewGoTextFaceSource(bytes.NewReader(fontBytes))
	if err != nil {
		log.Fatal(err)
	}

	u.faceSource = s

	u.smallFace = &text.GoTextFace{
		Source:    u.faceSource,
		Direction: 0,
		Size:      12,
	}

	u.normalFace = &text.GoTextFace{
		Source: u.faceSource,
		Size:   24,
	}

	u.bigFace = &text.GoTextFace{
		Source: u.faceSource,
		Size:   40,
	}

	u.AddLabel(LABEL_LAST_MESSAGE, "", u.normalFace, u.normalFace.Size, normalLabelColor)
	u.AddLabel(LABEL_TITLE, "Twitch Rat", u.bigFace, u.bigFace.Size, normalLabelColor)

	versionStr := "v0.0.0"
	if data, err := fileSystem.ReadFile("embed/version.txt"); err == nil {
		versionStr = "v" + string(data)
	}
	u.AddLabel(LABEL_VERSION, versionStr, u.smallFace, u.smallFace.Size, normalLabelColor)

	u.buttons = make([]button.Button, 0, MAX_BUTTONS)
	u.inputs = make([]input.Input, 0, MAX_INPUTS)

	u.addInput(INPUT_CHANNEL, INPUT_WIDTH, INPUT_HEIGHT, 24, "", "Twitch Channel")
	u.addButton(PLAY_BUTTON, BUTTON_WIDTH, BUTTON_HEIGHT, u.normalFace, "Play!", button.Enabled)
	u.addButton(BACK_BUTTON, BACK_BUTTON_WIDTH, BACK_BUTTON_HEIGHT, u.smallFace, "X", button.Enabled)
	u.addButton(OPTIONS_BUTTON, SMALL_BUTTON_WIDTH, SMALL_BUTTON_HEIGHT, u.smallFace, "Options", button.Enabled)
	u.addButton(ABOUT_BUTTON, SMALL_BUTTON_WIDTH, SMALL_BUTTON_HEIGHT, u.smallFace, "About", button.Enabled)

	u.AddLabel(LABEL_ABOUT_MESSAGE, u.aboutText, u.smallFace, u.smallFace.Size, normalLabelColor)
	u.addButton(SUBMENU_ABOUT_BACK_BUTTON, SMALL_BUTTON_WIDTH, SMALL_BUTTON_HEIGHT, u.smallFace, "Back", button.Enabled)
	u.addButton(SUBMENU_OPTION_BACK_BUTTON, SMALL_BUTTON_WIDTH, SMALL_BUTTON_HEIGHT, u.smallFace, "Back", button.Enabled)

	u.AddLabel(LABEL_LICENSE, u.licenseText, u.smallFace, u.smallFace.Size, normalLabelColor)
	u.addButton(ACCEPT_LICENSE_BUTTON, SMALL_BUTTON_WIDTH, SMALL_BUTTON_HEIGHT, u.smallFace, "Accept", button.Enabled)

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

func (u *uiImpl) layoutMainElements(cx, cy float64) {
	titleLabel := u.getLabel(LABEL_TITLE)
	ix, _ := titleLabel.Measure()
	px := cx - (ix / 2)
	py := cy
	titleLabel.Move(px, py)

	px = u.screenWidth - BACK_BUTTON_WIDTH
	py = 0
	u.moveButton(BACK_BUTTON, px, py)

	gapX := u.normalFace.Size * 0.5
	gapY := u.normalFace.Size * 1.5
	u.getLabel(LABEL_LAST_MESSAGE).Move(gapX, u.screenHeight-gapY)

	gapX = u.smallFace.Size * 0.5
	gapY = u.smallFace.Size * 1.5
	versionLabel := u.getLabel(LABEL_VERSION)
	cx, _ = versionLabel.Measure()
	versionLabel.Move(u.screenWidth-cx-gapX, u.screenHeight-gapY)
}

func (u *uiImpl) layoutMainMenuElements(cx, cy float64) {
	px := cx - (INPUT_WIDTH / 2)
	py := cy + (INPUT_HEIGHT / 2) + BUTTON_GAP
	u.moveInput(INPUT_CHANNEL, px, py)

	px = cx - (BUTTON_WIDTH / 2)
	py = py + INPUT_HEIGHT + BUTTON_GAP
	u.moveButton(PLAY_BUTTON, px, py)

	py = py + BUTTON_HEIGHT + BUTTON_GAP
	u.moveButton(OPTIONS_BUTTON, px, py)

	px = px + (BUTTON_WIDTH - SMALL_BUTTON_WIDTH)
	u.moveButton(ABOUT_BUTTON, px, py)
}

func (u *uiImpl) layoutLicenseElements(cx, cy float64) {
	py := cy + BUTTON_GAP
	w, h := u.getLabel(LABEL_LICENSE).Measure()
	u.moveLabel(LABEL_LICENSE, cx-(w/2), py)

	px := cx - (SMALL_BUTTON_WIDTH / 2)
	py = py + h + BUTTON_GAP
	u.moveButton(ACCEPT_LICENSE_BUTTON, px, py)
}

func (u *uiImpl) layoutAboutSubMenuElements(cx, cy float64) {
	py := cy + BUTTON_GAP
	w, h := u.getLabel(LABEL_ABOUT_MESSAGE).Measure()
	u.moveLabel(LABEL_ABOUT_MESSAGE, cx-(w/2), py)

	px := cx - (SMALL_BUTTON_WIDTH / 2)
	py = py + h + BUTTON_GAP
	u.moveButton(SUBMENU_ABOUT_BACK_BUTTON, px, py)
}

func (u *uiImpl) layoutOptionsSubMenuElements(cx, cy float64) {
	px := cx - (SMALL_BUTTON_WIDTH / 2)
	py := cy
	u.moveButton(SUBMENU_OPTION_BACK_BUTTON, px, py)
}

func (u *uiImpl) OnLayoutChange(width, height float64) {
	u.screenWidth = width
	u.screenHeight = height

	cx := u.screenWidth / 2
	cy := MENU_START

	u.layoutMainElements(cx, cy)
	cy += TITLE_TO_ELEMENTS_SEPARATION

	u.layoutMainMenuElements(cx, cy)
	u.layoutAboutSubMenuElements(cx, cy)
	u.layoutOptionsSubMenuElements(cx, cy)
	u.layoutLicenseElements(cx, cy)
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

func (u *uiImpl) addButton(id button.Id, w, h float64, face *text.GoTextFace, label string, state button.State) {
	u.buttons = append(u.buttons, button.New(id, w, h, label, face, u.audioPlayer, state))
}

func (u *uiImpl) SetButtonVisible(id button.Id, visible bool) {
	if b := u.getButton(id); b != nil {
		b.SetVisible(visible)
	}
}

func (u *uiImpl) ClickButton(id button.Id) {
	if b := u.getButton(id); b != nil {
		b.Click()
	}
}

func (u *uiImpl) SetButtonClickCallback(callback func(id button.Id)) {
	for i := range u.buttons {
		u.buttons[i].OnButtonClickCallback(callback)
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

func (ui uiImpl) IsInputEditing(id input.Id) bool {
	if i := ui.getInput(id); i != nil {
		return i.IsEditing()
	}
	return false
}

func (u *uiImpl) addInput(id input.Id, w, h float64, maxLength int, initialText string, placeHolder string) {
	u.inputs = append(u.inputs, input.New(id, w, h, maxLength, initialText, placeHolder, u.normalFace, u.audioPlayer))
}

func (u *uiImpl) updateInputs(mouseX, mouseY float64, leftPressed bool, elapsedTime int) {
	ebiten.SetCursorShape(ebiten.CursorShapeDefault)
	for i := range u.inputs {
		u.inputs[i].Update(mouseX, mouseY, leftPressed, u.keys, elapsedTime)
	}
}

func (u *uiImpl) AddLabel(id label.Id, text string, face *text.GoTextFace, lineHeight float64, color color.Color) {
	u.labels = append(u.labels, label.New(id, text, face, lineHeight, color))
}

func (u *uiImpl) getLabel(id label.Id) label.Label {
	for i := range u.labels {
		if u.labels[i].GetId() == id {
			return u.labels[i]
		}
	}
	return nil
}

func (u *uiImpl) moveLabel(id label.Id, x, y float64) {
	if l := u.getLabel(id); l != nil {
		l.Move(x, y)
	}
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

func New(audioPlayer audio.Player) UI {
	return &uiImpl{
		audioPlayer: audioPlayer,
	}
}
