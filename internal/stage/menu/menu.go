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

package menu

import (
	"embed"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/juan-medina/twitch-rat/internal/colors"
	"github.com/juan-medina/twitch-rat/internal/draw"
	"github.com/juan-medina/twitch-rat/internal/settings"
	"github.com/juan-medina/twitch-rat/internal/stage"
	"github.com/juan-medina/twitch-rat/internal/ui"
	"github.com/juan-medina/twitch-rat/internal/ui/button"
)

func (m *menu) Init() {
	channel := m.settings.GetValue("channel", "")
	m.ui.SetInputText(ui.INPUT_CHANNEL, channel)

	m.ui.OnButtonClick(m.onButtonClick)

	m.ui.SetButtonVisible(ui.PLAY_BUTTON, true)
	m.ui.SetInputVisible(ui.INPUT_CHANNEL, true)
	m.ui.SetLabelVisible(ui.LABEL_TITLE, true)

	m.ui.SetStatusMessage("Ready to Play!", colors.Yellow)
	m.firstScroll = 0
}

func (m *menu) End() {
	m.ui.SetButtonVisible(ui.PLAY_BUTTON, false)
	m.ui.SetInputVisible(ui.INPUT_CHANNEL, false)
	m.ui.SetLabelVisible(ui.LABEL_TITLE, false)

	m.settings.Save()
}

type menu struct {
	changer      stage.Changer
	ui           ui.UI
	settings     settings.Settings
	width        float32
	height       float32
	sewerMap     draw.Map
	firstScroll  float64
	secondScroll float64
}

func (m *menu) Update(elapsedTime int) {
	m.ui.Update(elapsedTime)
	w, _ := m.sewerMap.Size()

	m.firstScroll -= float64(elapsedTime) * 0.1
	if m.firstScroll < -w {
		m.firstScroll = 0
	}
	m.secondScroll = m.firstScroll + w
}

func (m *menu) Draw(screen *ebiten.Image) {
	m.sewerMap.Move(m.firstScroll, 0)
	m.sewerMap.Draw(screen)
	if m.secondScroll < float64(m.width) {
		m.sewerMap.Move(m.secondScroll, 0)
		m.sewerMap.Draw(screen)
	}
	m.ui.Draw(screen)
}

func (m *menu) OnLayoutChange(width, height float64) {
	m.width = float32(width)
	m.height = float32(height)
}

func (m *menu) onButtonClick(id button.Id) {
	switch id {
	case ui.PLAY_BUTTON:
		channel := m.ui.GetInputText(ui.INPUT_CHANNEL)
		if channel == "" {
			m.ui.SetStatusMessage("Please enter a channel name!", colors.DarkRed)
			return
		}
		m.settings.SetValue("channel", channel)
		m.settings.Save()

		m.ui.OnButtonClick(nil)
		m.changer.ChangeStage(stage.PLAYING)
	}
}

func New(changer stage.Changer, ui ui.UI, settings settings.Settings, fileSystem embed.FS) stage.Stage {
	return &menu{
		changer:  changer,
		ui:       ui,
		settings: settings,
		sewerMap: draw.NewMap(fileSystem, "embed/sprites/sewer/sewer.ldtk", 1, 4),
	}
}
