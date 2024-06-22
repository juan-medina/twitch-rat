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

package menu

import (
	"fmt"
	"regexp"
	"runtime"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/juan-medina/twitch-rat/internal/audio"
	"github.com/juan-medina/twitch-rat/internal/colors"
	"github.com/juan-medina/twitch-rat/internal/draw"
	"github.com/juan-medina/twitch-rat/internal/keys"
	"github.com/juan-medina/twitch-rat/internal/settings"
	"github.com/juan-medina/twitch-rat/internal/stage"
	"github.com/juan-medina/twitch-rat/internal/step"
	"github.com/juan-medina/twitch-rat/internal/ui"
	"github.com/juan-medina/twitch-rat/internal/ui/button"
	"github.com/juan-medina/twitch-rat/internal/ui/slider"
)

const (
	RAT_SPEED         = 400
	SCROLL_SPEED      = 250
	MENU_MUSIC        = "embed/music/menu.ogg"
	RAT_FRAME         = "rat_normal_run_%02d"
	INITIAL_FRAME     = 1
	END_FRAME         = 8
	RAT_SCALE         = 4
	RAT_Y_POS         = 768
	CHANNEL_REGEXP    = `^[a-zA-Z][a-zA-Z0-9\-_]*$`
	CHANNEL_NOT_EMPTY = "Please enter a channel name!"
	CHANNEL_NOT_VALID = "Channel name is invalid!"
)

type subMenu int

const (
	MAIN_MENU subMenu = iota
	OPTIONS_MENU
	ABOUT_MENU
	GAME_MODE_SETTINGS
)

func (m *menu) Init() {
	channel := m.settings.GetValue("channel", "")
	m.ui.SetInputText(ui.INPUT_CHANNEL, channel)

	m.ui.SetButtonClickCallback(m.onButtonClick)
	m.ui.SetSliderChangeCallback(m.onSliderChange)

	m.ui.SetLabelVisible(ui.LABEL_TITLE, true)

	m.ui.SetStatusMessage("Ready to Play!", colors.LightYellow)
	m.firstScroll = 0
	m.sewerMap.SetLevel(1)
	m.audioPlayer.PlaySong(MENU_MUSIC)
	m.changeSubMenu(MAIN_MENU)
}

func (m *menu) End() {
	m.ui.SetButtonVisible(ui.PLAY_BUTTON, false)
	m.ui.SetInputVisible(ui.INPUT_CHANNEL, false)
	m.ui.SetLabelVisible(ui.LABEL_TITLE, false)
	m.ui.SetButtonVisible(ui.OPTIONS_BUTTON, false)
	m.ui.SetButtonVisible(ui.ABOUT_BUTTON, false)
	m.ui.SetLabelVisible(ui.LABEL_DOWNLOAD, false)
	m.ui.SetRadioGroupVisible(ui.GAME_MODE_RADIO_GROUP, false)
	m.ui.SetButtonVisible(ui.SUBMENU_GAME_MODE_SETTINGS_BACK_BUTTON, false)
	m.ui.SetButtonVisible(ui.SUBMENU_GAME_MODE_SETTINGS_GO_BUTTON, false)

	if runtime.GOOS != "js" {
		m.ui.SetButtonVisible(ui.BACK_BUTTON, false)
	}

	m.settings.Save()
	m.audioPlayer.Stop()
}

type menu struct {
	changer        stage.Changer
	ui             ui.UI
	settings       settings.Settings
	width          float32
	height         float32
	sewerMap       draw.Map
	firstScroll    float64
	secondScroll   float64
	rats           draw.Sheet
	rat            draw.Sprite
	currentFrame   step.Value
	ratX           float64
	ratY           float64
	audioPlayer    audio.Player
	channelRegex   *regexp.Regexp
	currentSubMenu subMenu
}

