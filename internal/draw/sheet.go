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

package draw

import (
	"bytes"
	"embed"
	"encoding/json"
	"image"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
)

type Sheet interface {
	Load()
	Sprite(name string) Sprite
}

type SheetData struct {
	Frames []struct {
		Filename string `json:"filename"`
		Frame    struct {
			X int `json:"x"`
			Y int `json:"y"`
			W int `json:"w"`
			H int `json:"h"`
		} `json:"frame"`
		Rotated          bool `json:"rotated"`
		Trimmed          bool `json:"trimmed"`
		SpriteSourceSize struct {
			X int `json:"x"`
			Y int `json:"y"`
			W int `json:"w"`
			H int `json:"h"`
		} `json:"spriteSourceSize"`
		SourceSize struct {
			W int `json:"w"`
			H int `json:"h"`
		} `json:"sourceSize"`
		Pivot struct {
			X float64 `json:"x"`
			Y float64 `json:"y"`
		} `json:"pivot"`
	} `json:"frames"`
	Meta struct {
		App     string `json:"app"`
		Version string `json:"version"`
		Image   string `json:"image"`
		Format  string `json:"format"`
		Size    struct {
			W int `json:"w"`
			H int `json:"h"`
		} `json:"size"`
		Scale int `json:"scale"`
	} `json:"meta"`
}

type sheetImpl struct {
	data       SheetData
	fileName   string
	fileSystem embed.FS
	sprites    map[string]Sprite
	texture    *ebiten.Image
}

func (s *sheetImpl) Load() {
	if sheetBytes, err := s.fileSystem.ReadFile(s.fileName); err == nil {
		if err := json.Unmarshal([]byte(sheetBytes), &s.data); err == nil {
			basePath := s.fileName
			if strings.Contains(basePath, "/") {
				index := strings.LastIndex(basePath, "/")
				basePath = basePath[:index]
			} else {
				basePath = "."
			}
			textureFileName := basePath + "/" + s.data.Meta.Image

			if textureData, err := s.fileSystem.ReadFile(textureFileName); err == nil {
				if img, _, err := image.Decode(bytes.NewReader(textureData)); err == nil {
					s.texture = ebiten.NewImageFromImage(img)
				} else {
					panic(err)
				}
			} else {
				panic(err)
			}
			for _, frame := range s.data.Frames {
				spriteImage := s.texture.SubImage(image.Rect(frame.Frame.X, frame.Frame.Y, frame.Frame.X+frame.Frame.W, frame.Frame.Y+frame.Frame.H)).(*ebiten.Image)
				newSprite := NewSprite(spriteImage)
				newSprite.SetPivot(frame.Pivot.X, frame.Pivot.Y)
				s.sprites[frame.Filename] = newSprite
			}
		} else {
			panic(err)
		}
	} else {
		panic(err)
	}
}

func (s sheetImpl) Sprite(name string) Sprite {
	if s, ok := s.sprites[name]; !ok {
		panic("sprite not found: " + name)
	} else {
		return s
	}
}

func NewSheet(fileSystem embed.FS, fileName string) Sheet {
	s := sheetImpl{
		fileName:   fileName,
		fileSystem: fileSystem,
		data:       SheetData{},
		sprites:    make(map[string]Sprite),
	}
	s.Load()
	return &s
}
