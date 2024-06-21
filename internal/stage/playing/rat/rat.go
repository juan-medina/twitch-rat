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
	HEALTH_BAR_WIDTH       = 110
	HEALTH_BAR_HEIGHT      = 20
	CRIT_CHANGE            = 0.25
	MOD_VALUE              = 0.25
	HEALTH_MAX             = 50
	RAT_DAMAGE             = 15
	RAT_HEAL               = 15
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
	IsFullHealth() bool
	Cure(heal int) (amount int, over int, crit bool)
	Modify(points int) (value int, crit bool)
	AddFlyingText(text string, color colors.CustomColor)
	TimeSinceLastCommand() int
	ResetTimeSinceLastCommand()
	IsAttacking() bool
	GetHealth() int
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
	x, y, screenCenterX  float64
	vx                   float64
	sheet                draw.Sheet
	sprite               draw.Sprite
	animationStatus      animationStatus
	name                 string
	nameLabel            label.Label
	hpLabel              label.Label
	labelWidth           float64
	hpLabelWidth         float64
	hpLabelHeight        float64
	visible              bool
	facing               direction
	state                state
	destination          float64
	waitingTime          int
	target               Rat
	ui                   ui.UI
	audioPlayer          audio.Player
	color                colors.CustomColor
	barX                 float64
	barY                 float64
	redStart             float64
	redWidth             float64
	health               int
	flyingTextY          float64
	timeSinceLastCommand int
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

	vector.DrawFilledRect(screen, float32(r.barX), float32(r.barY), HEALTH_BAR_WIDTH, HEALTH_BAR_HEIGHT, colors.Green, false)
	vector.DrawFilledRect(screen, float32(r.redStart), float32(r.barY), float32(r.redWidth), HEALTH_BAR_HEIGHT, colors.Red, false)

	vector.StrokeRect(screen, float32(r.barX), float32(r.barY), HEALTH_BAR_WIDTH, HEALTH_BAR_HEIGHT, 3, colors.Black, false)
	r.hpLabel.Draw(screen)
}

func (r *ratImpl) calculateRedHPBar() {
	greenWidth := float64(HEALTH_BAR_WIDTH) * (float64(r.health) / HEALTH_MAX)
	r.redWidth = HEALTH_BAR_WIDTH - greenWidth
	r.redStart = r.barX + greenWidth
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

	r.hpLabel.Move(r.screenCenterX+r.x-(r.hpLabelWidth/2), r.barY+(HEALTH_BAR_HEIGHT/2)-(r.hpLabelHeight/2))

	r.flyingTextY = r.y - (rh / 2)
	r.calculateRedHPBar()
}

func (r *ratImpl) Update(elapsedTime int) {
	if !r.visible {
		return
	}

	r.timeSinceLastCommand += elapsedTime

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
					targetWasAttacking := r.target.IsAttacking()
					damage, over, crit := r.target.Hurt(RAT_DAMAGE)
					r.logDamage(damage, over, crit)
					if targetWasAttacking {
						r.ui.SetStatusMessage(r.target.GetColor().BBCoded(r.target.GetName())+" was interrupted", colors.Yellow)
					}
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
		if r.target != nil {
			if r.target.IsAlive() {
				r.updateDestinationToTarget()
				r.vx = 0
				if r.animationStatus.end {
					heal, over, crit := r.target.Cure(RAT_HEAL)
					r.logHeal(heal, over, crit)
					r.target = nil
					r.idle()
					r.audioPlayer.PlaySound(HEAL_SOUND)
				}
			} else {
				r.target = nil
				r.idle()
			}
		} else {
			r.idle()
			return
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
		if r.target != nil || r.state == running || r.state == healing {
			r.target = nil
		}
	}

	r.updateHP()
	r.vx = 0
	over = amount - (original - r.health)
	return
}

func (r *ratImpl) updateHP() {
	r.hpLabel.SetText(strconv.Itoa(r.health))
	r.hpLabelWidth, r.hpLabelHeight = r.hpLabel.Measure()
	r.calculateRedHPBar()
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
	r.timeSinceLastCommand = 0

	r.state = running
	r.SetAnimation(RUN_ANIM)
	r.target = otherRat
	r.waitingTime = 0
	r.updateDestinationToTarget()
	targetColor := otherRat.GetColor()
	damageStr := r.color.BBCoded(r.name) + " is attacking " +
		targetColor.BBCoded(r.target.GetName())
	r.ui.SetStatusMessage(damageStr, colors.White)
}

