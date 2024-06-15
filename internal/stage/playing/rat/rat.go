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

package rat

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/juan-medina/twitch-rat/internal/draw"
	"github.com/juan-medina/twitch-rat/internal/ui/label"
)

const (
	RAT_SCALE              = 4
	RAT_Y_POS              = 768
	LABEL_GAP              = 20
	ARENA_LEFT_X           = -750
	ARENA_RIGHT_X          = 670
	WALK_SPEED             = 0.05
	RUN_SPEED              = 0.2
	CLOSE_TO_OBJECTIVE     = 1.0
	WAIT_TO_WALK_AGAIN_MIN = 3000
	WAIT_TO_WALK_AGAIN_MAX = 4500
)

type Rat interface {
	Draw(screen *ebiten.Image)
	Update(elapsedTime int)
	SetCenter(screenCenterX float64)
	SetX(x float64)
	SetAnimation(animation string)
	GetName() string
	SetVisible(visible bool)
	IsVisible() bool
	RandomWalk()
}

type animation struct {
	pattern       string
	startFrame    int
	endFrame      int
	frameDuration int
	loop          bool
}

type animationStatus struct {
	currentFrame int
	currentTime  int
	animation    *animation
	end          bool
}

type direction int

const (
	left direction = iota
	right
)

type state int

const (
	idle state = iota
	walking
)

type ratImpl struct {
	x, y, screenCenterX float64
	vx                  float64
	sheet               draw.Sheet
	sprite              draw.Sprite
	animationStatus     animationStatus
	name                string
	nameLabel           label.Label
	labelWidth          float64
	visible             bool
	facing              direction
	state               state
	destination         float64
	waitingTime         int
}

func (r ratImpl) IsVisible() bool {
	return r.visible
}

func (r *ratImpl) SetVisible(visible bool) {
	r.visible = visible
}

func (r *ratImpl) Draw(screen *ebiten.Image) {
	if !r.visible {
		return
	}
	r.sprite.Draw(screen, r.screenCenterX+r.x, r.y, r.facing == left, false)
	r.nameLabel.Draw(screen)
}

func (r *ratImpl) SetCenter(screenCenterX float64) {
	r.screenCenterX = screenCenterX
	r.y = RAT_Y_POS
	r.moveLabel()
}

func (r *ratImpl) SetX(x float64) {
	r.x = x
	r.y = RAT_Y_POS
	r.moveLabel()
}

func (r *ratImpl) moveLabel() {
	_, rh := r.sprite.Size()
	y := r.y - rh - LABEL_GAP
	x := r.screenCenterX + r.x - (r.labelWidth / 2)
	r.nameLabel.Move(x, y)
}

func (r *ratImpl) Update(elapsedTime int) {
	if !r.visible {
		return
	}

	r.SetX(r.x + r.vx*float64(elapsedTime))

	if r.state == walking {
		diff := math.Abs(r.destination - r.x)
		if diff < CLOSE_TO_OBJECTIVE {
			r.idle()
		}
	} else if r.state == idle {
		if r.waitingTime > 0 {
			r.waitingTime -= elapsedTime
		} else {
			r.RandomWalk()
		}
	}
	if !r.animationStatus.end {
		r.animationStatus.currentTime += elapsedTime
		if r.animationStatus.currentTime >= r.animationStatus.animation.frameDuration {
			r.animationStatus.currentTime = 0
			r.animationStatus.currentFrame++
			if r.animationStatus.currentFrame > r.animationStatus.animation.endFrame {
				if r.animationStatus.animation.loop {
					r.animationStatus.currentFrame = r.animationStatus.animation.startFrame
				} else {
					r.animationStatus.end = true
					r.animationStatus.currentFrame = r.animationStatus.animation.endFrame
				}
			}
			r.sprite = r.sheet.Sprite(fmt.Sprintf(r.animationStatus.animation.pattern, r.animationStatus.currentFrame))
			r.sprite.SetScale(RAT_SCALE)
		}
	}

}

func (r ratImpl) GetName() string {
	if !r.visible {
		return ""
	}
	return r.name
}
func (r *ratImpl) SetAnimation(animation string) {
	if a, ok := animationMap[animation]; ok {
		r.animationStatus = animationStatus{
			currentFrame: a.startFrame,
			currentTime:  0,
			animation:    &a,
			end:          false,
		}
		r.sprite = r.sheet.Sprite(fmt.Sprintf(a.pattern, r.animationStatus.currentFrame))
		r.sprite.SetScale(RAT_SCALE)
	} else {
		panic("animation not found in set animation: " + animation)
	}
}
func (r *ratImpl) RandomWalk() {
	r.vx = WALK_SPEED
	r.destination = ARENA_LEFT_X + rand.Float64()*(ARENA_RIGHT_X-ARENA_LEFT_X)
	if r.destination < r.x {
		r.facing = left
	} else {
		r.facing = right
	}
	r.SetAnimation(WALK_ANIM)
	if r.facing == left {
		r.vx = -r.vx
	}
	r.state = walking
	r.waitingTime = 0
}
func (r *ratImpl) idle() {
	r.state = idle
	r.SetAnimation(IDLE_ANIM)
	r.vx = 0
	r.waitingTime = rand.Intn(WAIT_TO_WALK_AGAIN_MAX-WAIT_TO_WALK_AGAIN_MIN) + WAIT_TO_WALK_AGAIN_MIN
}

func New(sheet draw.Sheet, name string, face text.Face) Rat {
	label := label.New(0, name, face, 0, color.White)
	labelWidth, _ := label.Measure()
	label.SetVisible(true)
	return &ratImpl{
		name:       name,
		x:          0,
		y:          0,
		vx:         0,
		sheet:      sheet,
		nameLabel:  label,
		labelWidth: labelWidth,
		visible:    true,
		facing:     right,
	}
}
