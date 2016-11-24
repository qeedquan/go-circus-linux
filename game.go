package main

import (
	"image"
	"image/draw"
	"math"
	"math/rand"
	"unicode"

	"github.com/qeedquan/go-media/sdl"
	"github.com/qeedquan/go-media/sdl/sdlmixer"
)

type Rank struct {
	name  string
	score int
}

const (
	NUM_TITLE_BALLOONS     = 16
	STARTING_LIVES         = 5
	NUM_BARRIERS           = 3
	FASTEST_YM_OFF_BALLOON = 8
	MAX_YM                 = 32
	FPS                    = 1000 / 33
	FLYING_SPLAT_TIME      = 50
	NUM_ROWS               = 3
	FLYING_START_Y         = 192
	FLYING_START_YM        = -8
	BOUNCER_TIME           = 6
	GRAVITY                = 1
	SHOW_PLAYER_TIME       = 100
	LIMB_ANIMATION_TIME    = 8
	MAX_SCORE              = math.MaxInt32
)

const (
	NUM_BACKGROUND_CHANGES = 4
)

const (
	LEFT = iota
	RIGHT
)

const (
	GONE     = 0
	NORMAL   = 1
	POPPING0 = 2
	POPPING1 = 2
)

const (
	ACT_SEAL = iota
	ACT_BEAR
	NUM_ACTS
)

var backgroundChangeRects = [NUM_BACKGROUND_CHANGES]sdl.Rect{
	{424, 0, 88, 127},
	{256, 150, 153, 87},
	{26, 288, 57, 63},
	{580, 295, 44, 55},
}

var (
	numPlayers = 1
	barriers   = 0
	bouncy     = 0
	clearAll   = 0
	sfxVol     = 3
	musicVol   = 3

	coOp bool

	balloons      [2][NUM_ROWS][20]int
	balloonColors [2][NUM_ROWS]int

	lives [2]int
	score [2]int

	hasHighscore    int
	highscoreIndex  int
	showHighscores  bool
	highscoreEffect int
	hiscore         [8]Rank

	flyingActive       bool
	flyingSplat        int
	flyingDir          int
	flyingx, flyingy   int
	flyingxm, flyingym int
	flyingLeftArm      int
	flyingRightArm     int
	flyingLeftLeg      int
	flyingRightLeg     int
)

func intro() {
	for i := 0; i < 50; i++ {
		switch i {
		case 5:
			drawText(32, 176, "NEW BREED SOFTWARE")
		case 25:
			drawText(192, 288, "PRESENTS")
		}

		delay(30)
	}
}

func loop() {
	for {
		showHighscores = false
		title()
		if showHighscores {
			highScoreScreen()
		} else {
			game()
		}
	}
}

