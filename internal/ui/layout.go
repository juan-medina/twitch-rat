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

import "github.com/juan-medina/twitch-rat/internal/ui/label"

func (u *uiImpl) layoutCounter() {
	cx := u.screenWidth / 2
	cy := MENU_START

	countdownLabel := u.getLabel(LABEL_COUNTDOWN)
	ix, iy := countdownLabel.Measure()
	px := cx - ix
	py := cy
	countdownLabel.Move(px, py)

	py = py + iy + BUTTON_GAP
	w, _ := u.getLabel(LABEL_INSTRUCTIONS).Measure()
	u.moveLabel(LABEL_INSTRUCTIONS, cx-(w/2), py)

	cx = (u.screenWidth / 2) - (OPTION_PANEL_WIDTH / 2)
	cy = cy + 100
	u.movePanel(OPTIONS_PANEL, cx, cy)
}

func (u *uiImpl) layoutMainElements(cx, cy float64) {

	titleLabel := u.getLabel(LABEL_TITLE)
	ix, _ := titleLabel.Measure()
	px := cx - (ix / 2)
	py := cy
	titleLabel.Move(px, py)

	bb := u.getButton(BACK_BUTTON)
	bw, _ := bb.Size()
	px = u.screenWidth - bw
	py = 0
	u.moveButton(BACK_BUTTON, px, py)

	bb = u.getButton(IN_GAME_OPTIONS_BUTTON)
	bw, _ = bb.Size()
	px -= bw
	u.moveButton(IN_GAME_OPTIONS_BUTTON, px, py)

	gapX := float64(u.fontSmall.DefaultSize()) * 0.5
	gapY := float64(u.fontSmall.DefaultSize()) * 1.5

	py = u.screenHeight - gapY

	for i := 0; i < TOTAL_LAST_MESSAGES; i++ {
		labelId := LABEL_LAST_MESSAGE + label.Id(i)
		u.getLabel(labelId).Move(gapX, py)
		py -= u.fontSmall.DefaultSize()
	}

	gapX = u.fontSmall.DefaultSize() * 0.5
	gapY = u.fontSmall.DefaultSize() * 1.5
	versionLabel := u.getLabel(LABEL_VERSION)
	cx, _ = versionLabel.Measure()
	versionLabel.Move(u.screenWidth-cx-gapX, u.screenHeight-gapY)

	px = 450
	py = 0
	u.getInput(INPUT_DEBUG_USER).Move(px, py)

	px = px + INPUT_WIDTH + BUTTON_GAP
	u.getInput(INPUT_DEBUG_MESSAGE).Move(px, py)

	px = px + INPUT_WIDTH + BUTTON_GAP
	u.getButton(DEBUG_BUTTON).Move(px, py)
}

func (u *uiImpl) layoutMainMenuElements(cx, cy float64) {
	px := cx - (INPUT_WIDTH / 2)
	py := cy + (INPUT_HEIGHT / 2) + BUTTON_GAP*2
	u.moveInput(INPUT_CHANNEL, px, py)

	px = cx - (BUTTON_WIDTH / 2)
	py = py + INPUT_HEIGHT + BUTTON_GAP
	u.moveButton(PLAY_BUTTON, px, py)

	py = py + BUTTON_HEIGHT + BUTTON_GAP
	u.moveButton(OPTIONS_BUTTON, px, py)

	px = px + (BUTTON_WIDTH - SMALL_BUTTON_WIDTH)
	u.moveButton(ABOUT_BUTTON, px, py)

	lw, _ := u.getLabel(LABEL_DOWNLOAD).Measure()
	py = py + BUTTON_HEIGHT + BUTTON_GAP
	px = cx - (lw / 2)
	u.moveLabel(LABEL_DOWNLOAD, px, py)
}

func (u *uiImpl) layoutLicenseElements(cx, cy float64) {
	py := cy + BUTTON_GAP*3
	w, h := u.getLabel(LABEL_LICENSE).Measure()
	u.moveLabel(LABEL_LICENSE, cx-(w/2), py)

	px := cx - (SMALL_BUTTON_WIDTH / 2)
	py = py + h + BUTTON_GAP
	u.moveButton(ACCEPT_LICENSE_BUTTON, px, py)
}

func (u *uiImpl) layoutUpdate(cx, cy float64) {
	py := cy + BUTTON_GAP*6
	w, h := u.getLabel(LABEL_VERSION_UPDATE).Measure()
	u.moveLabel(LABEL_VERSION_UPDATE, cx-(w/2), py)

	px := cx + BUTTON_GAP
	py = py + h + BUTTON_GAP*4
	u.moveButton(CONTINUE_BUTTON, px, py)

	px = cx - (BUTTON_WIDTH) - BUTTON_GAP
	u.moveButton(DOWNLOAD_LATEST_BUTTON, px, py)
}

func (u *uiImpl) layoutAboutSubMenuElements(cx, cy float64) {
	py := cy + BUTTON_GAP*3
	w, h := u.getLabel(LABEL_ABOUT_MESSAGE).Measure()
	u.moveLabel(LABEL_ABOUT_MESSAGE, cx-(w/2), py)

	px := cx - (SMALL_BUTTON_WIDTH / 2)
	py = py + h + BUTTON_GAP
	u.moveButton(SUBMENU_ABOUT_BACK_BUTTON, px, py)
}

func (u *uiImpl) layoutOptionsSubMenuElements(cx, cy float64) {
	px := cx - SLIDER_WITH/2
	py := cy + BUTTON_GAP*3

	u.moveLabel(LABEL_OPTIONS_MUSIC_VOLUME, px, py)

	py = py + BUTTON_GAP*2
	u.moveSlider(MUSIC_VOLUME_SLIDER, px, py)

	py = py + BUTTON_GAP*2
	u.moveLabel(LABEL_OPTIONS_AUDIO_VOLUME, px, py)

	py = py + BUTTON_GAP*2
	u.moveSlider(AUDIO_VOLUME_SLIDER, px, py)

	px = cx - (SMALL_BUTTON_WIDTH / 2)
	py = py + SLIDER_HEIGHT + BUTTON_GAP
	u.moveButton(SUBMENU_OPTION_BACK_BUTTON, px, py)
}

func (u *uiImpl) layoutGameModeSettingsSubMenuElements(cx, cy float64) {
	py := cy + BUTTON_GAP*3

	// CENTER LABEL AND RADIO GROUP
	lw, lh := u.getLabel(LABEL_GAME_MODE).Measure()
	rw, _ := u.getRadioGroup(GAME_MODE_RADIO_GROUP).Measure()

	tw := lw + BUTTON_GAP + rw
	px := cx - (tw / 2)
	gy := (BUTTON_HEIGHT - lh) / 2
	u.moveLabel(LABEL_GAME_MODE, px, py+gy)

	px = px + lw + BUTTON_GAP
	u.moveRadioGroup(GAME_MODE_RADIO_GROUP, px, py)

	// BUTTONS
	py = py + BUTTON_HEIGHT + BUTTON_GAP
	px = cx - BUTTON_WIDTH - BUTTON_GAP
	u.moveButton(SUBMENU_GAME_MODE_SETTINGS_BACK_BUTTON, px, py)

	px = cx + (BUTTON_GAP)
	u.moveButton(SUBMENU_GAME_MODE_SETTINGS_GO_BUTTON, px, py)
}
