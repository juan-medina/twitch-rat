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
	"fmt"
	"image/color"
	"runtime"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/juan-medina/twitch-rat/internal/colors"
	"github.com/juan-medina/twitch-rat/internal/draw"
	"github.com/juan-medina/twitch-rat/internal/keys"
	"github.com/juan-medina/twitch-rat/internal/ui/button"
	"github.com/juan-medina/twitch-rat/internal/ui/button/imageButton"
	"github.com/juan-medina/twitch-rat/internal/ui/button/textButton"
	"github.com/juan-medina/twitch-rat/internal/ui/input"
	"github.com/juan-medina/twitch-rat/internal/ui/label"
	"github.com/juan-medina/twitch-rat/internal/ui/panel"
	"github.com/juan-medina/twitch-rat/internal/ui/radiogroup"
	"github.com/juan-medina/twitch-rat/internal/ui/slider"
)

type Widget interface {
	Draw(screen *ebiten.Image)
	Update(elapsedTime int, mouseX, mouseY float64, leftJustPressed bool, leftPressed bool, keys keys.Keys)
}

const (
	PLAY_BUTTON button.Id = iota
	BACK_BUTTON
	FULL_SCREEN_BUTTON
	WINDOWED_BUTTON
	OPTIONS_BUTTON
	ABOUT_BUTTON
	SUBMENU_ABOUT_BACK_BUTTON
	SUBMENU_OPTION_BACK_BUTTON
	ACCEPT_LICENSE_BUTTON
	DEBUG_BUTTON
	IN_GAME_OPTIONS_BUTTON
	DOWNLOAD_LATEST_BUTTON
	CONTINUE_BUTTON
	SUBMENU_GAME_MODE_SETTINGS_BACK_BUTTON
	SUBMENU_GAME_MODE_SETTINGS_GO_BUTTON
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
	LABEL_DOWNLOAD
	LABEL_VERSION_UPDATE
	LABEL_GAME_MODE
	LABEL_GAME_MODE_DESCRIPTION
	LABEL_JOIN_MODE
	LABEL_REJOIN_MODE
	LABEL_ATTACK_MODE
	LABEL_HEAL_MODE
	LABEL_HEALING_TO
	LABEL_RAT_HEALTH
	LABEL_RAT_DAMAGE
	LABEL_RAT_HEALING
	LABEL_LAST_MESSAGE
	LABEL_LAST_MESSAGE_LAST = LABEL_LAST_MESSAGE + TOTAL_LAST_MESSAGES
)

const (
	MUSIC_VOLUME_SLIDER slider.Id = iota
	AUDIO_VOLUME_SLIDER
	RAT_HEALTH_SLIDER
	RAT_DAMAGE_SLIDER
	RAT_HEALING_SLIDER
)

const (
	OPTIONS_PANEL panel.Id = iota
)

var (
	normalLabelColor = colors.White
)

const (
	GAME_MODE_RADIO_GROUP radiogroup.Id = iota
	JOIN_MODE_RADIO_GROUP
	REJOIN_MODE_RADIO_GROUP
	ATTACK_MODE_RADIO_GROUP
	HEAL_MODE_RADIO_GROUP
	HEALING_TO_RADIO_GROUP
)

func (u *uiImpl) createWidgets() {
	u.createMainComponents()

	u.createMainMenuUI()
	u.createAboutMenuUI()
	u.createOptionsMenuUI()
	u.createMatchSettingsMenuUI()

	u.createLicenseStageUI()
	u.createUpdateStageUI()
	u.createPlayingStageUI()

	u.createLastMessages()
	u.createDebugUI()
	u.SetLabelVisible(LABEL_VERSION, true)
}

func (u *uiImpl) createPlayingStageUI() {
	/*data, _ := u.fileSystem.ReadFile("embed/text/instructions.txt")
	instructions := string(data)*/

	u.addLabel(LABEL_COUNTDOWN, "30", u.fontBig, normalLabelColor)
	u.addLabel(LABEL_INSTRUCTIONS, "instructions", u.fontNormal, normalLabelColor)
	u.SetLabelBackgroundColor(LABEL_INSTRUCTIONS, colors.Black.NewWithAlpha(50), BACKGROUND_EXPAND)
}

