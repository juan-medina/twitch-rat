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
	"math"
	"math/rand"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/juan-medina/twitch-rat/internal/audio"
	"github.com/juan-medina/twitch-rat/internal/colors"
	"github.com/juan-medina/twitch-rat/internal/draw"
	"github.com/juan-medina/twitch-rat/internal/ui"
	"github.com/juan-medina/twitch-rat/internal/ui/label"
)

const (
	HIT_SOUND              = "embed/sounds/hit.ogg"
	DEAD_SOUND             = "embed/sounds/dead.ogg"
	HEAL_SOUND             = "embed/sounds/heal.ogg"
	RAT_SCALE              = 4
	RAT_Y_POS              = 768
	LABEL_GAP              = 25
	ARENA_LEFT_X           = -750
	ARENA_RIGHT_X          = 670
	WALK_SPEED             = 0.1
	RUN_SPEED              = 0.2
	CLOSE_TO_OBJECTIVE     = 8.0 * 4
	CLOSE_TO_OTHER_RAT     = 20.0 * 4
	WAIT_TO_WALK_AGAIN_MIN = 3000
	WAIT_TO_WALK_AGAIN_MAX = 4500
	WAIT_AFTER_HIT         = 1000
	HEALTH_BAR_GAP         = 22
	HEALTH_BAR_WIDTH       = 80
	HEALTH_BAR_HEIGHT      = 15
	HEALTH_MAX             = 100
	RAT_DAMAGE             = 20
	RAT_HEAL               = 20
	CRIT_CHANGE            = 0.25
	MOD_VALUE              = 0.25
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
	CanDoAction() bool
	Attack(otherRat Rat)
	GetX() float64
	Hurt(hit int) (amount int, over int, crit bool)
	IsAlive() bool
	ReSpawn(color colors.CustomColor)
	GetColor() colors.CustomColor
	Heal(otherRat Rat)
	Cure(heal int) (amount int, over int, crit bool)
	Modify(points int) (value int, crit bool)
	AddFlyingText(text string, color colors.CustomColor)
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
	running
	dead
	healing
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
	target              Rat
	ui                  ui.UI
	audioPlayer         audio.Player
	color               colors.CustomColor
	barX                float64
	barY                float64
	health              int
	flyingTextY         float64
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
	r.sprite.SetColor(r.color)
	r.sprite.Draw(screen, r.screenCenterX+r.x, r.y, r.facing == left, false)
	r.nameLabel.Draw(screen)

	greenWidth := float32(float64(HEALTH_BAR_WIDTH) * (float64(r.health) / HEALTH_MAX))
	redWidth := float32(HEALTH_BAR_WIDTH - greenWidth)
	redStart := float32(float32(r.barX) + greenWidth)

	vector.DrawFilledRect(screen, float32(r.barX), float32(r.barY), greenWidth, HEALTH_BAR_HEIGHT, colors.Green, false)
	vector.DrawFilledRect(screen, redStart, float32(r.barY), redWidth, HEALTH_BAR_HEIGHT, colors.Red, false)

	vector.StrokeRect(screen, float32(r.barX), float32(r.barY), HEALTH_BAR_WIDTH, HEALTH_BAR_HEIGHT, 2, colors.Black, false)
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

	r.barX = r.screenCenterX + r.x - HEALTH_BAR_WIDTH/2
	r.barY = y + HEALTH_BAR_GAP

	r.flyingTextY = r.y - (rh / 2)
}

