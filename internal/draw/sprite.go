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
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type Sprite interface {
	Draw(screen *ebiten.Image, x, y float64, flippedX, flippedY bool)
	SetScale(scale float64)
	Size() (float64, float64)
	SetPivot(x, y float64)
	SetColor(color color.Color)
	Filter(filter bool)
}

type spriteImpl struct {
	image       *ebiten.Image
	drawOptions *ebiten.DrawImageOptions
	scale       float64
	pivotX      float64
	pivotY      float64
	filter      bool
}

func (s *spriteImpl) Draw(screen *ebiten.Image, x, y float64, flippedX, flippedY bool) {
	ebitenScaleX := s.scale
	ebitenScaleY := s.scale
	if flippedX {
		ebitenScaleX = -ebitenScaleX
	}
	if flippedY {
		ebitenScaleY = -ebitenScaleY
	}

	if s.filter {
		s.drawOptions.Filter = ebiten.FilterLinear
	} else {
		s.drawOptions.Filter = ebiten.FilterNearest
	}

	s.drawOptions.GeoM.Reset()
	s.drawOptions.GeoM.Scale(ebitenScaleX, ebitenScaleY)
	pivotX := s.pivotX * float64(s.image.Bounds().Dx()) * ebitenScaleX
	pivotY := s.pivotY * float64(s.image.Bounds().Dy()) * ebitenScaleY
	s.drawOptions.GeoM.Translate(x-pivotX, y-pivotY)
	screen.DrawImage(s.image, s.drawOptions)
}

func (s *spriteImpl) SetScale(scale float64) {
	s.scale = scale
}

func (s spriteImpl) Size() (float64, float64) {
	return float64(s.image.Bounds().Dx()) * s.scale, float64(s.image.Bounds().Dy()) * s.scale
}

func (s *spriteImpl) SetPivot(x, y float64) {
	s.pivotX = x
	s.pivotY = y
}
func (s *spriteImpl) SetColor(color color.Color) {
	s.drawOptions.ColorScale.Reset()
	s.drawOptions.ColorScale.ScaleWithColor(color)
	_, _, _, a := color.RGBA()
	s.drawOptions.ColorScale.ScaleAlpha(float32(a) / 65536)
}

func (s *spriteImpl) Filter(filter bool) {
	s.filter = filter
}
func NewSprite(image *ebiten.Image) Sprite {
	return &spriteImpl{
		image:       image,
		drawOptions: &ebiten.DrawImageOptions{},
		scale:       1,
		pivotX:      0.5,
		pivotY:      0.5,
	}
}
