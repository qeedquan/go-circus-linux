package main

import (
	"bufio"
	"flag"
	"fmt"
	"image"
	"image/draw"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/qeedquan/go-media/sdl"
	"github.com/qeedquan/go-media/sdl/sdlmixer"
)

var (
	assets        = flag.String("assets", filepath.Join(sdl.GetBasePath(), "assets"), "assets directory")
	pref          = flag.String("pref", sdl.GetPrefPath("", "circus_linux"), "pref directory")
	fullscreen    = flag.Bool("fullscreen", false, "fullscreen")
	enableMusic   = flag.Bool("music", true, "enable music")
	enableSound   = flag.Bool("sound", true, "enable sound")
	infiniteLives = flag.Bool("inflives", false, "infinite lives")
)

// this game draws things incrementally, it doesn't erase everything
// and redraw every frames, meaning if we want a resizable window
// we have to render surface -> texture -> screen if we want sdl to
// be able to resize for us automagically. using render to texture
// kind of works, but it is glitches alot when resizing, sometimes
// it doesn't resize and just move the display top left.
var (
	screen  *Display
	rgba    *image.RGBA
	texture *sdl.Texture
)

func main() {
	runtime.LockOSThread()
	log.SetFlags(0)
	rand.Seed(time.Now().UnixNano())
	flag.Parse()
	loadOptions()
	initSDL()
	load()
	intro()
	loop()
	exit(0)
}

func initSDL() {
	err := sdl.Init(sdl.INIT_EVERYTHING &^ sdl.INIT_AUDIO)
	if err != nil {
		log.Fatal("sdl: ", err)
	}

	err = sdlmixer.OpenAudio(44100, sdlmixer.DEFAULT_FORMAT, 2, 4096)
	if err != nil {
		log.Println("sdl: ", err)
	}

	sdlmixer.AllocateChannels(128)

	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "best")
	wflag := sdl.WINDOW_RESIZABLE
	if *fullscreen {
		wflag |= sdl.WINDOW_FULLSCREEN_DESKTOP
	}
	width, height := 640, 480
	screen, err = newDisplay(width, height, wflag)
	if err != nil {
		log.Fatal("sdl: ", err)
	}

	screen.SetLogicalSize(width, height)
	screen.SetTitle("Circus Linux!")

	rgba = image.NewRGBA(image.Rect(0, 0, width, height))
	texture, err = screen.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, width, height)
	if err != nil {
		log.Fatal("sdl: ", err)
	}

	icon := loadSurface(dpath("images/icon.png"))
	screen.SetIcon(icon)
	icon.Free()

	sdl.StartTextInput()
}

type Display struct {
	*sdl.Window
	*sdl.Renderer
}

func newDisplay(w, h int, wflag sdl.WindowFlags) (*Display, error) {
	window, renderer, err := sdl.CreateWindowAndRenderer(w, h, wflag)
	if err != nil {
		return nil, err
	}
	return &Display{window, renderer}, nil
}

func load() {
	for i := 0; i < NUM_IMAGES; i++ {
		images[i] = loadImage(dpath(imageNames[i]))
		quitEvents()
		dest := image.Rect(0, 470, 640*i/NUM_IMAGES, 470+10)
		draw.Draw(rgba, dest, white, image.ZP, draw.Over)
		drawScreen()
	}

	for i := 0; i < NUM_SOUNDS; i++ {
		sounds[i] = loadSound(dpath(soundNames[i]))
	}

	for i := 0; i < NUM_MUSICS; i++ {
		musics[i] = loadMusic(dpath(musicNames[i]))
	}
}

func loadOptions() {
	sfxVol = 3
	musicVol = 3
	for i := range hiscore {
		hiscore[i] = Rank{"TUX", 100}
	}

	f, err := os.Open(spath("circuslinux.dat"))
	if err != nil {
		return
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	for s.Scan() {
		var name string
		var i, score int

		line := strings.TrimSpace(s.Text())
		n, err := fmt.Sscanf(line, "highscore%d=%d", &i, &score)
		if n == 2 && err == nil {
			if i >= 0 && i < len(hiscore) {
				hiscore[i].score = score
			}
			continue
		}

		n, err = fmt.Sscanf(line, "highscorer%d=%s", &i, &name)
		if n == 2 && err == nil {
			if i >= 0 && i < len(hiscore) {
				hiscore[i].name = ""
				l := 0
				for _, ch := range name {
					hiscore[i].name += string(ch)
					if l++; l == 3 {
						break
					}
				}
			}
			continue
		}

		n, err = fmt.Sscanf(line, "effects=%d", &sfxVol)
		if n == 1 && err == nil {
			continue
		}

		n, err = fmt.Sscanf(line, "music=%d", &musicVol)
		if n == 1 && err == nil {
			continue
		}
	}

	if sfxVol < 0 || sfxVol > 3 {
		sfxVol = 3
	}

	if musicVol < 0 || musicVol > 3 {
		musicVol = 3
	}
}

func saveOptions() {
	os.MkdirAll(spath(""), 0755)
	f, err := os.Create(spath("circuslinux.dat"))
	if err != nil {
		log.Print("options: ", err)
		return
	}
	defer f.Close()

	bw := bufio.NewWriter(f)
	for i, h := range hiscore {
		fmt.Fprintf(bw, "highscore%d=%d\n", i, h.score)
		fmt.Fprintf(bw, "highscorer%d=%s\n", i, h.name)
	}
	fmt.Fprintf(bw, "effects=%d\n", sfxVol)
	fmt.Fprintf(bw, "music=%d", musicVol)

	err = bw.Flush()
	if err != nil {
		log.Print("options: ", err)
	}

	cerr := f.Close()
	if err == nil && cerr != nil {
		log.Print("options: ", cerr)
	}
}