func (r *ratImpl) Update(elapsedTime int) {
	if !r.visible {
		return
	}

	r.SetX(r.x + r.vx*float64(elapsedTime))
	if x := r.GetX(); x < ARENA_LEFT_X || x > ARENA_RIGHT_X {
		r.RandomWalk()
	}

	switch r.state {
	case dead:
		if r.animationStatus.end {
			r.audioPlayer.PlaySound(DEAD_SOUND)
			r.SetVisible(false)
		}
	case walking:
		diff := math.Abs(r.destination - r.x)
		if diff < CLOSE_TO_OBJECTIVE {
			r.idle()
		}
	case running:
		if r.target != nil {
			if r.target.IsAlive() {
				r.updateDestinationToTarget()
				diff := math.Abs(r.destination - r.x)
				if diff < CLOSE_TO_OTHER_RAT {
					r.idle()
					r.SetAnimation(FIGHT_ANIM)
					r.waitingTime = WAIT_AFTER_HIT
					damage, over, crit := r.target.Hurt(RAT_DAMAGE)
					r.logDamage(damage, over, crit)
					r.audioPlayer.PlaySound(HIT_SOUND)
					r.target = nil
				}
			} else {
				r.idle()
			}
		}
	case idle:
		if r.waitingTime > 0 {
			r.waitingTime -= elapsedTime
		} else {
			r.RandomWalk()
		}
	case healing:
		if !r.target.IsAlive() {
			r.target = nil
			r.idle()
		}
		r.updateDestinationToTarget()
		r.vx = 0
		if r.animationStatus.end {
			heal, over, crit := r.target.Cure(RAT_HEAL)
			r.logHeal(heal, over, crit)
			r.target = nil
			r.idle()
			r.audioPlayer.PlaySound(HEAL_SOUND)
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
func (r *ratImpl) Hurt(hit int) (amount int, over int, crit bool) {
	amount, crit = r.Modify(hit)
	original := r.health
	r.health -= amount
	if r.health <= 0 {
		r.health = 0
		r.SetAnimation(DEAD_ANIM)
		r.state = dead
		r.waitingTime = 0
	} else {
		r.SetAnimation(HURT_ANIM)
		r.state = idle
		r.waitingTime = WAIT_AFTER_HIT
	}
	r.vx = 0
	over = amount - (original - r.health)
	return
}

func (r *ratImpl) idle() {
	r.state = idle
	r.SetAnimation(IDLE_ANIM)
	r.vx = 0
	r.waitingTime = rand.Intn(WAIT_TO_WALK_AGAIN_MAX-WAIT_TO_WALK_AGAIN_MIN) + WAIT_TO_WALK_AGAIN_MIN
}

func (r ratImpl) CanDoAction() bool {
	return r.state == idle || r.state == walking
}

func (r *ratImpl) Attack(otherRat Rat) {
	r.state = running
	r.SetAnimation(RUN_ANIM)
	r.target = otherRat
	r.waitingTime = 0
	r.updateDestinationToTarget()
	targetColor := otherRat.GetColor()
	damageStr := r.color.Tag() + r.name +
		colors.White.Tag() + " is attacking " +
		targetColor.Tag() + r.target.GetName()
	r.ui.SetStatusMessage(damageStr, colors.White)
}

func (r *ratImpl) Heal(otherRat Rat) {
	r.state = healing
	r.SetAnimation(HEAL_ANIM)
	r.target = otherRat
	r.waitingTime = 0
	r.updateDestinationToTarget()
	r.vx = 0
	targetColor := otherRat.GetColor()
	healingStr := ""
	if r.name != r.target.GetName() {
		healingStr = r.color.Tag() + r.name +
			colors.White.Tag() + " is healing " +
			targetColor.Tag() + r.target.GetName()
	} else {
		healingStr = r.color.Tag() + r.name +
			colors.White.Tag() + " is healing himself"
	}

	r.ui.SetStatusMessage(healingStr, colors.White)
}

func (r *ratImpl) updateDestinationToTarget() {
	r.vx = RUN_SPEED
	r.destination = r.target.GetX()
	if r.destination < r.x {
		r.facing = left
	} else {
		r.facing = right
	}
	if r.facing == left {
		r.vx = -r.vx
	}
}

func (r ratImpl) GetX() float64 {
	return r.x
}

func (r ratImpl) IsAlive() bool {
	return r.health > 0
}

func (r *ratImpl) ReSpawn(color colors.CustomColor) {
	r.nameLabel.SetColor(color)
	r.color = color
	r.health = HEALTH_MAX
	r.x = 0
	r.RandomWalk()
	r.visible = true
}

func (r ratImpl) GetColor() colors.CustomColor {
	return r.color
}

func (r *ratImpl) Cure(heal int) (amount int, over int, crit bool) {
	amount, crit = r.Modify(heal)
	original := r.health
	r.health += amount
	if r.health >= HEALTH_MAX {
		r.health = HEALTH_MAX
	}
	over = heal - (r.health - original)
	return
}

func (r ratImpl) logHeal(amount, over int, crit bool) {
	targetColor := r.target.GetColor()
	var healStr = ""
	if r.name == r.target.GetName() {
		healStr = r.color.Tag() + r.name +
			colors.Yellow.Tag() + " heal himself"
	} else {
		healStr = r.color.Tag() + r.name +
			colors.Yellow.Tag() + " heal " +
			targetColor.Tag() + r.target.GetName()
	}

	healStr += colors.Yellow.Tag() + " by " + colors.Green.Tag()

	if crit {
		healStr += "*"
	}

	healStr += strconv.Itoa(amount)

	if crit {
		healStr += "* CRITICAL"
	}

	if over > 0 {
		healStr += colors.Yellow.Tag() + " (" +
			colors.Blue.Tag() + strconv.Itoa(over) +
			colors.Yellow.Tag() + " over heal)"
	}

	r.ui.SetStatusMessage(healStr, colors.Yellow)

	healStr = strconv.Itoa(amount)
	if crit {
		healStr = "*" + healStr + "*"
	}
	r.target.AddFlyingText(healStr, colors.Green)
}

func (r ratImpl) logDamage(amount, over int, crit bool) {
	targetColor := r.target.GetColor()
	var damageStr = ""

	damageStr = r.color.Tag() + r.name +
		colors.Yellow.Tag() + " hurt " +
		targetColor.Tag() + r.target.GetName()

	damageStr += colors.Yellow.Tag() + " by " + colors.Red.Tag()

	if crit {
		damageStr += "*"
	}

	damageStr += strconv.Itoa(amount)

	if crit {
		damageStr += "* CRITICAL"
	}

	if over > 0 {
		damageStr += colors.Yellow.Tag() + " (" +
			colors.Blue.Tag() + strconv.Itoa(over) +
			colors.Yellow.Tag() + " over kill)"
	}
	if !r.target.IsAlive() {
		damageStr += colors.Yellow.Tag() + " and" +
			colors.Red.Tag() + " kill " +
			colors.Yellow.Tag() + "it"
	}

	r.ui.SetStatusMessage(damageStr, colors.Yellow)

	damageStr = strconv.Itoa(amount)
	if crit {
		damageStr = "*" + damageStr + "*"
	}
	r.target.AddFlyingText(damageStr, colors.Red)
}

func (r ratImpl) Modify(points int) (value int, crit bool) {
	value = points
	randomAmount := rand.Float64() * MOD_VALUE
	value += int(randomAmount * float64(points))

	if rand.Float64() < CRIT_CHANGE {
		value *= 2
		crit = true
	}
	return
}

func (r ratImpl) AddFlyingText(text string, color colors.CustomColor) {
	r.ui.AddFlyingText(text, color, r.screenCenterX+r.x, r.flyingTextY)
}

func New(audioPlayer audio.Player, sheet draw.Sheet, ui ui.UI, font draw.Font, name string, ratColor colors.CustomColor) Rat {
	label := label.NewLabel(0, name, font, font.DefaultSize(), ratColor)
	labelWidth, _ := label.Measure()
	label.SetVisible(true)

	return &ratImpl{
		audioPlayer: audioPlayer,
		name:        name,
		x:           0,
		y:           0,
		vx:          0,
		sheet:       sheet,
		nameLabel:   label,
		labelWidth:  labelWidth,
		visible:     true,
		facing:      right,
		ui:          ui,
		color:       ratColor,
		health:      HEALTH_MAX,
	}
}
