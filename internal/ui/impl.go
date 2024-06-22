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

package ui

import (
	"embed"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/juan-medina/twitch-rat/internal/audio"
	"github.com/juan-medina/twitch-rat/internal/draw"
	"github.com/juan-medina/twitch-rat/internal/keys"
	"github.com/juan-medina/twitch-rat/internal/ui/button"
	"github.com/juan-medina/twitch-rat/internal/ui/input"
	"github.com/juan-medina/twitch-rat/internal/ui/label"
	"github.com/juan-medina/twitch-rat/internal/ui/panel"
	"github.com/juan-medina/twitch-rat/internal/ui/radiogroup"
	"github.com/juan-medina/twitch-rat/internal/ui/scores"
	"github.com/juan-medina/twitch-rat/internal/ui/slider"
	"github.com/juan-medina/twitch-rat/internal/version"
)

type uiImpl struct {
	screenWidth   float64
	screenHeight  float64
	fileSystem    embed.FS
	buttons       []button.Button
	inputs        []input.Input
	labels        []label.Label
	sliders       []slider.Slider
	panels        []panel.Panel
	radioGroups   []radiogroup.RadioGroup
	keys          keys.Keys
	audioPlayer   audio.Player
	fontVerySmall draw.Font
	fontSmall     draw.Font
	fontNormal    draw.Font
	fontBig       draw.Font
	flyingTexts   []flyingText
	scores        scores.Scores
	sheet         draw.Sheet
	widgets       []Widget
	version       version.Version
}

func (u *uiImpl) Init(fileSystem embed.FS, keys keys.Keys, sheet draw.Sheet) {
	u.sheet = sheet
	u.fileSystem = fileSystem
	u.keys = keys

	u.labels = make([]label.Label, 0, MAX_LABELS)
	u.widgets = make([]Widget, 0, MAX_WIDGETS)
	u.buttons = make([]button.Button, 0, MAX_BUTTONS)
	u.inputs = make([]input.Input, 0, MAX_INPUTS)
	u.panels = make([]panel.Panel, 0, MAX_PANELS)
	u.radioGroups = make([]radiogroup.RadioGroup, 0, MAX_RADIO_GROUPS)
	u.widgets = append(u.widgets, u.scores)

	u.scores.Init()

	u.createWidgets()
}

func (u *uiImpl) Draw(screen *ebiten.Image) {
	u.drawWidgets(screen)
	u.drawFlyingText(screen)
}

func (u *uiImpl) Update(elapsedTime int) {
	xe, ye := ebiten.CursorPosition()
	x, y := float64(xe), float64(ye)

	justPressed := inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
	pressed := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)

	ebiten.SetCursorShape(ebiten.CursorShapeDefault)

	u.updateWidgets(elapsedTime, x, y, justPressed, pressed)
	u.updateFlyingTexts(elapsedTime)
}

func (u *uiImpl) OnLayoutChange(width, height float64) {
	u.screenWidth = width
	u.screenHeight = height

	cx := u.screenWidth / 2
	cy := MENU_START

	u.layoutMainElements(cx, cy)
	cy += TITLE_TO_ELEMENTS_SEPARATION

	u.layoutMainMenuElements(cx, cy)
	u.layoutAboutSubMenuElements(cx, cy)
	u.layoutOptionsSubMenuElements(cx, cy)
	u.layoutGameModeSettingsSubMenuElements(cx, cy)

	u.layoutLicenseElements(cx, cy)
	u.layoutUpdate(cx, cy)
	u.layoutCounter()

	u.scores.Move(SCORE_X, SCORE_Y)
}
