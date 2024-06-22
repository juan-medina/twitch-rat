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
	"github.com/prgra/bbcode"
)

const (
	CLICK_SENT_DELAY = 200
	CLICK_SOUND      = "embed/sounds/click.ogg"
)

type linkState int

const (
	None linkState = iota
	Enabled
	Hover
	Pressed
)

var (
	linkEnabledColor     = colors.Yellow
	linkHoverColor       = colors.Purple
	linkPressedColor     = colors.Violet
	linkUnderliningColor = colors.Violet
)

type linePart struct {
	text            string
	color           colors.CustomColor
	link            string
	w, h            float64
	linkState       linkState
	timeToSendClick int
	hasAlpha        bool
}

type line struct {
	parts []linePart
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
	audio            audio.Player
	lines            []line
}

func (l labelBMPF) GetId() Id {
	return l.id
}

func (l *labelBMPF) Draw(screen *ebiten.Image) {
	if l.visible {
		if l.hasBackground {
			vector.DrawFilledRect(screen, float32(l.bgX), float32(l.bgY), float32(l.bgW), float32(l.bgH), l.backgroundColor, false)
		}
		var currentX float64
		currentY := l.y
		lineHeight := l.font.GetLineHeight()
		_, _, _, a := l.color.RGBA()
		originalAlpha := uint8(a >> 8)
		for _, line := range l.lines {
			var color color.Color
			currentX = l.x
			for _, part := range line.parts {
				if part.color != nil {
					if part.hasAlpha {
						color = part.color
					} else {
						color = part.color.NewWithAlpha(originalAlpha)
					}

				} else {
					color = l.color
				}
				l.font.Draw(screen, part.text, currentX, currentY, l.lineHeight, color)
				if part.link != "" {
					vector.StrokeLine(screen, float32(currentX), float32(currentY+l.font.DefaultSize()), float32(currentX+part.w), float32(currentY+l.font.DefaultSize()), 2, linkUnderliningColor, false)
				}
				currentX += part.w
			}
			currentY += lineHeight
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
	l.parse()
	l.calculateBackground()
}

type parseElements struct {
	text  string
	code  string
	param string
}

func (l labelBMPF) parseBBCode(text string) []parseElements {
	var result []parseElements

	decoded := bbcode.Parse(text)
	lastStart := 0
	for _, code := range decoded.BBCodes {
		begin := code.Pos - 1
		end := code.Pos + code.Len - 1
		if (code.Pos - 1) > lastStart {
			result = append(result, parseElements{text: decoded.NewString[lastStart:begin]})
		}
		if code.Len > 0 {
			lastStart = end
			result = append(result, parseElements{text: code.Text, code: code.Name, param: code.Param})
		}
	}
	if lastStart < len(decoded.NewString) {
		result = append(result, parseElements{text: decoded.NewString[lastStart:]})
	}

	return result
}

func (l *labelBMPF) parse() {
	l.lines = make([]line, 0)
	newLines := strings.Split(l.text, "\n")
	var currentColor colors.CustomColor
	for _, nl := range newLines {
		newParts := make([]linePart, 0)
		for _, e := range l.parseBBCode(nl) {
			link := ""
			linkState := None
			hasAlpha := false
			currentColor = nil
			if e.code == "color" {
				colorLen := len(e.param)
				if colorLen > 7 {
					hasAlpha = true
				}
				currentColor = colors.FromHtml(e.param)
			} else if e.code == "url" {
				link = e.param
				linkState = Enabled
			}
			w, h := l.font.Measure(e.text, l.lineHeight)
			part := linePart{
				text:      e.text,
				w:         w,
				h:         h,
				hasAlpha:  hasAlpha,
				color:     currentColor,
				link:      link,
				linkState: linkState,
			}
			newParts = append(newParts, part)
		}

		l.lines = append(l.lines, line{parts: newParts})
	}
}

func (l labelBMPF) GetText() string {
	return l.text
}

func (l labelBMPF) Measure() (width float64, height float64) {
	width, height = l.measureWithoutBackground()

	width += l.expandBackground * 2
	height += l.expandBackground * 2
	return
}

func (l labelBMPF) measureWithoutBackground() (width float64, height float64) {
	var currentX float64
	currentY := 0.0
	lineHeight := l.font.GetLineHeight()
	for _, line := range l.lines {
		currentX = 0.0
		for _, part := range line.parts {
			currentX += part.w
		}
		currentY += lineHeight
		if currentX > width {
			width = currentX
		}
	}
	height = currentY

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
		w, h := l.measureWithoutBackground()
		l.bgX = l.x - l.expandBackground
		l.bgY = l.y - l.expandBackground
		l.bgW = w + l.expandBackground*2
		l.bgH = h + l.expandBackground*2
	}
}

func (l *labelBMPF) Update(elapsedTime int, mouseX, mouseY float64, leftJustPressed bool, leftPressed bool, keys keys.Keys) {
	if !l.visible {
		return
	}

	var currentX float64
	currentY := l.y
	lineHeight := l.font.GetLineHeight()
	for i := range l.lines {
		currentX = l.x
		for j := range l.lines[i].parts {
			if l.lines[i].parts[j].link != "" {
				if mouseX >= currentX &&
					mouseX <= currentX+l.lines[i].parts[j].w &&
					mouseY >= currentY &&
					mouseY <= currentY+l.lines[i].parts[j].h {
					ebiten.SetCursorShape(ebiten.CursorShapePointer)
					if leftJustPressed {
						if l.lines[i].parts[j].timeToSendClick == 0 && l.lines[i].parts[j].linkState != Pressed {
							l.lines[i].parts[j].linkState = Pressed
							l.lines[i].parts[j].timeToSendClick = CLICK_SENT_DELAY
							l.lines[i].parts[j].color = linkPressedColor
							return
						}
					} else {
						l.lines[i].parts[j].linkState = Hover
						l.lines[i].parts[j].color = linkHoverColor
					}
				} else {
					l.lines[i].parts[j].linkState = Enabled
					l.lines[i].parts[j].color = linkEnabledColor
				}
				if l.lines[i].parts[j].timeToSendClick != 0 {
					l.lines[i].parts[j].timeToSendClick -= elapsedTime
					if l.lines[i].parts[j].timeToSendClick <= 0 {
						l.lines[i].parts[j].timeToSendClick = 0
						l.lines[i].parts[j].linkState = Enabled
						l.lines[i].parts[j].color = linkEnabledColor
						l.audio.PlaySound(CLICK_SOUND)
						l.newTab(l.lines[i].parts[j].link)
					}
				}
			}
			currentX += l.lines[i].parts[j].w
		}
		currentY += lineHeight
	}
}

func NewLabel(id Id, text string, font draw.Font, lineHeight float64, color color.Color, audio audio.Player) Label {
	l := labelBMPF{
		id:         id,
		text:       "",
		visible:    false,
		lineHeight: lineHeight,
		font:       font,
		color:      color,
		audio:      audio,
	}
	l.SetText(text)
	return &l
}
