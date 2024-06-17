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

package playing

import (
	"fmt"
	"image/color"
	"math/rand"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/juan-medina/twitch-rat/internal/audio"
	"github.com/juan-medina/twitch-rat/internal/chat"
	"github.com/juan-medina/twitch-rat/internal/colors"
	"github.com/juan-medina/twitch-rat/internal/draw"
	"github.com/juan-medina/twitch-rat/internal/keys"
	"github.com/juan-medina/twitch-rat/internal/settings"
	"github.com/juan-medina/twitch-rat/internal/stage"
	"github.com/juan-medina/twitch-rat/internal/stage/playing/rat"
	"github.com/juan-medina/twitch-rat/internal/ui"
	"github.com/juan-medina/twitch-rat/internal/ui/button"
)

const (
	GAME_MUSIC      = "embed/music/game.ogg"
	RAT_SPAWN_POINT = -64
	MAX_RATS        = 2
	SPAWN_COMMAND   = "rat"
	ATTACK_COMMAND  = "attack"
)

func (p *playing) Init() {
	p.ui.SetButtonClickCallback(p.onButtonClick)

	p.ui.SetButtonVisible(ui.BACK_BUTTON, true)

	p.ui.SetStatusMessage("Connecting..", colors.White)

	p.channel = p.settings.GetValue("channel", "")
	p.chat.OnEvent(p.onChatEvent)
	p.chat.Connect(p.channel)
	p.sewerMap.SetLevel(0)
	p.audioPlayer.PlaySong(GAME_MUSIC)

	if debug := p.settings.GetBoolValue("debug", false); debug {
		p.ui.SetInputVisible(ui.INPUT_DEBUG_USER, true)
		p.ui.SetInputVisible(ui.INPUT_DEBUG_MESSAGE, true)
		p.ui.SetButtonVisible(ui.DEBUG_BUTTON, true)
	}
	p.herd = p.herd[:0]
}

func (p *playing) End() {
	p.ui.SetButtonVisible(ui.BACK_BUTTON, false)
	p.chat.Disconnect()
	p.ui.SetStatusMessage("Disconnected", colors.White)
	p.audioPlayer.Stop()

	if debug := p.settings.GetBoolValue("debug", false); debug {
		p.ui.SetInputVisible(ui.INPUT_DEBUG_USER, false)
		p.ui.SetInputVisible(ui.INPUT_DEBUG_MESSAGE, false)
		p.ui.SetButtonVisible(ui.DEBUG_BUTTON, false)
	}
}

type playing struct {
	changer       stage.Changer
	ui            ui.UI
	settings      settings.Settings
	eventsChan    chan chat.Event
	chat          chat.Chat
	channel       string
	currentWidth  float64
	currentHeight float64
	sewerMap      draw.Map
	audioPlayer   audio.Player
	rats          draw.Sheet
	herd          []rat.Rat
	font          draw.Font
}

func (p *playing) Update(elapsedTime int, keys keys.Keys) {
	p.ui.Update(elapsedTime)
	select {
	case event := <-p.eventsChan:
		switch event.Type_ {
		case chat.Connect:
			p.ui.SetStatusMessage("Connected to "+p.channel, colors.White)
		case chat.Disconnect:
			p.ui.SetStatusMessage("Disconnected", colors.White)
		case chat.Message:
			if event.Message != "" {
				if event.Message[0] == '!' {
					p.processCommand(event.Message, event.Sender, event.UserColor)
				}
			}
		}
	default:
		/// no new event
	}

	if !p.ui.IsInputEditing(ui.INPUT_CHANNEL) {
		if keys.IsDownNoRepeat(ebiten.KeyEscape) {
			p.ui.ClickButton(ui.BACK_BUTTON)
		}
	}
	for i := range p.herd {
		p.herd[i].Update(elapsedTime)
	}
}

func (p *playing) processCommand(message string, user string, userColor color.Color) {
	lenStr := len(message)

	firstSpace := strings.Index(message, " ")
	var endCommand int
	if firstSpace == -1 {
		endCommand = lenStr
	} else {
		endCommand = firstSpace
	}
	command := strings.ToLower(message[1:endCommand])

	args := ""
	if firstSpace != -1 && lenStr > firstSpace+1 {
		args = message[firstSpace+1:]
	}

	if command == SPAWN_COMMAND {
		findRat := p.getRat(user)
		if findRat == nil {
			rat := rat.New(p.audioPlayer, p.rats, p.ui, p.font, user, userColor)
			rat.RandomWalk()
			rat.SetX(RAT_SPAWN_POINT)
			rat.SetCenter((p.currentWidth / 2))
			p.herd = append(p.herd, rat)
			p.ui.SetStatusMessage(fmt.Sprintf("%s spawned", user), colors.White)
		} else {
			if !findRat.IsAlive() {
				findRat.SetX(RAT_SPAWN_POINT)
				findRat.ReSpawn(userColor)
				findRat.SetCenter((p.currentWidth / 2))
				p.ui.SetStatusMessage(fmt.Sprintf("%s re-spawned", user), colors.White)
			}
		}
	} else if command == ATTACK_COMMAND {
		if userRat := p.getRat(user); userRat != nil {
			if userRat.CanAttack() {
				if targetRat := p.getRat(args); targetRat != nil {
					targetName := targetRat.GetName()
					userName := userRat.GetName()
					if targetName != userName {
						userRat.Attack(targetRat)
					}
				}
			}
		}
	}
}
func (p playing) getRat(name string) rat.Rat {
	if name == "" {
		return nil
	}
	lower := strings.ToLower(name)
	if lower[0] == '@' {
		lower = lower[1:]
	}
	for i := range p.herd {
		ratName := strings.ToLower(p.herd[i].GetName())
		if ratName == lower {
			return p.herd[i]
		}
	}
	return nil
}

