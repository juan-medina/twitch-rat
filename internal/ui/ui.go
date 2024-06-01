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
	"bytes"
	"embed"
	"image/color"
	"io/fs"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/juan-medina/twitch-rat/internal/keys"
)

var (
	yellow      = color.RGBA64{0xFFFF, 0xFFFF, 0x0000, 0xFFFF}
	lightPurple = color.RGBA64{0xFFFF, 0x0000, 0xFFFF, 0xFFFF}
	darkPurple  = color.RGBA64{0x6666, 0x0000, 0x6666, 0xFFFF}
	purple      = color.RGBA64{0x8888, 0x0000, 0x8888, 0xFFFF}
	gray        = color.RGBA64{0x6666, 0x6666, 0x6666, 0xFFFF}
	darkGray    = color.RGBA64{0x3333, 0x3333, 0x3333, 0xFFFF}
	white       = color.RGBA64{0xFFFF, 0xFFFF, 0xFFFF, 0xFFFF}
	black       = color.RGBA64{0x0000, 0x0000, 0x0000, 0xFFFF}
)

type UI interface {
	Init(fileSystem embed.FS, keys keys.Keys, screenWidth, screenHeight int)
	Update(elapsedTime int)
	Draw(screen *ebiten.Image)
	SetStatusMessage(message string)
	OnButtonClick(callback func(id ButtonId))
	EnableButton(id ButtonId)
	DisableButton(id ButtonId)
	GetInputText(id InputId) string
	SetInputText(id InputId, text string)
}

type uiImpl struct {
	screenWidth   int
	screenHeight  int
	fileSystem    embed.FS
	faceSource    *text.GoTextFaceSource
	normalFace    *text.GoTextFace
	lastMessage   string
	lastMessageDO text.DrawOptions
	buttons       []button
	onButtonClick func(id ButtonId)
	inputs        []input
	keys          keys.Keys
}

// init implements UI.
func (u *uiImpl) Init(fileSystem embed.FS, keys keys.Keys, width int, height int) {
	u.fileSystem = fileSystem
	u.screenWidth = width
	u.screenHeight = height
	u.keys = keys

	fontBytes, err := fs.ReadFile(u.fileSystem, "embed/fonts/default.ttf")
	if err != nil {
		log.Fatal(err)
	}

	s, err := text.NewGoTextFaceSource(bytes.NewReader(fontBytes))
	if err != nil {
		log.Fatal(err)
	}

	u.faceSource = s

	u.normalFace = &text.GoTextFace{
		Source: u.faceSource,
		Size:   24,
	}

	u.lastMessageDO.GeoM.Reset()
	gapX := u.normalFace.Size
	gapY := u.normalFace.Size * 1.5
	u.lastMessageDO.GeoM.Translate(gapX, float64(u.screenHeight)-gapY)
	u.lastMessage = "Ready!"

	u.buttons = make([]button, 0, MAX_BUTTONS)

	bx := float64((u.screenWidth / 2) - (BUTTON_WIDTH / 2))
	by := float64((u.screenHeight / 2) - (BUTTON_HEIGHT / 2))
	pxLeft := bx - (BUTTON_WIDTH / 2) - BUTTON_GAP
	u.addButton(CONNECT_BUTTON, pxLeft, by, BUTTON_WIDTH, BUTTON_HEIGHT, "CONNECT", buttonEnabled)
	pxRight := bx + (BUTTON_WIDTH / 2) + BUTTON_GAP
	u.addButton(DISCONNECT_BUTTON, pxRight, by, BUTTON_WIDTH, BUTTON_HEIGHT, "DISCONNECT", buttonDisabled)

	u.inputs = make([]input, 0, MAX_INPUTS)
	u.addInput(INPUT_CHANNEL, pxLeft, by-BUTTON_GAP-BUTTON_HEIGHT, INPUT_WIDTH, INPUT_HEIGHT, 25, "", "Input Channel Name")
}

func (u *uiImpl) Draw(screen *ebiten.Image) {
	text.Draw(screen, u.lastMessage, u.normalFace, &u.lastMessageDO)
	for _, b := range u.buttons {
		b.draw(screen)

	}

	for _, i := range u.inputs {
		i.draw(screen)
	}
}

func (u *uiImpl) Update(elapsedTime int) {
	xe, ye := ebiten.CursorPosition()
	x, y := float64(xe), float64(ye)
	pressed := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)

	u.updateButtons(x, y, pressed, elapsedTime)
	u.updateInputs(elapsedTime)
}

func (u *uiImpl) SetStatusMessage(message string) {
	u.lastMessage = message
}

func New() UI {
	return &uiImpl{
		onButtonClick: func(id ButtonId) {},
	}
}
