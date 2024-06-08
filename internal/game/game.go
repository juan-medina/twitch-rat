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
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/juan-medina/twitch-rat/internal/keys"
	"github.com/juan-medina/twitch-rat/internal/settings"
	"github.com/juan-medina/twitch-rat/internal/stage"
	"github.com/juan-medina/twitch-rat/internal/stage/menu"
	"github.com/juan-medina/twitch-rat/internal/stage/playing"
	"github.com/juan-medina/twitch-rat/internal/ui"
)

const (
	WIDTH  = 1920.0
	HEIGHT = 1080.0
	TITLE  = "Twitch Rat"

	APPLICATION = "twitch-rats"
)

type GameState int

const (
	LOADING GameState = iota
	RUNNING
)

type game struct {
	fileSystem     embed.FS
	lastUpdateTime time.Time
	ui             ui.UI
	keys           keys.Keys
	settings       settings.Settings
	state          GameState
	stages         map[stage.Id]stage.Stage
	currentStage   stage.Id
	currentWidth   float64
	currentHeight  float64
}

func (g *game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return WIDTH, HEIGHT
}
func (g *game) LayoutF(outsideWidth, outsideHeight float64) (screenWidth, screenHeight float64) {
	currentHeight := HEIGHT
	currentWidth := (HEIGHT * outsideWidth) / outsideHeight

	if currentWidth != g.currentWidth || currentHeight != g.currentHeight {
		g.onLayoutChange(currentWidth, currentHeight)
	}
	return currentWidth, currentHeight
}

func (g *game) onLayoutChange(width, height float64) {
	g.currentWidth = width
	g.currentHeight = height
	if g.state == RUNNING {
		g.ui.OnLayoutChange(width, height)
		if g.currentStage != stage.NONE {
			g.stages[g.currentStage].OnLayoutChange(width, height)
		}
	}
}

func (g *game) init() {
	g.settings.Init()
	g.keys.Init()
	g.ui.Init(g.fileSystem, g.keys)
	g.ui.OnLayoutChange(g.currentWidth, g.currentHeight)

	g.lastUpdateTime = time.Now()

	g.addStage(stage.MENU, menu.New(g, g.ui, g.settings))
	g.addStage(stage.PLAYING, playing.New(g, g.ui, g.settings))
	g.ChangeStage(stage.MENU)

	g.state = RUNNING
}

func (g *game) addStage(id stage.Id, st stage.Stage) {
	g.stages[id] = st
}

func (g *game) getElapsedTime() int {
	elapsedTime := int(time.Since(g.lastUpdateTime).Milliseconds())
	g.lastUpdateTime = time.Now()
	return elapsedTime
}

func (g *game) Update() error {
	elapsedTime := g.getElapsedTime()

	if g.state == LOADING {
		g.init()
		return nil
	}

	g.keys.Update(elapsedTime)

	if g.currentStage != stage.NONE {
		g.stages[g.currentStage].Update(elapsedTime)
	}

	return nil
}
func (g *game) Draw(screen *ebiten.Image) {
	if g.currentStage != stage.NONE {
		g.stages[g.currentStage].Draw(screen)
	}
}

func (g *game) Run() error {
	ebiten.SetWindowSize(WIDTH, HEIGHT)
	ebiten.SetWindowTitle(TITLE)
	ebiten.SetTPS(60)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetFullscreen(false)

	return ebiten.RunGame(g)
}

func (g *game) ChangeStage(id stage.Id) {
	if g.currentStage != stage.NONE {
		g.stages[g.currentStage].End()
	}
	g.currentStage = id
	g.stages[g.currentStage].Init()
	g.stages[g.currentStage].OnLayoutChange(g.currentWidth, g.currentHeight)
}

func New(er embed.FS) *game {
	g := game{
		fileSystem:   er,
		state:        LOADING,
		ui:           ui.New(),
		keys:         keys.New(),
		settings:     settings.New(APPLICATION),
		stages:       make(map[stage.Id]stage.Stage),
		currentStage: stage.NONE,
	}

	return &g
}
