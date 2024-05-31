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

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/juan-medina/twitch-rat/internal/chat"
	"github.com/juan-medina/twitch-rat/internal/ui"
)

const (
	WIDTH  = 1920
	HEIGHT = 1080
	TITLE  = "Twitch Rat"
)

type game struct {
	eventsChan  chan chat.Event
	fileSystem  embed.FS
	initialized bool
	channel     string
	ui          ui.UI
	chat        chat.Chat
}

func (g game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return WIDTH, HEIGHT
}

func (g *game) init() {
	g.ui.Init(g.fileSystem, WIDTH, HEIGHT)
	g.ui.OnButtonClick(g.OnButtonClick)
	g.chat.OnEvent(g.onChatEvent)
	g.initialized = true
}
func (g *game) OnButtonClick(id ui.ButtonId) {
	if id == ui.CONNECT_BUTTON {
		g.ui.DisableButton(ui.CONNECT_BUTTON)
		g.chat.Connect(g.channel)
	}
}

func (g *game) Update() error {
	if !g.initialized {
		g.init()
	}

	g.ui.Update()

	select {
	case event := <-g.eventsChan:
		switch event.Type_ {
		case chat.Connect:
			g.ui.SetStatusMessage("Connected to " + g.channel)
		case chat.Message:
			g.ui.SetStatusMessage(fmt.Sprintf("%s: %s", event.Sender, event.Message))
		}
	default:
		return nil
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
		channel:     "harukakaribu",
	}

	return &g
}
