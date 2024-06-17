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
	"embed"
	"fmt"
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/juan-medina/twitch-rat/internal/audio"
	"github.com/juan-medina/twitch-rat/internal/colors"
	"github.com/juan-medina/twitch-rat/internal/draw"
	"github.com/juan-medina/twitch-rat/internal/keys"
	"github.com/juan-medina/twitch-rat/internal/step"
	"github.com/juan-medina/twitch-rat/internal/ui/button"
	"github.com/juan-medina/twitch-rat/internal/ui/input"
	"github.com/juan-medina/twitch-rat/internal/ui/label"
	"github.com/juan-medina/twitch-rat/internal/ui/slider"
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
	GetLabelColor(id label.Id) color.Color
	SetLabelVisible(id label.Id, visible bool)

	SetSliderVisible(id slider.Id, visible bool)
	SetSliderValue(id slider.Id, value float64)
	SetSliderChangeCallback(callback func(id slider.Id, value float64))
	AddFlyingText(text string, color colors.CustomColor, x, y float64)
}

var (
	normalLabelColor = colors.White
)

type flyingText struct {
	text    string
	color   colors.CustomColor
	alpha   step.Value
	x, y    float64
	vy      float64
	visible bool
}

type uiImpl struct {
	screenWidth  float64
	screenHeight float64
	fileSystem   embed.FS
	buttons      []button.Button
	inputs       []input.Input
	labels       []label.Label
	sliders      []slider.Slider
	keys         keys.Keys
	audioPlayer  audio.Player
	licenseText  string
	aboutText    string
	fontSmall    draw.Font
	fontNormal   draw.Font
	fontBig      draw.Font
	flyingTexts  []flyingText
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
	SLIDER_WITH                  = 300
	SLIDER_HEIGHT                = 30
	TOTAL_LAST_MESSAGES          = 10
	MAX_FLYING_TEXTS             = 20
	FLYING_TEXT_VY               = 0.10
	FLYING_TIME_FULL             = 1000
	FLYING_TIME_TO_VANISH        = 500
	FLYING_TIME_TO_FREE          = 1
)

const (
	PLAY_BUTTON button.Id = iota
	BACK_BUTTON
	OPTIONS_BUTTON
	ABOUT_BUTTON
	SUBMENU_ABOUT_BACK_BUTTON
	SUBMENU_OPTION_BACK_BUTTON
	ACCEPT_LICENSE_BUTTON
	DEBUG_BUTTON
)

const (
	INPUT_CHANNEL input.Id = iota
	INPUT_DEBUG_USER
	INPUT_DEBUG_MESSAGE
)

const (
	LABEL_VERSION label.Id = iota
	LABEL_TITLE
	LABEL_LICENSE
	LABEL_ABOUT_MESSAGE
	LABEL_OPTIONS_MUSIC_VOLUME
	LABEL_OPTIONS_AUDIO_VOLUME
	LABEL_LAST_MESSAGE
	LABEL_LAST_MESSAGE_LAST = LABEL_LAST_MESSAGE + TOTAL_LAST_MESSAGES
)

const (
	MUSIC_VOLUME_SLIDER slider.Id = iota
	AUDIO_VOLUME_SLIDER
)