func title() {
	defer sdlmixer.HaltMusic()

	draw.Draw(rgba, rgba.Bounds(), images[IMG_TITLE], image.ZP, draw.Src)

	hilight := sdl.Rect{-1, -1, -1, -1}
	oldHilight := hilight

	var x, y, xm, ym, col, bumped [NUM_TITLE_BALLOONS]int
	for i := 0; i < NUM_TITLE_BALLOONS; i++ {
		x[i] = rand.Intn(640 - 32)
		y[i] = rand.Intn(480 - 32)

		xm[i] = rand.Intn(5) + 1
		if rand.Intn(2) == 0 {
			xm[i] = -xm[i]
		}

		ym[i] = rand.Intn(5) + 1
		if rand.Intn(2) == 0 {
			ym[i] = -ym[i]
		}

		col[i] = rand.Intn(8)*6 + IMG_BALLOON_RED_LEFT_0
		bumped[i] = 0
	}

	textx, textxm, textImg, textTime := -640, 36, 0, 0
	highscoreEffect = 0
	frame := uint(0)

	for {
		frame++

		for {
			ev := sdl.PollEvent()
			if ev == nil {
				break
			}

			switch ev := ev.(type) {
			case sdl.QuitEvent:
				exit(0)

			case sdl.KeyDownEvent:
				switch ev.Sym {
				case sdl.K_ESCAPE:
					exit(0)
				}

			case sdl.MouseButtonDownEvent:
				switch {
				// one players
				case ev.X >= 16 && ev.X <= 16+238 &&
					ev.Y >= 283 && ev.Y <= 283+27:
					numPlayers = 1
					return

				// two players
				case ev.X >= 16 && ev.X <= 16+264 &&
					ev.Y >= 310 && ev.Y <= 310+27:
					numPlayers = 2
					return

				// start two player coop
				case ev.X >= 16 && ev.X <= 16+356 &&
					ev.Y >= 337 && ev.Y <= 337+27:
					numPlayers = 2
					coOp = true
					return

				// toggle barriers
				case ev.X <= 207 &&
					ev.Y >= 371 && ev.Y <= 371+27:
					barriers = 1 - barriers
					playSound(SND_TEETER2 - barriers)

				// bouncy balloons
				case ev.X <= 374 &&
					ev.Y >= 398 && ev.Y <= 398+27:
					bouncy = 1 - bouncy
					playSound(SND_TEETER2 - bouncy)

				// clear all
				case ev.X <= 234 &&
					ev.Y >= 425 && ev.Y <= 425+27:
					clearAll = 1 - clearAll
					playSound(SND_TEETER2 - bouncy)

				// set sfx volume
				case ev.X >= 559 && ev.X <= 559+73 &&
					ev.Y >= 284 && ev.Y <= 284+25:
					sfxVol = (sfxVol + 1) % 4
					sdlmixer.Volume(-1, sfxVol*sdlmixer.MAX_VOLUME/3)
					playSound(SND_POP)

				// set music volume
				case ev.X >= 512 && ev.X <= 512+121 &&
					ev.Y >= 336 && ev.Y <= 336+52:
					musicVol = (musicVol + 1) % 4
					sdlmixer.VolumeMusic(musicVol * sdlmixer.MAX_VOLUME / 3)

				// hiscore
				case ev.X >= 440 && ev.X <= 440+195 &&
					ev.Y >= 398 && ev.Y <= 398+29:
					showHighscores = true
					playSound(SND_HIGHSCORE)
					return

				// exit
				case ev.X >= 535 && ev.X <= 535+100 &&
					ev.Y >= 429 && ev.Y <= 429+29:
					exit(1)
				}

			case sdl.MouseMotionEvent:
				switch {
				// one players
				case ev.X >= 16 && ev.X <= 16+238 &&
					ev.Y >= 283 && ev.Y <= 283+27:
					hilight = sdl.Rect{16, 283, 238, 27}

				// two players
				case ev.X >= 16 && ev.X <= 16+264 &&
					ev.Y >= 310 && ev.Y <= 310+27:
					hilight = sdl.Rect{16, 310, 264, 27}

				// start two player coop
				case ev.X >= 16 && ev.X <= 16+356 &&
					ev.Y >= 337 && ev.Y <= 337+27:
					hilight = sdl.Rect{16, 337, 356, 27}

				// toggle barriers
				case ev.X <= 207 &&
					ev.Y >= 371 && ev.Y <= 371+27:
					hilight = sdl.Rect{0, 371, 207, 27}

				// bouncy balloons
				case ev.X <= 374 &&
					ev.Y >= 398 && ev.Y <= 398+27:
					hilight = sdl.Rect{0, 398, 374, 27}

				// clear all
				case ev.X <= 234 &&
					ev.Y >= 425 && ev.Y <= 425+27:
					hilight = sdl.Rect{0, 425, 234, 27}

				// set sfx volume
				case ev.X >= 559 && ev.X <= 559+73 &&
					ev.Y >= 284 && ev.Y <= 284+25:
					hilight = sdl.Rect{559, 284, 73, 52}

				// set music volume
				case ev.X >= 512 && ev.X <= 512+121 &&
					ev.Y >= 336 && ev.Y <= 336+52:
					hilight = sdl.Rect{512, 336, 121, 52}

				// hiscore
				case ev.X >= 440 && ev.X <= 440+195 &&
					ev.Y >= 398 && ev.Y <= 398+29:
					hilight = sdl.Rect{440, 398, 195, 29}

				// exit
				case ev.X >= 535 && ev.X <= 535+100 &&
					ev.Y >= 429 && ev.Y <= 429+29:
					hilight = sdl.Rect{535, 429, 100, 29}

				default:
					hilight = sdl.Rect{-1, -1, -1, -1}
				}
			}
		}

		sdl.Delay(30)

		// erase the hilights
		if hilight != oldHilight {
			if oldHilight.X != -1 {
				erase(oldHilight.X, oldHilight.Y, oldHilight.W, oldHilight.H, IMG_TITLE)
			}
			oldHilight = hilight
		}

		// erase balloons
		for i := 0; i < NUM_TITLE_BALLOONS; i++ {
			erase(int32(x[i]), int32(y[i]), 32, 32, IMG_TITLE)
		}

		// erase credits
		erase(0, 252, 640, 32, IMG_TITLE)

		// move the balloons
		for i := 0; i < NUM_TITLE_BALLOONS; i++ {
			x[i] += xm[i]
			y[i] += ym[i]

			if frame%3 == 0 {
				ym[i] += 1
				if ym[i] > 16 {
					ym[i] = 16
				}
			}

			bumped[i] = 0
		}

		// make the balloons bounce into each other
		for i := 0; i < NUM_TITLE_BALLOONS; i++ {
			for j := 0; j < NUM_TITLE_BALLOONS; j++ {
				if i != j && bumped[j] == 0 && bumped[i] == 0 {
					if x[i] > x[j]-32 &&
						x[i] < x[j]+32 &&
						y[i] > y[j]-32 &&
						y[i] < y[j]+32 {

						x[i] -= xm[i] * 2 / 3
						y[i] -= ym[i] * 2 / 3

						xm[i], xm[j] = xm[j], xm[i]
						ym[i], ym[j] = ym[j], ym[i]

						bumped[i], bumped[j] = 1, 1
					}
				}
			}
		}

		// keep the balloons in bound
		for i := 0; i < NUM_TITLE_BALLOONS; i++ {
			if x[i] <= 0 {
				xm[i] = rand.Intn(5) + 1
				x[i] = 0
			} else if x[i] >= 640-32 {
				xm[i] = -(rand.Intn(5) + 1)
				x[i] = 640 - 32
			}

			if y[i] <= 0 {
				ym[i] = rand.Intn(5) + 1
				y[i] = 0
			} else if y[i] >= 480-32 {
				ym[i] = -ym[i]
				y[i] = 480 - 32
			}
		}

		// handle credits
		if textx < 0 {
			textx += textxm
			if textx >= 0 {
				textx = 0
				textxm = 0
			} else {
				textxm--
			}
		} else if textx == 0 && textTime < 100 {
			textTime++
			textxm = 0
		} else if textx < 640 {
			textx += textxm
			textxm++
		} else {
			textImg = (textImg + 1) % 3
			textx = -640
			textxm = 36
			textTime = 0
		}

		// draw the hilights
		if hilight.X != -1 {
			r := image.Rect(int(hilight.X), int(hilight.Y), int(hilight.X+hilight.W), int(hilight.Y+hilight.H))
			sp := image.Pt(int(hilight.X), int(hilight.Y-275))
			draw.Draw(rgba, r, images[IMG_TITLE_HIGHLIGHTS], sp, draw.Over)
		}

		// draw whether the options are on
		if barriers != 0 {
			drawObj(0, 376, IMG_LIGHT_ON)
		} else {
			drawObj(0, 376, IMG_LIGHT_OFF)
		}

		if bouncy != 0 {
			drawObj(0, 404, IMG_LIGHT_ON)
		} else {
			drawObj(0, 404, IMG_LIGHT_OFF)
		}

		if clearAll != 0 {
			drawObj(0, 429, IMG_LIGHT_ON)
		} else {
			drawObj(0, 429, IMG_LIGHT_OFF)
		}

		for i := 0; i < 3; i++ {
			if sfxVol > i {
				drawObj(583+i*16, 311, IMG_LIGHT_ON)
			} else {
				drawObj(583+i*16, 311, IMG_LIGHT_OFF)
			}

			if musicVol > i {
				drawObj(583+i*16, 363, IMG_LIGHT_ON)
			} else {
				drawObj(583+i*16, 363, IMG_LIGHT_OFF)
			}
		}

		// draw the balloons
		for i := 0; i < NUM_TITLE_BALLOONS; i++ {
			n := 0
			if xm[i] > 0 {
				n = 2
			}
			drawObj(x[i], y[i], col[i]+n+rand.Intn(2))
		}

		// draw the credits
		drawObj(textx, 252, IMG_PROGRAMMER+textImg)

		drawScreen()

		streamMusic(MUS_TITLE)
	}
}

