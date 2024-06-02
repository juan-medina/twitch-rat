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

package game

import (
	"embed"
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/juan-medina/twitch-rat/internal/chat"
	"github.com/juan-medina/twitch-rat/internal/keys"
	"github.com/juan-medina/twitch-rat/internal/settings"
	"github.com/juan-medina/twitch-rat/internal/ui"
)

const (
	WIDTH       = 1920
	HEIGHT      = 1080
	TITLE       = "Twitch Rat"
	APPLICATION = "twitch-rats"
)

type game struct {
	eventsChan     chan chat.Event
	fileSystem     embed.FS
	initialized    bool
	channel        string
	lastUpdateTime time.Time
	ui             ui.UI
	keys           keys.Keys
	chat           chat.Chat
	settings       settings.Settings
}

func (g game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return WIDTH, HEIGHT
}

func (g *game) init() {
	g.settings.Init()
	g.chat.OnEvent(g.onChatEvent)
	g.keys.Init()
	g.ui.Init(g.fileSystem, g.keys, WIDTH, HEIGHT)
	g.ui.OnButtonClick(g.OnButtonClick)

	g.channel = g.settings.GetValue("channel", "")
	g.ui.SetInputText(ui.INPUT_CHANNEL, g.channel)
	g.initialized = true
	g.lastUpdateTime = time.Now()
}
func (g *game) OnButtonClick(id ui.ButtonId) {
	switch id {
	case ui.DISCONNECT_BUTTON:
		g.ui.SetStatusMessage("Disconnecting...")
		g.ui.DisableButton(ui.DISCONNECT_BUTTON)
		g.chat.Disconnect()
	case ui.CONNECT_BUTTON:
		channel := g.ui.GetInputText(ui.INPUT_CHANNEL)
		if channel == "" {
			g.ui.SetStatusMessage("Please enter a channel name")
			return
		}
		g.channel = channel
		g.settings.SetValue("channel", g.channel)
		g.settings.Save()
		g.ui.SetStatusMessage("Connecting...")
		g.ui.DisableButton(ui.CONNECT_BUTTON)
		g.chat.Connect(g.channel)
	}
}

func (g *game) Update() error {
	if !g.initialized {
		g.init()
	}

	elapsedTime := time.Since(g.lastUpdateTime)
	g.lastUpdateTime = time.Now()
	elapsedMillis := int(elapsedTime.Milliseconds())

	g.ui.Update(elapsedMillis)
	g.keys.Update(elapsedMillis)

	select {
	case event := <-g.eventsChan:
		switch event.Type_ {
		case chat.Connect:
			g.ui.SetStatusMessage("Connected to " + g.channel)
			g.ui.EnableButton(ui.DISCONNECT_BUTTON)
		case chat.Disconnect:
			g.ui.SetStatusMessage("Ready!")
			g.ui.EnableButton(ui.CONNECT_BUTTON)
		case chat.Message:
			g.ui.SetStatusMessage(fmt.Sprintf("%s: %s", event.Sender, event.Message))
		}
	default:
		/// no new event
	}

	return nil
}
func (g *game) Draw(screen *ebiten.Image) {
	g.ui.Draw(screen)
}

func (g *game) Run() error {
	ebiten.SetWindowSize(WIDTH, HEIGHT)
	ebiten.SetWindowTitle(TITLE)
	ebiten.SetTPS(60)

	return ebiten.RunGame(g)
}

func (g *game) onChatEvent(e chat.Event) {
	g.eventsChan <- e
}

func New(er embed.FS) *game {
	g := game{
		eventsChan:  make(chan chat.Event, 10),
		fileSystem:  er,
		initialized: false,
		ui:          ui.New(),
		chat:        chat.New(),
		keys:        keys.New(),
		settings:    settings.New(APPLICATION),
		channel:     "",
	}

	return &g
}
