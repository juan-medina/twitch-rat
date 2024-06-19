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
