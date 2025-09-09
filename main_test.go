package main

import (
	"testing"
)

// TestCheckCollision verifies collision detection between hitboxes.
func TestCheckCollision(t *testing.T) {
	tests := []struct {
		name     string
		p, o     Hitbox
		expected bool
	}{
		{
			name:     "Overlapping boxes",
			p:        Hitbox{0, 0, 10, 10},
			o:        Hitbox{5, 5, 10, 10},
			expected: true,
		},
		{
			name:     "Non-overlapping boxes",
			p:        Hitbox{0, 0, 10, 10},
			o:        Hitbox{20, 20, 10, 10},
			expected: false,
		},
		{
			name:     "Touching edge",
			p:        Hitbox{0, 0, 10, 10},
			o:        Hitbox{10, 0, 10, 10},
			expected: false, // touching edges is not overlapping
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := checkCollision(tt.p, tt.o)
			if got != tt.expected {
				t.Errorf("checkCollision() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// TestGameReset ensures reset restores game state.
func TestGameReset(t *testing.T) {
	g := &Game{}
	g.score = 10
	g.highScore = 20
	g.logs = append(g.logs, Log{x: 100, y: 100})
	g.gameOver = true

	g.reset()

	if g.score != 0 {
		t.Errorf("expected score=0, got %d", g.score)
	}
	if g.gameOver {
		t.Errorf("expected gameOver=false, got true")
	}
	if g.logs != nil && len(g.logs) != 0 {
		t.Errorf("expected logs to be empty, got %v", g.logs)
	}
	if g.playerX == 0 || g.playerY == 0 {
		t.Errorf("expected player position initialized, got (%f, %f)", g.playerX, g.playerY)
	}
}

// TestSpawnLog ensures spawnLog creates a log at correct position with given velocity.
func TestSpawnLog(t *testing.T) {
	g := &Game{}
	velocity := -5.0
	g.spawnLog(velocity)

	if len(g.logs) != 1 {
		t.Fatalf("expected 1 log, got %d", len(g.logs))
	}
	log := g.logs[0]
	if log.x != screenW-chocoLogW {
		t.Errorf("expected x=%d, got %f", screenW-chocoLogW, log.x)
	}
	if log.y != screenH-chocoLogH-35 {
		t.Errorf("expected y=%d, got %f", screenH-chocoLogH-35, log.y)
	}
	if log.velocity != velocity {
		t.Errorf("expected velocity=%f, got %f", velocity, log.velocity)
	}
}

// TestLayout verifies that Layout initializes positions correctly.
func TestLayout(t *testing.T) {
	g := &Game{}
	w, h := g.Layout(800, 600)

	if w != screenW || h != screenH {
		t.Errorf("expected (%d,%d), got (%d,%d)", screenW, screenH, w, h)
	}
	expectedPlayerX := screenW/2 - playerW/2
	expectedPlayerY := screenH - playerH - 65
	if g.playerX != float64(expectedPlayerX) || g.playerY != float64(expectedPlayerY) {
		t.Errorf("expected player pos (%d,%d), got (%f,%f)",
			expectedPlayerX, expectedPlayerY, g.playerX, g.playerY)
	}
	expectedLogX := screenW - chocoLogW
	expectedLogY := screenH - chocoLogH - 35
	if g.chocoLogX != float64(expectedLogX) || g.chocoLogY != float64(expectedLogY) {
		t.Errorf("expected log pos (%d,%d), got (%f,%f)",
			expectedLogX, expectedLogY, g.chocoLogX, g.chocoLogY)
	}
}
