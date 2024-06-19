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

package label

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	ebitenText "github.com/hajimehoshi/ebiten/v2/text/v2"
)

type labelTTFImpl struct {
	id              Id
	text            string
	textDrawOptions ebitenText.DrawOptions
	face            ebitenText.Face
	visible         bool
	lineHeight      float64
}

func (l labelTTFImpl) GetId() Id {
	return l.id
}

func (l *labelTTFImpl) Draw(screen *ebiten.Image) {
	if l.visible {
		ebitenText.Draw(screen, l.text, l.face, &l.textDrawOptions)
	}
}

func (l *labelTTFImpl) Move(x float64, y float64) {
	l.textDrawOptions.GeoM.Reset()
	l.textDrawOptions.GeoM.Translate(x, y)
}

func (l *labelTTFImpl) SetText(text string) {
	l.text = text
}

func (l labelTTFImpl) GetText() string {
	return l.text
}

func (l labelTTFImpl) Measure() (float64, float64) {
	return ebitenText.Measure(l.text, l.face, l.lineHeight)
}

func (l *labelTTFImpl) SetColor(color color.Color) {
	l.textDrawOptions.ColorScale.Reset()
	l.textDrawOptions.ColorScale.ScaleWithColor(color)
	_, _, _, a := color.RGBA()
	l.textDrawOptions.ColorScale.ScaleAlpha(float32(a) / 65536)
}
func (l labelTTFImpl) GetColor() color.Color {
	r := uint8(l.textDrawOptions.ColorScale.R() * 255)
	g := uint8(l.textDrawOptions.ColorScale.G() * 255)
	b := uint8(l.textDrawOptions.ColorScale.B() * 255)
	a := uint8(l.textDrawOptions.ColorScale.A() * 255)
	return color.RGBA{r, g, b, a}
}

func (l *labelTTFImpl) SetAlpha(alpha float32) {
	l.textDrawOptions.ColorScale.ScaleAlpha(alpha)
}

func (l *labelTTFImpl) SetVisible(visible bool) {
	l.visible = visible
}

func NewLabelTTF(id Id, text string, face ebitenText.Face, lineHeight float64, color color.Color) Label {
	textDrawOptions := ebitenText.DrawOptions{}
	textDrawOptions.ColorScale.ScaleWithColor(color)
	textDrawOptions.Filter = ebiten.FilterNearest
	textDrawOptions.LayoutOptions.LineSpacing = lineHeight
	return &labelTTFImpl{
		id:              id,
		text:            text,
		textDrawOptions: textDrawOptions,
		face:            face,
		lineHeight:      lineHeight,
	}
}
