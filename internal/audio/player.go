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
	Update()
}

type playerImpl struct {
	audioContext *ebitenAudio.Context
	musicPlayer  *ebitenAudio.Player
	fileSystem   embed.FS
	songs        map[string]*vorbis.Stream
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
			p.musicPlayer.SetVolume(0.20)
			p.musicPlayer.Play()
		}
	}
}

func (p *playerImpl) StopCurrentSong() {
	p.musicPlayer.SetPosition(0)
	p.musicPlayer.Pause()
	p.musicPlayer.Close()
	p.musicPlayer = nil
}

func (p *playerImpl) Update() {
	if p.musicPlayer != nil {
		if !p.musicPlayer.IsPlaying() {
			p.musicPlayer.SetPosition(0)
			p.musicPlayer.Play()
		}
	}
}

func NewPlayer(fileSystem embed.FS) Player {
	return &playerImpl{
		audioContext: ebitenAudio.NewContext(44100),
		fileSystem:   fileSystem,
		songs:        make(map[string]*vorbis.Stream),
	}
}
