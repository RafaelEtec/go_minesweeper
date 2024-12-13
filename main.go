package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math/rand/v2"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	SCREEN_WIDTH  = 488
	SCREEN_HEIGHT = 488

	FRAME_OX     = 0
	FRAME_OY     = 0
	FRAME_WIDTH  = 32
	FRAME_HEIGHT = 32

	BOARD_WIDTH  = 488
	BOARD_HEIGHT = 488

	COLS = 15
	ROWS = 15
	W    = 32

	STARTING_BOMBS = 50

	MESSAGE_MOVES   = "Moves: %d"
	MESSAGE_DEFEAT  = "You Die!"
	MESSAGE_VICTORY = "You Win!"
)

var (
	STARTED = false
)

type Game struct {
	state   int
	message string

	board *Board
}

// state: 0 = stopped
//        1 = playing
//        2 = won

type Board struct {
	tiles [][]*Tile
}

type Tile struct {
	Img *ebiten.Image

	neighborCount int
	isBomb        bool
	isRevealed    bool

	X int
	Y int

	R int
	C int
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return SCREEN_WIDTH, SCREEN_HEIGHT
}

func (g *Game) Update() error {

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{180, 180, 180, 255})

	opts := ebiten.DrawImageOptions{}

	drawTiles(g, opts, screen)
}

func main() {
	game := &Game{
		state:   1,
		message: "",
		board: &Board{
			tiles: [][]*Tile{
				{{}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}},
				{{}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}},
				{{}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}},
				{{}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}},
				{{}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}},
				{{}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}},
				{{}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}},
				{{}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}},
				{{}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}},
				{{}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}},
				{{}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}},
				{{}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}},
				{{}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}},
				{{}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}},
				{{}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}},
			},
		},
	}

	createBoard(game)

	randomizeBombs(game)
	countBombs(game)

	printBoard(game)

	ebiten.SetWindowSize(SCREEN_WIDTH*2, SCREEN_HEIGHT*2)
	ebiten.SetWindowTitle("MINESWEEPER by Rafael Goulart")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func countBombs(g *Game) {
	bomb, err := ebitenutil.NewImageFromURL("https://github.com/RafaelEtec/go_minesweeper/blob/27dc2e25bf4362beb80684bc2c91c56963481388/assets/images/skull.png?raw=true")
	if err != nil {
		log.Fatal(err)
	}

	for r := 0; r < ROWS; r++ {
		for c := 0; c < COLS; c++ {
			tile := g.board.tiles[r][c]

			if tile.isBomb {
				tile.Img = bomb
				continue
			}

			total := 0
			for roff := -1; roff <= 1; roff++ {
				i := tile.R + roff
				if i < 0 || i >= ROWS {
					continue
				}

				for coff := -1; coff <= 1; coff++ {
					j := tile.C + coff
					if j < 0 || j >= COLS {
						continue
					}

					neighbor := g.board.tiles[i][j]
					if neighbor.isBomb {
						total++
					}
				}
			}
			tile.neighborCount = total
		}
	}
}

func randomizeBombs(g *Game) {
	for n := 0; n < STARTING_BOMBS; {
		r := rand.IntN(ROWS)
		c := rand.IntN(COLS)

		if !g.board.tiles[r][c].isBomb {
			g.board.tiles[r][c].isBomb = true
			n++
		}
	}
}

func drawTiles(g *Game, opts ebiten.DrawImageOptions, screen *ebiten.Image) {
	for r := 0; r < ROWS; r++ {
		for c := 0; c < COLS; c++ {
			tile := g.board.tiles[r][c]

			fox, foy, fw, fh := FRAME_OX, FRAME_OY, FRAME_WIDTH, FRAME_HEIGHT
			foy += 32 * tile.neighborCount
			fh *= tile.neighborCount + 1

			rspace := r/3 + 2
			cspace := c/3 + 2

			opts.GeoM.Translate(float64(c)*32+float64(cspace), float64(r)*32+float64(rspace))
			screen.DrawImage(
				tile.Img.SubImage(
					image.Rect(fox, foy, fw, fh),
				).(*ebiten.Image),
				&opts,
			)

			opts.GeoM.Reset()
		}
	}
}

func createBoard(g *Game) {
	blank, err := ebitenutil.NewImageFromURL("https://github.com/RafaelEtec/go_minesweeper/blob/27dc2e25bf4362beb80684bc2c91c56963481388/assets/images/blank.png?raw=true")
	if err != nil {
		log.Fatal(err)
	}

	for r := 0; r < ROWS; r++ {
		for c := 0; c < COLS; c++ {
			g.board.tiles[r][c].isRevealed = false
			g.board.tiles[r][c].isBomb = false
			g.board.tiles[r][c].neighborCount = 0
			g.board.tiles[r][c].Img = blank
			g.board.tiles[r][c].X = c * W
			g.board.tiles[r][c].Y = r * W
			g.board.tiles[r][c].R = r
			g.board.tiles[r][c].C = c
		}
	}
}

func printBoard(g *Game) {
	for i := 0; i < ROWS; i++ {
		for j := 0; j < COLS; j++ {
			if j == ROWS-1 {
				fmt.Print(g.board.tiles[i][j].neighborCount)
			} else {
				fmt.Print(g.board.tiles[i][j].neighborCount, "-")
			}
			time.Sleep(time.Millisecond * 1)
		}
		fmt.Println("")
	}
	fmt.Println("")
}
