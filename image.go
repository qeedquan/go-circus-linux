package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"
	"math/rand"
	"os"
	"unicode"

	"github.com/qeedquan/go-media/image/imageutil"
	"github.com/qeedquan/go-media/sdl"
	"github.com/qeedquan/go-media/sdl/sdlimage"
	"github.com/qeedquan/go-media/sdl/sdlimage/sdlcolor"
)

const (
	IMG_TITLE = iota
	IMG_TITLE_HIGHLIGHTS
	IMG_LIGHT_ON
	IMG_LIGHT_OFF
	IMG_PROGRAMMER
	IMG_GRAPHICS
	IMG_MUSIC

	IMG_BACKGROUND_0
	IMG_BACKGROUND_1

	IMG_BALLOON_RED_LEFT_0
	IMG_BALLOON_RED_LEFT_1
	IMG_BALLOON_RED_RIGHT_0
	IMG_BALLOON_RED_RIGHT_1
	IMG_BALLOON_RED_DIE_0
	IMG_BALLOON_RED_DIE_1

	IMG_BALLOON_ORANGE_LEFT_0
	IMG_BALLOON_ORANGE_LEFT_1
	IMG_BALLOON_ORANGE_RIGHT_0
	IMG_BALLOON_ORANGE_RIGHT_1
	IMG_BALLOON_ORANGE_DIE_0
	IMG_BALLOON_ORANGE_DIE_1

	IMG_BALLOON_YELLOW_LEFT_0
	IMG_BALLOON_YELLOW_LEFT_1
	IMG_BALLOON_YELLOW_RIGHT_0
	IMG_BALLOON_YELLOW_RIGHT_1
	IMG_BALLOON_YELLOW_DIE_0
	IMG_BALLOON_YELLOW_DIE_1

	IMG_BALLOON_GREEN_LEFT_0
	IMG_BALLOON_GREEN_LEFT_1
	IMG_BALLOON_GREEN_RIGHT_0
	IMG_BALLOON_GREEN_RIGHT_1
	IMG_BALLOON_GREEN_DIE_0
	IMG_BALLOON_GREEN_DIE_1

	IMG_BALLOON_CYAN_LEFT_0
	IMG_BALLOON_CYAN_LEFT_1
	IMG_BALLOON_CYAN_RIGHT_0
	IMG_BALLOON_CYAN_RIGHT_1
	IMG_BALLOON_CYAN_DIE_0
	IMG_BALLOON_CYAN_DIE_1

	IMG_BALLOON_BLUE_LEFT_0
	IMG_BALLOON_BLUE_LEFT_1
	IMG_BALLOON_BLUE_RIGHT_0
	IMG_BALLOON_BLUE_RIGHT_1
	IMG_BALLOON_BLUE_DIE_0
	IMG_BALLOON_BLUE_DIE_1

	IMG_BALLOON_PURPLE_LEFT_0
	IMG_BALLOON_PURPLE_LEFT_1
	IMG_BALLOON_PURPLE_RIGHT_0
	IMG_BALLOON_PURPLE_RIGHT_1
	IMG_BALLOON_PURPLE_DIE_0
	IMG_BALLOON_PURPLE_DIE_1

	IMG_BALLOON_WHITE_LEFT_0
	IMG_BALLOON_WHITE_LEFT_1
	IMG_BALLOON_WHITE_RIGHT_0
	IMG_BALLOON_WHITE_RIGHT_1
	IMG_BALLOON_WHITE_DIE_0
	IMG_BALLOON_WHITE_DIE_1

	IMG_CLOWN_BODY_LEFT
	IMG_CLOWN_BODY_RIGHT

	IMG_CLOWN_BODY_UPSIDE_DOWN

	IMG_CLOWN_LEFT_ARM_0
	IMG_CLOWN_LEFT_ARM_1
	IMG_CLOWN_LEFT_ARM_2

	IMG_CLOWN_RIGHT_ARM_0
	IMG_CLOWN_RIGHT_ARM_1
	IMG_CLOWN_RIGHT_ARM_2

	IMG_CLOWN_LEFT_LEG_0
	IMG_CLOWN_LEFT_LEG_1

	IMG_CLOWN_LEFT_LEG_0_UPSIDE_DOWN
	IMG_CLOWN_LEFT_LEG_1_UPSIDE_DOWN

	IMG_CLOWN_RIGHT_LEG_0
	IMG_CLOWN_RIGHT_LEG_1

	IMG_CLOWN_RIGHT_LEG_0_UPSIDE_DOWN
	IMG_CLOWN_RIGHT_LEG_1_UPSIDE_DOWN

	IMG_TEETER_TOTTER_LEFT_0
	IMG_TEETER_TOTTER_LEFT_1
	IMG_TEETER_TOTTER_LEFT_2
	IMG_TEETER_TOTTER_LEFT_3

	IMG_TEETER_TOTTER_RIGHT_0
	IMG_TEETER_TOTTER_RIGHT_1
	IMG_TEETER_TOTTER_RIGHT_2
	IMG_TEETER_TOTTER_RIGHT_3

	IMG_BOUNCER_0
	IMG_BOUNCER_1
	IMG_BARRIER

	IMG_TIMES
	IMG_NUMBERS_0
	IMG_NUMBERS_1
	IMG_LETTERS
	IMG_FUZZ
	IMG_CLOWN_HEAD
	IMG_CLOWN_HEAD_OH
	IMG_SADCLOWN_0
	IMG_SADCLOWN_1
	IMG_SADCLOWN_2
	IMG_ENTER_INITIALS

	IMG_HIGHSCORE_TOP
	IMG_HIGHSCORE_LEFT

	IMG_SEAL_0
	IMG_SEAL_1

	IMG_BEACHBALL_0
	IMG_BEACHBALL_1
	IMG_BEACHBALL_2

	IMG_BEAR_RIGHT_0
	IMG_BEAR_RIGHT_1

	IMG_BEAR_LEFT_0
	IMG_BEAR_LEFT_1

	NUM_IMAGES
)

