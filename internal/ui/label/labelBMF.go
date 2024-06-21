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
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/juan-medina/twitch-rat/internal/audio"
	"github.com/juan-medina/twitch-rat/internal/colors"
	"github.com/juan-medina/twitch-rat/internal/draw"
	"github.com/juan-medina/twitch-rat/internal/keys"
)

const (
	CLICK_SENT_DELAY = 200
	CLICK_SOUND      = "embed/sounds/click.ogg"
)

type linkState int

const (
	Enabled linkState = iota
	Hover
	Pressed
)

var (
	linkEnabledColor = colors.Yellow
	linkHoverColor   = colors.Purple
	linkPressedColor = colors.Violet
)

type linkInfo struct {
	start           int
	end             int
	url             string
	x, y, w, h      float64
	state           linkState
	timeToSendClick int
}

func (l linkInfo) Hit(x, y float64) bool {
	return x >= l.x && x <= l.x+l.w && y >= l.y && y <= l.y+l.h
}

type labelBMPF struct {
	id               Id
	text             string
	visible          bool
	lineHeight       float64
	font             draw.Font
	color            color.Color
	x, y             float64
	hasBackground    bool
	backgroundColor  color.Color
	bgX, bgY         float64
	bgW, bgH         float64
	expandBackground float64
	links            []linkInfo
	audio            audio.Player
}

func (l labelBMPF) GetId() Id {
	return l.id
}

func (l *labelBMPF) Draw(screen *ebiten.Image) {
	if l.visible {
		if l.hasBackground {
			vector.DrawFilledRect(screen, float32(l.bgX), float32(l.bgY), float32(l.bgW), float32(l.bgH), l.backgroundColor, false)
		}
		l.font.Draw(screen, l.text, l.x, l.y, l.lineHeight, l.color)
		for _, link := range l.links {
			var color colors.CustomColor
			switch link.state {
			case Enabled:
				color = linkEnabledColor
			case Hover:
				color = linkHoverColor
			case Pressed:
				color = linkPressedColor
			}
			l.font.Draw(screen, link.url, link.x+l.x, link.y+l.y, l.lineHeight, color)
		}
	}
}

func (l *labelBMPF) Move(x float64, y float64) {
	l.x = x
	l.y = y
	l.calculateBackground()
}

func (l *labelBMPF) SetText(text string) {
	l.text = text
	l.calculateBackground()
}

func (l labelBMPF) GetText() string {
	return l.text
}

func (l labelBMPF) Measure() (width float64, height float64) {
	width, height = l.font.Measure(l.text, l.lineHeight)
	width += l.expandBackground * 2
	height += l.expandBackground * 2
	return
}

func (l *labelBMPF) SetColor(color color.Color) {
	l.color = color
}
func (l labelBMPF) GetColor() color.Color {
	return l.color
}

func (l *labelBMPF) SetAlpha(alpha float32) {
	r, g, b, _ := l.color.RGBA()
	l.color = color.RGBA{uint8(r), uint8(g), uint8(b), uint8(alpha * 255)}
}

func (l *labelBMPF) SetVisible(visible bool) {
	l.visible = visible
}

func (l *labelBMPF) SetBackgroundColor(color color.Color, expand float64) {
	l.hasBackground = color != nil
	l.backgroundColor = color
	l.expandBackground = expand
	l.calculateBackground()
}

func (l *labelBMPF) calculateBackground() {
	if l.hasBackground {
		w, h := l.font.Measure(l.text, l.lineHeight)
		l.bgX = l.x - l.expandBackground
		l.bgY = l.y - l.expandBackground
		l.bgW = w + l.expandBackground*2
		l.bgH = h + l.expandBackground*2
	}
}

func (l *labelBMPF) ParseLinks() {
	text := l.stripColorTags()
	currentY := 0.0
	for _, link := range strings.Split(text, "\n") {
		start := strings.Index(link, "https://")
		if start != -1 {
			end := strings.Index(link[start:], " ")
			if end == -1 {
				end = len(link)
			} else {
				end += start
			}
			found := link[start:end]
			preLink := link[:start]
			preLinkWidth, _ := l.font.Measure(preLink, l.lineHeight)
			linkWidth, _ := l.font.Measure(found, l.lineHeight)

			l.links = append(l.links, linkInfo{
				start: start,
				end:   end,
				url:   found,
				x:     preLinkWidth,
				y:     currentY,
				w:     linkWidth,
				h:     l.font.GetLineHeight(),
			})
		}

		currentY += l.font.GetLineHeight()
	}
}
func (l labelBMPF) stripColorTags() string {
	result := make([]rune, 0, len(l.text))
	skip := 0
	for _, rune := range l.text {
		if rune == colors.TEXT_TAG {
			skip = 8
		} else {
			if skip > 0 {
				skip--
			} else {
				result = append(result, rune)
			}
		}
	}

	return string(result)
}

func (l *labelBMPF) Update(elapsedTime int, mouseX, mouseY float64, leftJustPressed bool, leftPressed bool, keys keys.Keys) {
	if !l.visible {
		return
	}
	for i, link := range l.links {
		if link.Hit(mouseX-l.x, mouseY-l.y) {
			if l.links[i].state == Hover {
				ebiten.SetCursorShape(ebiten.CursorShapePointer)
				if leftJustPressed {
					if l.links[i].timeToSendClick == 0 && l.links[i].state != Pressed {
						l.links[i].state = Pressed
						l.links[i].timeToSendClick = CLICK_SENT_DELAY
						return
					}
				} else {
					l.links[i].state = Hover
				}
			} else {
				l.links[i].state = Hover
			}
		} else {
			l.links[i].state = Enabled
		}
		if l.links[i].timeToSendClick != 0 {
			l.links[i].timeToSendClick -= elapsedTime
			if l.links[i].timeToSendClick <= 0 {
				l.links[i].timeToSendClick = 0
				l.links[i].state = Enabled
				l.audio.PlaySound(CLICK_SOUND)
				l.newTab(l.links[i].url)
			}
		}
	}
}

func NewLabel(id Id, text string, font draw.Font, lineHeight float64, color color.Color, audio audio.Player) Label {
	return &labelBMPF{
		id:         id,
		text:       text,
		visible:    false,
		lineHeight: lineHeight,
		font:       font,
		color:      color,
		links:      make([]linkInfo, 0, 10),
		audio:      audio,
	}
}
