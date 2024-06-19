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

package license

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/juan-medina/twitch-rat/internal/colors"
	"github.com/juan-medina/twitch-rat/internal/keys"
	"github.com/juan-medina/twitch-rat/internal/stage"
	"github.com/juan-medina/twitch-rat/internal/ui"
	"github.com/juan-medina/twitch-rat/internal/ui/button"
)

func (l *license) Init() {
	l.ui.SetButtonClickCallback(l.onButtonClick)

	l.ui.SetLabelVisible(ui.LABEL_TITLE, true)
	l.ui.SetButtonVisible(ui.ACCEPT_LICENSE_BUTTON, true)
	l.ui.SetLabelVisible(ui.LABEL_LICENSE, true)
	l.ui.SetButtonVisible(ui.BACK_BUTTON, true)

	l.ui.SetStatusMessage("Showing license..", colors.LightYellow)
}

func (l *license) End() {
	l.ui.SetButtonClickCallback(nil)

	l.ui.SetButtonVisible(ui.BACK_BUTTON, false)
	l.ui.SetLabelVisible(ui.LABEL_TITLE, false)
	l.ui.SetButtonVisible(ui.ACCEPT_LICENSE_BUTTON, false)
	l.ui.SetLabelVisible(ui.LABEL_LICENSE, false)
}

type license struct {
	changer stage.Changer
	ui      ui.UI
	width   float32
	height  float32
}

func (l *license) Update(elapsedTime int, keys keys.Keys) {
	l.ui.Update(elapsedTime)
	if keys.IsDownNoRepeat(ebiten.KeyEnter) || keys.IsDownNoRepeat(ebiten.KeyNumpadEnter) {
		l.ui.ClickButton(ui.ACCEPT_LICENSE_BUTTON)
	}
	if keys.IsDownNoRepeat(ebiten.KeyEscape) {
		l.ui.ClickButton(ui.BACK_BUTTON)
	}
}

func (l *license) Draw(screen *ebiten.Image) {
	vector.DrawFilledRect(screen, 0, 0, l.width, l.height, colors.Teal, false)
	l.ui.Draw(screen)
}

func (l *license) OnLayoutChange(width, height float64) {
	l.width = float32(width)
	l.height = float32(height)
}

func (l *license) onButtonClick(id button.Id) {
	switch id {
	case ui.ACCEPT_LICENSE_BUTTON:
		l.changer.ChangeStage(stage.MENU)
	case ui.BACK_BUTTON:
		l.changer.ChangeStage(stage.EXIT)
	}
}

func New(changer stage.Changer, ui ui.UI) stage.Stage {
	m := license{
		changer: changer,
		ui:      ui,
	}
	return &m
}