func highScoreScreen() {
	defer sdlmixer.HaltMusic()

	draw.Draw(rgba, rgba.Bounds(), black, image.ZP, draw.Src)
	drawObj(0, 0, IMG_HIGHSCORE_TOP)

	// draw scores
	height := images[IMG_HIGHSCORE_TOP].Bounds().Dy()
	for i, h := range hiscore {
		y := height + i*32
		drawNumber(32, y+16, h.score, IMG_NUMBERS_0+(i%2))
		drawText(224, y+16, h.name)

		// barriers
		if i >= 4 {
			drawObj(336, y+24, IMG_LIGHT_ON)
		} else {
			drawObj(336, y+24, IMG_LIGHT_OFF)
		}

		// bouncy balloons
		if i == 2 || i == 3 || i == 6 || i == 7 {
			drawObj(444, y+24, IMG_LIGHT_ON)
		} else {
			drawObj(444, y+24, IMG_LIGHT_OFF)
		}

		// clear all
		if i%2 == 1 {
			drawObj(564, y+24, IMG_LIGHT_ON)
		} else {
			drawObj(564, y+24, IMG_LIGHT_OFF)
		}
	}

	for {
		for {
			ev := sdl.PollEvent()
			if ev == nil {
				break
			}
			switch ev.(type) {
			case *sdl.QuitEvent:
				exit(0)
			case *sdl.KeyDownEvent:
				return
			}
		}

		sdl.Delay(30)
		drawScreen()
		streamMusic(MUS_HISCORESCREEN)
	}
}