func (u *uiImpl) Init(fileSystem embed.FS, keys keys.Keys) {
	u.fileSystem = fileSystem
	u.keys = keys

	data, _ := fileSystem.ReadFile("embed/text/license.txt")
	u.licenseText = string(data)

	data, _ = fileSystem.ReadFile("embed/text/about.txt")
	u.aboutText = string(data)

	for i := 0; i < TOTAL_LAST_MESSAGES; i++ {
		id := LABEL_LAST_MESSAGE + label.Id(i)
		u.addLabel(id, "", u.fontSmall, normalLabelColor)
	}

	u.addLabel(LABEL_TITLE, "Twitch Rat", u.fontBig, normalLabelColor)

	versionStr := "0.0.0.0"
	if data, err := fileSystem.ReadFile("embed/version.txt"); err == nil {
		versionStr = string(data)
	}
	parts := strings.Split(versionStr, ".")
	versionStr = fmt.Sprintf("%sv%s%s%s.%s%s%s.%s%s%s.%s%s%s",
		colors.Blue.Tag(),
		colors.Green.Tag(), parts[0], colors.White.Tag(),
		colors.Yellow.Tag(), parts[1], colors.White.Tag(),
		colors.Orange.Tag(), parts[2], colors.White.Tag(),
		colors.Red.Tag(), parts[3], colors.White.Tag())
	u.addLabel(LABEL_VERSION, versionStr, u.fontSmall, normalLabelColor)

	u.buttons = make([]button.Button, 0, MAX_BUTTONS)
	u.inputs = make([]input.Input, 0, MAX_INPUTS)

	u.addInput(INPUT_CHANNEL, INPUT_WIDTH, INPUT_HEIGHT, 24, "", "Twitch Channel", u.fontNormal)
	u.addButton(PLAY_BUTTON, BUTTON_WIDTH, BUTTON_HEIGHT, u.fontNormal, "Play!", button.Enabled)
	u.addButton(BACK_BUTTON, BACK_BUTTON_WIDTH, BACK_BUTTON_HEIGHT, u.fontSmall, "X", button.Enabled)
	u.addButton(OPTIONS_BUTTON, SMALL_BUTTON_WIDTH, SMALL_BUTTON_HEIGHT, u.fontSmall, "Options", button.Enabled)
	u.addButton(ABOUT_BUTTON, SMALL_BUTTON_WIDTH, SMALL_BUTTON_HEIGHT, u.fontSmall, "About", button.Enabled)

	u.addLabel(LABEL_ABOUT_MESSAGE, u.aboutText, u.fontSmall, normalLabelColor)
	u.addButton(SUBMENU_ABOUT_BACK_BUTTON, SMALL_BUTTON_WIDTH, SMALL_BUTTON_HEIGHT, u.fontSmall, "Back", button.Enabled)
	u.addButton(SUBMENU_OPTION_BACK_BUTTON, SMALL_BUTTON_WIDTH, SMALL_BUTTON_HEIGHT, u.fontSmall, "Back", button.Enabled)

	u.addLabel(LABEL_LICENSE, u.licenseText, u.fontSmall, normalLabelColor)
	u.addButton(ACCEPT_LICENSE_BUTTON, SMALL_BUTTON_WIDTH, SMALL_BUTTON_HEIGHT, u.fontSmall, "Accept", button.Enabled)

	u.addLabel(LABEL_OPTIONS_MUSIC_VOLUME, "Music Volume", u.fontNormal, normalLabelColor)
	u.addSlider(MUSIC_VOLUME_SLIDER, SLIDER_WITH, SLIDER_HEIGHT, u.fontSmall, normalLabelColor)

	u.addLabel(LABEL_OPTIONS_AUDIO_VOLUME, "Audio Volume", u.fontNormal, normalLabelColor)
	u.addSlider(AUDIO_VOLUME_SLIDER, SLIDER_WITH, SLIDER_HEIGHT, u.fontSmall, normalLabelColor)

	u.addInput(INPUT_DEBUG_USER, INPUT_WIDTH, INPUT_HEIGHT, 24, "", "User", u.fontNormal)
	u.addInput(INPUT_DEBUG_MESSAGE, INPUT_WIDTH, INPUT_HEIGHT, 24, "", "Message", u.fontNormal)
	u.addButton(DEBUG_BUTTON, BUTTON_WIDTH, BUTTON_HEIGHT, u.fontNormal, "Debug", button.Enabled)

	for i := 0; i < TOTAL_LAST_MESSAGES; i++ {
		u.SetLabelVisible(LABEL_LAST_MESSAGE+label.Id(i), true)
	}

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

	for _, s := range u.sliders {
		s.Draw(screen)
	}

	for _, f := range u.flyingTexts {
		f.draw(screen, u.fontSmall)
	}
}

func (f flyingText) draw(screen *ebiten.Image, font draw.Font) {
	if !f.visible {
		return
	}
	font.Draw(screen, f.text, f.x, f.y, font.DefaultSize(), f.color)
}