var imageNames = []string{
	"images/title/title.png",
	"images/title/title-highlights.png",
	"images/title/light-on.png",
	"images/title/light-off.png",
	"images/title/programming.png",
	"images/title/graphics.png",
	"images/title/music.png",

	"images/backgrounds/background0.png",
	"images/backgrounds/background1.png",

	"images/balloons/red-left-0.png",
	"images/balloons/red-left-1.png",
	"images/balloons/red-right-0.png",
	"images/balloons/red-right-1.png",
	"images/balloons/red-die-0.png",
	"images/balloons/red-die-1.png",

	"images/balloons/orange-left-0.png",
	"images/balloons/orange-left-1.png",
	"images/balloons/orange-right-0.png",
	"images/balloons/orange-right-1.png",
	"images/balloons/orange-die-0.png",
	"images/balloons/orange-die-1.png",

	"images/balloons/yellow-left-0.png",
	"images/balloons/yellow-left-1.png",
	"images/balloons/yellow-right-0.png",
	"images/balloons/yellow-right-1.png",
	"images/balloons/yellow-die-0.png",
	"images/balloons/yellow-die-1.png",

	"images/balloons/green-left-0.png",
	"images/balloons/green-left-1.png",
	"images/balloons/green-right-0.png",
	"images/balloons/green-right-1.png",
	"images/balloons/green-die-0.png",
	"images/balloons/green-die-1.png",

	"images/balloons/cyan-left-0.png",
	"images/balloons/cyan-left-1.png",
	"images/balloons/cyan-right-0.png",
	"images/balloons/cyan-right-1.png",
	"images/balloons/cyan-die-0.png",
	"images/balloons/cyan-die-1.png",

	"images/balloons/blue-left-0.png",
	"images/balloons/blue-left-1.png",
	"images/balloons/blue-right-0.png",
	"images/balloons/blue-right-1.png",
	"images/balloons/blue-die-0.png",
	"images/balloons/blue-die-1.png",

	"images/balloons/purple-left-0.png",
	"images/balloons/purple-left-1.png",
	"images/balloons/purple-right-0.png",
	"images/balloons/purple-right-1.png",
	"images/balloons/purple-die-0.png",
	"images/balloons/purple-die-1.png",

	"images/balloons/white-left-0.png",
	"images/balloons/white-left-1.png",
	"images/balloons/white-right-0.png",
	"images/balloons/white-right-1.png",
	"images/balloons/white-die-0.png",
	"images/balloons/white-die-1.png",

	"images/clowns/body-left.png",
	"images/clowns/body-right.png",

	"images/clowns/body-upside-down.png",

	"images/clowns/left-arm-0.png",
	"images/clowns/left-arm-1.png",
	"images/clowns/left-arm-2.png",

	"images/clowns/right-arm-0.png",
	"images/clowns/right-arm-1.png",
	"images/clowns/right-arm-2.png",

	"images/clowns/left-leg-0.png",
	"images/clowns/left-leg-1.png",

	"images/clowns/left-leg-0-upside-down.png",
	"images/clowns/left-leg-1-upside-down.png",

	"images/clowns/right-leg-0.png",
	"images/clowns/right-leg-1.png",

	"images/clowns/right-leg-0-upside-down.png",
	"images/clowns/right-leg-1-upside-down.png",

	"images/teeter-totter/left-0.png",
	"images/teeter-totter/left-1.png",
	"images/teeter-totter/left-2.png",
	"images/teeter-totter/left-3.png",

	"images/teeter-totter/right-0.png",
	"images/teeter-totter/right-1.png",
	"images/teeter-totter/right-2.png",
	"images/teeter-totter/right-3.png",

	"images/bouncers/bouncer-0.png",
	"images/bouncers/bouncer-1.png",
	"images/bouncers/barrier.png",

	"images/status/times.png",
	"images/status/numbers-0.png",
	"images/status/numbers-1.png",
	"images/status/letters.png",
	"images/status/fuzz.png",
	"images/status/clown-head.png",
	"images/status/clown-head-oh.png",
	"images/status/sadclown-0.png",
	"images/status/sadclown-1.png",
	"images/status/sadclown-2.png",
	"images/status/enter-initials.png",

	"images/highscore/top.png",
	"images/highscore/left.png",

	"images/acts/seal-0.png",
	"images/acts/seal-1.png",

	"images/acts/beachball-0.png",
	"images/acts/beachball-1.png",
	"images/acts/beachball-2.png",

	"images/acts/bear-right-0.png",
	"images/acts/bear-right-1.png",
	"images/acts/bear-left-0.png",
	"images/acts/bear-left-1.png",
}

