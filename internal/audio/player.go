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

package audio

import (
	"bytes"
	"embed"

	ebitenAudio "github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
)

type Player interface {
	LoadSong(song string)
	PlaySong(song string)
	StopCurrentSong()
	ChangeSongVolume(volume float64)

	LoadSound(sound string)
	PlaySound(sound string)
	StopAllSounds()
	ChangeSoundVolume(volume float64)

	Stop()

	Update()
}

type playerImpl struct {
	audioContext *ebitenAudio.Context
	musicPlayer  *ebitenAudio.Player
	fileSystem   embed.FS
	songs        map[string]*vorbis.Stream
	sounds       map[string]*ebitenAudio.Player
	songVolume   float64
	soundVolume  float64
}

func (p *playerImpl) LoadSong(song string) {
	if _, ok := p.songs[song]; ok {
		return
	}
	if songData, err := p.fileSystem.ReadFile(song); err == nil {
		if stream, err := vorbis.DecodeWithoutResampling(bytes.NewReader(songData)); err == nil {
			p.songs[song] = stream
		} else {
			panic(err)
		}
	} else {
		panic(err)
	}
}

func (p *playerImpl) PlaySong(song string) {
	if stream, ok := p.songs[song]; !ok {
		panic("song not found")
	} else {
		if p.musicPlayer != nil {
			p.StopCurrentSong()
		}
		var err error
		if p.musicPlayer, err = p.audioContext.NewPlayer(stream); err != nil {
			panic(err)
		} else {
			p.musicPlayer.SetPosition(0)
			p.musicPlayer.SetVolume(p.songVolume)
			p.musicPlayer.Play()
		}
	}
}

func (p *playerImpl) StopCurrentSong() {
	if p.musicPlayer != nil {
		p.musicPlayer.SetPosition(0)
		p.musicPlayer.Pause()
		p.musicPlayer.Close()
		p.musicPlayer = nil
	}
}

func (p *playerImpl) Update() {
	if p.musicPlayer != nil {
		if !p.musicPlayer.IsPlaying() {
			p.musicPlayer.SetPosition(0)
			p.musicPlayer.Play()
		}
	}
}

func (p *playerImpl) LoadSound(sound string) {
	if _, ok := p.sounds[sound]; ok {
		return
	}
	if soundData, err := p.fileSystem.ReadFile(sound); err == nil {
		if stream, err := vorbis.DecodeWithoutResampling(bytes.NewReader(soundData)); err == nil {
			if p.sounds[sound], err = p.audioContext.NewPlayer(stream); err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	} else {
		panic(err)
	}
}

func (p *playerImpl) PlaySound(sound string) {
	if player, ok := p.sounds[sound]; !ok {
		panic("sound not found")
	} else {
		if player.IsPlaying() {
			player.Pause()
		}
		player.SetPosition(0)
		player.SetVolume(p.soundVolume)
		player.Play()
	}
}

func (p *playerImpl) StopAllSounds() {
	for _, player := range p.sounds {
		if player != nil {
			if player.IsPlaying() {
				player.Pause()
				player.SetPosition(0)
			}
		}
	}
}

func (p *playerImpl) Stop() {
	p.StopCurrentSong()
	p.StopAllSounds()
}

func (p *playerImpl) ChangeSongVolume(volume float64) {
	p.songVolume = volume
	if p.musicPlayer != nil {
		p.musicPlayer.SetVolume(volume)
	}
}

func (p *playerImpl) ChangeSoundVolume(volume float64) {
	p.soundVolume = volume
	for _, player := range p.sounds {
		player.SetVolume(volume)
	}
}

func NewPlayer(fileSystem embed.FS, songVolume float64, soundVolume float64) Player {
	return &playerImpl{
		audioContext: ebitenAudio.NewContext(44100),
		fileSystem:   fileSystem,
		songs:        make(map[string]*vorbis.Stream),
		sounds:       make(map[string]*ebitenAudio.Player),
		songVolume:   songVolume,
		soundVolume:  soundVolume,
	}
}