func (u *uiImpl) Update(elapsedTime int) {
	xe, ye := ebiten.CursorPosition()
	x, y := float64(xe), float64(ye)

	justPressed := inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
	pressed := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)

	ebiten.SetCursorShape(ebiten.CursorShapeDefault)
	u.updateButtons(x, y, justPressed, elapsedTime)
	u.updateInputs(x, y, justPressed, elapsedTime)
	u.updateSliders(x, y, justPressed, pressed, elapsedTime)
	u.updateFlyingTexts(elapsedTime)
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

func (u *uiImpl) SetStatusMessage(message string, textColor color.Color) {
	for i := TOTAL_LAST_MESSAGES - 1; i > 0; i-- {
		prevLabelId := LABEL_LAST_MESSAGE + label.Id(i-1)
		currentLabelId := LABEL_LAST_MESSAGE + label.Id(i)

		u.SetLabelText(currentLabelId, u.GetLabelText(prevLabelId))
		prevColor := u.GetLabelColor(prevLabelId)
		r, g, b, _ := prevColor.RGBA()
		newAlpha := uint32(255 - (i * (255 - 64) / (TOTAL_LAST_MESSAGES - 1)))
		u.SetLabelColor(currentLabelId, color.RGBA{uint8(r), uint8(g), uint8(b), uint8(newAlpha)})
	}

	u.SetLabelText(LABEL_LAST_MESSAGE, message)
	u.SetLabelColor(LABEL_LAST_MESSAGE, textColor)
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

	gapX := float64(u.fontSmall.DefaultSize()) * 0.5
	gapY := float64(u.fontSmall.DefaultSize()) * 1.5

	py = u.screenHeight - gapY

	for i := 0; i < TOTAL_LAST_MESSAGES; i++ {
		labelId := LABEL_LAST_MESSAGE + label.Id(i)
		u.getLabel(labelId).Move(gapX, py)
		py -= u.fontSmall.DefaultSize()
	}

	gapX = u.fontSmall.DefaultSize() * 0.5
	gapY = u.fontSmall.DefaultSize() * 1.5
	versionLabel := u.getLabel(LABEL_VERSION)
	cx, _ = versionLabel.Measure()
	versionLabel.Move(u.screenWidth-cx-gapX, u.screenHeight-gapY)

	px = 0
	py = 0
	u.getInput(INPUT_DEBUG_USER).Move(px, py)

	px = px + INPUT_WIDTH + BUTTON_GAP
	u.getInput(INPUT_DEBUG_MESSAGE).Move(px, py)

	px = px + INPUT_WIDTH + BUTTON_GAP
	u.getButton(DEBUG_BUTTON).Move(px, py)
}

func (u *uiImpl) layoutMainMenuElements(cx, cy float64) {
	px := cx - (INPUT_WIDTH / 2)
	py := cy + (INPUT_HEIGHT / 2) + BUTTON_GAP*2
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
	py := cy + BUTTON_GAP*3
	w, h := u.getLabel(LABEL_LICENSE).Measure()
	u.moveLabel(LABEL_LICENSE, cx-(w/2), py)

	px := cx - (SMALL_BUTTON_WIDTH / 2)
	py = py + h + BUTTON_GAP
	u.moveButton(ACCEPT_LICENSE_BUTTON, px, py)
}

func (u *uiImpl) layoutAboutSubMenuElements(cx, cy float64) {
	py := cy + BUTTON_GAP*3
	w, h := u.getLabel(LABEL_ABOUT_MESSAGE).Measure()
	u.moveLabel(LABEL_ABOUT_MESSAGE, cx-(w/2), py)

	px := cx - (SMALL_BUTTON_WIDTH / 2)
	py = py + h + BUTTON_GAP
	u.moveButton(SUBMENU_ABOUT_BACK_BUTTON, px, py)
}

