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

package ui

import (
	"embed"
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
	"github.com/juan-medina/twitch-rat/internal/ui/button/imageButton"
	"github.com/juan-medina/twitch-rat/internal/ui/button/textButton"
	"github.com/juan-medina/twitch-rat/internal/ui/input"
	"github.com/juan-medina/twitch-rat/internal/ui/label"
	"github.com/juan-medina/twitch-rat/internal/ui/panel"
	"github.com/juan-medina/twitch-rat/internal/ui/scores"
	"github.com/juan-medina/twitch-rat/internal/ui/slider"
)

type Widget interface {
	Draw(screen *ebiten.Image)
	Update(elapsedTime int, mouseX, mouseY float64, leftJustPressed bool, leftPressed bool, keys keys.Keys)
}
type UI interface {
	Init(fileSystem embed.FS, keys keys.Keys, sheet draw.Sheet)
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
	SetLabelBackgroundColor(id label.Id, color color.Color, expand float64)

	SetSliderVisible(id slider.Id, visible bool)
	SetSliderValue(id slider.Id, value float64)
	SetSliderChangeCallback(callback func(id slider.Id, value float64))
	AddFlyingText(text string, color colors.CustomColor, x, y float64)

	SetScoreVisible(visible bool)
	AddScoreEntry(data scores.ScoreData)
	StartsCore()

	SetPanelVisible(id panel.Id, visible bool)
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
	screenWidth   float64
	screenHeight  float64
	fileSystem    embed.FS
	buttons       []button.Button
	inputs        []input.Input
	labels        []label.Label
	sliders       []slider.Slider
	panels        []panel.Panel
	keys          keys.Keys
	audioPlayer   audio.Player
	fontVerySmall draw.Font
	fontSmall     draw.Font
	fontNormal    draw.Font
	fontBig       draw.Font
	flyingTexts   []flyingText
	scores        scores.Scores
	sheet         draw.Sheet
	widgets       []Widget
}

const (
	MAX_WIDGETS                  = 200
	MAX_BUTTONS                  = 10
	MENU_START                   = 250.0
	BUTTON_WIDTH                 = 180.0
	BUTTON_HEIGHT                = 50.0
	SMALL_BUTTON_WIDTH           = 85.0
	SMALL_BUTTON_HEIGHT          = 25.0
	BUTTON_GAP                   = 20.0
	TITLE_TO_ELEMENTS_SEPARATION = 50.0
	MAX_INPUTS                   = 1
	INPUT_WIDTH                  = BUTTON_WIDTH*2.0 + BUTTON_GAP*2.0
	INPUT_HEIGHT                 = 50
	MAX_LABELS                   = 10
	SLIDER_WITH                  = 300
	SLIDER_HEIGHT                = 30
	TOTAL_LAST_MESSAGES          = 15
	MAX_FLYING_TEXTS             = 20
	FLYING_TEXT_VY               = 0.10
	FLYING_TIME_FULL             = 1000
	FLYING_TIME_TO_VANISH        = 500
	FLYING_TIME_TO_FREE          = 1
	SCORE_X                      = 10
	SCORE_Y                      = 10
	BACKGROUND_EXPAND            = 5
	MAX_PANELS                   = 5
	OPTION_PANEL_WIDTH           = 450
	OPTION_PANEL_HEIGHT          = 220
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
	IN_GAME_OPTIONS_BUTTON
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
	LABEL_COUNTDOWN
	LABEL_INSTRUCTIONS
	LABEL_LAST_MESSAGE
	LABEL_LAST_MESSAGE_LAST = LABEL_LAST_MESSAGE + TOTAL_LAST_MESSAGES
)

const (
	MUSIC_VOLUME_SLIDER slider.Id = iota
	AUDIO_VOLUME_SLIDER
)

const (
	OPTIONS_PANEL panel.Id = iota
)

