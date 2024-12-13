package main

import (
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	SCREEN_WIDTH  = 294
	SCREEN_HEIGHT = 308

	FRAME_OX     = 0
	FRAME_OY     = 0
	FRAME_WIDTH  = 32
	FRAME_HEIGHT = 32

	BOARD_WIDTH  = 288
	BOARD_HEIGHT = 288

	COLS = 20
	ROWS = 20
	W    = 32

	STARTING_BOMBS = 30

	MESSAGE_MOVES   = "Moves: %d"
	MESSAGE_DEFEAT  = "You Die!"
	MESSAGE_VICTORY = "You Win!"
)

var (
	STARTED = false
	X       = -1
	Y       = -1
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
	tiles [20][20]*Tile
}

type Tile struct {
	Img           *ebiten.Image
	neighborCount int

	isBomb     bool
	isRevealed bool
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
			tiles: [20][20]*Tile{},
		},
	}

	createBoard(game)

	//randomizeBombs(game)

	ebiten.SetWindowSize(SCREEN_WIDTH*2, SCREEN_HEIGHT*2)
	ebiten.SetWindowTitle("MINESWEEPER by Rafael Goulart")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
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
	tile, err := ebitenutil.NewImageFromFile("/assets/images/clues.png")
	if err != nil {
		log.Fatal(err)
	}

	for r := 0; r < ROWS; r++ {
		for c := 0; c < COLS; c++ {
			g.board.tiles[r][c].isRevealed = false
			g.board.tiles[r][c].isBomb = false
			g.board.tiles[r][c].neighborCount = 0
			g.board.tiles[r][c].Img = tile
		}
	}
}
