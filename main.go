package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	mrand "math/rand"
	"os"
	"time"

	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

var (
	player           *ebiten.Image
	bg               *ebiten.Image
	chocoLog         *ebiten.Image
	jumpHeight       float64
	landingHeight    float64
	jumpUpComplete   bool
	jumpDownComplete bool
	jumpComplete     bool    = true
	logSpeed         float64 = -4
	myFont           font.Face
)

const (
	screenW, screenH                             = 1000, 667
	playerW, playerH                             = 128, 128
	chocoLogW, chocoLogH                         = 128, 128
	playerHitboxW, playerHitboxH                 = 88, 100
	playerHitboxOffsetX, playerHitboxOffsetY     = 40, 25
	chocoLogHitboxW, chocoLogHitboxH             = 40, 60
	chocoLogHitboxOffsetX, chocoLogHitboxOffsetY = 30, 50
	normalFontSize                               = 24
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

	logs       []Log
	spawnTimer int
	spawnSpeed int
	gameOver   bool
	score      int
	highScore  int
}

type Log struct {
	x, y     float64
	velocity float64
}

type Hitbox struct {
	x, y, w, h float64
}

func (g *Game) Update() error {

	if g.gameOver {
		if inpututil.IsKeyJustPressed(ebiten.KeyR) {
			g.reset()
		}
		return nil
	}

	// Handle jump.
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) && jumpComplete {
		jumpHeight = -150
		g.playerVelocity = 0.5
		jumpUpComplete = false
		jumpDownComplete = false
		jumpComplete = false
	}

	if !jumpUpComplete && jumpHeight < g.playerVelocity {
		g.playerVelocity -= 4
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

	// Spawn logs
	g.spawnTimer--
	if g.spawnTimer <= 0 {
		temp, err := rand.Int(rand.Reader, big.NewInt(60))
		if err != nil {
			panic(err)
		}
		num := temp.Int64()

		logSpeed -= 0.5
		g.spawnLog(logSpeed)

		g.spawnTimer = 120 + int(num) // Spawns logs every 2-3 seconds
	}

	// Update logs
	for i := range g.logs {
		g.logs[i].x += g.logs[i].velocity
		if g.logs[i].x+chocoLogW < g.playerX && g.logs[i].x+chocoLogW > g.playerX+g.logs[i].velocity {
			g.score++
			if g.score > g.highScore {
				g.highScore = g.score
			}
		}
	}

	// Remove logs off screen
	newLogs := g.logs[:0]
	for _, logObj := range g.logs {
		if logObj.x+chocoLogW > 0 {
			newLogs = append(newLogs, logObj)
		}
	}

	g.logs = newLogs

	// Check collisions
	for _, logObj := range g.logs {
		if checkCollision(
			Hitbox{
				x: g.playerX + playerHitboxOffsetX,
				y: g.playerY + playerHitboxOffsetY,
				w: playerHitboxW,
				h: playerHitboxH,
			},
			Hitbox{
				x: logObj.x + chocoLogHitboxOffsetX,
				y: logObj.y + chocoLogHitboxOffsetY,
				w: chocoLogHitboxW,
				h: chocoLogHitboxH,
			},
		) {
			log.Println("Collision detected!")
			jumpComplete = true
			g.gameOver = true
		}
	}
	return nil
}

func (g *Game) spawnLog(velocity float64) {
	g.logs = append(g.logs, Log{
		x:        screenW - chocoLogW,
		y:        screenH - chocoLogH - 35,
		velocity: velocity,
	})
}

func (g *Game) reset() {
	temp, err := rand.Int(rand.Reader, big.NewInt(60))
	if err != nil {
		panic(err)
	}
	num := temp.Int64()
	g.logs = nil
	g.spawnTimer = 60 + int(num)
	g.playerVelocity = 0
	g.playerX = screenW/2 - playerW/2
	g.playerY = screenH - playerH - 65
	g.gameOver = false
	g.score = 0
	logSpeed = -4

}

func checkCollision(p Hitbox, o Hitbox) bool {
	return p.x < o.x+o.w &&
		p.x+p.w > o.x &&
		p.y < o.y+o.h &&
		p.y+p.h > o.y
}

func loadFont() font.Face {
	ttfBytes, err := os.ReadFile("assets/comicJungle/comicJungle.ttf")
	if err != nil {
		log.Fatal(err)
	}
	tt, err := opentype.Parse(ttfBytes)
	if err != nil {
		log.Fatal(err)
	}
	const dpi = 72
	face, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    36,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	return face
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw Background
	bgOp := &ebiten.DrawImageOptions{}
	screen.DrawImage(bg, bgOp)

	// Logs
	for _, logObj := range g.logs {
		clOp := &ebiten.DrawImageOptions{}
		clOp.GeoM.Translate(logObj.x, logObj.y)
		screen.DrawImage(chocoLog, clOp)
	}

	if g.gameOver {
		ebitenutil.DebugPrint(screen, "GAME OVER! Press R to restart")
	}

	// Draw Scores
	scoreText := fmt.Sprintf("SCORE: %d", g.score)
	highScoreText := fmt.Sprintf("HIGH SCORE: %d", g.highScore)
	text.Draw(screen, scoreText, myFont, 20, 40, colornames.Black)
	text.Draw(screen, highScoreText, myFont, 20, 80, colornames.Black)

	// Draw player
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(g.playerX, g.playerY)
	screen.DrawImage(player, op)
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
	myFont = loadFont()

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
	mrand.Seed(time.Now().UnixNano())
	ebiten.SetWindowSize(screenW, screenH)
	ebiten.SetWindowTitle("ChocoJump")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
