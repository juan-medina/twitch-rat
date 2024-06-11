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

package draw

import (
	"embed"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/ldtkgo"
	renderer "github.com/solarlune/ldtkgo/renderer/ebitengine"
)

type Map interface {
	Move(x, y float64)
	Draw(screen *ebiten.Image)
	Size() (float64, float64)
}

type mapImpl struct {
	fileSystem           embed.FS
	fileName             string
	level                int
	world                *ldtkgo.Project
	mapRender            *renderer.Renderer
	mapImage             *ebiten.Image
	mapDrawingOptions    ebiten.DrawImageOptions
	renderDrawingOptions *renderer.DrawOptions
	x, y                 float64
	scale                float64
}

func (m mapImpl) Draw(screen *ebiten.Image) {
	screen.DrawImage(m.mapImage, &m.mapDrawingOptions)
}

func (m *mapImpl) Move(x float64, y float64) {
	m.x = x
	m.y = y
	m.mapDrawingOptions.GeoM.Reset()
	m.mapDrawingOptions.GeoM.Scale(m.scale, m.scale)
	m.mapDrawingOptions.GeoM.Translate(m.x, m.y)
}

func (m *mapImpl) Size() (float64, float64) {
	return float64(m.world.Levels[m.level].Width) * m.scale,
		float64(m.world.Levels[m.level].Height) * m.scale
}

func NewMap(fileSystem embed.FS, fileName string, level int, scale float64) Map {
	m := mapImpl{
		fileSystem: fileSystem,
		fileName:   fileName,
		level:      level,
		scale:      scale,
	}

	var err error
	if m.world, err = ldtkgo.Open(fileName, fileSystem); err != nil {
		panic(err)
	} else {
		basePath := fileName

		if strings.Contains(basePath, "/") {
			index := strings.LastIndex(basePath, "/")
			basePath = basePath[:index]
		}

		for i, _ := range m.world.Tilesets {
			tilePath := m.world.Tilesets[i].Path
			if !strings.Contains(tilePath, "/") {
				m.world.Tilesets[i].Path = basePath + "/" + tilePath
			}
		}

		if m.mapRender, err = renderer.New(m.fileSystem, m.world); err != nil {
			panic(err)
		}
	}

	levelData := m.world.Levels[level]
	m.mapImage = ebiten.NewImage(levelData.Width, levelData.Height)
	m.renderDrawingOptions = renderer.NewDefaultDrawOptions()
	m.renderDrawingOptions.BackgroundColorFill = false
	m.mapRender.Render(levelData, m.mapImage, m.renderDrawingOptions)

	return &m
}
