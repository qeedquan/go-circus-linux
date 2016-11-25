package main

import (
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/qeedquan/go-media/sdl"
)

func dpath(path ...string) string {
	return filepath.Join(*assets, filepath.Join(path...))
}

func spath(path ...string) string {
	return filepath.Join(*pref, filepath.Join(path...))
}

func quitEvents() {
	for {
		ev := sdl.PollEvent()
		if ev == nil {
			break
		}
		switch ev := ev.(type) {
		case sdl.QuitEvent:
			os.Exit(0)
		case sdl.KeyDownEvent:
			switch ev.Sym {
			case sdl.K_ESCAPE:
				os.Exit(0)
			}
		}
	}
}

func delay(ms time.Duration) {
	ticker := time.NewTicker(ms * time.Millisecond)
	defer ticker.Stop()

loop:
	for {
		select {
		case <-ticker.C:
			break loop
		default:
			quitEvents()
			drawScreen()
		}
	}
}

func clamp(x, a, b int) int {
	if x < a {
		return a
	}
	if x > b {
		return b
	}
	return x
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func getInitials() string {
	user, err := user.Current()
	if err != nil {
		return "TUX"
	}

	s := ""
	l := 0
	for _, ch := range user.Name {
		s += string(ch)
		if l++; l == 3 {
			break
		}
	}

	for ; l < 3; l++ {
		s += " "
	}
	return s
}

func exit(status int) {
	saveOptions()
	os.Exit(status)
}
