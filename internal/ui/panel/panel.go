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

package panel

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Id int
type Panel interface {
	GetId() Id
	Draw(screen *ebiten.Image)
	Move(x, y float64)
	Measure() (float64, float64)
	SetColor(color color.Color)
	SetVisible(visible bool)
}

type panelImp struct {
	id         Id
	visible    bool
	color      color.Color
	x, y, w, h float64
}

func (p panelImp) GetId() Id {
	return p.id
}

func (p *panelImp) Draw(screen *ebiten.Image) {
	if p.visible {
		vector.DrawFilledRect(screen, float32(p.x), float32(p.y), float32(p.w), float32(p.h), p.color, false)
	}
}

func (p *panelImp) Move(x float64, y float64) {
	p.x = x
	p.y = y
}

func (p panelImp) Measure() (width float64, height float64) {
	return p.w, p.h
}

func (p *panelImp) SetColor(color color.Color) {
	p.color = color
}

func (p *panelImp) SetVisible(visible bool) {
	p.visible = visible
}

func New(id Id, width float64, height float64, color color.Color) Panel {
	return &panelImp{
		id:      id,
		visible: false,
		color:   color,
		w:       width,
		h:       height,
	}
}