func (u *uiImpl) Init(fileSystem embed.FS, keys keys.Keys, sheet draw.Sheet) {
	u.sheet = sheet
	u.fileSystem = fileSystem
	u.keys = keys

	u.labels = make([]label.Label, 0, MAX_LABELS)
	u.widgets = make([]Widget, 0, MAX_WIDGETS)
	u.buttons = make([]button.Button, 0, MAX_BUTTONS)
	u.inputs = make([]input.Input, 0, MAX_INPUTS)
	u.panels = make([]panel.Panel, 0, MAX_PANELS)
	u.widgets = append(u.widgets, u.scores)

	u.scores.Init()

	data, _ := fileSystem.ReadFile("embed/text/license.txt")
	licenseText := string(data)

	data, _ = fileSystem.ReadFile("embed/text/about.txt")
	aboutText := string(data)

	data, _ = fileSystem.ReadFile("embed/text/instructions.txt")
	instructions := string(data)

	u.addLabel(LABEL_TITLE, "Twitch Rat", u.fontBig, normalLabelColor)

	versionStr := "0.0.0.0"
	if data, err := fileSystem.ReadFile("embed/version.txt"); err == nil {
		versionStr = string(data)
	}
	parts := strings.Split(versionStr, ".")
	versionStr = colors.Blue.BBCoded("v") +
		colors.Green.BBCoded(parts[0]) + "." +
		colors.Yellow.BBCoded(parts[1]) + "." +
		colors.Orange.BBCoded(parts[2]) + "." +
		colors.Red.BBCoded(parts[3])
	u.addLabel(LABEL_VERSION, versionStr, u.fontSmall, normalLabelColor)

	u.addInput(INPUT_CHANNEL, INPUT_WIDTH, INPUT_HEIGHT, 24, "", "Twitch Channel", u.fontNormal)
	u.addTextButton(PLAY_BUTTON, BUTTON_WIDTH, BUTTON_HEIGHT, u.fontNormal, "Play!", button.Enabled)
	u.addImageButton(BACK_BUTTON, "exitLeft", button.Enabled)
	u.addImageButton(IN_GAME_OPTIONS_BUTTON, "gear", button.Enabled)

	u.addTextButton(OPTIONS_BUTTON, SMALL_BUTTON_WIDTH, SMALL_BUTTON_HEIGHT, u.fontSmall, "Options", button.Enabled)
	u.addTextButton(ABOUT_BUTTON, SMALL_BUTTON_WIDTH, SMALL_BUTTON_HEIGHT, u.fontSmall, "About", button.Enabled)

	u.addLabel(LABEL_ABOUT_MESSAGE, aboutText, u.fontSmall, normalLabelColor)

	u.SetLabelBackgroundColor(LABEL_ABOUT_MESSAGE, colors.Black.NewWithAlpha(50), BACKGROUND_EXPAND)
	u.addTextButton(SUBMENU_ABOUT_BACK_BUTTON, SMALL_BUTTON_WIDTH, SMALL_BUTTON_HEIGHT, u.fontSmall, "Back", button.Enabled)

	u.addLabel(LABEL_LICENSE, licenseText, u.fontSmall, normalLabelColor)

	u.addTextButton(ACCEPT_LICENSE_BUTTON, SMALL_BUTTON_WIDTH, SMALL_BUTTON_HEIGHT, u.fontSmall, "Accept", button.Enabled)

	u.addInput(INPUT_DEBUG_USER, INPUT_WIDTH, INPUT_HEIGHT, 24, "", "User", u.fontNormal)
	u.addInput(INPUT_DEBUG_MESSAGE, INPUT_WIDTH, INPUT_HEIGHT, 24, "", "Message", u.fontNormal)
	u.addTextButton(DEBUG_BUTTON, BUTTON_WIDTH, BUTTON_HEIGHT, u.fontNormal, "Debug", button.Enabled)

	for i := 0; i < TOTAL_LAST_MESSAGES; i++ {
		id := LABEL_LAST_MESSAGE + label.Id(i)
		u.addLabel(id, "", u.fontSmall, normalLabelColor)
	}

	for i := 0; i < TOTAL_LAST_MESSAGES; i++ {
		u.SetLabelVisible(LABEL_LAST_MESSAGE+label.Id(i), true)
	}

	u.addLabel(LABEL_COUNTDOWN, "30", u.fontBig, normalLabelColor)
	u.addLabel(LABEL_INSTRUCTIONS, instructions, u.fontNormal, normalLabelColor)
	u.SetLabelBackgroundColor(LABEL_INSTRUCTIONS, colors.Black.NewWithAlpha(50), BACKGROUND_EXPAND)

	u.addPanel(OPTIONS_PANEL, OPTION_PANEL_WIDTH, OPTION_PANEL_HEIGHT, colors.Black.NewWithAlpha(50))
	u.addTextButton(SUBMENU_OPTION_BACK_BUTTON, SMALL_BUTTON_WIDTH, SMALL_BUTTON_HEIGHT, u.fontSmall, "Back", button.Enabled)
	u.addLabel(LABEL_OPTIONS_MUSIC_VOLUME, "Music Volume", u.fontNormal, normalLabelColor)
	u.addSlider(MUSIC_VOLUME_SLIDER, SLIDER_WITH, SLIDER_HEIGHT, u.fontSmall, normalLabelColor)

	u.addLabel(LABEL_OPTIONS_AUDIO_VOLUME, "Audio Volume", u.fontNormal, normalLabelColor)
	u.addSlider(AUDIO_VOLUME_SLIDER, SLIDER_WITH, SLIDER_HEIGHT, u.fontSmall, normalLabelColor)

	u.SetLabelVisible(LABEL_VERSION, true)
}

func (u *uiImpl) Draw(screen *ebiten.Image) {
	for _, w := range u.widgets {
		w.Draw(screen)
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

	for _, w := range u.widgets {
		w.Update(elapsedTime, x, y, justPressed, pressed, u.keys)
	}

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
		newAlpha := uint32(255 - (i * (255 - 128) / (TOTAL_LAST_MESSAGES - 1)))
		u.SetLabelColor(currentLabelId, color.RGBA{uint8(r), uint8(g), uint8(b), uint8(newAlpha)})
	}

	u.SetLabelText(LABEL_LAST_MESSAGE, message)
	u.SetLabelColor(LABEL_LAST_MESSAGE, textColor)
}