func (m *menu) Update(elapsedTime int, keys keys.Keys) {
	m.ui.Update(elapsedTime)
	w, _ := m.sewerMap.Size()

	m.firstScroll -= float64(elapsedTime) * SCROLL_SPEED / 1000
	if m.firstScroll < -w {
		m.firstScroll = 0
	}
	m.secondScroll = m.firstScroll + w

	if m.currentFrame.Update(elapsedTime) {
		m.updateRatFrame()
	}

	if m.currentSubMenu == MAIN_MENU {
		if !m.ui.IsInputEditing(ui.INPUT_CHANNEL) {
			if keys.IsDownNoRepeat(ebiten.KeyEnter) || keys.IsDownNoRepeat(ebiten.KeyNumpadEnter) {
				m.ui.ClickButton(ui.PLAY_BUTTON)
			}
			if runtime.GOOS != "js" {
				if keys.IsDownNoRepeat(ebiten.KeyEscape) {
					m.ui.ClickButton(ui.BACK_BUTTON)
				}
			}
		}
	} else {
		if m.currentSubMenu == GAME_MODE_SETTINGS {
			if keys.IsDownNoRepeat(ebiten.KeyEnter) || keys.IsDownNoRepeat(ebiten.KeyNumpadEnter) {
				m.ui.ClickButton(ui.SUBMENU_GAME_MODE_SETTINGS_GO_BUTTON)
			}
		}
		if keys.IsDownNoRepeat(ebiten.KeyEscape) {
			switch m.currentSubMenu {
			case OPTIONS_MENU:
				m.ui.ClickButton(ui.SUBMENU_OPTION_BACK_BUTTON)
			case ABOUT_MENU:
				m.ui.ClickButton(ui.SUBMENU_ABOUT_BACK_BUTTON)
			case GAME_MODE_SETTINGS:
				m.ui.ClickButton(ui.SUBMENU_GAME_MODE_SETTINGS_BACK_BUTTON)
			}
		}
	}
}

func (m *menu) updateRatFrame() {
	frame := fmt.Sprintf(RAT_FRAME, int(m.currentFrame.GetValue()))
	m.rat = m.rats.Sprite(frame)
	m.rat.SetScale(RAT_SCALE)
}

func (m *menu) Draw(screen *ebiten.Image) {
	m.sewerMap.Move(m.firstScroll, 0)
	m.sewerMap.Draw(screen)
	if m.secondScroll < float64(m.width) {
		m.sewerMap.Move(m.secondScroll, 0)
		m.sewerMap.Draw(screen)
	}

	m.rat.Draw(screen, m.ratX, m.ratY, false, false)

	m.ui.Draw(screen)
}

func (m *menu) OnLayoutChange(width, height float64) {
	m.width = float32(width)
	m.height = float32(height)

	cx := float64(m.width / 2)
	m.ratX = cx
	m.ratY = RAT_Y_POS
}

func (m *menu) onButtonClick(id button.Id) {
	switch id {
	case ui.PLAY_BUTTON:
		channel := m.ui.GetInputText(ui.INPUT_CHANNEL)
		if channel == "" {
			m.ui.SetStatusMessage(CHANNEL_NOT_EMPTY, colors.Red)
			return
		} else {
			if !m.channelRegex.MatchString(channel) {
				m.ui.SetStatusMessage(CHANNEL_NOT_VALID, colors.Red)
				return
			}
		}
		m.settings.SetValue("channel", channel)
		m.settings.Save()
		m.changeSubMenu(GAME_MODE_SETTINGS)
	case ui.BACK_BUTTON:
		m.changer.ChangeStage(stage.EXIT)
	case ui.OPTIONS_BUTTON:
		m.changeSubMenu(OPTIONS_MENU)
	case ui.ABOUT_BUTTON:
		m.changeSubMenu(ABOUT_MENU)
	case ui.SUBMENU_ABOUT_BACK_BUTTON:
		m.changeSubMenu(MAIN_MENU)
	case ui.SUBMENU_OPTION_BACK_BUTTON:
		m.changeSubMenu(MAIN_MENU)
	case ui.SUBMENU_GAME_MODE_SETTINGS_BACK_BUTTON:
		m.changeSubMenu(MAIN_MENU)
	case ui.SUBMENU_GAME_MODE_SETTINGS_GO_BUTTON:
		m.ui.SetButtonClickCallback(nil)
		m.ui.SetSliderChangeCallback(nil)
		m.changer.ChangeStage(stage.PLAYING)
	}
}