func (u *uiImpl) layoutOptionsSubMenuElements(cx, cy float64) {
	px := cx - SLIDER_WITH/2
	py := cy + BUTTON_GAP*3

	u.moveLabel(LABEL_OPTIONS_MUSIC_VOLUME, px, py)

	py = py + BUTTON_GAP*2
	u.moveSlider(MUSIC_VOLUME_SLIDER, px, py)

	py = py + BUTTON_GAP*2
	u.moveLabel(LABEL_OPTIONS_AUDIO_VOLUME, px, py)

	py = py + BUTTON_GAP*2
	u.moveSlider(AUDIO_VOLUME_SLIDER, px, py)

	px = cx - (SMALL_BUTTON_WIDTH / 2)
	py = py + SLIDER_HEIGHT + BUTTON_GAP
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

func (u *uiImpl) addButton(id button.Id, w, h float64, font draw.Font, label string, state button.State) {
	u.buttons = append(u.buttons, button.New(id, w, h, label, font, font.DefaultSize(), u.audioPlayer, state))
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

func (u *uiImpl) addInput(id input.Id, w, h float64, maxLength int, initialText string, placeHolder string, font draw.Font) {
	u.inputs = append(u.inputs, input.New(id, w, h, maxLength, initialText, placeHolder, font, font.DefaultSize(), u.audioPlayer))
}

func (u *uiImpl) updateInputs(mouseX, mouseY float64, leftPressed bool, elapsedTime int) {
	for i := range u.inputs {
		u.inputs[i].Update(mouseX, mouseY, leftPressed, u.keys, elapsedTime)
	}
}

func (u *uiImpl) addLabel(id label.Id, text string, font draw.Font, color color.Color) {
	u.labels = append(u.labels, label.NewLabel(id, text, font, font.DefaultSize(), color))
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

func (u *uiImpl) GetLabelText(id label.Id) string {
	if l := u.getLabel(id); l != nil {
		return l.GetText()
	}
	return ""
}

func (u *uiImpl) SetLabelColor(id label.Id, color color.Color) {
	if l := u.getLabel(id); l != nil {
		l.SetColor(color)
	}
}

func (u uiImpl) GetLabelColor(id label.Id) color.Color {
	if l := u.getLabel(id); l != nil {
		return l.GetColor()
	}
	return colors.White
}

func (ui *uiImpl) SetLabelVisible(id label.Id, visible bool) {
	if l := ui.getLabel(id); l != nil {
		l.SetVisible(visible)
	}
}

func (ui *uiImpl) getSlider(id slider.Id) slider.Slider {
	for i := range ui.sliders {
		if ui.sliders[i].GetId() == id {
			return ui.sliders[i]
		}
	}
	return nil
}

func (ui *uiImpl) moveSlider(id slider.Id, x, y float64) {
	if s := ui.getSlider(id); s != nil {
		s.Move(x, y)
	}
}

func (ui *uiImpl) SetSliderVisible(id slider.Id, visible bool) {
	if s := ui.getSlider(id); s != nil {
		s.SetVisible(visible)
	}
}
func (ui *uiImpl) SetSliderValue(id slider.Id, value float64) {
	if s := ui.getSlider(id); s != nil {
		s.SetValue(value)
	}
}

func (ui *uiImpl) addSlider(id slider.Id, w, h float64, font draw.Font, labelColor color.Color) {
	ui.sliders = append(ui.sliders, slider.New(id, w, h, font, font.DefaultSize(), labelColor))
}

func (u *uiImpl) updateSliders(mouseX, mouseY float64, leftJustPressed bool, leftPressed bool, elapsedTime int) {
	for i := range u.sliders {
		u.sliders[i].Update(mouseX, mouseY, leftJustPressed, leftPressed, elapsedTime)
	}
}

func (u *uiImpl) SetSliderChangeCallback(callback func(id slider.Id, value float64)) {
	for i := range u.sliders {
		u.sliders[i].OnValueChangeCallback(callback)
	}
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

func New(audioPlayer audio.Player, fontSmall draw.Font, fontNormal draw.Font, fontBig draw.Font) UI {
	return &uiImpl{
		audioPlayer: audioPlayer,
		fontSmall:   fontSmall,
		fontNormal:  fontNormal,
		fontBig:     fontBig,
		flyingTexts: make([]flyingText, 0, MAX_FLYING_TEXTS),
	}
}