func (u *uiImpl) layoutCounter() {
	cx := u.screenWidth / 2
	cy := MENU_START

	countdownLabel := u.getLabel(LABEL_COUNTDOWN)
	ix, iy := countdownLabel.Measure()
	px := cx - ix
	py := cy
	countdownLabel.Move(px, py)

	py = py + iy + BUTTON_GAP
	w, _ := u.getLabel(LABEL_INSTRUCTIONS).Measure()
	u.moveLabel(LABEL_INSTRUCTIONS, cx-(w/2), py)

	cx = (u.screenWidth / 2) - (OPTION_PANEL_WIDTH / 2)
	cy = cy + 100
	u.movePanel(OPTIONS_PANEL, cx, cy)
}

func (u *uiImpl) layoutMainElements(cx, cy float64) {

	titleLabel := u.getLabel(LABEL_TITLE)
	ix, _ := titleLabel.Measure()
	px := cx - (ix / 2)
	py := cy
	titleLabel.Move(px, py)

	bb := u.getButton(BACK_BUTTON)
	bw, _ := bb.Size()
	px = u.screenWidth - bw
	py = 0
	u.moveButton(BACK_BUTTON, px, py)

	bb = u.getButton(IN_GAME_OPTIONS_BUTTON)
	bw, _ = bb.Size()
	px -= bw
	u.moveButton(IN_GAME_OPTIONS_BUTTON, px, py)

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

	px = 450
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
	u.layoutCounter()

	u.scores.Move(SCORE_X, SCORE_Y)
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

func (u *uiImpl) addTextButton(id button.Id, w, h float64, font draw.Font, label string, state button.State) {
	tb := textButton.New(id, w, h, label, font, font.DefaultSize(), u.audioPlayer, state)
	u.buttons = append(u.buttons, tb)
	u.widgets = append(u.widgets, tb)
}

func (u *uiImpl) addImageButton(id button.Id, spriteName string, state button.State) {
	ib := imageButton.New(id, u.sheet.Sprite(spriteName), u.audioPlayer, state)
	u.buttons = append(u.buttons, ib)
	u.widgets = append(u.widgets, ib)
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
	i := input.New(id, w, h, maxLength, initialText, placeHolder, font, font.DefaultSize(), u.audioPlayer)
	u.inputs = append(u.inputs, i)
	u.widgets = append(u.widgets, i)
}

func (u *uiImpl) addLabel(id label.Id, text string, font draw.Font, color color.Color) {
	l := label.NewLabel(id, text, font, font.DefaultSize(), color, u.audioPlayer)
	u.labels = append(u.labels, l)
	u.widgets = append(u.widgets, l)
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
		if id == LABEL_COUNTDOWN {
			u.layoutCounter()
		}
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
	s := slider.New(id, w, h, font, font.DefaultSize(), labelColor)
	ui.sliders = append(ui.sliders, s)
	ui.widgets = append(ui.widgets, s)
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
func (u *uiImpl) AddScoreEntry(data scores.ScoreData) {
	u.scores.Add(data)
}

func (u *uiImpl) SetScoreVisible(visible bool) {
	u.scores.SetVisible(visible)
	if visible {
		u.scores.Reset()
	}
}

func (u *uiImpl) StartsCore() {
	u.scores.Start()
}

func (u *uiImpl) SetLabelBackgroundColor(id label.Id, color color.Color, expand float64) {
	if l := u.getLabel(id); l != nil {
		l.SetBackgroundColor(color, expand)
	}
}

func (u *uiImpl) addPanel(id panel.Id, w, h float64, color color.Color) {
	p := panel.New(id, w, h, color)
	u.panels = append(u.panels, p)
	u.widgets = append(u.widgets, p)
}

func (u *uiImpl) getPanel(id panel.Id) panel.Panel {
	for i := range u.panels {
		if u.panels[i].GetId() == id {
			return u.panels[i]
		}
	}
	return nil
}

func (u *uiImpl) SetPanelVisible(id panel.Id, visible bool) {
	if p := u.getPanel(id); p != nil {
		p.SetVisible(visible)
	}
}

func (u *uiImpl) movePanel(id panel.Id, x, y float64) {
	if p := u.getPanel(id); p != nil {
		p.Move(x, y)
	}
}

func New(audioPlayer audio.Player, fontVerySmall draw.Font, fontSmall draw.Font, fontNormal draw.Font, fontBig draw.Font) UI {
	return &uiImpl{
		audioPlayer:   audioPlayer,
		fontVerySmall: fontVerySmall,
		fontSmall:     fontSmall,
		fontNormal:    fontNormal,
		fontBig:       fontBig,
		flyingTexts:   make([]flyingText, 0, MAX_FLYING_TEXTS),
		scores:        scores.New(fontSmall, fontNormal),
	}
}
