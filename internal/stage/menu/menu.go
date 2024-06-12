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
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/juan-medina/twitch-rat/internal/colors"
	"github.com/juan-medina/twitch-rat/internal/draw"
	"github.com/juan-medina/twitch-rat/internal/settings"
	"github.com/juan-medina/twitch-rat/internal/stage"
	"github.com/juan-medina/twitch-rat/internal/step"
	"github.com/juan-medina/twitch-rat/internal/ui"
	"github.com/juan-medina/twitch-rat/internal/ui/button"
)

const (
	RAT_SPEED    = 400
	SCROLL_SPEED = 250
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
	m.sewerMap.SetLevel(1)
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
	rats         draw.Sheet
	rat          draw.Sprite
	currentFrame step.Value
	ratX         float64
	ratY         float64
}

func (m *menu) Update(elapsedTime int) {
	m.ui.Update(elapsedTime)
	w, _ := m.sewerMap.Size()

	m.firstScroll -= float64(elapsedTime) * SCROLL_SPEED / 1000
	if m.firstScroll < -w {
		m.firstScroll = 0
	}
	m.secondScroll = m.firstScroll + w

	if m.currentFrame.Update(elapsedTime) {
		frame := fmt.Sprintf("rat_run_%02d", int(m.currentFrame.GetValue()))
		m.rat = m.rats.Sprite(frame)
		m.rat.SetScale(4)
	}
}

func (m *menu) Draw(screen *ebiten.Image) {
	m.sewerMap.Move(m.firstScroll, 0)
	m.sewerMap.Draw(screen)
	if m.secondScroll < float64(m.width) {
		m.sewerMap.Move(m.secondScroll, 0)
		m.sewerMap.Draw(screen)
	}

	m.rat.Draw(screen, m.ratX, m.ratY)

	m.ui.Draw(screen)
}

func (m *menu) OnLayoutChange(width, height float64) {
	m.width = float32(width)
	m.height = float32(height)

	cx := float64(m.width / 2)
	ratW, _ := m.rat.Size()
	m.ratX = cx - (ratW / 2)
	m.ratY = 640
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

func New(changer stage.Changer, ui ui.UI, settings settings.Settings, rats draw.Sheet, sewerMap draw.Map) stage.Stage {
	m := menu{
		changer:      changer,
		ui:           ui,
		settings:     settings,
		sewerMap:     sewerMap,
		rats:         rats,
		currentFrame: step.NewLoopValue(1, 8, RAT_SPEED),
	}
	m.rat = m.rats.Sprite("rat_run_01")
	m.rat.SetScale(4)
	return &m
}
