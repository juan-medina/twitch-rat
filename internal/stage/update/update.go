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

package update

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/juan-medina/twitch-rat/internal/colors"
	"github.com/juan-medina/twitch-rat/internal/keys"
	"github.com/juan-medina/twitch-rat/internal/stage"
	"github.com/juan-medina/twitch-rat/internal/ui"
	"github.com/juan-medina/twitch-rat/internal/ui/button"
	"github.com/juan-medina/twitch-rat/internal/version"
)

func (u *update) Init() {
	u.ui.SetButtonClickCallback(u.onButtonClick)

	u.ui.SetLabelVisible(ui.LABEL_TITLE, true)
	u.ui.SetLabelVisible(ui.LABEL_VERSION_UPDATE, true)
	u.ui.SetButtonVisible(ui.DOWNLOAD_LATEST_BUTTON, true)
	u.ui.SetButtonVisible(ui.CONTINUE_BUTTON, true)
	u.ui.SetButtonVisible(ui.BACK_BUTTON, true)

	u.ui.SetStatusMessage("Showing update..", colors.LightYellow)
	u.ui.SetLabelText(ui.LABEL_VERSION_UPDATE, fmt.Sprintf(ui.VERSION_OUTDATED_STRING, u.version.Latest().Bbcode))
}

func (u *update) End() {
	u.ui.SetButtonClickCallback(nil)
	u.ui.SetLabelVisible(ui.LABEL_TITLE, false)

	u.ui.SetLabelVisible(ui.LABEL_VERSION_UPDATE, false)
	u.ui.SetButtonVisible(ui.DOWNLOAD_LATEST_BUTTON, false)
	u.ui.SetButtonVisible(ui.CONTINUE_BUTTON, false)
	u.ui.SetButtonVisible(ui.BACK_BUTTON, false)
}

type update struct {
	changer stage.Changer
	ui      ui.UI
	version version.Version
}

func (u *update) Update(elapsedTime int, keys keys.Keys) {
	u.ui.Update(elapsedTime)
	if keys.IsDownNoRepeat(ebiten.KeyEnter) || keys.IsDownNoRepeat(ebiten.KeyNumpadEnter) {
		u.ui.ClickButton(ui.CONTINUE_BUTTON)
	} else if keys.IsDownNoRepeat(ebiten.KeyEscape) {
		u.ui.ClickButton(ui.BACK_BUTTON)
	}
}

func (u *update) Draw(screen *ebiten.Image) {
	u.ui.Draw(screen)
}

func (u *update) OnLayoutChange(width, height float64) {
}

func (u *update) onButtonClick(id button.Id) {
	switch id {
	case ui.CONTINUE_BUTTON:
		u.changer.ChangeStage(stage.MENU)
	case ui.DOWNLOAD_LATEST_BUTTON:
		u.download()
	case ui.BACK_BUTTON:
		u.changer.ChangeStage(stage.EXIT)
	}
}

func New(changer stage.Changer, ui ui.UI, version version.Version) stage.Stage {
	m := update{
		changer: changer,
		ui:      ui,
		version: version,
	}
	return &m
}
