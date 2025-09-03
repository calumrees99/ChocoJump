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
	chocoLog         *ebiten.Image
	jumpHeight       float64
	landingHeight    float64
	jumpUpComplete   bool
	jumpDownComplete bool
	jumpComplete     bool = true
)

const (
	screenW, screenH                             = 1000, 667
	playerW, playerH                             = 128, 128
	chocoLogW, chocoLogH                         = 128, 128
	playerHitboxW, playerHitboxH                 = 88, 100
	playerHitboxOffsetX, playerHitboxOffsetY     = 40, 25
	chocoLogHitboxW, chocoLogHitboxH             = 40, 60
	chocoLogHitboxOffsetX, chocoLogHitboxOffsetY = 30, 50
)

type Game struct {

	// Player position.
	playerX, playerY float64

	// Player velocity.
	playerVelocity float64

	// Choco log position
	chocoLogX, chocoLogY float64

	// Choco log velocity
	chocoLogVelocity float64
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
		g.playerVelocity -= 3
		if jumpHeight > g.playerVelocity {
			jumpUpComplete = true
		}
	}

	if jumpUpComplete {
		jumpHeight = 0
		if jumpHeight > g.playerVelocity {
			g.playerVelocity += 3
			if jumpHeight < g.playerVelocity {
				jumpDownComplete = true
			}
		}
	}

	if jumpUpComplete && jumpDownComplete {
		jumpComplete = true
	}

	// Apply player velocity.
	g.playerY += g.playerVelocity

	// Move chocoLog
	if 0 < 1 {
		g.chocoLogVelocity -= 4
		g.chocoLogX += g.chocoLogVelocity
	}

	// Check collision between player and chocoLog
	if checkCollision(
		g.playerX+playerHitboxOffsetX, g.playerY+playerHitboxOffsetY, playerHitboxW, playerHitboxH,
		g.chocoLogX+chocoLogHitboxOffsetX, g.chocoLogY+chocoLogHitboxOffsetY, chocoLogHitboxW, chocoLogHitboxH,
	) {
		log.Println("Collision detected!")
		// Example: stop chocoLog
		g.chocoLogVelocity = 0
	}
	return nil
}

func checkCollision(px, py, pw, ph, ox, oy, ow, oh float64) bool {
	return px < ox+ow &&
		px+pw > ox &&
		py < oy+oh &&
		py+ph > oy
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw Background
	bgOp := &ebiten.DrawImageOptions{}
	screen.DrawImage(bg, bgOp)

	// Draw player
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(g.playerX, g.playerY)
	screen.DrawImage(player, op)

	// Draw chocoLog
	clOp := &ebiten.DrawImageOptions{}
	clOp.GeoM.Translate(g.chocoLogX, g.chocoLogY)
	screen.DrawImage(chocoLog, clOp)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (w, h int) {
	g.playerX = screenW/2 - playerW/2
	g.playerY = screenH - playerH - 65
	g.chocoLogX = screenW - chocoLogW
	g.chocoLogY = screenH - chocoLogH - 35
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
	chocoLog, _, err = ebitenutil.NewImageFromFile("assets/chocoLog.png")
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
