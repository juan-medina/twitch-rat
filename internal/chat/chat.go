/*
 * Copyright (c) 2023 Juan Antonio Medina Iglesias
 *
 *  Permission is hereby granted, free of charge, to any person obtaining a copy
 *  of this software and associated documentation files (the "Software"), to deal
 *  in the Software without restriction, including without limitation the rights
 *  to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *  copies of the Software, and to permit persons to whom the Software is
 *  furnished to do so, subject to the following conditions:
 *
 *  The above copyright notice and this permission notice shall be included in
 *  all copies or substantial portions of the Software.
 *
 *  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *  IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *  FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *  AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *  LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 *  THE SOFTWARE.
 */

package chat

import (
	"image/color"
	"strconv"
)

type EventType int

const (
	Connect EventType = iota
	Disconnect
	Message
)

type Event struct {
	Type_     EventType
	Message   string
	Sender    string
	UserColor color.Color
}

type Chat interface {
	Connect(channel string)
	Disconnect()
	OnEvent(callback func(e Event))
}

var registeredImpl func() Chat = func() Chat {
	panic("chat implementation not registered")
}

func registerImpl(impl func() Chat) {
	registeredImpl = impl
}

func paseHtmlColor(colorStr string) color.Color {
	result := color.RGBA{0, 0, 0, 255}

	if len(colorStr) > 0 && colorStr[0] == '#' {
		colorStr = colorStr[1:]
	}

	if len(colorStr) == 6 {
		if r, err := strconv.ParseUint(colorStr[0:2], 16, 8); err == nil {
			if g, err := strconv.ParseUint(colorStr[2:4], 16, 8); err == nil {
				if b, err := strconv.ParseUint(colorStr[4:6], 16, 8); err == nil {
					result = color.RGBA{uint8(r), uint8(g), uint8(b), 255}
				}
			}
		}
	}

	return result
}

func New() Chat {
	return registeredImpl()
}