func (r *ratImpl) Heal(otherRat Rat) {
	r.timeSinceLastCommand = 0

	r.state = healing
	r.SetAnimation(HEAL_ANIM)
	r.target = otherRat
	r.waitingTime = 0
	r.updateDestinationToTarget()
	r.vx = 0

	healingStr := ""
	targetColor := otherRat.GetColor()
	if r.name != r.target.GetName() {
		healingStr = r.color.BBCoded(r.name) + " is healing " +
			targetColor.BBCoded(r.target.GetName())
	} else {
		healingStr = r.color.BBCoded(r.name) + " is healing himself"
	}

	r.ui.SetStatusMessage(healingStr, colors.White)
}

func (r ratImpl) IsFullHealth() bool {
	return r.health == HEALTH_MAX
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
	r.timeSinceLastCommand = 0
	r.nameLabel.SetColor(color)
	r.color = color
	r.health = HEALTH_MAX
	r.updateHP()
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
	r.hpLabel.SetText(strconv.Itoa(r.health))
	over = heal - (r.health - original)
	return
}

func (r ratImpl) logHeal(amount, over int, crit bool) {
	targetColor := r.target.GetColor()
	var healStr = ""
	if r.name == r.target.GetName() {
		healStr = r.color.BBCoded(r.name) + " heal himself"
	} else {
		healStr = r.color.BBCoded(r.name) + " heal " +
			targetColor.BBCoded(r.target.GetName())
	}

	strAmount := ""
	if crit {
		strAmount += "*"
	}
	strAmount += strconv.Itoa(amount)
	if crit {
		strAmount += "* CRITICAL"
	}

	healStr += " by " + colors.Green.BBCoded(strAmount)

	if over > 0 {
		healStr += " " +
			colors.Blue.BBCoded("("+strconv.Itoa(over)+" over heal)")
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

	damageStr = r.color.BBCoded(r.name) + " damage " +
		targetColor.BBCoded(r.target.GetName())

	strAmount := ""
	if crit {
		strAmount += "*"
	}
	strAmount += strconv.Itoa(amount)
	if crit {
		strAmount += "* CRITICAL"
	}

	damageStr += " by " + colors.Red.BBCoded(strAmount)

	if over > 0 {
		damageStr += " " +
			colors.Blue.BBCoded("("+strconv.Itoa(over)+" over kill)")
	} else {
		if !r.target.IsAlive() {
			damageStr += " " + colors.Blue.BBCoded("(killed)")
		}
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

func (r ratImpl) TimeSinceLastCommand() int {
	return r.timeSinceLastCommand
}

func (r *ratImpl) ResetTimeSinceLastCommand() {
	r.timeSinceLastCommand = 0
}

func (r *ratImpl) IsAttacking() bool {
	return (r.state == running || r.state == healing) && r.target != nil
}
func (r ratImpl) GetHealth() int {
	return r.health
}

func New(audioPlayer audio.Player, sheet draw.Sheet, ui ui.UI, fontVerySmall draw.Font, fontSmall draw.Font, name string, ratColor colors.CustomColor) Rat {
	nameLabel := label.NewLabel(0, name, fontSmall, fontSmall.DefaultSize(), ratColor, nil)
	labelWidth, _ := nameLabel.Measure()
	nameLabel.SetVisible(true)

	hpLabel := label.NewLabel(0, strconv.Itoa(HEALTH_MAX), fontVerySmall, fontVerySmall.DefaultSize(), colors.White, nil)
	hpLabelWidth, hpLabelHeight := hpLabel.Measure()
	hpLabel.SetVisible(true)

	return &ratImpl{
		audioPlayer:   audioPlayer,
		name:          name,
		x:             0,
		y:             0,
		vx:            0,
		sheet:         sheet,
		nameLabel:     nameLabel,
		hpLabel:       hpLabel,
		labelWidth:    labelWidth,
		hpLabelWidth:  hpLabelWidth,
		hpLabelHeight: hpLabelHeight,
		visible:       true,
		facing:        right,
		ui:            ui,
		color:         ratColor,
		health:        HEALTH_MAX,
	}
}
