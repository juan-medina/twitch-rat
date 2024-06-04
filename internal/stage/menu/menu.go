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
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/juan-medina/twitch-rat/internal/settings"
	"github.com/juan-medina/twitch-rat/internal/stage"
	"github.com/juan-medina/twitch-rat/internal/ui"
)

func (m *menu) Init() {
	channel := m.settings.GetValue("channel", "")
	m.ui.SetInputText(ui.INPUT_CHANNEL, channel)

	m.ui.OnButtonClick(m.onButtonClick)

	m.ui.SetButtonVisible(ui.PLAY_BUTTON, true)
	m.ui.SetInputVisible(ui.INPUT_CHANNEL, true)

	m.ui.SetStatusMessage("Ready to Play!")
}

func (m *menu) End() {
	m.ui.SetButtonVisible(ui.PLAY_BUTTON, false)
	m.ui.SetInputVisible(ui.INPUT_CHANNEL, false)

	m.settings.Save()
}

type menu struct {
	changer  stage.Changer
	ui       ui.UI
	settings settings.Settings
}

func (m *menu) Update(elapsedTime int) {
	m.ui.Update(elapsedTime)
}

func (m *menu) Draw(screen *ebiten.Image) {
	m.ui.Draw(screen)
}

func (m *menu) onButtonClick(id ui.ButtonId) {
	switch id {
	case ui.PLAY_BUTTON:
		channel := m.ui.GetInputText(ui.INPUT_CHANNEL)
		m.settings.SetValue("channel", channel)
		m.settings.Save()

		m.ui.OnButtonClick(nil)
		m.changer.ChangeStage(stage.PLAYING)
	}
}

func New(changer stage.Changer, ui ui.UI, settings settings.Settings) stage.Stage {
	return &menu{
		changer:  changer,
		ui:       ui,
		settings: settings,
	}
}
