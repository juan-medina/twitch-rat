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

package playing

import (
	"fmt"
	"math/rand"
	"strconv"
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
	"github.com/juan-medina/twitch-rat/internal/ui/slider"
)

const (
	GAME_MUSIC       = "embed/music/game.ogg"
	COUNTDOWN_SOUND  = "embed/sounds/%d.ogg"
	GO_SOUND         = "embed/sounds/go.ogg"
	TICK_SOUND       = "embed/sounds/tick.ogg"
	START_SOUND      = "embed/sounds/start.ogg"
	RAT_SPAWN_POINT  = -64
	MAX_RATS         = 40
	SPAWN_COMMAND    = "rat"
	ATTACK_COMMAND   = "attack"
	HEAL_COMMAND     = "heal"
	SECOND           = 1000
	COUNTDOWN_LENGTH = 30 * SECOND
	GO_VANISH        = 0.5 * SECOND
	TIME_TO_AUTO     = 30 * SECOND
)

func (p *playing) Init() {
	p.ui.SetButtonClickCallback(p.onButtonClick)
	p.ui.SetSliderChangeCallback(p.onSliderChange)

	p.ui.SetButtonVisible(ui.BACK_BUTTON, true)
	p.ui.SetButtonVisible(ui.IN_GAME_OPTIONS_BUTTON, true)

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
	p.ui.SetLabelVisible(ui.LABEL_COUNTDOWN, true)
	p.ui.SetLabelText(ui.LABEL_COUNTDOWN, strconv.FormatInt(COUNTDOWN_LENGTH/SECOND, 10))
	p.ui.SetLabelVisible(ui.LABEL_INSTRUCTIONS, true)
	p.herd = p.herd[:0]
	p.searchSlice = p.searchSlice[:0]
	p.status = connecting
	p.ui.SetScoreVisible(true)
}

func (p *playing) End() {
	p.ui.SetButtonVisible(ui.BACK_BUTTON, false)
	p.ui.SetButtonVisible(ui.IN_GAME_OPTIONS_BUTTON, false)
	p.chat.Disconnect()
	p.ui.SetStatusMessage("Disconnected", colors.White)
	p.audioPlayer.Stop()

	if debug := p.settings.GetBoolValue("debug", false); debug {
		p.ui.SetInputVisible(ui.INPUT_DEBUG_USER, false)
		p.ui.SetInputVisible(ui.INPUT_DEBUG_MESSAGE, false)
		p.ui.SetButtonVisible(ui.DEBUG_BUTTON, false)
	}
	p.ui.SetLabelVisible(ui.LABEL_COUNTDOWN, false)
	p.ui.SetLabelVisible(ui.LABEL_INSTRUCTIONS, false)
	p.ui.SetScoreVisible(false)

	p.ui.SetButtonVisible(ui.SUBMENU_OPTION_BACK_BUTTON, false)
	p.ui.SetSliderVisible(ui.MUSIC_VOLUME_SLIDER, false)
	p.ui.SetLabelVisible(ui.LABEL_OPTIONS_MUSIC_VOLUME, false)
	p.ui.SetSliderVisible(ui.AUDIO_VOLUME_SLIDER, false)
	p.ui.SetLabelVisible(ui.LABEL_OPTIONS_AUDIO_VOLUME, false)
	p.ui.SetPanelVisible(ui.OPTIONS_PANEL, false)
}

type status int

const (
	connecting status = iota
	counting
	endCountdown
	fighting
)

type playing struct {
	changer        stage.Changer
	ui             ui.UI
	settings       settings.Settings
	eventsChan     chan chat.Event
	chat           chat.Chat
	channel        string
	currentWidth   float64
	currentHeight  float64
	sewerMap       draw.Map
	audioPlayer    audio.Player
	rats           draw.Sheet
	herd           []rat.Rat
	searchSlice    []rat.Rat
	fontVerySmall  draw.Font
	fontSmall      draw.Font
	countdown      int
	status         status
	optionsVisible bool
}

