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
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/juan-medina/twitch-rat/internal/audio"
	"github.com/juan-medina/twitch-rat/internal/colors"
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

type UI interface {
	Init(fileSystem embed.FS, keys keys.Keys, sheet draw.Sheet)
	Update(elapsedTime int)

	Draw(screen *ebiten.Image)

	SetStatusMessage(message string, color color.Color)

	EnableButton(id button.Id)
	DisableButton(id button.Id)
	SetButtonVisible(id button.Id, visible bool)
	ClickButton(id button.Id)
	SetButtonClickCallback(callback func(id button.Id))

	GetInputText(id input.Id) string
	SetInputText(id input.Id, text string)
	SetInputVisible(id input.Id, visible bool)
	IsInputEditing(id input.Id) bool

	OnLayoutChange(width, height float64)

	SetLabelText(id label.Id, text string)
	SetLabelColor(id label.Id, color color.Color)
	GetLabelColor(id label.Id) color.Color
	SetLabelVisible(id label.Id, visible bool)
	SetLabelBackgroundColor(id label.Id, color color.Color, expand float64)

	SetSliderVisible(id slider.Id, visible bool)
	SetSliderValue(id slider.Id, value float64)
	SetSliderChangeCallback(callback func(id slider.Id, value float64))
	AddFlyingText(text string, color colors.CustomColor, x, y float64)

	SetScoreVisible(visible bool)
	AddScoreEntry(data scores.ScoreData)
	StartsCore()

	SetPanelVisible(id panel.Id, visible bool)
	SetRadioGroupVisible(id radiogroup.Id, visible bool)
	SelectRadioGroup(id radiogroup.Id, index int)
}

func New(version version.Version, audioPlayer audio.Player, fontVerySmall draw.Font, fontSmall draw.Font, fontNormal draw.Font, fontBig draw.Font) UI {
	return &uiImpl{
		audioPlayer:   audioPlayer,
		fontVerySmall: fontVerySmall,
		fontSmall:     fontSmall,
		fontNormal:    fontNormal,
		fontBig:       fontBig,
		flyingTexts:   make([]flyingText, 0, MAX_FLYING_TEXTS),
		scores:        scores.New(fontSmall, fontNormal),
		version:       version,
	}
}