func pauseScreen() int {
	sdlmixer.PauseMusic()
	defer sdlmixer.ResumeMusic()

	drawFuzz(224, 224, 192, 32)
	drawText(224, 224, "PAUSED")

	rc := 0
loop:
	for {
		for {
			ev := sdl.PollEvent()
			if ev == nil {
				break
			}
			switch ev := ev.(type) {
			case *sdl.QuitEvent:
				exit(0)
			case *sdl.KeyDownEvent:
				switch ev.Sym {
				case sdl.K_SPACE, sdl.K_TAB, sdl.K_p:
					break loop
				case sdl.K_ESCAPE:
					rc = 1
					break loop
				}
			}
		}
		drawScreen()
	}

	erase(224, 224, 192, 32, IMG_BACKGROUND_0)
	return rc
}

func game() {
	var barrierx [NUM_BARRIERS]int
	var bouncers [2]int

	hasHighscore = -1
	act := rand.Intn(NUM_ACTS)
	actx, actxm := 0, 0
	acty, actym := 0, 0
	frame := 0
	backgroundFrame := 0
	highscoreIndex = (barriers*4 + bouncy*2 + clearAll)
	fire := false
	showPlayer := SHOW_PLAYER_TIME
	mousex, mousey := 0, 0

	flyingActive = false
	flyingSplat = 0

	teeterSide := LEFT
	oldTeeterx := 0
	teeterx, teeterxm, teeterxmm := 0, 0, 0
	teeterRoll, teeterSound := 0, 0

	player := 0
	for i := 0; i < 2; i++ {
		lives[i] = STARTING_LIVES
		score[i] = 0

		for y := 0; y < NUM_ROWS; y++ {
			resetBalloons(i, y)
		}

		for y := 0; y < NUM_ROWS; y++ {
			balloonColors[i][y] = y * 2
		}
	}

	for i := range barrierx {
		barrierx[i] = i * 128
	}

	doPause := false

	draw.Draw(rgba, rgba.Bounds(), images[IMG_BACKGROUND_0], image.ZP, draw.Src)
loop:
	for {
		lastTime := sdl.GetTicks()
		frame++

		if doPause {
			if pauseScreen() != 0 {
				break loop
			}
			doPause = false
		}

		// animate background
		if frame%5 == 0 {
			backgroundFrame = (backgroundFrame + 1) % 2
			updateBackground(backgroundFrame)
		}

		// erase teeter-totter
		erase(int32(teeterx), 444, 96, 36, IMG_BACKGROUND_0+backgroundFrame)

		// erase flying clown
		if flyingActive || flyingSplat != 0 {
			erase(int32(flyingx), int32(flyingy), 32, 32, IMG_BACKGROUND_0+backgroundFrame)
		}

		// erase bouncers
		for i := 0; i < 2; i++ {
			erase(608*int32(i), 448, 32, 32, IMG_BACKGROUND_0+backgroundFrame)
		}

		// erase act
		erase(152, 347, 48, 48, IMG_BACKGROUND_0+backgroundFrame)

		// erase balloons
		for y := 0; y < NUM_ROWS; y++ {
			erase(0, int32(y)*32+32, 640, 32, IMG_BACKGROUND_0+backgroundFrame)
		}

		// erase barriers
		if barriers != 0 {
			erase(0, NUM_ROWS*32+32, 640, 32, IMG_BACKGROUND_0+backgroundFrame)
		}

		// erase lives status
		erase(512, 0, 128, 32, IMG_BACKGROUND_0+backgroundFrame)

		// erase score status
		erase(0, 0, 192, 32, IMG_BACKGROUND_0+backgroundFrame)

		// move teeter totter
		oldTeeterx = teeterx
		teeterx += teeterxm
		teeterxm += teeterxmm

		teeterxm = clamp(teeterxm, -32, 32)
		teeterx = clamp(teeterx, 32, 512)

		fire = false

		for {
			ev := sdl.PollEvent()
			if ev == nil {
				break
			}

			keys := sdl.GetKeyboardState()
			if keys[sdl.SCANCODE_LEFT] != 0 {
				teeterxmm = -2
				if teeterxm > 0 {
					teeterxm = 0
				}
			} else if keys[sdl.SCANCODE_RIGHT] != 0 {
				teeterxmm = 2
				if teeterxm < 0 {
					teeterxm = 0
				}
			} else {
				teeterxmm = 0
				teeterxm = 0
			}

			switch ev := ev.(type) {
			case sdl.QuitEvent:
				exit(0)

			case sdl.KeyDownEvent:
				switch ev.Sym {
				case sdl.K_ESCAPE:
					break loop
				case sdl.K_SPACE, sdl.K_TAB, sdl.K_p:
					doPause = true
				case sdl.K_RETURN:
					fire = true
				}

			case sdl.MouseMotionEvent:
				teeterx = int(ev.X - 48)
				mousex, mousey = int(ev.X), int(ev.Y)
				teeterx = clamp(teeterx, 32, 512)

			case sdl.MouseButtonDownEvent:
				fire = true
			}
		}

		// handle fire
		if fire {
			// swap teeter-totter side
			if flyingActive {
				teeterSide = 1 - teeterSide
			}

			// activate a new clown
			if !flyingActive && flyingSplat == 0 {
				newClown()

				if showPlayer > 0 {
					showPlayer = 1
				}
			}
		}

		// handle high score effect
		if highscoreEffect > 0 {
			highscoreEffect--
		}

		// handle the barrel
		if teeterx > oldTeeterx {
			if teeterRoll++; teeterRoll > 3 {
				teeterRoll = 0
			}
		} else if teeterx < oldTeeterx {
			if teeterRoll--; teeterRoll < 0 {
				teeterRoll = 3
			}
		}

		// handle bouncers
		for i := range bouncers {
			if bouncers[i] > 0 {
				bouncers[i]--
			}
		}

		// handle barriers
		if barriers != 0 {
			for i := range barrierx {
				barrierx[i] += 8
				if barrierx[i] > 640 {
					barrierx[i] = -64
				}
			}
		}

		// handle balloons
		any := false
		for y := 0; y < NUM_ROWS; y++ {
			// handle popping
			some := false

			for x := range balloons[player][y] {
				if balloons[player][y][x] == POPPING0 {
					balloons[player][y][x] = GONE
				} else if balloons[player][y][x] > POPPING0 {
					balloons[player][y][x]--
				}

				if balloons[player][y][x] == NORMAL {
					some = true
					any = true
				}
			}

			// all balloons popped? add more?
			if !some && flyingy > NUM_ROWS*32+64 && clearAll == 0 {
				resetBalloons(player, y)
				balloonColors[player][y]++

				if balloonColors[player][y] > 7 {
					balloonColors[player][y] = 0
				}

				addPlayerScores(player, y)
			}

			// move balloons
			if frame%4 == 0 {
				if y%2 == 0 {
					// left
					i := balloons[player][y][0]
					for x := 0; x < 19; x++ {
						balloons[player][y][x] = balloons[player][y][x+1]
					}
					balloons[player][y][19] = i
				} else {
					// right
					i := balloons[player][y][19]
					for x := 18; x >= 0; x-- {
						balloons[player][y][x+1] = balloons[player][y][x]
					}
					balloons[player][y][0] = i
				}
			}
		}

		// reset all balloons?
		if clearAll != 0 && !any && flyingy > NUM_ROWS*32+64 {
			for y := 0; y < 3; y++ {
				resetBalloons(player, y)
				balloonColors[player][y]++

				if balloonColors[player][y] > 7 {
					balloonColors[player][y] = 0
				}

				addPlayerScores(player, y)
			}
		}

		// handle flying clown
		if flyingActive {
			// move the clown
			flyingx += flyingxm
			flyingy += flyingym

			// bounce off top
			if flyingy < 32 {
				flyingy = 32
				flyingym = 0
			}

			// bounce off bouncers
			if flyingy > 416 && flyingx < 32 || flyingx > 576 {
				flyingy = 416
				flyingym = -abs(flyingym)

				// make bouncer squish and launch clown torwards center
				if flyingx < 32 {
					bouncers[0] = BOUNCER_TIME
					flyingxm = 8
				} else {
					bouncers[1] = BOUNCER_TIME
					flyingxm = -8
				}

				// give a point for bouncing
				addScore(player, 1)

				playSound(SND_BOUNCE)
			}

			// bounce off barriers
			if barriers != 0 {
				for i := 0; i < NUM_BARRIERS; i++ {
					if flyingy >= (NUM_ROWS*32) && flyingy <= (NUM_ROWS*32)+64 &&
						flyingx >= barrierx[i]-32 && flyingx <= barrierx[i]+64 {
						if flyingy <= (NUM_ROWS*32)+32 {
							flyingy = NUM_ROWS * 32
							flyingym = -abs(flyingym)
						} else {
							flyingy = (NUM_ROWS * 32) + 64
							flyingym = abs(flyingym)
						}

						playSound(SND_BOUNCE)
					}
				}
			}

			// Bounce off teeter-totter or splat
			if flyingy > 448 {
				flyingy = 448

				// Did we hit the teeter-totter?

				if (teeterSide == RIGHT && (flyingx >= teeterx && flyingx <= teeterx+96)) ||
					(teeterSide == LEFT && (flyingx >= teeterx-32 && flyingx <= teeterx+64)) {
					// Yes!  Bounce other the clown:

					flyingym = -(abs(flyingx-(teeterx+48-16)) / 3) - 16
					flyingy = 432

					if teeterSide == LEFT {
						flyingx = teeterx + 64
						flyingDir = LEFT
						teeterSide = RIGHT
					} else {
						flyingx = teeterx
						flyingDir = RIGHT
						teeterSide = LEFT
					}

					// Randomly pick a X direction:

					if rand.Intn(2) == 0 {
						if flyingxm != 0 {
							flyingxm = 0
						} else {
							flyingxm = -4
						}
					}

					// Randomly change X direction (sign):

					if rand.Intn(2) == 0 {
						flyingxm = -flyingxm
					}

					// Give a point for bouncing:
					addScore(player, 1)

					/* Play teeter-totter bounce sound: */

					playSound(SND_TEETER1 + teeterSound)
					teeterSound = 1 - teeterSound
				} else {
					// No!  Splat the flying clown!

					flyingActive = false
					flyingSplat = FLYING_SPLAT_TIME
					if !*infiniteLives {
						lives[player]--
					}
					playSound(SND_SPLAT)
				}
			}

			// bounce off balloons
			x := (flyingx + 16) / 32
			y := (flyingy / 32) - 1

			if (frame/2)%2 != 0 {
				if (y % 2) == 0 {
					x = (flyingx / 32)
				} else {
					x = (flyingx / 32) + 1
				}
			}

			if y >= 0 && y < NUM_ROWS {
				if balloons[player][y][x] == NORMAL {
					balloons[player][y][x] = POPPING1
					playSound(SND_POP)

					addScore(player, y+1)

					// Bounce horizontally:
					if (flyingx % 32) < 16 {
						flyingxm = -4
					} else {
						flyingxm = 4
					}

					// Bounce vertically:
					if bouncy == 1 {
						flyingym = -flyingym

						if flyingym > FASTEST_YM_OFF_BALLOON {
							flyingym = FASTEST_YM_OFF_BALLOON
						}
					}
				}
			}

			// bounce off sides
			if flyingx < 0 {
				flyingx = 0
				flyingxm = abs(flyingxm)
			} else if flyingx > 608 {
				flyingx = 608
				flyingxm = -abs(flyingxm)
			}

			/* Deal with gravity: */

			flyingym = flyingym + GRAVITY

			if flyingym > MAX_YM {
				flyingym = MAX_YM
			}
			if flyingym < -MAX_YM {
				flyingym = -MAX_YM
			}
		}

		// count splats down
		if flyingSplat > 0 {
			flyingSplat--

			// if out of clowns, show game over while clown is splat
			if lives[player] == 0 {
				if numPlayers == 1 {
					drawFuzz(176, 224, 288, 32)
					drawText(176, 224, "GAME OVER")
				} else if numPlayers == 2 {
					drawFuzz(176, 192, 288, 96)
					drawText(224, 192, "PLAYER")
					drawNumber(304, 224, player+1, IMG_NUMBERS_0+player)
					drawText(176, 256, "GAME OVER")
				}
			}

			if flyingSplat == 0 {
				// switch players
				if numPlayers == 2 {
					if coOp {
						// copy balloons in coop mode
						for y := 0; y < NUM_ROWS; y++ {
							for x := range balloons[player][y] {
								balloons[1-player][y][x] = balloons[player][y][x]
							}
						}
					}

					// swap player
					player = 1 - player
					if lives[player] == 0 {
						player = 1 - player
						erase(0, 0, 640, 480, IMG_BACKGROUND_0+backgroundFrame)
						drawScreen()
					}
				}

				// erase game over display
				if numPlayers == 1 {
					erase(176, 224, 288, 32, IMG_BACKGROUND_0+backgroundFrame)
				} else {
					erase(176, 192, 288, 96, IMG_BACKGROUND_0+backgroundFrame)
				}

				// show which player is playing now
				showPlayer = SHOW_PLAYER_TIME

				if lives[player] == 0 {
					break loop
				}
			}
		}

		// change limb positions
		if flyingActive || flyingSplat != 0 {
			if frame%LIMB_ANIMATION_TIME == 0 {
				flyingLeftArm = rand.Intn(3)
				flyingRightArm = rand.Intn(3)
				flyingLeftLeg = rand.Intn(2)
				flyingRightLeg = rand.Intn(2)
			}
		}

		// draw act
		switch act {
		case ACT_SEAL:
			erase(148, int32(acty), 32, 32, IMG_BACKGROUND_0+backgroundFrame)

			acty += actym
			actym++

			if acty >= 315 {
				acty = 315
				actym = -10
			}

			drawObj(148, acty, IMG_BEACHBALL_0+((frame/4)%3))
			drawObj(152, 347, IMG_SEAL_0+((frame/4)%2))

		case ACT_BEAR:
			erase(int32(actx), 340, 48, 96, IMG_BACKGROUND_0+backgroundFrame)
			actx += actxm
			if act <= 64 {
				actx = 64
				actxm = 4
			} else if actx >= 524 {
				actx = 524
				actxm = -4
			}

			if actxm > 0 {
				drawObj(actx, 340, IMG_BEAR_RIGHT_0+((frame/4)%2))
			} else {
				drawObj(actx, 340, IMG_BEAR_LEFT_0+((frame/4)%2))
			}
		}

		// draw balloons
		for y := 0; y < NUM_ROWS; y++ {
			for x := 0; x < 20; x++ {
				drawBalloon(player, x, y, (frame/2)%2)
			}
		}

		// draw barriers
		if barriers != 0 {
			for i := 0; i < NUM_BARRIERS; i++ {
				drawObj(barrierx[i], NUM_ROWS*32+32, IMG_BARRIER)
			}
		}

		// draw teeter-totter
		switch teeterSide {
		case LEFT:
			drawObj(teeterx, 448, IMG_TEETER_TOTTER_LEFT_0+teeterRoll)
		case RIGHT:
			drawObj(teeterx, 448, IMG_TEETER_TOTTER_RIGHT_0+teeterRoll)
		}

		// draw clown on teeter-totter
		drawClown(teeterx+64-(teeterSide*64), 444, teeterSide, 1, 1, 1, 1)

		// draw flying clown
		if flyingActive {
			drawClown(flyingx, flyingy, flyingDir, flyingLeftArm, flyingRightArm, flyingLeftLeg, flyingRightLeg)
		}

		// draw splat clown
		if flyingSplat != 0 {
			drawClown(flyingx, flyingy, 2, -1, -1, flyingLeftLeg+2, flyingRightLeg+2)
		}

		// draw bouncers
		for i := range bouncers {
			x, y := 608*i, 448
			if bouncers[i] == 0 {
				drawObj(x, y, IMG_BOUNCER_0)
			} else {
				drawObj(x, y, IMG_BOUNCER_1)
			}
		}

		// draw lives status
		drawFuzz(512, 0, 128, 32)

		if mousex < 500 || mousex > 556 || mousey > 44 {
			// not near head show normal clown face
			drawObj(512, 0, IMG_CLOWN_HEAD)
		} else {
			// near head, show excited face
			drawObj(512, 0, IMG_CLOWN_HEAD_OH)
		}

		drawObj(544, 0, IMG_TIMES)
		drawNumber(576, 0, lives[player], IMG_NUMBERS_0+player)

		// draw score status
		drawFuzz(0, 0, 192, 32)
		drawNumber(0, 0, score[player], IMG_NUMBERS_0+player)

		// draw "Player X Ready" message
		if showPlayer > 0 && lives[player] > 0 {
			if numPlayers == 1 {
				drawFuzz(240, 224, 160, 32)
				drawText(240, 224, "READY")
			} else if numPlayers == 2 {
				drawFuzz(224, 192, 192, 96)
				drawText(224, 192, "PLAYER")
				drawNumber(304, 224, player+1, IMG_NUMBERS_0+player)
				drawText(240, 256, "READY")
			}

			showPlayer--

			if showPlayer == 0 {
				if numPlayers == 1 {
					erase(240, 224, 160, 32, IMG_BACKGROUND_0+backgroundFrame)
				} else if numPlayers == 2 {
					erase(224, 192, 192, 96, IMG_BACKGROUND_0+backgroundFrame)
				}
			}
		}

		drawScreen()
		nowTime := sdl.GetTicks()
		if nowTime < lastTime+FPS {
			sdl.Delay(lastTime + FPS - nowTime)
		}

		streamMusic(MUS_GAME)
	}

	drawScreen()
	sdlmixer.HaltMusic()
	sdlmixer.HaltChannel(-1)

	// show scores
	drawFuzz(0, 0, 640, 480)
	delay(300)

	draw.Draw(rgba, rgba.Bounds(), black, image.ZP, draw.Src)

	// draw "Final Score(s)" text
	if numPlayers == 1 {
		drawText(144, 144, "FINAL SCORE")
	} else {
		drawText(128, 144, "FINAL SCORES")
	}

	// show players score
	if numPlayers == 1 {
		drawNumber(224, 272, score[0], IMG_NUMBERS_0)
		if hasHighscore == 0 {
			drawText(208, 0, "HISCORE")
			drawObj(207, 64, IMG_ENTER_INITIALS)
		}
	} else {
		// show player 1 score
		drawText(0, 208, "PLAYER")
		drawText(208, 208, "ONE")
		drawNumber(0, 272, score[0], IMG_NUMBERS_0)

		if hasHighscore == 0 {
			drawText(0, 0, "HISCORE")
			drawObj(0, 64, IMG_ENTER_INITIALS)
		}

		// show player 2 score
		drawText(336, 208, "PLAYER")
		drawText(544, 208, "TWO")
		drawNumber(336, 272, score[1], IMG_NUMBERS_1)

		if hasHighscore == 1 {
			drawText(416, 0, "HISCORE")
			drawObj(415, 64, IMG_ENTER_INITIALS)
		}
	}

loop1:
	for {
		frame++
		for {
			ev := sdl.PollEvent()
			if ev == nil {
				break
			}
			switch ev := ev.(type) {
			case sdl.QuitEvent:
				exit(0)

			case sdl.KeyDownEvent:
				switch ev.Sym {
				case sdl.K_ESCAPE:
					break loop1
				case sdl.K_BACKSPACE:
					if hasHighscore != -1 {
						h := &hiscore[highscoreIndex]
						l := len(h.name)
						if l > 0 {
							h.name = h.name[:len(h.name)-1]
						}
						playSound(SND_POP)
					}
				case sdl.K_RETURN:
					playSound(SND_HIGHSCORE)
					break loop1
				}

			case sdl.TextInputEvent:
				for _, ch := range string(ev.Text[:]) {
					if ch == 0 {
						break
					}

					ch = unicode.ToUpper(ch)
					if !('A' <= ch && ch <= 'Z') {
						continue
					}

					if hasHighscore != -1 {
						h := &hiscore[highscoreIndex]
						l := len(h.name)
						if l < 3 {
							h.name += string(ch)
						}
						playSound(SND_HIGHSCORE)
					}
				}

			case sdl.MouseButtonDownEvent:
				break loop1
			}
		}

		// update text
		if hasHighscore != -1 {
			if hasHighscore == 0 {
				if numPlayers == 1 {
					r := image.Rect(272, 32, 272+96, 32+32)
					draw.Draw(rgba, r, black, image.ZP, draw.Src)
					drawText(272, 32, hiscore[highscoreIndex].name)
				} else {
					r := image.Rect(0, 32, 96, 32+32)
					draw.Draw(rgba, r, black, image.ZP, draw.Src)
					drawText(0, 32, hiscore[highscoreIndex].name)
				}
			} else {
				r := image.Rect(544, 32, 544+96, 32+32)
				draw.Draw(rgba, r, black, image.ZP, draw.Src)
				drawText(544, 32, hiscore[highscoreIndex].name)
			}
		}

		drawObj(512, 320, IMG_SADCLOWN_0+((frame/5)%3))
		drawScreen()

		if hasHighscore == -1 {
			streamMusic(MUS_GAMEOVER)
		} else {
			streamMusic(MUS_HISCORE)
		}

		sdl.Delay(30)
	}

	// did a name get entered?
	if hasHighscore != -1 {
		if hiscore[highscoreIndex].name == "" {
			hiscore[highscoreIndex].name = getInitials()
		}
	}

	sdlmixer.HaltMusic()
}

