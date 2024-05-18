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
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const (
	WIDTH  = 1920
	HEIGHT = 1080
	TITLE  = "Twitch Rat"
)

type twitch_message struct {
	message string
	sender  string
}

type game struct {
	messageChan      chan twitch_message
	faceSource       *text.GoTextFaceSource
	normalFace       *text.GoTextFace
	normalFaceHeight float64
	fileSystem       embed.FS
	initialized      bool
	lastMessage      string
	lastMessageDO    text.DrawOptions
}

func (g game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return WIDTH, HEIGHT
}

func (g *game) init() {
	// Load font
	fontBytes, err := fs.ReadFile(g.fileSystem, "embed/fonts/default.ttf")
	if err != nil {
		panic(err)
	}

	s, err := text.NewGoTextFaceSource(bytes.NewReader(fontBytes))
	if err != nil {
		log.Fatal(err)
	}

	g.faceSource = s

	g.normalFace = &text.GoTextFace{
		Source: g.faceSource,
		Size:   24,
	}
	g.normalFaceHeight = g.normalFace.Size * 1.5

	g.initialized = true

	g.lastMessage = "Loading..."
}

func (g *game) Update() error {

	if !g.initialized {
		g.init()
	}

	select {
	case msg := <-g.messageChan:
		g.lastMessage = fmt.Sprintf("%s: %s", msg.sender, msg.message)
	default:
		// No new messages
	}

	w, h := text.Measure(g.lastMessage, g.normalFace, g.normalFaceHeight)
	g.lastMessageDO.GeoM.Reset()
	g.lastMessageDO.GeoM.Translate(WIDTH/2, HEIGHT/2)
	g.lastMessageDO.GeoM.Translate(-w/2, -h/2)

	return nil
}
func (g *game) Draw(screen *ebiten.Image) {
	text.Draw(screen, g.lastMessage, g.normalFace, &g.lastMessageDO)
}

func (g game) pre_chat() {
	ebiten.SetWindowSize(WIDTH, HEIGHT)
	ebiten.SetWindowTitle(TITLE)
	ebiten.SetTPS(60)
}

func (g *game) post_chat() error {
	return ebiten.RunGame(g)
}

func New(er embed.FS) *game {
	g := game{
		messageChan: make(chan twitch_message, 10),
		fileSystem:  er,
		initialized: false,
	}

	return &g
}