func (p *playing) Update(elapsedTime int, keys keys.Keys) {
	if p.status == counting {
		currentCountdownSeconds := p.countdown / SECOND
		p.countdown -= elapsedTime
		countDownSeconds := p.countdown / SECOND
		if countDownSeconds != currentCountdownSeconds {
			if countDownSeconds == 0 {
				p.ui.SetLabelText(ui.LABEL_COUNTDOWN, "GO!")
				p.audioPlayer.PlaySound(GO_SOUND)
				p.audioPlayer.PlaySound(START_SOUND)
				p.ui.SetLabelVisible(ui.LABEL_INSTRUCTIONS, false)
				p.ui.SetStatusMessage("Game "+
					colors.Blue.Tag()+"started."+
					colors.Yellow.Tag()+" Attack any rat with"+
					colors.Red.Tag()+" !attack"+
					colors.Yellow.Tag()+" or heal yourself or any rat with"+
					colors.Green.Tag()+" !heal",
					colors.Yellow)
				p.ui.SetStatusMessage("If you"+
					colors.Red.Tag()+" die"+
					colors.Yellow.Tag()+" you can re-spawn using"+
					colors.Red.Tag()+" !rat"+
					colors.Yellow.Tag()+" again",
					colors.Yellow)
				p.status = endCountdown
				p.countdown = GO_VANISH
				p.ui.StartsCore()
			} else {
				p.ui.SetLabelText(ui.LABEL_COUNTDOWN, fmt.Sprintf("%d", countDownSeconds))
				if countDownSeconds >= 1 && countDownSeconds <= 10 {
					p.audioPlayer.PlaySound(fmt.Sprintf(COUNTDOWN_SOUND, countDownSeconds))
					p.audioPlayer.PlaySound(TICK_SOUND)
				} else {
					p.audioPlayer.PlaySound(TICK_SOUND)
				}
			}
		}
	} else if p.status == endCountdown {
		p.countdown -= elapsedTime
		if p.countdown <= 0 {
			p.ui.SetLabelVisible(ui.LABEL_COUNTDOWN, false)
			p.status = fighting
			p.ResetTimeSinceLastCommand()
		}
	}

	p.ui.Update(elapsedTime)
	select {
	case event := <-p.eventsChan:
		switch event.Type_ {
		case chat.Connect:
			p.ui.SetStatusMessage("Connected to "+colors.Green.Tag()+p.channel, colors.White)
			p.countdown = COUNTDOWN_LENGTH
			p.status = counting
			p.ui.SetStatusMessage("Game starting in"+
				colors.Green.Tag()+" "+strconv.FormatInt(COUNTDOWN_LENGTH/SECOND, 10)+
				colors.Yellow.Tag()+" seconds."+
				colors.Yellow.Tag()+" Join at any time with"+
				colors.Red.Tag()+" !rat",
				colors.Yellow)
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
			if p.optionsVisible {
				p.ui.ClickButton(ui.SUBMENU_OPTION_BACK_BUTTON)
			} else {
				p.ui.ClickButton(ui.BACK_BUTTON)
			}
		}
	}

	for i := range p.herd {
		p.herd[i].Update(elapsedTime)
		if p.status == fighting {
			timeSinceLastCommand := p.herd[i].TimeSinceLastCommand()
			if p.herd[i].CanDoAction() {
				origin := p.herd[i]
				if timeSinceLastCommand > TIME_TO_AUTO {
					var target rat.Rat
					if rand.Intn(100) < 30 {
						if target = p.findTargetToHeal(); target != nil {
							p.herd[i].Heal(target)
						} else {
							if target = p.findTargetToAttack(origin); target != nil {
								p.herd[i].Attack(target)
							}
						}
					} else {
						if target = p.findTargetToAttack(origin); target != nil {
							p.herd[i].Attack(target)
						}
					}

				}
			}
		}
	}
}

func (p *playing) ResetTimeSinceLastCommand() {
	for i := range p.herd {
		p.herd[i].ResetTimeSinceLastCommand()
	}
}
func (p playing) findTargetToAttack(origin rat.Rat) rat.Rat {
	p.searchSlice = p.searchSlice[:0]
	for i := range p.herd {
		if origin != p.herd[i] {
			if p.herd[i].IsVisible() && p.herd[i].IsAlive() {
				p.searchSlice = append(p.searchSlice, p.herd[i])
			}
		}
	}

	if len(p.searchSlice) == 0 {
		return nil
	}

	r := rand.Intn(len(p.searchSlice))
	return p.searchSlice[r]
}

func (p playing) findTargetToHeal() rat.Rat {
	p.searchSlice = p.searchSlice[:0]
	for i := range p.herd {
		if p.herd[i].IsAlive() && !p.herd[i].IsFullHealth() {
			p.searchSlice = append(p.searchSlice, p.herd[i])
		}
	}

	if len(p.searchSlice) == 0 {
		return nil
	}

	r := rand.Intn(len(p.searchSlice))
	return p.searchSlice[r]
}

var lastRandomName = 0

