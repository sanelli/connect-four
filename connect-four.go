package main

import (
	"connect-four/board"
	"image/color"
	"math"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

func main() {

	const TILE_SIZE = int32(100)

	gameBoard := board.MakeConnectFourBoard()

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	ttf.Init()
	defer ttf.Quit()

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
	white := sdl.Color{R: 255, G: 255, B: 255, A: 255}

	// Create surface for empty
	emtpySurface, _ := sdl.CreateRGBSurface(0, TILE_SIZE, TILE_SIZE, 32, 0, 0, 0, 0)
	emtpySurface.FillRect(&sdl.Rect{X: 0, Y: 0, W: TILE_SIZE, H: TILE_SIZE}, sdl.MapRGBA(emtpySurface.Format, blue.R, blue.G, blue.B, blue.A))
	FillCircle(emtpySurface, 50, 50, 30, black)
	emptyTexture, _ := renderer.CreateTextureFromSurface(emtpySurface)
	emtpySurface.Free()
	defer emptyTexture.Destroy()

	// Create surface for player 1
	player1Surface, _ := sdl.CreateRGBSurface(0, TILE_SIZE, TILE_SIZE, 32, 0, 0, 0, 0)
	player1Surface.FillRect(&sdl.Rect{X: 0, Y: 0, W: TILE_SIZE, H: TILE_SIZE}, sdl.MapRGBA(player1Surface.Format, blue.R, blue.G, blue.B, blue.A))
	FillCircle(player1Surface, 50, 50, 30, red)
	player1Texture, _ := renderer.CreateTextureFromSurface(player1Surface)
	player1Surface.Free()
	defer player1Texture.Destroy()

	// Create surface for player 2
	player2Surface, _ := sdl.CreateRGBSurface(0, TILE_SIZE, TILE_SIZE, 32, 0, 0, 0, 0)
	player2Surface.FillRect(&sdl.Rect{X: 0, Y: 0, W: TILE_SIZE, H: TILE_SIZE}, sdl.MapRGBA(player2Surface.Format, blue.R, blue.G, blue.B, blue.A))
	FillCircle(player2Surface, 50, 50, 30, yellow)
	player2Texture, _ := renderer.CreateTextureFromSurface(player2Surface)
	player2Surface.Free()
	defer player2Texture.Destroy()

	// Create surface for player 1 token
	player1TokenSurface, _ := sdl.CreateRGBSurface(0, TILE_SIZE, TILE_SIZE, 32, 0, 0, 0, 0)
	player1TokenSurface.FillRect(&sdl.Rect{X: 0, Y: 0, W: TILE_SIZE, H: TILE_SIZE}, sdl.MapRGBA(player1TokenSurface.Format, black.R, black.G, black.B, black.A))
	FillCircle(player1TokenSurface, 50, 50, 30, red)
	player1TokenTexture, _ := renderer.CreateTextureFromSurface(player1TokenSurface)
	player1TokenSurface.Free()
	defer player1TokenTexture.Destroy()

	// Create surface for player 2 token
	player2TokenSurface, _ := sdl.CreateRGBSurface(0, TILE_SIZE, TILE_SIZE, 32, 0, 0, 0, 0)
	player2TokenSurface.FillRect(&sdl.Rect{X: 0, Y: 0, W: TILE_SIZE, H: TILE_SIZE}, sdl.MapRGBA(player2TokenSurface.Format, black.R, black.G, black.B, black.A))
	FillCircle(player2TokenSurface, 50, 50, 30, yellow)
	player2TokenTexture, _ := renderer.CreateTextureFromSurface(player2TokenSurface)
	player2TokenSurface.Free()
	defer player2TokenTexture.Destroy()

	// Footer texture
	footerFont, _ := ttf.OpenFont("./BebasNeue-Regular.ttf", 30)
	defer footerFont.Close()

	footerSurface, _ := footerFont.RenderUTF8Solid("Left/Righ: Select column; Down: drop token; N: New game; Q: Exit", white)
	footerTexture, _ := renderer.CreateTextureFromSurface(footerSurface)
	footerRect := sdl.Rect{X: 0, Y: 0, W: footerSurface.W, H: footerSurface.H}
	footerSurface.Free()
	defer footerTexture.Destroy()

	// Winner textures
	winnerFont, _ := ttf.OpenFont("./BebasNeue-Regular.ttf", 52)
	defer winnerFont.Close()

	winner1Surface, _ := footerFont.RenderUTF8Solid("Player 1 won!", red)
	winner1Texture, _ := renderer.CreateTextureFromSurface(winner1Surface)
	winner1Rect := sdl.Rect{X: 0, Y: 0, W: winner1Surface.W, H: winner1Surface.H}
	winner1Surface.Free()
	defer winner1Texture.Destroy()

	winner2Surface, _ := footerFont.RenderUTF8Solid("Player 2 won!", yellow)
	winner2Texture, _ := renderer.CreateTextureFromSurface(winner2Surface)
	winner2Rect := sdl.Rect{X: 0, Y: 0, W: winner2Surface.W, H: winner2Surface.H}
	winner2Surface.Free()
	defer winner2Texture.Destroy()

	noWinnerSurface, _ := footerFont.RenderUTF8Solid("Nobody won!", white)
	noWinnerTexture, _ := renderer.CreateTextureFromSurface(noWinnerSurface)
	noWinnerRect := sdl.Rect{X: 0, Y: 0, W: noWinnerSurface.W, H: noWinnerSurface.H}
	noWinnerSurface.Free()
	defer noWinnerTexture.Destroy()

	// Sample rendering
	running := true
	changed := true
	selectedColumn := 3
	srcRect := sdl.Rect{X: 0, Y: 0, W: TILE_SIZE, H: TILE_SIZE}
	winner := 0
	for running {

		start := sdl.GetPerformanceCounter()

		if changed {
			renderer.SetDrawColor(black.R, black.G, black.B, black.A)
			renderer.FillRect(&sdl.Rect{X: 0, Y: 0, W: 700, H: 750})

			switch winner {
			case 0:
				dtsRect := sdl.Rect{X: int32(selectedColumn) * TILE_SIZE, Y: 0, W: TILE_SIZE, H: TILE_SIZE}
				switch gameBoard.CurrentPlayer() {
				case 1:
					renderer.Copy(player1TokenTexture, &srcRect, &dtsRect)
				case 2:
					renderer.Copy(player2TokenTexture, &srcRect, &dtsRect)
				}
			case 1:
				renderer.Copy(winner1Texture, &winner1Rect, &sdl.Rect{X: (700 - winner1Rect.W) / 2, Y: 25, W: winner1Rect.W, H: winner1Rect.H})
			case 2:
				renderer.Copy(winner2Texture, &winner2Rect, &sdl.Rect{X: (700 - winner2Rect.W) / 2, Y: 25, W: winner2Rect.W, H: winner2Rect.H})
			case -1:
				renderer.Copy(noWinnerTexture, &noWinnerRect, &sdl.Rect{X: (700 - noWinnerRect.W) / 2, Y: 25, W: noWinnerRect.W, H: noWinnerRect.H})
			}

			for row := 0; row < 6; row++ {
				for column := 0; column < 7; column++ {
					value := gameBoard.GetContent(row, column)
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

			renderer.Copy(footerTexture, &footerRect, &sdl.Rect{X: 0, Y: 710, W: footerRect.W, H: footerRect.H})
			renderer.Present()

			changed = false
		}

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.KeyboardEvent:
				if t.Type == sdl.KEYUP {
					if t.Keysym.Sym == sdl.K_q {
						running = false
					} else if t.Keysym.Sym == sdl.K_n {
						changed = true
						selectedColumn = 3
						winner = 0
						gameBoard = board.MakeConnectFourBoard()
					} else if winner == 0 && t.Keysym.Sym == sdl.K_RIGHT {
						changed = true
						selectedColumn = selectedColumn + 1
						if selectedColumn > 6 {
							selectedColumn = 6
						}
					} else if winner == 0 && t.Keysym.Sym == sdl.K_LEFT {
						changed = true
						selectedColumn = selectedColumn - 1
						if selectedColumn < 0 {
							selectedColumn = 0
						}
					} else if winner == 0 && t.Keysym.Sym == sdl.K_DOWN || t.Keysym.Sym == sdl.K_SPACE {
						changed = true
						gameBoard.Play(selectedColumn)
						selectedColumn = 3
						winner = gameBoard.Winner()
					}
				}
			}
		}

		end := sdl.GetPerformanceCounter()
		elapsedMS := float64(end-start) / float64(sdl.GetPerformanceFrequency()) * 1000.0
		sdl.Delay(uint32(math.Trunc(math.Max(0.00, math.Floor(float64(16.666)-elapsedMS)))))
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
