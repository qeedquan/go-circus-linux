package main

import (
	"log"

	"github.com/qeedquan/go-media/sdl/sdlmixer"
)

const (
	SND_POP = iota
	SND_BOUNCE
	SND_TEETER1
	SND_TEETER2
	SND_SPLAT
	SND_APPLAUSE
	SND_CHEERING
	SND_HIGHSCORE
	SND_KEYPRESS
	NUM_SOUNDS
)

var soundNames = [NUM_SOUNDS]string{
	"sounds/pop.wav",
	"sounds/bounce.wav",
	"sounds/teeter1.wav",
	"sounds/teeter2.wav",
	"sounds/splat.wav",
	"sounds/applause.wav",
	"sounds/cheering.wav",
	"sounds/wahoo.wav",
	"sounds/keypress.wav",
}

const (
	MUS_TITLE = iota
	MUS_GAME
	MUS_GAMEOVER
	MUS_HISCORE
	MUS_HISCORESCREEN
	NUM_MUSICS
)

var musicNames = [...]string{
	"music/finally.ogg",
	"music/klovninarki.ogg",
	"music/kaupunki.ogg",
	"music/hiscore.ogg",
	"music/hiscreen.ogg",
}

var (
	sounds [NUM_SOUNDS]*sdlmixer.Chunk
	musics [NUM_MUSICS]*sdlmixer.Music
)

func loadSound(name string) *sdlmixer.Chunk {
	chunk, err := sdlmixer.LoadWAV(name)
	if err != nil {
		log.Print("sound: ", err)
	}
	return chunk
}

func loadMusic(name string) *sdlmixer.Music {
	music, err := sdlmixer.LoadMUS(name)
	if err != nil {
		log.Print("music: ", err)
	}
	return music
}

func streamMusic(id int) {
	if *enableMusic && sdlmixer.PlayingMusic() == 0 {
		musics[id].Play(0)
		sdlmixer.VolumeMusic(musicVol * sdlmixer.MAX_VOLUME / 3)
	}
}

func playSound(id int) {
	if *enableSound && sounds[id] != nil {
		sounds[id].PlayChannel(-1, 0)
	}
}
