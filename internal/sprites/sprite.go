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

package sprites

import "github.com/hajimehoshi/ebiten/v2"

type Sprite interface {
	Draw(screen *ebiten.Image, x, y float64)
	SetScale(scale float64)
	Size() (float64, float64)
}

type spriteImpl struct {
	image       *ebiten.Image
	drawOptions *ebiten.DrawImageOptions
	scale       float64
}

func (s *spriteImpl) Draw(screen *ebiten.Image, x, y float64) {
	s.drawOptions.GeoM.Reset()
	s.drawOptions.GeoM.Scale(s.scale, s.scale)
	s.drawOptions.GeoM.Translate(x, y)
	screen.DrawImage(s.image, s.drawOptions)
}
func (s *spriteImpl) SetScale(scale float64) {
	s.scale = scale
}

func (s spriteImpl) Size() (float64, float64) {
	return float64(s.image.Bounds().Dx()) * s.scale, float64(s.image.Bounds().Dy()) * s.scale
}

func NewSprite(image *ebiten.Image) Sprite {
	return &spriteImpl{
		image:       image,
		drawOptions: &ebiten.DrawImageOptions{},
		scale:       1,
	}
}
