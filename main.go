package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var (
	player           *ebiten.Image
	bg               *ebiten.Image
	chocolateLog     *ebiten.Image
	jumpHeight       float64
	landingHeight    float64
	jumpUpComplete   bool
	jumpDownComplete bool
	jumpComplete     bool = true
)

const (
	screenW, screenH = 1000, 667
	playerW, playerH = 128, 128
)

type Game struct {

	// Player position.
	playerX, playerY float64

	// Player velocity.
	playerVelocity float64
}

func (g *Game) Update() error {

	// Handle jump.
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) && jumpComplete {
		jumpHeight = -150
		g.playerVelocity = 0.5
		jumpUpComplete = false
		jumpDownComplete = false
		jumpComplete = false
	}

	if !jumpUpComplete && jumpHeight < g.playerVelocity {
		g.playerVelocity -= 2
		if jumpHeight > g.playerVelocity {
			jumpUpComplete = true
		}
	}

	if jumpUpComplete {
		jumpHeight = 0
		if jumpHeight > g.playerVelocity {
			g.playerVelocity += 2
			if jumpHeight < g.playerVelocity {
				jumpDownComplete = true
			}
		}
	}

	if jumpUpComplete && jumpDownComplete {
		jumpComplete = true
	}

	// Apply velocity.
	g.playerY += g.playerVelocity
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw Background
	bgOp := &ebiten.DrawImageOptions{}
	screen.DrawImage(bg, bgOp)

	// Draw player
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(g.playerX, g.playerY)
	screen.DrawImage(player, op)

	// Draw ChocolateLog
	clOp := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, 0)
	screen.DrawImage(chocolateLog, clOp)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (w, h int) {
	g.playerX = screenW/2 - playerW/2
	g.playerY = screenH - playerH - 65
	return screenW, screenH
}

func init() {
	var err error
	bg, _, err = ebitenutil.NewImageFromFile("assets/background.png")
	if err != nil {
		log.Fatal(err)
	}
	player, _, err = ebitenutil.NewImageFromFile("assets/player.png")
	if err != nil {
		log.Fatal(err)
	}
	chocolateLog, _, err = ebitenutil.NewImageFromFile("assets/chocolateLog.png")
	if err != nil {
		log.Fatal(err)
	}

}
func main() {
	ebiten.SetWindowSize(screenW, screenH)
	ebiten.SetWindowTitle("ChocoJump")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