func (u *uiImpl) createOptionsMenuUI() {
	u.addPanel(OPTIONS_PANEL, OPTION_PANEL_WIDTH, OPTION_PANEL_HEIGHT, colors.Black.NewWithAlpha(50))
	u.addTextButton(SUBMENU_OPTION_BACK_BUTTON, SMALL_BUTTON_WIDTH, SMALL_BUTTON_HEIGHT, u.fontSmall, "Back", button.Enabled)
	u.addLabel(LABEL_OPTIONS_MUSIC_VOLUME, "Music Volume", u.fontNormal, normalLabelColor)
	u.addSlider(MUSIC_VOLUME_SLIDER, 0, 100, SLIDER_WITH, SLIDER_HEIGHT, u.fontSmall, normalLabelColor)

	u.addLabel(LABEL_OPTIONS_AUDIO_VOLUME, "Audio Volume", u.fontNormal, normalLabelColor)
	u.addSlider(AUDIO_VOLUME_SLIDER, 0, 100, SLIDER_WITH, SLIDER_HEIGHT, u.fontSmall, normalLabelColor)
}

func (u *uiImpl) createDebugUI() {
	u.addInput(INPUT_DEBUG_USER, INPUT_WIDTH, INPUT_HEIGHT, 24, "", "User", u.fontNormal)
	u.addInput(INPUT_DEBUG_MESSAGE, INPUT_WIDTH, INPUT_HEIGHT, 24, "", "Message", u.fontNormal)
	u.addTextButton(DEBUG_BUTTON, BUTTON_WIDTH, BUTTON_HEIGHT, u.fontNormal, "Debug", button.Enabled)
}

func (u *uiImpl) createLicenseStageUI() {
	data, _ := u.fileSystem.ReadFile("embed/text/license.txt")
	licenseText := string(data)

	u.addLabel(LABEL_LICENSE, licenseText, u.fontSmall, normalLabelColor)
	u.addTextButton(ACCEPT_LICENSE_BUTTON, SMALL_BUTTON_WIDTH, SMALL_BUTTON_HEIGHT, u.fontSmall, "Accept", button.Enabled)
}

func (u *uiImpl) createMainMenuUI() {
	u.addInput(INPUT_CHANNEL, INPUT_WIDTH, INPUT_HEIGHT, 24, "", "Twitch Channel", u.fontNormal)
	u.addTextButton(PLAY_BUTTON, BUTTON_WIDTH, BUTTON_HEIGHT, u.fontNormal, "Play!", button.Enabled)

	u.addLabel(LABEL_DOWNLOAD, "[url=https://juan-medina.com/twitch-rat/twitch-rat.zip]Download Desktop Version[/url]", u.fontSmall, normalLabelColor)

	u.addTextButton(OPTIONS_BUTTON, SMALL_BUTTON_WIDTH, SMALL_BUTTON_HEIGHT, u.fontSmall, "Options", button.Enabled)
	u.addTextButton(ABOUT_BUTTON, SMALL_BUTTON_WIDTH, SMALL_BUTTON_HEIGHT, u.fontSmall, "About", button.Enabled)
}

func (u *uiImpl) createAboutMenuUI() {
	u.SetLabelBackgroundColor(LABEL_ABOUT_MESSAGE, colors.Black.NewWithAlpha(50), BACKGROUND_EXPAND)
	u.addTextButton(SUBMENU_ABOUT_BACK_BUTTON, SMALL_BUTTON_WIDTH, SMALL_BUTTON_HEIGHT, u.fontSmall, "Back", button.Enabled)

	data, _ := u.fileSystem.ReadFile("embed/text/about.txt")
	aboutText := string(data)
	u.addLabel(LABEL_ABOUT_MESSAGE, aboutText, u.fontSmall, normalLabelColor)
}