func (m *menu) changeSubMenu(subMenu subMenu) {
	m.currentSubMenu = subMenu

	m.ui.SetInputVisible(ui.INPUT_CHANNEL, false)
	m.ui.SetButtonVisible(ui.PLAY_BUTTON, false)
	m.ui.SetButtonVisible(ui.OPTIONS_BUTTON, false)
	m.ui.SetButtonVisible(ui.ABOUT_BUTTON, false)
	m.ui.SetLabelVisible(ui.LABEL_ABOUT_MESSAGE, false)
	m.ui.SetButtonVisible(ui.SUBMENU_ABOUT_BACK_BUTTON, false)
	m.ui.SetButtonVisible(ui.SUBMENU_OPTION_BACK_BUTTON, false)
	m.ui.SetSliderVisible(ui.MUSIC_VOLUME_SLIDER, false)
	m.ui.SetLabelVisible(ui.LABEL_OPTIONS_MUSIC_VOLUME, false)
	m.ui.SetSliderVisible(ui.AUDIO_VOLUME_SLIDER, false)
	m.ui.SetLabelVisible(ui.LABEL_OPTIONS_AUDIO_VOLUME, false)
	m.ui.SetLabelVisible(ui.LABEL_DOWNLOAD, false)
	m.ui.SetRadioGroupVisible(ui.GAME_MODE_RADIO_GROUP, false)
	m.ui.SetButtonVisible(ui.SUBMENU_GAME_MODE_SETTINGS_BACK_BUTTON, false)
	m.ui.SetButtonVisible(ui.SUBMENU_GAME_MODE_SETTINGS_GO_BUTTON, false)

	switch m.currentSubMenu {
	case MAIN_MENU:
		m.ui.SetInputVisible(ui.INPUT_CHANNEL, true)
		m.ui.SetButtonVisible(ui.PLAY_BUTTON, true)
		m.ui.SetButtonVisible(ui.OPTIONS_BUTTON, true)
		m.ui.SetButtonVisible(ui.ABOUT_BUTTON, true)
		if runtime.GOOS != "js" {
			m.ui.SetButtonVisible(ui.BACK_BUTTON, true)
		} else {
			m.ui.SetLabelVisible(ui.LABEL_DOWNLOAD, true)
		}
	case ABOUT_MENU:
		m.ui.SetLabelVisible(ui.LABEL_ABOUT_MESSAGE, true)
		m.ui.SetButtonVisible(ui.SUBMENU_ABOUT_BACK_BUTTON, true)
	case OPTIONS_MENU:
		m.ui.SetButtonVisible(ui.SUBMENU_OPTION_BACK_BUTTON, true)
		m.ui.SetSliderVisible(ui.MUSIC_VOLUME_SLIDER, true)
		m.ui.SetLabelVisible(ui.LABEL_OPTIONS_MUSIC_VOLUME, true)
		m.ui.SetSliderVisible(ui.AUDIO_VOLUME_SLIDER, true)
		m.ui.SetLabelVisible(ui.LABEL_OPTIONS_AUDIO_VOLUME, true)
		song := m.settings.GetFloatValue("song_volume", 0.2)
		m.ui.SetSliderValue(ui.MUSIC_VOLUME_SLIDER, song)
		sound := m.settings.GetFloatValue("sound_volume", 0.5)
		m.ui.SetSliderValue(ui.AUDIO_VOLUME_SLIDER, sound)
	case GAME_MODE_SETTINGS:
		m.ui.SetRadioGroupVisible(ui.GAME_MODE_RADIO_GROUP, true)
		m.ui.SetButtonVisible(ui.SUBMENU_GAME_MODE_SETTINGS_BACK_BUTTON, true)
		m.ui.SetButtonVisible(ui.SUBMENU_GAME_MODE_SETTINGS_GO_BUTTON, true)
		m.ui.SelectRadioGroup(ui.GAME_MODE_RADIO_GROUP, 0)
	}
}

func (m *menu) onSliderChange(id slider.Id, value float64) {
	switch id {
	case ui.MUSIC_VOLUME_SLIDER:
		m.settings.SetFloatValue("song_volume", value)
		m.settings.Save()
		m.audioPlayer.ChangeSongVolume(value)
	case ui.AUDIO_VOLUME_SLIDER:
		m.settings.SetFloatValue("sound_volume", value)
		m.settings.Save()
		m.audioPlayer.ChangeSoundVolume(value)
	}
}

func New(changer stage.Changer, ui ui.UI, settings settings.Settings, rats draw.Sheet, sewerMap draw.Map, audioPlayer audio.Player) stage.Stage {
	audioPlayer.LoadSong(MENU_MUSIC)
	m := menu{
		changer:      changer,
		ui:           ui,
		settings:     settings,
		sewerMap:     sewerMap,
		rats:         rats,
		currentFrame: step.NewLoopValue(INITIAL_FRAME, END_FRAME, RAT_SPEED),
		audioPlayer:  audioPlayer,
		channelRegex: regexp.MustCompile(CHANNEL_REGEXP),
	}
	m.updateRatFrame()
	return &m
}