var (
	images [NUM_IMAGES]image.Image
	white  = image.NewUniform(color.RGBA{0xFF, 0xFF, 0xFF, 0xFF})
	black  = image.NewUniform(color.RGBA{0x0, 0x0, 0x0, 0xFF})
)

func loadImage(name string) image.Image {
	f, err := os.Open(name)
	if err != nil {
		log.Fatal("image: ", err)
	}
	defer f.Close()

	m, _, err := image.Decode(f)
	if err != nil {
		log.Fatal("image: ", name, ": ", err)
	}

	return imageutil.ColorKey(m, sdl.Color{255, 255, 255, 255})
}

func loadSurface(name string) *sdl.Surface {
	surface, _ := sdlimage.LoadSurfaceImage(loadImage(name))
	return surface
}

func drawText(x, y int, text string) {
	for i, ch := range text {
		ch = unicode.ToUpper(ch)
		if 'A' <= ch && ch <= 'Z' {
			px := x + i*32
			r := image.Rect(px, y, px+32, y+32)
			sp := image.Pt(int(ch-'A')*32, 0)
			draw.Draw(rgba, r, images[IMG_LETTERS], sp, draw.Over)
		}
	}
}

func drawBalloon(player, x, y, off int) {
	if balloons[player][y][x] == GONE {
		return
	}

	img := IMG_BALLOON_RED_LEFT_0 + balloonColors[player][y]*6
	if highscoreEffect != 0 {
		img = IMG_BALLOON_RED_LEFT_0 + rand.Intn(8)*6
	}

	if balloons[player][y][x] != NORMAL {
		img += 4
		if balloons[player][y][x] == POPPING0 {
			img++
		}
	} else {
		img += rand.Intn(2)
		img += (y % 2) * 2
	}

	px := x * 32
	if off == 1 {
		if y%2 == 0 {
			px -= 16
		} else {
			px += 16
		}
	}
	py := y*32 + 32
	dest := image.Rect(px, py, px+32, py+32)

	draw.Draw(rgba, dest, images[img], image.ZP, draw.Over)
}

func drawObj(x, y, pict int) {
	b := images[pict].Bounds()
	draw.Draw(rgba, image.Rect(x, y, x+b.Dx(), y+b.Dy()), images[pict], image.ZP, draw.Over)
}

func drawClown(x, y, side, leftArm, rightArm, leftLeg, rightLeg int) {
	drawObj(x, y, IMG_CLOWN_BODY_LEFT+side)

	if leftArm != -1 {
		drawObj(x, y, IMG_CLOWN_LEFT_ARM_0+leftArm)
	}

	if rightArm != -1 {
		drawObj(x, y, IMG_CLOWN_RIGHT_ARM_0+rightArm)
	}

	drawObj(x, y, IMG_CLOWN_LEFT_LEG_0+leftLeg)
	drawObj(x, y, IMG_CLOWN_RIGHT_LEG_0+rightLeg)
}

func drawNumber(x, y, v, img int) {
	str := fmt.Sprint(v)
	for i, ch := range str {
		px := x + i*32
		r := image.Rect(px, y, px+32, y+32)
		sp := image.Pt(int(ch-'0')*32, 0)
		draw.Draw(rgba, r, images[img], sp, draw.Over)
	}
}

func drawFuzz(x, y, w, h int) {
	b := images[IMG_FUZZ].Bounds()
	iw, ih := b.Dx(), b.Dy()
	for yy := y; yy < y+h; yy += ih {
		for xx := x; xx < x+w; xx += iw {
			sw, sh := iw, ih
			if xx+sw > x+w {
				sw = x + w - xx
			}

			if yy+sh > y+h {
				sh = y + h - yy
			}

			r := image.Rect(xx, yy, xx+sw, yy+sh)
			draw.Draw(rgba, r, images[IMG_FUZZ], image.ZP, draw.Over)
		}
	}
}

func drawScreen() {
	texture.Update(nil, rgba.Pix, rgba.Stride)
	screen.SetDrawColor(sdlcolor.Black)
	screen.Clear()
	screen.Copy(texture, nil, nil)
	screen.Present()
}

func erase(x, y, w, h int32, bkgd int) {
	draw.Draw(rgba, image.Rect(int(x), int(y), int(x+w), int(y+h)), images[bkgd], image.Pt(int(x), int(y)), draw.Src)
}
