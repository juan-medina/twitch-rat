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

package game

import (
	"bytes"
	"embed"
	"image"
	"image/color"
	"runtime"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/juan-medina/twitch-rat/internal/audio"
	"github.com/juan-medina/twitch-rat/internal/colors"
	"github.com/juan-medina/twitch-rat/internal/draw"
	"github.com/juan-medina/twitch-rat/internal/keys"
	"github.com/juan-medina/twitch-rat/internal/settings"
	"github.com/juan-medina/twitch-rat/internal/stage"
	"github.com/juan-medina/twitch-rat/internal/stage/license"
	"github.com/juan-medina/twitch-rat/internal/stage/menu"
	"github.com/juan-medina/twitch-rat/internal/stage/playing"
	"github.com/juan-medina/twitch-rat/internal/stage/update"
	"github.com/juan-medina/twitch-rat/internal/step"
	"github.com/juan-medina/twitch-rat/internal/ui"
	"github.com/juan-medina/twitch-rat/internal/version"
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
	FADING_OUT
	FADING_IN
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
	nextStage      stage.Id
	currentWidth   float64
	currentHeight  float64
	valueToChange  step.Value
	audioPlayer    audio.Player
	fontVerySmall  draw.Font
	fontSmall      draw.Font
	fontNormal     draw.Font
	fontBig        draw.Font
	version        version.Version
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
	g.version.Init()
	g.settings.Init()
	g.keys.Init()
	g.fontVerySmall.Init(g.fileSystem, "embed/fonts/pixeloid/pixeloid_very_small.fnt")
	g.fontSmall.Init(g.fileSystem, "embed/fonts/pixeloid/pixeloid_small.fnt")
	g.fontNormal.Init(g.fileSystem, "embed/fonts/pixeloid/pixeloid_normal.fnt")
	g.fontBig.Init(g.fileSystem, "embed/fonts/pixeloid/pixeloid_big.fnt")
	sewerMap := draw.NewMap(g.fileSystem, "embed/sprites/sewer/sewer.ldtk", 4)
	rats := draw.NewSheet(g.fileSystem, "embed/sprites/rats/rats.json")

	g.ui.Init(g.fileSystem, g.keys, rats)
	g.ui.OnLayoutChange(g.currentWidth, g.currentHeight)

	song := g.settings.GetIntValue(settings.SONG_VOLUME, settings.DEFAULT_SONG_VOLUME)
	sound := g.settings.GetIntValue(settings.SOUND_VOLUME, settings.DEFAULT_SOUND_VOLUME)
	g.audioPlayer.ChangeSongVolume(song)
	g.audioPlayer.ChangeSoundVolume(sound)

	g.lastUpdateTime = time.Now()

	g.addStage(stage.LICENSE, license.New(g.version, g, g.ui))
	g.addStage(stage.UPDATE, update.New(g, g.ui, g.version))
	g.addStage(stage.MENU, menu.New(g, g.ui, g.settings, rats, sewerMap, g.audioPlayer))
	g.addStage(stage.PLAYING, playing.New(g, g.ui, g.settings, sewerMap, rats, g.audioPlayer, g.fontVerySmall, g.fontSmall))

	if debug := g.settings.GetBoolValue(settings.DEBUG, settings.DEFAULT_DEBUG); !debug {
		g.changeStage(stage.LICENSE)
	} else {
		g.changeStage(stage.PLAYING)
	}

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

	g.audioPlayer.Update()

	g.keys.Update()

	if g.currentStage != stage.NONE {
		if g.state == RUNNING {
			g.stages[g.currentStage].Update(elapsedTime, g.keys)
		} else {
			g.valueToChange.Update(elapsedTime)
			if g.valueToChange.IsAtEnd() {
				if g.state == FADING_OUT {
					if g.nextStage == stage.EXIT {
						if g.audioPlayer != nil {
							g.audioPlayer.Stop()
						}
						return g.exit()
					}
					g.changeStage(g.nextStage)
					g.stages[g.currentStage].Update(0, g.keys)
					g.state = FADING_IN
					g.valueToChange.Reset()
				} else if g.state == FADING_IN {
					g.state = RUNNING
				}
			}
		}
	}

	return nil
}
func (g *game) Draw(screen *ebiten.Image) {
	if g.currentStage != stage.NONE {
		g.stages[g.currentStage].Draw(screen)
	}
	if g.state == FADING_OUT {
		vector.DrawFilledRect(screen, 0, 0, float32(g.currentWidth), float32(g.currentHeight), color.RGBA{0, 0, 0, uint8(g.valueToChange.GetValue())}, true)
	} else if g.state == FADING_IN {
		vector.DrawFilledRect(screen, 0, 0, float32(g.currentWidth), float32(g.currentHeight), color.RGBA{0, 0, 0, 255 - uint8(g.valueToChange.GetValue())}, true)
	}
}

func (g *game) Run() error {
	ebiten.SetWindowSize(WIDTH, HEIGHT)
	ebiten.SetWindowTitle(TITLE)
	ebiten.SetTPS(60)
	ebiten.SetVsyncEnabled(true)
	ebiten.SetRunnableOnUnfocused(true)
	if runtime.GOOS != "js" {
		if iconData, err := g.fileSystem.ReadFile("embed/icon/rat.png"); err == nil {
			if img, _, err := image.Decode(bytes.NewReader(iconData)); err == nil {
				ebiten.SetWindowIcon([]image.Image{img})
			} else {
				panic(err)
			}
		} else {
			panic(err)
		}
		ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	} else {
		ebiten.SetFullscreen(false)
	}

	return ebiten.RunGame(g)
}

func (g *game) ChangeStage(id stage.Id) {
	g.nextStage = id
	g.state = FADING_OUT
	g.valueToChange.Reset()
}

func (g *game) changeStage(id stage.Id) {
	if g.currentStage != stage.NONE {
		g.stages[g.currentStage].End()
	}
	g.currentStage = id
	g.stages[g.currentStage].Init()
	g.stages[g.currentStage].OnLayoutChange(g.currentWidth, g.currentHeight)
}

func (g *game) DrawFinalScreen(screen ebiten.FinalScreen, offscreen *ebiten.Image, geoM ebiten.GeoM) {
	opt := ebiten.DrawImageOptions{GeoM: geoM}
	opt.Filter = ebiten.FilterNearest
	screen.Fill(colors.Teal)
	screen.DrawImage(offscreen, &opt)
}

func New(er embed.FS) *game {
	audioPlayer := audio.NewPlayer(er, 1, 1)
	fontVerySmall := draw.NewFont()
	fontSmall := draw.NewFont()
	fontNormal := draw.NewFont()
	fontBig := draw.NewFont()
	version := version.New(er)
	g := game{
		fileSystem:    er,
		state:         LOADING,
		ui:            ui.New(version, audioPlayer, fontVerySmall, fontSmall, fontNormal, fontBig),
		keys:          keys.New(),
		settings:      settings.New(APPLICATION),
		stages:        make(map[stage.Id]stage.Stage),
		currentStage:  stage.NONE,
		valueToChange: step.NewFromToPauseValue(0, 255, 200, 100),
		audioPlayer:   audioPlayer,
		fontVerySmall: fontVerySmall,
		fontSmall:     fontSmall,
		fontNormal:    fontNormal,
		fontBig:       fontBig,
		version:       version,
	}

	return &g
}