func (p *playing) Draw(screen *ebiten.Image) {
	p.sewerMap.Draw(screen)
	for i := range p.herd {
		p.herd[i].Draw(screen)
	}
	p.ui.Draw(screen)
}

func (p *playing) OnLayoutChange(width, height float64) {
	p.currentWidth = width
	p.currentHeight = height

	w, _ := p.sewerMap.Size()
	dx := (p.currentWidth - w) / 2
	p.sewerMap.Move(dx, 0)
	for i := range p.herd {
		p.herd[i].SetCenter((p.currentWidth / 2))
	}
}

var (
	labelColors []color.RGBA = []color.RGBA{
		{0xFF, 0x00, 0x00, 0xFF}, // Red
		{0x00, 0xFF, 0x00, 0xFF}, // Green
		{0xFF, 0xFF, 0x00, 0xFF}, // Yellow
		{0xFF, 0x00, 0xFF, 0xFF}, // Magenta
		{0x00, 0xFF, 0xFF, 0xFF}, // Cyan
		{0xFF, 0x7F, 0x00, 0xFF}, // Gold
		{0x7F, 0x00, 0xFF, 0xFF}, // Purple
		{0xFF, 0x00, 0x7F, 0xFF}, // Orange
		{0x00, 0x80, 0x00, 0xFF}, // Dark Green
		{0x80, 0x00, 0x80, 0xFF}, // Maroon
		{0x00, 0x80, 0x80, 0xFF}, // Dark Red
		{0x00, 0xFF, 0x80, 0xFF}, // Dark Yellow
		{0xFF, 0x00, 0x80, 0xFF}, // Dark Magenta
		{0x00, 0x80, 0xFF, 0xFF}, // Dark Cyan
		{0x80, 0x00, 0xFF, 0xFF}, // Dark Blue
		{0x80, 0x80, 0x00, 0xFF}, // Dark Brown
		{0x80, 0x80, 0x80, 0xFF}, // Dark Gray
		{0x80, 0xC0, 0x80, 0xFF}, // Dark Olive
		{0x00, 0x00, 0x80, 0xFF}, // Dark Navy
		{0x00, 0x00, 0xFF, 0xFF}, // Black
	}
	nextLabelColor = rand.Intn(len(labelColors))
)

func (p *playing) onButtonClick(id button.Id) {
	switch id {
	case ui.BACK_BUTTON:
		p.ui.SetButtonClickCallback(nil)
		if debug := p.settings.GetBoolValue("debug", false); !debug {
			p.changer.ChangeStage(stage.MENU)
		} else {
			p.changer.ChangeStage(stage.EXIT)
		}
	case ui.DEBUG_BUTTON:
		user := p.ui.GetInputText(ui.INPUT_DEBUG_USER)
		message := p.ui.GetInputText(ui.INPUT_DEBUG_MESSAGE)
		userColor := labelColors[nextLabelColor]
		nextLabelColor = (nextLabelColor + 1) % len(labelColors)
		p.onChatEvent(chat.Event{Type_: chat.Message, Sender: user, Message: message, UserColor: userColor})
	}
}
func (p *playing) onChatEvent(e chat.Event) {
	p.eventsChan <- e
}

func New(changer stage.Changer, ui ui.UI, settings settings.Settings, sewerMap draw.Map, rats draw.Sheet, audioPlayer audio.Player, font draw.Font) stage.Stage {
	audioPlayer.LoadSong(GAME_MUSIC)
	audioPlayer.LoadSound(rat.HIT_SOUND)
	audioPlayer.LoadSound(rat.DEAD_SOUND)
	var chatInterface chat.Chat
	if debug := settings.GetBoolValue("debug", false); !debug {
		chatInterface = chat.New()
	} else {
		chatInterface = chat.NewDebug()
	}
	p := playing{
		changer:     changer,
		settings:    settings,
		ui:          ui,
		eventsChan:  make(chan chat.Event, 10),
		chat:        chatInterface,
		sewerMap:    sewerMap,
		audioPlayer: audioPlayer,
		rats:        rats,
		herd:        make([]rat.Rat, 0, MAX_RATS),
		font:        font,
	}
	return &p
}
