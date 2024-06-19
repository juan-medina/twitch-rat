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

package colors

import (
	"fmt"
	"strconv"
)

const (
	TEXT_TAG = 'ͼ'
)

type CustomColor interface {
	RGBA() (r, g, b, a uint32)
	Tag() string
	NewWithAlpha(a uint8) CustomColor
}

type customColorImpl struct {
	r, g, b, a uint8
	str        string
}

var (
	Black       = New(0x00, 0x00, 0x00, 0xFF)
	White       = New(0xFF, 0xFF, 0xFF, 0xFF)
	Magenta     = New(0xFF, 0x00, 0xFF, 0xFF)
	LightGray   = New(0xC8, 0xC8, 0xC8, 0xFF)
	Gray        = New(0x82, 0x82, 0x82, 0xFF)
	DarkGray    = New(0x50, 0x50, 0x50, 0xFF)
	Yellow      = New(0xCF, 0xCF, 0x00, 0xFF)
	DarkYellow  = New(0x4E, 0x4E, 0x00, 0xFF)
	LightYellow = New(0xFD, 0xFD, 0x00, 0xFF)
	Gold        = New(0xFF, 0xCB, 0x00, 0xFF)
	Orange      = New(0xFF, 0xA1, 0x00, 0xFF)
	Pink        = New(0xFF, 0x6D, 0xC2, 0xFF)
	Red         = New(0xE6, 0x29, 0x37, 0xFF)
	Maroon      = New(0xBE, 0x21, 0x37, 0xFF)
	Green       = New(0x00, 0xE4, 0x30, 0xFF)
	Lime        = New(0x9E, 0x3E, 0x2F, 0xFF)
	DarkGreen   = New(0x75, 0x2C, 0x2C, 0xFF)
	SkyBlue     = New(0x66, 0xBF, 0xFF, 0xFF)
	Blue        = New(0x79, 0xF1, 0xF1, 0xFF)
	DarkBlue    = New(0x52, 0xAC, 0xFF, 0xFF)
	Purple      = New(0xC8, 0x7A, 0xFF, 0xFF)
	Violet      = New(0x87, 0x3C, 0xBE, 0xFF)
	DarkPurple  = New(0x70, 0x1F, 0x7E, 0xFF)
	Beige       = New(0xD3, 0xB0, 0x83, 0xFF)
	Brown       = New(0x7F, 0x6A, 0x4F, 0xFF)
	DarkBrown   = New(0x4C, 0x3F, 0x2F, 0xFF)
	Teal        = New(0x00, 0x25, 0x56, 0xFF)
	DarkTeal    = New(0x00, 0x1A, 0x3D, 0xFF)
)

func (c customColorImpl) RGBA() (r, g, b, a uint32) {
	r = uint32(c.r)
	r |= r << 8
	g = uint32(c.g)
	g |= g << 8
	b = uint32(c.b)
	b |= b << 8
	a = uint32(c.a)
	a |= a << 8
	return
}

func (c *customColorImpl) NewWithAlpha(a uint8) CustomColor {
	return New(c.r, c.g, c.b, a)
}
func (c customColorImpl) Tag() string {
	return c.tagged(TEXT_TAG)
}
func (c customColorImpl) HTML() string {
	return c.tagged('#')
}
func (c customColorImpl) tagged(r rune) string {
	return string(r) + c.str
}

func FromHtml(colorStr string) CustomColor {
	result := Black

	if len(colorStr) > 0 && colorStr[0] == '#' {
		colorStr = colorStr[1:]
	}

	if len(colorStr) == 6 || len(colorStr) == 8 {
		if r, err := strconv.ParseUint(colorStr[0:2], 16, 8); err == nil {
			if g, err := strconv.ParseUint(colorStr[2:4], 16, 8); err == nil {
				if b, err := strconv.ParseUint(colorStr[4:6], 16, 8); err == nil {
					var a uint64 = 255
					if len(colorStr) == 8 {
						if a, err = strconv.ParseUint(colorStr[6:8], 16, 8); err != nil {
							return result
						}
					}
					result = New(uint8(r), uint8(g), uint8(b), uint8(a))
				}
			}
		}
	}

	return result
}

func toStr(r, g, b, a uint8) string {
	return fmt.Sprintf("%02X%02X%02X%02X", r, g, b, a)
}

func New(r, g, b, a uint8) CustomColor {
	return &customColorImpl{
		r:   r,
		g:   g,
		b:   b,
		a:   a,
		str: toStr(r, g, b, a),
	}
}
