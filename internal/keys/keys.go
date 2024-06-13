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

package keys

import "github.com/hajimehoshi/ebiten/v2"

type Keys interface {
	Init()
	Update()
	IsDown(key ebiten.Key) bool
	IsDownNoRepeat(key ebiten.Key) bool
	LastInputChar() rune
	SwallowKey(key ebiten.Key)
}

var (
	emptyRune = []rune{}
)

const (
	REPEAT_DELAY = 200
)

type keyStatus struct {
	wasDown         bool
	isDown          bool
	isDownNotRepeat bool
}

type keyManagerImpl struct {
	lastInputChar rune
	status        []keyStatus
}

func (k *keyManagerImpl) Init() {
	k.lastInputChar = 0
}

func (k *keyManagerImpl) Update() {
	k.lastInputChar = 0
	rune := ebiten.AppendInputChars(emptyRune)
	if len(rune) == 1 {
		k.lastInputChar = rune[0]
	}

	for i := 0; i < int(ebiten.KeyMax); i++ {
		k.status[i].isDown = ebiten.IsKeyPressed(ebiten.Key(i))
		if k.status[i].isDown {
			if k.status[i].wasDown {
				k.status[i].isDownNotRepeat = false
			} else {
				k.status[i].isDownNotRepeat = true
				k.status[i].wasDown = true
			}
		} else {
			k.status[i].isDownNotRepeat = false
			k.status[i].wasDown = false
		}
	}
}

func (k keyManagerImpl) IsDown(key ebiten.Key) bool {
	return k.status[key].isDown
}

func (k keyManagerImpl) IsDownNoRepeat(key ebiten.Key) bool {
	return k.status[key].isDownNotRepeat
}

func (k keyManagerImpl) LastInputChar() rune {
	return k.lastInputChar
}

func (k *keyManagerImpl) SwallowKey(key ebiten.Key) {
	k.status[key].isDownNotRepeat = false
}

func New() Keys {
	return &keyManagerImpl{
		lastInputChar: 0,
		status:        make([]keyStatus, ebiten.KeyMax),
	}
}
