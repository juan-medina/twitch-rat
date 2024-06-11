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

package playing

import (
	"embed"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/juan-medina/twitch-rat/internal/chat"
	"github.com/juan-medina/twitch-rat/internal/colors"
	"github.com/juan-medina/twitch-rat/internal/settings"
	"github.com/juan-medina/twitch-rat/internal/stage"
	"github.com/juan-medina/twitch-rat/internal/ui"
	"github.com/juan-medina/twitch-rat/internal/ui/button"
	"github.com/solarlune/ldtkgo"
	renderer "github.com/solarlune/ldtkgo/renderer/ebitengine"
)

func (p *playing) Init() {
	p.ui.OnButtonClick(p.onButtonClick)

	p.ui.SetButtonVisible(ui.BACK_BUTTON, true)

	p.ui.SetStatusMessage("Connecting..", colors.White)

	p.channel = p.settings.GetValue("channel", "")
	p.chat.OnEvent(p.onChatEvent)
	p.chat.Connect(p.channel)
}

func (p *playing) End() {
	p.ui.SetButtonVisible(ui.BACK_BUTTON, false)
	p.chat.Disconnect()
}

type playing struct {
	changer       stage.Changer
	ui            ui.UI
	settings      settings.Settings
	eventsChan    chan chat.Event
	chat          chat.Chat
	channel       string
	fileSystem    embed.FS
	currentWidth  float64
	currentHeight float64
	sewerMap      *ldtkgo.Project
	mapRender     *renderer.Renderer
	mapIMage      *ebiten.Image
	mapDO         ebiten.DrawImageOptions
	renderOptions *renderer.DrawOptions
}

func (p *playing) Update(elapsedTime int) {
	p.ui.Update(elapsedTime)
	select {
	case event := <-p.eventsChan:
		switch event.Type_ {
		case chat.Connect:
			p.ui.SetStatusMessage("Connected to "+p.channel, colors.White)
		case chat.Disconnect:
		case chat.Message:
			p.ui.SetStatusMessage(fmt.Sprintf("%s: %s", event.Sender, event.Message), colors.White)
		}
	default:
		/// no new event
	}
}

func (p *playing) Draw(screen *ebiten.Image) {
	screen.DrawImage(p.mapIMage, &p.mapDO)
	p.ui.Draw(screen)
}

func (p *playing) OnLayoutChange(width, height float64) {
	p.currentWidth = width
	p.currentHeight = height

	level := p.sewerMap.Levels[0]
	p.mapDO.GeoM.Reset()
	p.mapDO.GeoM.Scale(4, 4)
	dx := (p.currentWidth - float64(level.Width*4)) / 2
	p.mapDO.GeoM.Translate(dx, 0)

}

func (p *playing) onButtonClick(id button.Id) {
	switch id {
	case ui.BACK_BUTTON:
		p.ui.OnButtonClick(nil)
		p.changer.ChangeStage(stage.MENU)
	}
}
func (p *playing) onChatEvent(e chat.Event) {
	p.eventsChan <- e
}

func New(changer stage.Changer, ui ui.UI, settings settings.Settings, fileSystem embed.FS) stage.Stage {
	p := playing{
		changer:    changer,
		settings:   settings,
		ui:         ui,
		eventsChan: make(chan chat.Event, 10),
		chat:       chat.New(),
		fileSystem: fileSystem,
	}

	var err error
	if p.sewerMap, err = ldtkgo.Open("embed/sprites/sewer/sewer.ldtk", p.fileSystem); err != nil {
		panic(err)
	} else {
		p.sewerMap.Tilesets[0].Path = "embed/sprites/sewer/" + p.sewerMap.Tilesets[0].Path
		path := p.sewerMap.Tilesets[0].Path
		fmt.Println(path)
		if p.mapRender, err = renderer.New(p.fileSystem, p.sewerMap); err != nil {
			panic(err)
		}
	}
	level := p.sewerMap.Levels[0]
	p.mapIMage = ebiten.NewImage(level.Width, level.Height)
	p.renderOptions = renderer.NewDefaultDrawOptions()
	p.renderOptions.BackgroundColorFill = false
	p.mapRender.Render(level, p.mapIMage, p.renderOptions)

	return &p
}
