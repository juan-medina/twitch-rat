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
	"bytes"
	"embed"
	"encoding/xml"
	"image"
	"image/color"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/juan-medina/twitch-rat/internal/colors"
)

type FontDef struct {
	XMLName xml.Name `xml:"font"`
	Text    string   `xml:",chardata"`
	Info    struct {
		Text     string `xml:",chardata"`
		Face     string `xml:"face,attr"`
		Size     int    `xml:"size,attr"`
		Bold     int    `xml:"bold,attr"`
		Italic   int    `xml:"italic,attr"`
		Charset  string `xml:"charset,attr"`
		Unicode  int    `xml:"unicode,attr"`
		StretchH int    `xml:"stretchH,attr"`
		Smooth   int    `xml:"smooth,attr"`
		Aa       int    `xml:"aa,attr"`
		Padding  string `xml:"padding,attr"`
		Spacing  string `xml:"spacing,attr"`
		Outline  int    `xml:"outline,attr"`
	} `xml:"info"`
	Common struct {
		Text       string `xml:",chardata"`
		LineHeight int    `xml:"lineHeight,attr"`
		Base       int    `xml:"base,attr"`
		ScaleW     int    `xml:"scaleW,attr"`
		ScaleH     int    `xml:"scaleH,attr"`
		Pages      int    `xml:"pages,attr"`
		Packed     int    `xml:"packed,attr"`
		AlphaChnl  int    `xml:"alphaChnl,attr"`
		RedChnl    int    `xml:"redChnl,attr"`
		GreenChnl  int    `xml:"greenChnl,attr"`
		BlueChnl   int    `xml:"blueChnl,attr"`
	} `xml:"common"`
	Pages struct {
		Text string `xml:",chardata"`
		Page []struct {
			Text string `xml:",chardata"`
			ID   int    `xml:"id,attr"`
			File string `xml:"file,attr"`
		} `xml:"page"`
	} `xml:"pages"`
	Chars struct {
		Text  string `xml:",chardata"`
		Count string `xml:"count,attr"`
		Char  []struct {
			Text     string `xml:",chardata"`
			ID       int    `xml:"id,attr"`
			X        int    `xml:"x,attr"`
			Y        int    `xml:"y,attr"`
			Width    int    `xml:"width,attr"`
			Height   int    `xml:"height,attr"`
			Xoffset  int    `xml:"xoffset,attr"`
			Yoffset  int    `xml:"yoffset,attr"`
			Xadvance int    `xml:"xadvance,attr"`
			Page     int    `xml:"page,attr"`
			Chnl     int    `xml:"chnl,attr"`
		} `xml:"char"`
	} `xml:"chars"`
}

func (f *FontDef) load(fileSystem embed.FS, fileName string) {
	data, err := fileSystem.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	err = xml.Unmarshal(data, f)
	if err != nil {
		panic(err)
	}
}

type Font interface {
	Init(fileSystem embed.FS, fileName string)
	Draw(screen *ebiten.Image, text string, x, y float64, size float64, color color.Color)
	Measure(text string, size float64) (float64, float64)
	DefaultSize() float64
}

type runeDef struct {
	id       int
	x        int
	y        int
	width    int
	height   int
	xOffset  int
	yOffset  int
	xAdvance int
	page     int
	channel  int
	sprite   Sprite
}

type fontImpl struct {
	fontDef  FontDef
	pages    []*ebiten.Image
	runes    map[rune]runeDef
	spacingX float64
	spacingY float64
}

func (f *fontImpl) Init(fileSystem embed.FS, fileName string) {
	f.fontDef.load(fileSystem, fileName)
	basePath := fileName
	if strings.Contains(basePath, "/") {
		index := strings.LastIndex(basePath, "/")
		basePath = basePath[:index]
	} else {
		basePath = "."
	}
	spacingStr := f.fontDef.Info.Spacing
	spacing := strings.Split(spacingStr, ",")
	if len(spacing) != 2 {
		panic("invalid spacing: " + spacingStr)
	}
	f.spacingX, _ = strconv.ParseFloat(spacing[0], 64)
	f.spacingY, _ = strconv.ParseFloat(spacing[1], 64)

	f.loadPages(fileSystem, basePath)
	f.processChars()
}
func (f *fontImpl) loadPages(fileSystem embed.FS, basePath string) {
	f.pages = make([]*ebiten.Image, f.fontDef.Common.Pages)
	for i, page := range f.fontDef.Pages.Page {
		textureFileName := basePath + "/" + page.File
		if textureData, err := fileSystem.ReadFile(textureFileName); err == nil {
			if img, _, err := image.Decode(bytes.NewReader(textureData)); err == nil {
				f.pages[i] = ebiten.NewImageFromImage(img)
			} else {
				panic(err)
			}
		} else {
			panic(err)
		}
	}
}