func (p playing) getRandomRatName() (name string) {

	r := rand.Intn(100)
	if r < 50 {
		name = "rat"
	} else {
		name = "mouse"
	}

	lastRandomName += rand.Intn(10) + 1

	name = fmt.Sprintf("%s%d", name, lastRandomName)
	return
}

func (p *playing) processCommand(message string, user string, userColor colors.CustomColor) {
	if userColor == colors.Black {
		userColor = labelColors[nextLabelColor]
		nextLabelColor = (nextLabelColor + 1) % len(labelColors)
	}

	if user == "" {
		user = p.getRandomRatName()
	}

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
			rat := rat.New(p.audioPlayer, p.rats, p.ui, p.fontVerySmall, p.fontSmall, user, userColor)
			rat.RandomWalk()
			rat.SetX(RAT_SPAWN_POINT)
			rat.SetCenter((p.currentWidth / 2))
			p.herd = append(p.herd, rat)
			p.ui.SetStatusMessage(fmt.Sprintf("%s%s%s join the fight!", userColor.Tag(), user, colors.White.Tag()), colors.White)
			p.ui.AddScoreEntry(rat)
		} else {
			if !findRat.IsAlive() {
				findRat.SetX(RAT_SPAWN_POINT)
				findRat.ReSpawn(userColor)
				findRat.SetCenter((p.currentWidth / 2))
				p.ui.SetStatusMessage(fmt.Sprintf("%s%s%s rejoin", userColor.Tag(), user, colors.White.Tag()), colors.White)
			}
		}
	} else if command == ATTACK_COMMAND && p.status == fighting {
		if userRat := p.getRat(user); userRat != nil {
			if userRat.CanDoAction() {
				if targetRat := p.getRat(args); targetRat != nil {
					if targetRat.IsAlive() {
						targetName := targetRat.GetName()
						userName := userRat.GetName()
						if targetName != userName {
							userRat.Attack(targetRat)
						}
					}
				}
			}
		}
	} else if command == HEAL_COMMAND && p.status == fighting {
		if userRat := p.getRat(user); userRat != nil {
			if userRat.CanDoAction() {
				targetRat := userRat
				if args != "" {
					targetRat = p.getRat(args)
				}
				if targetRat != nil {
					if targetRat.IsAlive() {
						userRat.Heal(targetRat)
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
	labelColors []colors.CustomColor = []colors.CustomColor{
		colors.New(0xFF, 0x00, 0x00, 0xFF), // Red
		colors.New(0x00, 0xFF, 0x00, 0xFF), // Green
		colors.New(0xFF, 0xFF, 0x00, 0xFF), // Yellow
		colors.New(0xFF, 0x00, 0xFF, 0xFF), // Magenta
		colors.New(0x00, 0xFF, 0xFF, 0xFF), // Cyan
		colors.New(0xFF, 0x7F, 0x00, 0xFF), // Gold
		colors.New(0x7F, 0x00, 0xFF, 0xFF), // Purple
		colors.New(0xFF, 0x00, 0x7F, 0xFF), // Orange
		colors.New(0x00, 0x80, 0x00, 0xFF), // Dark Green
		colors.New(0x80, 0x00, 0x80, 0xFF), // Maroon
		colors.New(0x00, 0x80, 0x80, 0xFF), // Dark Red
		colors.New(0x00, 0xFF, 0x80, 0xFF), // Dark Yellow
		colors.New(0xFF, 0x00, 0x80, 0xFF), // Dark Magenta
		colors.New(0x00, 0x80, 0xFF, 0xFF), // Dark Cyan
		colors.New(0x80, 0x00, 0xFF, 0xFF), // Dark Blue
		colors.New(0x80, 0x80, 0x00, 0xFF), // Dark Brown
		colors.New(0x80, 0x80, 0x80, 0xFF), // Dark Gray
		colors.New(0x80, 0xC0, 0x80, 0xFF), // Dark Olive
		colors.New(0x00, 0x00, 0x80, 0xFF), // Dark Navy
		colors.New(0x00, 0x00, 0xFF, 0xFF), // Black
	}
	nextLabelColor = rand.Intn(len(labelColors))
)

func (p *playing) onButtonClick(id button.Id) {
	switch id {
	case ui.BACK_BUTTON:
		p.ui.SetButtonClickCallback(nil)
		p.ui.SetSliderChangeCallback(nil)
		if debug := p.settings.GetBoolValue("debug", false); !debug {
			p.changer.ChangeStage(stage.MENU)
		} else {
			p.changer.ChangeStage(stage.EXIT)
		}
	case ui.DEBUG_BUTTON:
		user := p.ui.GetInputText(ui.INPUT_DEBUG_USER)
		message := p.ui.GetInputText(ui.INPUT_DEBUG_MESSAGE)
		p.onChatEvent(chat.Event{Type_: chat.Message, Sender: user, Message: message, UserColor: colors.Black})
	case ui.IN_GAME_OPTIONS_BUTTON:
		p.optionsMenu(true)
	case ui.SUBMENU_OPTION_BACK_BUTTON:
		p.optionsMenu(false)
	}
}

func (p *playing) optionsMenu(enable bool) {
	p.optionsVisible = enable
	p.ui.SetButtonVisible(ui.SUBMENU_OPTION_BACK_BUTTON, p.optionsVisible)
	p.ui.SetSliderVisible(ui.MUSIC_VOLUME_SLIDER, p.optionsVisible)
	p.ui.SetLabelVisible(ui.LABEL_OPTIONS_MUSIC_VOLUME, p.optionsVisible)
	p.ui.SetSliderVisible(ui.AUDIO_VOLUME_SLIDER, p.optionsVisible)
	p.ui.SetLabelVisible(ui.LABEL_OPTIONS_AUDIO_VOLUME, p.optionsVisible)
	p.ui.SetPanelVisible(ui.OPTIONS_PANEL, p.optionsVisible)
	song := p.settings.GetFloatValue("song_volume", 0.2)
	p.ui.SetSliderValue(ui.MUSIC_VOLUME_SLIDER, song)
	sound := p.settings.GetFloatValue("sound_volume", 0.5)
	p.ui.SetSliderValue(ui.AUDIO_VOLUME_SLIDER, sound)
	if p.optionsVisible {
		p.ui.DisableButton(ui.BACK_BUTTON)
		p.ui.DisableButton(ui.IN_GAME_OPTIONS_BUTTON)
		p.ui.SetLabelVisible(ui.LABEL_INSTRUCTIONS, false)
		p.ui.SetLabelVisible(ui.LABEL_COUNTDOWN, false)
	} else {
		p.ui.EnableButton(ui.BACK_BUTTON)
		p.ui.EnableButton(ui.IN_GAME_OPTIONS_BUTTON)
		if p.status != fighting {
			p.ui.SetLabelVisible(ui.LABEL_INSTRUCTIONS, true)
			p.ui.SetLabelVisible(ui.LABEL_COUNTDOWN, true)
		}
	}
}
func (p *playing) onChatEvent(e chat.Event) {
	p.eventsChan <- e
}

func (p *playing) onSliderChange(id slider.Id, value float64) {
	switch id {
	case ui.MUSIC_VOLUME_SLIDER:
		p.settings.SetFloatValue("song_volume", value)
		p.settings.Save()
		p.audioPlayer.ChangeSongVolume(value)
	case ui.AUDIO_VOLUME_SLIDER:
		p.settings.SetFloatValue("sound_volume", value)
		p.settings.Save()
		p.audioPlayer.ChangeSoundVolume(value)
	}
}

func New(changer stage.Changer, ui ui.UI, settings settings.Settings, sewerMap draw.Map, rats draw.Sheet, audioPlayer audio.Player, fontVerySmall draw.Font, fontSmall draw.Font) stage.Stage {
	audioPlayer.LoadSong(GAME_MUSIC)
	audioPlayer.LoadSound(rat.HIT_SOUND)
	audioPlayer.LoadSound(rat.DEAD_SOUND)
	audioPlayer.LoadSound(rat.HEAL_SOUND)
	audioPlayer.LoadSound(GO_SOUND)
	audioPlayer.LoadSound(TICK_SOUND)
	audioPlayer.LoadSound(START_SOUND)
	for i := 1; i <= 10; i++ {
		audioPlayer.LoadSound(fmt.Sprintf(COUNTDOWN_SOUND, i))
	}
	var chatInterface chat.Chat
	if debug := settings.GetBoolValue("debug", false); !debug {
		chatInterface = chat.New()
	} else {
		chatInterface = chat.NewDebug()
	}
	p := playing{
		changer:       changer,
		settings:      settings,
		ui:            ui,
		eventsChan:    make(chan chat.Event, 10),
		chat:          chatInterface,
		sewerMap:      sewerMap,
		audioPlayer:   audioPlayer,
		rats:          rats,
		herd:          make([]rat.Rat, 0, MAX_RATS),
		searchSlice:   make([]rat.Rat, 0, MAX_RATS),
		fontVerySmall: fontVerySmall,
		fontSmall:     fontSmall,
	}
	return &p
}