func updateBackground(which int) {
	for _, br := range backgroundChangeRects {
		r := image.Rect(int(br.X), int(br.Y), int(br.X+br.W), int(br.Y+br.H))
		sp := image.Pt(int(br.X), int(br.Y))
		draw.Draw(rgba, r, images[IMG_BACKGROUND_0+which], sp, draw.Over)
	}
}

func newClown() {
	flyingActive = true
	flyingSplat = 0

	flyingx = 608 * rand.Intn(2)
	flyingy = FLYING_START_Y

	if flyingx > 0 {
		flyingDir = RIGHT
	} else {
		flyingDir = LEFT
	}

	flyingxm = 0
	flyingym = FLYING_START_YM

	flyingLeftArm = rand.Intn(3)
	flyingRightArm = rand.Intn(3)
	flyingLeftLeg = rand.Intn(2)
	flyingRightLeg = rand.Intn(2)
}

func addPlayerScores(player, y int) {
	switch y {
	case 0:
		addScore(player, 1000)
		lives[player]++
		playSound(SND_CHEERING)
	case 1:
		addScore(player, 250)
		playSound(SND_APPLAUSE)
	case 2:
		addScore(player, 100)
		playSound(SND_APPLAUSE)
	}
}

func addScore(player, inc int) {
	score[player] += inc
	if score[player] > MAX_SCORE {
		score[player] = MAX_SCORE
	}

	if score[player] >= hiscore[highscoreIndex].score {
		// did they just get this high score?

		if hasHighscore != player {
			hasHighscore = player
			playSound(SND_HIGHSCORE)
			highscoreEffect = 50
		}

		hiscore[highscoreIndex].score = score[player]
	}
}

func resetBalloons(player, row int) {
	for i := range balloons[player][row] {
		balloons[player][row][i] = NORMAL
	}
}
