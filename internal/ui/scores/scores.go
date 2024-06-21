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

package scores

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/juan-medina/twitch-rat/internal/colors"
	"github.com/juan-medina/twitch-rat/internal/draw"
	"github.com/juan-medina/twitch-rat/internal/keys"
)

const (
	HEAD_SPRITE        = "rat_head_small"
	ENTRIES_GAP_Y      = 8
	ENTRIES_GAP_X      = 24
	MAX_ENTRIES        = 40
	SCORE_WIDTH        = 400
	SECOND             = 1000
	COUNTDOWN_LENGTH   = 11 * SECOND
	MAX_SCORES_DISPLAY = 24
)

type ScoreData interface {
	GetColor() colors.CustomColor
	GetName() string
	GetHealth() int
}

type Scores interface {
	Init()

	Update(elapsedTime int, mouseX, mouseY float64, leftJustPressed bool, leftPressed bool, keys keys.Keys)
	Draw(screen *ebiten.Image)
	SetVisible(visible bool)

	Move(x, y float64)

	Add(data ScoreData)
	Reset()

	Start()
}

type scoreEntry struct {
	data         ScoreData
	score        int
	pointX       float64
	text         string
	addText      string
	addX         float64
	delayAddTime int
	delayValue   int
}
type scoresImpl struct {
	visible       bool
	entries       []scoreEntry
	x, y          float64
	fontSmall     draw.Font
	fontNormal    draw.Font
	scoreWidth    float64
	timeX         float64
	timeY         float64
	started       bool
	countdown     int
	countdownText string
	limitText     string
}

func (s *scoresImpl) Init() {
}

func (s *scoresImpl) Update(elapsedTime int, mouseX, mouseY float64, leftJustPressed bool, leftPressed bool, keys keys.Keys) {
	if !s.visible {
		return
	}

	secondsLogic := -1
	calculateScore := false
	if s.started {

		currentCountdownSeconds := s.countdown / SECOND
		s.countdown -= elapsedTime
		countDownSeconds := s.countdown / SECOND

		if currentCountdownSeconds != countDownSeconds {
			secondsLogic = int(countDownSeconds)
			s.countdownText = fmt.Sprintf("%02d/10s", secondsLogic)
		}
		if countDownSeconds <= 0 {
			s.countdown = COUNTDOWN_LENGTH
			s.countdownText = "10/10s"
			calculateScore = true
		}
	}

	maxName := "ZZZZZZZZZZ"
	maxNameLength := len(maxName)

	maxScoreLength := 0
	maxScore := ""

	needSort := false
	for i := range s.entries {
		if calculateScore {
			health := s.entries[i].data.GetHealth()
			toAdd := 5 * health
			if toAdd > 0 {
				s.entries[i].delayValue = toAdd
				s.entries[i].delayAddTime = 2 * SECOND
				s.entries[i].addText = fmt.Sprintf("+ %d", toAdd)
			}
		} else if s.entries[i].delayAddTime > 0 {
			s.entries[i].delayAddTime -= elapsedTime
			if s.entries[i].delayAddTime <= 0 {
				s.entries[i].score += s.entries[i].delayValue
				s.entries[i].delayValue = 0
				s.entries[i].delayAddTime = 0
				s.entries[i].addText = ""
				needSort = true
			}
		}

		s.entries[i].text = strconv.Itoa(s.entries[i].score)
		name := s.entries[i].data.GetName()
		if len(name) > maxNameLength {
			maxNameLength = len(name)
			maxName = name
		}
		if len(s.entries[i].text) > maxScoreLength {
			maxScoreLength = len(s.entries[i].text)
			maxScore = s.entries[i].text
		}
	}

	if needSort {
		sort.Slice(s.entries, func(i, j int) bool {
			return s.entries[i].score > s.entries[j].score
		})
	}

	mameWidth, _ := s.fontSmall.Measure(maxName, s.fontSmall.DefaultSize())
	scoreWidth, _ := s.fontSmall.Measure(maxScore, s.fontSmall.DefaultSize())
	for i := range s.entries {
		entryScoreWidth, _ := s.fontSmall.Measure(s.entries[i].text, s.fontSmall.DefaultSize())
		s.entries[i].pointX = s.x + (mameWidth + ENTRIES_GAP_X) + (scoreWidth - entryScoreWidth)
		s.entries[i].addX = s.entries[i].pointX + entryScoreWidth + ENTRIES_GAP_X/2
	}

	s.scoreWidth = mameWidth + ENTRIES_GAP_X + scoreWidth

	timeWidth, _ := s.fontSmall.Measure(s.countdownText, s.fontSmall.DefaultSize())
	s.timeX = s.x + s.scoreWidth - timeWidth
	s.timeY = s.y + s.fontNormal.DefaultSize()
}

func (s scoresImpl) Draw(screen *ebiten.Image) {
	if !s.visible {
		return
	}

	startX := s.x
	startY := s.y

	scoreText := "Score"

	if !s.started {
		scoreText = "Rats"
	}

	s.fontNormal.Draw(screen, scoreText, startX-(s.x/2), startY, s.fontNormal.DefaultSize(), colors.White)

	startY += s.fontNormal.DefaultSize() + ENTRIES_GAP_Y

	vector.StrokeLine(screen, 0, float32(startY), float32(startX+s.scoreWidth+(s.x/2)), float32(startY), 2, colors.White, false)

	startY += ENTRIES_GAP_Y * 2

	textHeight := s.fontSmall.DefaultSize()

	for i, e := range s.entries {
		if i >= MAX_SCORES_DISPLAY {
			break
		}
		s.fontSmall.Draw(screen, e.data.GetName(), startX, startY, textHeight, e.data.GetColor())

		if s.started {
			s.fontSmall.Draw(screen, e.text, e.pointX, startY, textHeight, colors.White)
			s.fontSmall.Draw(screen, e.addText, e.addX, startY, textHeight, colors.Green)
		}
		startY += textHeight + ENTRIES_GAP_Y
	}

	if len(s.entries) > 0 {
		vector.StrokeLine(screen, 0, float32(startY), float32(s.x+s.scoreWidth+(s.x/2)), float32(startY), 2, colors.White, false)
		s.fontSmall.Draw(screen, s.limitText, 6, startY-6, textHeight, colors.White)
	}

	if s.started {
		s.fontSmall.Draw(screen, s.countdownText, s.timeX, s.timeY, textHeight, colors.White)
	}
}

func (s *scoresImpl) SetVisible(visible bool) {
	s.visible = visible
}

func (s *scoresImpl) Add(data ScoreData) {
	s.entries = append(s.entries, scoreEntry{data: data, score: 0})
	if len(s.entries) > MAX_SCORES_DISPLAY {
		s.limitText = fmt.Sprintf("%d/%d rats", MAX_SCORES_DISPLAY, len(s.entries))
	} else {
		s.limitText = fmt.Sprintf("%d rats", len(s.entries))
	}
}

func (s *scoresImpl) Move(x float64, y float64) {
	s.x = x
	s.y = y
}

func (s *scoresImpl) Reset() {
	s.entries = s.entries[:0]
	s.started = false
	s.limitText = ""
}

func (s *scoresImpl) Start() {
	s.started = true
	s.countdown = COUNTDOWN_LENGTH
	s.countdownText = "10/10s"
}

func New(fontSmall draw.Font, fontNormal draw.Font) Scores {
	return &scoresImpl{
		entries:    make([]scoreEntry, 0, MAX_ENTRIES),
		fontSmall:  fontSmall,
		fontNormal: fontNormal,
		limitText:  "",
	}
}
