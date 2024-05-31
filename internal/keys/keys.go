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
	Update(elapsed int)
	IsKeyDown(key ebiten.Key) bool
	IsKeyRepeat(key ebiten.Key) bool
	LastRepeatedKey() ebiten.Key
}

type keyStatus struct {
	isDown           bool
	repeatDown       bool
	timeToNextRepeat int
}

const (
	REPEAT_DELAY = 200
)

type keyManagerImpl struct {
	status          map[ebiten.Key]keyStatus
	lastRepeatedKey ebiten.Key
}

func (k *keyManagerImpl) Init() {
	for i := ebiten.Key(0); i < ebiten.KeyMax; i++ {
		k.status[i] = keyStatus{
			isDown:     false,
			repeatDown: false,
		}
	}
}

func (k *keyManagerImpl) Update(elapsed int) {
	k.lastRepeatedKey = ebiten.KeyMax
	for id := range k.status {
		currentStatus := k.status[id]
		currentlyEbitenDown := ebiten.IsKeyPressed(id)
		currentStatus.isDown = currentlyEbitenDown
		if currentlyEbitenDown {
			if currentStatus.timeToNextRepeat <= 0 {
				currentStatus.timeToNextRepeat = REPEAT_DELAY
				currentStatus.repeatDown = true
			} else {
				currentStatus.timeToNextRepeat -= elapsed
				currentStatus.repeatDown = false
			}
		} else {
			currentStatus.timeToNextRepeat = 0
			currentStatus.repeatDown = false
		}

		if currentStatus.repeatDown {
			k.lastRepeatedKey = id
		}

		k.status[id] = currentStatus
	}
}

func (k *keyManagerImpl) IsKeyDown(key ebiten.Key) bool {
	return k.status[key].isDown
}

func (k *keyManagerImpl) IsKeyRepeat(key ebiten.Key) bool {
	return k.status[key].repeatDown
}

func (k keyManagerImpl) LastRepeatedKey() ebiten.Key {
	return k.lastRepeatedKey
}

func New() Keys {
	return &keyManagerImpl{
		status:          make(map[ebiten.Key]keyStatus, ebiten.KeyMax),
		lastRepeatedKey: ebiten.KeyMax,
	}
}
