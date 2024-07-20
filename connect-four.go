package main

import (
	"connect-four/board"
	"image/color"
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

func main() {

	const TILE_SIZE = int32(100)

	board := board.MakeConnectFourBoard()

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("Connect 4", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 700, 750, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	renderer, _ := sdl.CreateRenderer(window, -1, 0)

	// Colors
	red := sdl.Color{R: 255, G: 0, B: 0, A: 255}
	yellow := sdl.Color{R: 255, G: 255, B: 0, A: 255}
	blue := sdl.Color{R: 0, G: 0, B: 255, A: 255}
	black := sdl.Color{R: 0, G: 0, B: 0, A: 255}

	// Create surface for empty
	emtpySurface, _ := sdl.CreateRGBSurface(0, TILE_SIZE, TILE_SIZE, 32, 0, 0, 0, 0)
	emtpySurface.FillRect(&sdl.Rect{X: 0, Y: 0, W: TILE_SIZE, H: TILE_SIZE}, sdl.MapRGBA(emtpySurface.Format, blue.R, blue.G, blue.B, blue.A))
	FillCircle(emtpySurface, 50, 50, 30, black)
	emptyTexture, _ := renderer.CreateTextureFromSurface(emtpySurface)

	// Create surface for player 1
	player1Surface, _ := sdl.CreateRGBSurface(0, TILE_SIZE, TILE_SIZE, 32, 0, 0, 0, 0)
	player1Surface.FillRect(&sdl.Rect{X: 0, Y: 0, W: TILE_SIZE, H: TILE_SIZE}, sdl.MapRGBA(emtpySurface.Format, blue.R, blue.G, blue.B, blue.A))
	FillCircle(player1Surface, 50, 50, 30, red)
	player1Texture, _ := renderer.CreateTextureFromSurface(player1Surface)

	// Create surface for player 2
	player2Surface, _ := sdl.CreateRGBSurface(0, TILE_SIZE, TILE_SIZE, 32, 0, 0, 0, 0)
	player2Surface.FillRect(&sdl.Rect{X: 0, Y: 0, W: TILE_SIZE, H: TILE_SIZE}, sdl.MapRGBA(emtpySurface.Format, blue.R, blue.G, blue.B, blue.A))
	FillCircle(player2Surface, 50, 50, 30, yellow)
	player2Texture, _ := renderer.CreateTextureFromSurface(player2Surface)

	// Create surface for player 1 token
	player1TokenSurface, _ := sdl.CreateRGBSurface(0, TILE_SIZE, TILE_SIZE, 32, 0, 0, 0, 0)
	player1TokenSurface.FillRect(&sdl.Rect{X: 0, Y: 0, W: TILE_SIZE, H: TILE_SIZE}, sdl.MapRGBA(emtpySurface.Format, black.R, black.G, black.B, black.A))
	FillCircle(player1TokenSurface, 50, 50, 30, red)
	player1TokenTexture, _ := renderer.CreateTextureFromSurface(player1TokenSurface)

	// Create surface for player 2 token
	player2TokenSurface, _ := sdl.CreateRGBSurface(0, TILE_SIZE, TILE_SIZE, 32, 0, 0, 0, 0)
	player2TokenSurface.FillRect(&sdl.Rect{X: 0, Y: 0, W: TILE_SIZE, H: TILE_SIZE}, sdl.MapRGBA(emtpySurface.Format, black.R, black.G, black.B, black.A))
	FillCircle(player2TokenSurface, 50, 50, 30, yellow)
	player2TokenTexture, _ := renderer.CreateTextureFromSurface(player2TokenSurface)

	// Sample rendering
	running := true
	changed := true
	selectedColumn := 3
	srcRect := sdl.Rect{X: 0, Y: 0, W: TILE_SIZE, H: TILE_SIZE}
	for running {

		if changed {
			renderer.SetDrawColor(black.R, black.G, black.B, black.A)
			renderer.FillRect(&sdl.Rect{X: 0, Y: 0, W: 700, H: 750})

			if board.Winner() == 0 {
				dtsRect := sdl.Rect{X: int32(selectedColumn) * TILE_SIZE, Y: 0, W: TILE_SIZE, H: TILE_SIZE}
				switch board.CurrentPlayer() {
				case 1:
					renderer.Copy(player1TokenTexture, &srcRect, &dtsRect)
				case 2:
					renderer.Copy(player2TokenTexture, &srcRect, &dtsRect)
				}
			}

			for row := 0; row < 6; row++ {
				for column := 0; column < 7; column++ {
					value := board.GetContent(row, column)
					dtsRect := sdl.Rect{X: int32(column) * TILE_SIZE, Y: 6*TILE_SIZE - int32(row)*TILE_SIZE, W: TILE_SIZE, H: TILE_SIZE}
					switch value {
					case 0:
						renderer.Copy(emptyTexture, &srcRect, &dtsRect)
					case 1:
						renderer.Copy(player1Texture, &srcRect, &dtsRect)
					case 2:
						renderer.Copy(player2Texture, &srcRect, &dtsRect)
					}
				}
			}

			renderer.Present()

			changed = false
		}

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent: // NOTE: Please use `*sdl.QuitEvent` for `v0.4.x` (current version).
				println("Quit")
				running = false
			case *sdl.KeyboardEvent:
				if t.Type == sdl.KEYUP {
					if t.Keysym.Sym == sdl.K_RIGHT {
						changed = true
						selectedColumn = selectedColumn + 1
						if selectedColumn > 6 {
							selectedColumn = 6
						}
					} else if t.Keysym.Sym == sdl.K_LEFT {
						changed = true
						selectedColumn = selectedColumn - 1
						if selectedColumn < 0 {
							selectedColumn = 0
						}
					} else if t.Keysym.Sym == sdl.K_DOWN || t.Keysym.Sym == sdl.K_SPACE {
						changed = true
						board.Play(selectedColumn)
						selectedColumn = 3
					}
				}
			}
		}

		// TODO: Remove this delay
		sdl.Delay(33)
	}
}

func FillCircle(surface *sdl.Surface, centreX int, centreY int, radius int, c sdl.Color) {

	clr := color.RGBA{R: c.R, G: c.G, B: c.B, A: c.A}
	for x := centreX - radius; x <= (centreX + radius); x++ {
		for y := centreY - radius; y <= (centreY + radius); y++ {
			distance := math.Sqrt(float64((centreX-x)*(centreX-x) + (centreY-y)*(centreY-y)))
			if distance <= float64(radius) {
				surface.Set(x, y, clr)
			}
		}
	}
}
