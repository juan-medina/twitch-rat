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
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	whiteImage    = ebiten.NewImage(3, 3)
	whiteSubImage = whiteImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
)

func init() {
	b := whiteImage.Bounds()
	pix := make([]byte, 4*b.Dx()*b.Dy())
	for i := range pix {
		pix[i] = 0xff
	}
	whiteImage.WritePixels(pix)
}

type Direction int

const (
	Horizontal Direction = iota
	Vertical
)

func DrawGradientRect(dst *ebiten.Image, x, y, width, height float32, clr1 color.Color, clr2 color.Color, direction Direction, antialias bool) {
	var path vector.Path
	path.MoveTo(x, y)
	path.LineTo(x, y+height)
	path.LineTo(x+width, y+height)
	path.LineTo(x+width, y)
	vs, is := path.AppendVerticesAndIndicesForFilling(nil, nil)

	r1, g1, b1, a1 := clr1.RGBA()
	r2, g2, b2, a2 := clr2.RGBA()

	var color1Vertex []uint16
	var color2Vertex []uint16

	switch direction {
	case Horizontal:
		color1Vertex = []uint16{0, 1}
		color2Vertex = []uint16{2, 3}
	case Vertical:
		color1Vertex = []uint16{0, 3}
		color2Vertex = []uint16{1, 2}
	}

	for _, i := range color1Vertex {
		vs[i].SrcX = 1
		vs[i].SrcY = 1
		vs[i].ColorR = float32(r1) / 0xffff
		vs[i].ColorG = float32(g1) / 0xffff
		vs[i].ColorB = float32(b1) / 0xffff
		vs[i].ColorA = float32(a1) / 0xffff
	}

	for _, i := range color2Vertex {
		vs[i].SrcX = 1
		vs[i].SrcY = 1
		vs[i].ColorR = float32(r2) / 0xffff
		vs[i].ColorG = float32(g2) / 0xffff
		vs[i].ColorB = float32(b2) / 0xffff
		vs[i].ColorA = float32(a2) / 0xffff
	}

	op := &ebiten.DrawTrianglesOptions{}
	op.ColorScaleMode = ebiten.ColorScaleModePremultipliedAlpha
	op.AntiAlias = antialias
	dst.DrawTriangles(vs, is, whiteSubImage, op)
}