func (u *uiImpl) createUpdateStageUI() {
	u.addLabel(LABEL_VERSION_UPDATE, fmt.Sprintf(VERSION_OUTDATED_STRING, u.version.Current().Bbcode), u.fontNormal, normalLabelColor)
	var downloadText string
	if runtime.GOOS == "js" {
		downloadText = "Refresh"
	} else {
		downloadText = "Download"
	}
	u.addTextButton(DOWNLOAD_LATEST_BUTTON, BUTTON_WIDTH, BUTTON_HEIGHT, u.fontNormal, downloadText, button.Enabled)
	u.addTextButton(CONTINUE_BUTTON, BUTTON_WIDTH, BUTTON_HEIGHT, u.fontNormal, "Continue", button.Enabled)
}

func (u *uiImpl) createLastMessages() {
	for i := 0; i < TOTAL_LAST_MESSAGES; i++ {
		id := LABEL_LAST_MESSAGE + label.Id(i)
		u.addLabel(id, "", u.fontSmall, normalLabelColor)
	}

	for i := 0; i < TOTAL_LAST_MESSAGES; i++ {
		u.SetLabelVisible(LABEL_LAST_MESSAGE+label.Id(i), true)
	}
}

func (u *uiImpl) createMainComponents() {
	u.addLabel(LABEL_TITLE, "Twitch Rat", u.fontBig, normalLabelColor)
	u.addLabel(LABEL_VERSION, u.version.Current().Bbcode, u.fontSmall, normalLabelColor)
	u.addImageButton(IN_GAME_OPTIONS_BUTTON, "gear", button.Enabled)
	u.addImageButton(BACK_BUTTON, "exitLeft", button.Enabled)
	u.addImageButton(FULL_SCREEN_BUTTON, "larger", button.Enabled)
	u.addImageButton(WINDOWED_BUTTON, "smaller", button.Enabled)
	u.SetButtonVisible(FULL_SCREEN_BUTTON, true)
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

func (ui *uiImpl) addSlider(id slider.Id, min, max int, w, h float64, font draw.Font, labelColor color.Color) {
	s := slider.New(id, min, max, w, h, font, font.DefaultSize(), labelColor)
	ui.sliders = append(ui.sliders, s)
	ui.widgets = append(ui.widgets, s)
}

func (u *uiImpl) addPanel(id panel.Id, w, h float64, color color.Color) {
	p := panel.New(id, w, h, color)
	u.panels = append(u.panels, p)
	u.widgets = append(u.widgets, p)
}

func (u *uiImpl) getButton(id button.Id) button.Button {
	for i, b := range u.buttons {
		if b.GetId() == id {
			return u.buttons[i]
		}
	}
	return nil
}

func (u *uiImpl) getInput(id input.Id) input.Input {
	for i := range u.inputs {
		if u.inputs[i].GetId() == id {
			return u.inputs[i]
		}
	}
	return nil
}

func (u *uiImpl) getLabel(id label.Id) label.Label {
	for i := range u.labels {
		if u.labels[i].GetId() == id {
			return u.labels[i]
		}
	}
	return nil
}

func (u *uiImpl) getPanel(id panel.Id) panel.Panel {
	for i := range u.panels {
		if u.panels[i].GetId() == id {
			return u.panels[i]
		}
	}
	return nil
}

func (ui *uiImpl) getSlider(id slider.Id) slider.Slider {
	for i := range ui.sliders {
		if ui.sliders[i].GetId() == id {
			return ui.sliders[i]
		}
	}
	return nil
}

func (u *uiImpl) movePanel(id panel.Id, x, y float64) {
	if p := u.getPanel(id); p != nil {
		p.Move(x, y)
	}
}

func (u *uiImpl) changeButtonState(id button.Id, state button.State) {
	if b := u.getButton(id); b != nil {
		b.ChangeState(state)
	}
}

func (u *uiImpl) moveButton(id button.Id, x, y float64) {
	if b := u.getButton(id); b != nil {
		b.Move(x, y)
	}
}

func (ui *uiImpl) moveInput(id input.Id, x, y float64) {
	if i := ui.getInput(id); i != nil {
		i.Move(x, y)
	}
}

func (u *uiImpl) moveLabel(id label.Id, x, y float64) {
	if l := u.getLabel(id); l != nil {
		l.Move(x, y)
	}
}

func (ui *uiImpl) moveSlider(id slider.Id, x, y float64) {
	if s := ui.getSlider(id); s != nil {
		s.Move(x, y)
	}
}

func (u *uiImpl) EnableButton(id button.Id) {
	u.changeButtonState(id, button.Enabled)
}

func (u *uiImpl) DisableButton(id button.Id) {
	u.changeButtonState(id, button.Disabled)
}

func (u *uiImpl) SetButtonVisible(id button.Id, visible bool) {
	if b := u.getButton(id); b != nil {
		b.SetVisible(visible)
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

func (ui *uiImpl) SetSliderVisible(id slider.Id, visible bool) {
	if s := ui.getSlider(id); s != nil {
		s.SetVisible(visible)
	}
}
func (ui *uiImpl) SetSliderValue(id slider.Id, value int) {
	if s := ui.getSlider(id); s != nil {
		s.SetValue(value)
	}
}

func (ui uiImpl) GetSliderValue(id slider.Id) int {
	if s := ui.getSlider(id); s != nil {
		return s.GetValue()
	}
	return 0
}

func (u *uiImpl) SetSliderChangeCallback(callback func(id slider.Id, value int)) {
	for i := range u.sliders {
		u.sliders[i].OnValueChangeCallback(callback)
	}
}

func (u *uiImpl) SetLabelBackgroundColor(id label.Id, color color.Color, expand float64) {
	if l := u.getLabel(id); l != nil {
		l.SetBackgroundColor(color, expand)
	}
}

func (u *uiImpl) SetPanelVisible(id panel.Id, visible bool) {
	if p := u.getPanel(id); p != nil {
		p.SetVisible(visible)
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
	u.getButton(FULL_SCREEN_BUTTON).OnButtonClickCallback(u.toggleFullScreenClick)
	u.getButton(WINDOWED_BUTTON).OnButtonClickCallback(u.toggleFullScreenClick)
}
func (u *uiImpl) toggleFullScreenClick(id button.Id) {
	willBeFullScreen := !ebiten.IsFullscreen()
	ebiten.SetFullscreen(willBeFullScreen)
	u.SetButtonVisible(FULL_SCREEN_BUTTON, !willBeFullScreen)
	u.SetButtonVisible(WINDOWED_BUTTON, willBeFullScreen)
}

func (u *uiImpl) drawWidgets(screen *ebiten.Image) {
	for _, w := range u.widgets {
		w.Draw(screen)
	}
}

func (u *uiImpl) updateWidgets(elapsedTime int, x float64, y float64, justPressed bool, pressed bool) {
	for _, w := range u.widgets {
		w.Update(elapsedTime, x, y, justPressed, pressed, u.keys)
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

func (u *uiImpl) addRadioGroup(id radiogroup.Id, w, h float64, font draw.Font, options ...string) {
	rg := radiogroup.New(id, w, h, font, u.audioPlayer, options...)
	u.radioGroups = append(u.radioGroups, rg)
	u.widgets = append(u.widgets, rg)
}

func (u *uiImpl) getRadioGroup(id radiogroup.Id) radiogroup.RadioGroup {
	for i := range u.radioGroups {
		if u.radioGroups[i].GetId() == id {
			return u.radioGroups[i]
		}
	}
	return nil
}

func (u *uiImpl) SetRadioGroupVisible(id radiogroup.Id, visible bool) {
	if rg := u.getRadioGroup(id); rg != nil {
		rg.SetVisible(visible)
	}
}

func (u *uiImpl) SelectRadioGroup(id radiogroup.Id, index int) {
	if rg := u.getRadioGroup(id); rg != nil {
		rg.SetSelected(index)
	}
}
func (u *uiImpl) GetRadioGroupSelection(id radiogroup.Id) int {
	if rg := u.getRadioGroup(id); rg != nil {
		return rg.GetSelected()
	}
	return -1
}

func (u *uiImpl) SetRadioChangeCallback(callback func(id radiogroup.Id, index int)) {
	for i := range u.radioGroups {
		u.radioGroups[i].OnChange(callback)
	}
}

func (u *uiImpl) moveRadioGroup(id radiogroup.Id, x, y float64) {
	if rg := u.getRadioGroup(id); rg != nil {
		rg.Move(x, y)
	}
}

func (u *uiImpl) createMatchSettingsMenuUI() {
	u.addLabel(LABEL_GAME_MODE, "Game Mode:", u.fontNormal, normalLabelColor)
	u.addRadioGroup(GAME_MODE_RADIO_GROUP, BUTTON_WIDTH*3, BUTTON_HEIGHT, u.fontNormal, "AFK", "Battle", "Custom")

	u.addLabel(LABEL_GAME_MODE_DESCRIPTION, "AA", u.fontSmall, normalLabelColor)
	u.SetLabelBackgroundColor(LABEL_GAME_MODE_DESCRIPTION, colors.Black.NewWithAlpha(50), BACKGROUND_EXPAND)

	u.addLabel(LABEL_JOIN_MODE, "Joining Mode:", u.fontSmall, normalLabelColor)
	u.addRadioGroup(JOIN_MODE_RADIO_GROUP, SMALL_RADIO_OPTION_SIZE*2, SMALL_BUTTON_HEIGHT, u.fontSmall, "Any Chatter", "With !rat")

	u.addLabel(LABEL_REJOIN_MODE, "Allow Re-join:", u.fontSmall, normalLabelColor)
	u.addRadioGroup(REJOIN_MODE_RADIO_GROUP, SMALL_RADIO_OPTION_SIZE*2, SMALL_BUTTON_HEIGHT, u.fontSmall, "Yes", "No")

	u.addLabel(LABEL_ATTACK_MODE, "Attack Mode:", u.fontSmall, normalLabelColor)
	u.addRadioGroup(ATTACK_MODE_RADIO_GROUP, SMALL_RADIO_OPTION_SIZE*2, SMALL_BUTTON_HEIGHT, u.fontSmall, "Auto", "With !attack")

	u.addLabel(LABEL_HEAL_MODE, "Heal Mode:", u.fontSmall, normalLabelColor)
	u.addRadioGroup(HEAL_MODE_RADIO_GROUP, SMALL_RADIO_OPTION_SIZE*2, SMALL_BUTTON_HEIGHT, u.fontSmall, "Auto", "With !heal")

	u.addLabel(LABEL_HEALING_TO, "Healing To:", u.fontSmall, normalLabelColor)
	u.addRadioGroup(HEALING_TO_RADIO_GROUP, SMALL_RADIO_OPTION_SIZE*3, SMALL_BUTTON_HEIGHT, u.fontSmall, "Anyone", "Self", "Others")

	u.addLabel(LABEL_RAT_HEALTH, "Rat Health:", u.fontSmall, normalLabelColor)
	u.addSlider(RAT_HEALTH_SLIDER, 25, 200, SMALL_RADIO_OPTION_SIZE*3, SMALL_BUTTON_HEIGHT, u.fontSmall, normalLabelColor)

	u.addLabel(LABEL_RAT_DAMAGE, "Rat Damage:", u.fontSmall, normalLabelColor)
	u.addSlider(RAT_DAMAGE_SLIDER, 10, 80, SMALL_RADIO_OPTION_SIZE*3, SMALL_BUTTON_HEIGHT, u.fontSmall, normalLabelColor)

	u.addLabel(LABEL_RAT_HEALING, "Rat Healing:", u.fontSmall, normalLabelColor)
	u.addSlider(RAT_HEALING_SLIDER, 7, 60, SMALL_RADIO_OPTION_SIZE*3, SMALL_BUTTON_HEIGHT, u.fontSmall, normalLabelColor)

	u.addTextButton(SUBMENU_GAME_MODE_SETTINGS_BACK_BUTTON, BUTTON_WIDTH, BUTTON_HEIGHT, u.fontNormal, "Back", button.Enabled)
	u.addTextButton(SUBMENU_GAME_MODE_SETTINGS_GO_BUTTON, BUTTON_WIDTH, BUTTON_HEIGHT, u.fontNormal, "Go!", button.Enabled)
}