func (f *fontImpl) processChars() {
	for _, char := range f.fontDef.Chars.Char {
		if char.Page >= len(f.pages) {
			panic("page not found for char: " + strconv.Itoa(char.ID))
		}
		spriteImage := f.pages[char.Page].SubImage(image.Rect(char.X, char.Y, char.X+char.Width, char.Y+char.Height)).(*ebiten.Image)
		sprite := NewSprite(spriteImage)
		sprite.SetPivot(0, 0)
		f.runes[rune(char.ID)] = runeDef{
			id:       char.ID,
			x:        char.X,
			y:        char.Y,
			width:    char.Width,
			height:   char.Height,
			xOffset:  char.Xoffset,
			yOffset:  char.Yoffset,
			xAdvance: char.Xadvance,
			page:     char.Page,
			channel:  char.Chnl,
			sprite:   sprite,
		}
	}
}

func (f *fontImpl) Draw(screen *ebiten.Image, text string, x, y float64, size float64, color color.Color) {
	_, _, _, originalAlpha := color.RGBA()
	originalAlpha = originalAlpha >> 8
	currentPosX := x
	currentPosY := y
	scale := size / float64(f.fontDef.Common.LineHeight)
	var currentChar runeDef
	skip := 0
	for i, r := range text {
		if skip > 0 {
			skip--
			continue
		}
		if r == '\n' {
			currentPosX = x
			currentPosY += float64(f.fontDef.Common.LineHeight) * scale
			currentPosY += f.spacingY * scale
			continue
		}
		if r == colors.TEXT_TAG {
			colorStr := text[i+2 : i+10]
			color = colors.FromHtml(colorStr).NewWithAlpha(uint8(originalAlpha))
			skip = 8
			continue
		} else if char, ok := f.runes[r]; ok {
			currentChar = char
		} else {
			currentChar = f.runes[-1]
		}
		currentChar.sprite.SetScale(scale)
		currentChar.sprite.SetColor(color)
		destX := currentPosX + (float64(currentChar.xOffset) * scale)
		destY := currentPosY + (float64(currentChar.yOffset) * scale)
		currentChar.sprite.Draw(screen, destX, destY, false, false)
		currentPosX += float64(currentChar.xAdvance) * scale
		currentPosX += f.spacingX * scale

	}
}
func (f fontImpl) Measure(text string, size float64) (float64, float64) {
	//
	x := 0.0
	y := 0.0
	maxX := 0.0
	maxY := 0.0
	currentPosX := x
	currentPosY := y
	scale := size / float64(f.fontDef.Common.LineHeight)
	var currentChar runeDef
	skip := 0
	for _, r := range text {
		if skip > 0 {
			skip--
			continue
		}
		if r == '\n' {
			currentPosX = x
			currentPosY += float64(f.fontDef.Common.LineHeight) * scale
			currentPosY += f.spacingY * scale
			continue
		}
		if r == colors.TEXT_TAG {
			skip = 8
			continue
		} else if char, ok := f.runes[r]; ok {
			currentChar = char
		} else {
			currentChar = f.runes[-1]
		}
		currentPosX += float64(currentChar.xAdvance) * scale
		currentPosX += f.spacingX * scale

		if currentPosX > maxX {
			maxX = currentPosX
		}
		if currentPosY > maxY {
			maxY = currentPosY
		}
	}
	maxY += float64(f.fontDef.Common.LineHeight) * scale
	return maxX, maxY
}

func (f fontImpl) DefaultSize() float64 {
	return float64(f.fontDef.Common.LineHeight)
}

func NewFont() Font {
	return &fontImpl{
		fontDef: FontDef{},
		runes:   make(map[rune]runeDef),
	}
}
