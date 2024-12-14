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
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	SCREEN_WIDTH  = 480
	SCREEN_HEIGHT = 494

	FRAME_OX     = 0
	FRAME_OY     = 0
	FRAME_WIDTH  = 32
	FRAME_HEIGHT = 32

	BOARD_WIDTH  = 480
	BOARD_HEIGHT = 494

	COLS = 15
	ROWS = 15
	W    = 32

	STARTING_BOMBS = 38

	MESSAGE_FLAGS   = "Flags: %d"
	MESSAGE_DEFEAT  = "You're Dead!"
	MESSAGE_VICTORY = "You Win!"
	MESSAGE_TIME    = ""
)

var (
	STARTED = false
)

type Game struct {
	state   int
	message string

	flags int

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
	isFlag        bool

	X int
	Y int

	R int
	C int
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return SCREEN_WIDTH, SCREEN_HEIGHT
}

func (g *Game) Update() error {
	if g.state == 1 {
		handleMouse(g)
	}

	handleVictory(g)

	options(g)

	return nil
}

func handleVictory(g *Game) {
	won := true

	for r := 0; r < ROWS; r++ {
		for c := 0; c < COLS; c++ {
			tile := g.board.tiles[r][c]

			if tile.isBomb && !tile.isFlag {
				won = false
			}
		}
	}

	if won {
		g.state = 2
		g.message = MESSAGE_VICTORY
	}
}

func options(g *Game) {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		restart(g)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		gameOver(g)
	}
}

func restart(g *Game) {
	g.flags = STARTING_BOMBS
	g.state = 1
	g.message = ""

	createBoard(g)
	randomizeBombs(g)
	countBombs(g)
}

func handleMouse(g *Game) {
	if inpututil.IsMouseButtonJustPressed(0) {
		x, y := ebiten.CursorPosition()
		r, c := checkPosition(g, x, y)

		if r != -1 && c != -1 {
			tile := g.board.tiles[r][c]
			reveal(g, tile)

			if tile.isBomb {
				gameOver(g)
			}
		}
	}

	if inpututil.IsMouseButtonJustPressed(2) {
		x, y := ebiten.CursorPosition()
		r, c := checkPosition(g, x, y)

		if r != -1 && c != -1 {
			tile := g.board.tiles[r][c]
			placeFlag(g, tile)
		}
	}
}

func placeFlag(g *Game, tile *Tile) {
	if !tile.isRevealed {
		if !tile.isFlag {
			if g.flags > 0 {
				tile.isFlag = true
				g.flags -= 1
			}
		} else {
			tile.isFlag = false
			g.flags += 1
		}
	}
}

func floodFill(g *Game, tile *Tile) {
	for roff := -1; roff <= 1; roff++ {
		for coff := -1; coff <= 1; coff++ {
			i := tile.R + roff
			j := tile.C + coff

			if i > -1 && i < ROWS && j > -1 && j < COLS {
				neighbor := g.board.tiles[i][j]
				if !neighbor.isBomb && !neighbor.isRevealed {
					reveal(g, neighbor)
				}
			}
		}
	}
}

func reveal(g *Game, tile *Tile) {
	tile.isRevealed = true
	if tile.neighborCount == 0 {
		floodFill(g, tile)
	}
}

func checkPosition(g *Game, x int, y int) (int, int) {
	for r := 0; r < ROWS; r++ {
		for c := 0; c < COLS; c++ {
			tile := g.board.tiles[r][c]
			if x > tile.X && x < tile.X+W && y > tile.Y && y < tile.Y+W {
				return r, c
			}
		}
	}
	return -1, -1
}

func gameOver(g *Game) {
	g.state = 0
	g.message = MESSAGE_DEFEAT

	for r := 0; r < ROWS; r++ {
		for c := 0; c < COLS; c++ {
			g.board.tiles[r][c].isRevealed = true
		}
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{180, 180, 180, 255})

	opts := ebiten.DrawImageOptions{}

	drawTiles(g, opts, screen)
	drawStats(g, screen)
}

func main() {
	game := &Game{
		state:   1,
		flags:   STARTING_BOMBS,
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
	for r := 0; r < ROWS; r++ {
		for c := 0; c < COLS; c++ {
			tile := g.board.tiles[r][c]

			if tile.isBomb {
				tile.neighborCount = 10
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
			offset := tile.neighborCount
			if !tile.isRevealed {
				offset = 9
			}

			if tile.isFlag {
				offset = 11
			}

			fox, foy, fw, fh := FRAME_OX, FRAME_OY, FRAME_WIDTH, FRAME_HEIGHT
			foy += 32 * offset
			fh *= offset + 1

			opts.GeoM.Translate(float64(c)*W, float64(r)*W)
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

func drawStats(g *Game, screen *ebiten.Image) {
	flags := fmt.Sprintf(MESSAGE_FLAGS, g.flags)

	ebitenutil.DebugPrintAt(screen, g.message, 1, SCREEN_HEIGHT-16)
	ebitenutil.DebugPrintAt(screen, flags, 420, SCREEN_HEIGHT-16)
}

func createBoard(g *Game) {
	tiles, err := ebitenutil.NewImageFromURL("https://github.com/RafaelEtec/go_minesweeper/blob/03c9bb37cee39d27eb2b2cb70cbef9bfbee2fe15/assets/images/tiles_new.png?raw=true")
	if err != nil {
		log.Fatal(err)
	}

	for r := 0; r < ROWS; r++ {
		for c := 0; c < COLS; c++ {
			g.board.tiles[r][c].isRevealed = false
			g.board.tiles[r][c].isBomb = false
			g.board.tiles[r][c].isFlag = false
			g.board.tiles[r][c].neighborCount = 0
			g.board.tiles[r][c].Img = tiles
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
