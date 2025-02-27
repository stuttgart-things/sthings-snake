package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/stuttgart-things/homerun-library"

	tl "github.com/JoelOtter/termloop"
	"github.com/charmbracelet/huh"
	"github.com/fatih/color"
)

var (
	banner = `

	█▀ ▀█▀ █░█ █ █▄░█ █▀▀ █▀ ▄▄ █▀ █▄░█ ▄▀█ █▄▀ █▀▀
	▄█ ░█░ █▀█ █ █░▀█ █▄█ ▄█ ░░ ▄█ █░▀█ █▀█ █░█ ██▄

	`
	insecure       = true
	dt             = time.Now()
	homerunAddr    = os.Getenv("HOMERUN_ADDR")
	homerunToken   = os.Getenv("HOMERUN_TOKEN")
	logPath        = os.Getenv("LOG_PATH")
	severityPreFix = os.Getenv("HOMERUN_SEVERITY_PREFIX")
)

type Coordinates struct {
	X, Y int
}

type Snake struct {
	body      []Coordinates
	direction string
	tickCount int
	growth    int
}

type Food struct {
	*tl.Entity
	placed     bool
	playerName string
}

const (
	LevelWidth  = 80
	LevelHeight = 24
)

func NewFood(playerName string) *Food {
	return &Food{
		Entity:     tl.NewEntityFromCanvas(2, 2, tl.CanvasFromString("O")),
		placed:     false,
		playerName: playerName, // ✅ Ensure this is assigned
	}
}

func (f *Food) Tick(event tl.Event) {
	// Check if food has been placed, if not, place the food
	if !f.placed {
		width, height := game.Screen().Size()
		if width > 0 && height > 0 {
			f.PlaceFood(width, height)
			f.placed = true
		}
	}
}

func (f *Food) PlaceFood(levelWidth, levelHeight int) {
	rand.Seed(time.Now().UnixNano())
	foodX := rand.Intn(LevelWidth-4) + 2
	foodY := rand.Intn(LevelHeight-4) + 2

	f.SetPosition(foodX, foodY)
}

func (f *Food) Draw(screen *tl.Screen) {
	// Draw food after it has been placed
	if f.placed {
		f.Entity.Draw(screen)
	}
}

func (f *Food) AtPosition(x, y int) bool {
	foodX, foodY := f.Position()
	if (x == foodX || x == foodX-1 || x == foodX+1) && y == foodY {
		sendNotificationToHomerun(f.playerName, score)
		return true
	}
	return false
}

func drawWalls(screen *tl.Screen) {
	// Top and bottom walls
	for x := 0; x < LevelWidth; x++ {
		screen.RenderCell(x, 0, &tl.Cell{Fg: tl.ColorWhite, Ch: '-'})             // Top wall
		screen.RenderCell(x, LevelHeight-1, &tl.Cell{Fg: tl.ColorWhite, Ch: '-'}) // Bottom wall
	}
	// Left and right walls
	for y := 0; y < LevelHeight; y++ {
		screen.RenderCell(0, y, &tl.Cell{Fg: tl.ColorWhite, Ch: '|'})            // Left wall
		screen.RenderCell(LevelWidth-1, y, &tl.Cell{Fg: tl.ColorWhite, Ch: '|'}) // Right wall
	}
}

func (snake *Snake) CollidesWithWalls() bool {
	head := snake.body[0]
	return head.X < 1 || head.Y < 1 || head.X >= LevelWidth-1 || head.Y >= LevelHeight-1
}

func (snake *Snake) CollidesWithSelf() bool {
	head := snake.body[0]
	for _, segment := range snake.body[1:] {
		if head.X == segment.X && head.Y == segment.Y {
			return true
		}
	}
	return false
}

func GameOver() {
	log.Println("Game Over!")
	playerName := showMenu() // Immediately show the menu
	startGame(playerName)    // Restart game loop with new player
}

func startGame(playerName string) {
	game = tl.NewGame()
	game.Screen().SetFps(30)

	level := tl.NewBaseLevel(tl.Cell{
		Bg: tl.ColorBlack,
		Fg: tl.ColorWhite,
		Ch: ' ',
	})

	snake := NewSnake(20, 20)
	food = NewFood(playerName)

	// Place initial food
	foodX := rand.Intn(LevelWidth-4) + 2
	foodY := rand.Intn(LevelHeight-4) + 2
	food.SetPosition(foodX, foodY)
	food.placed = true

	level.AddEntity(snake)
	level.AddEntity(food)

	// Render Banner on the Right Side
	bannerLines := []string{
		"█▀ ▀█▀ █░█ █ █▄░█ █▀▀ █▀",
		"▄█ ░█░ █▀█ █ █░▀█ █▄█ ▄█",
	}

	for i, line := range bannerLines {
		xPos := LevelWidth - len(line) + 12 // Align to the right with some margin
		bannerText := tl.NewText(xPos, i, line, tl.ColorYellow, tl.ColorBlack)
		level.AddEntity(bannerText)
	}

	// Display Score
	scoreText = tl.NewText(1, 0, "Score: 0", tl.ColorWhite, tl.ColorBlack)
	level.AddEntity(scoreText)

	game.Screen().SetLevel(level)
	game.Start()
}

func NewSnake(x, y int) *Snake {
	snake := &Snake{
		direction: "right",
		tickCount: 0,
		growth:    0,
	}
	// Initialize snake with 3 segments
	for i := 0; i < 3; i++ {
		snake.body = append(snake.body, Coordinates{X: x - i*2, Y: y})
	}
	return snake
}

func (snake *Snake) Draw(screen *tl.Screen) {
	drawWalls(screen)
	for _, segment := range snake.body {
		screen.RenderCell(segment.X, segment.Y, &tl.Cell{Fg: tl.ColorGreen, Ch: '■'})
	}
}

func showFinalScreen() {
	// Set up a blank level to display end game information
	blankLevel := tl.NewBaseLevel(tl.Cell{
		Bg: tl.ColorBlack, // Background color of the level
		Ch: ' ',           // Character to fill the screen with
	})
	game.Screen().SetLevel(blankLevel)

	// Create the final score message
	finalMessage := fmt.Sprintf("Final Score: %d", score)
	messageLength := len(finalMessage)
	startX := (LevelWidth / 2) - (messageLength / 2)
	startY := LevelHeight/2 - 1 // Positioned slightly above center for multiple lines

	finalScoreText := tl.NewText(startX, startY, finalMessage, tl.ColorWhite, tl.ColorBlack)
	blankLevel.AddEntity(finalScoreText)

	// Instructions for restarting or quitting
	restartMessage := "Press CTRL+C to QUIT"
	restartX := (LevelWidth / 2) - (len(restartMessage) / 2)
	restartY := startY + 2

	restartText := tl.NewText(restartX, restartY, restartMessage, tl.ColorWhite, tl.ColorBlack)
	blankLevel.AddEntity(restartText)

	game.Screen().Draw()
}

var score int

func (snake *Snake) Tick(event tl.Event) {
	// Handle direction change input
	if event.Type == tl.EventKey {
		switch event.Key {
		case tl.KeyArrowRight:
			if snake.direction != "left" {
				snake.direction = "right"
			}
		case tl.KeyArrowLeft:
			if snake.direction != "right" {
				snake.direction = "left"
			}
		case tl.KeyArrowUp:
			if snake.direction != "down" {
				snake.direction = "up"
			}
		case tl.KeyArrowDown:
			if snake.direction != "up" {
				snake.direction = "down"
			}
		}
	}

	// Update snake every two ticks
	snake.tickCount++
	if snake.tickCount >= 2 {
		snake.tickCount = 0
		newHead := snake.body[0]
		// Move head based on the current direction
		switch snake.direction {
		case "right":
			newHead.X += 2
		case "left":
			newHead.X -= 2
		case "up":
			newHead.Y -= 1
		case "down":
			newHead.Y += 1
		}

		// Check for food collision
		if food.AtPosition(newHead.X, newHead.Y) {
			snake.growth += 1
			food.placed = false
			score++
			scoreText.SetText(fmt.Sprintf("Score: %d", score))
			// sendNotificationToHomerun("pat")
		}

		// Grow the snake if needed
		if snake.growth > 0 {
			snake.body = append([]Coordinates{newHead}, snake.body...)
			snake.growth--
		} else {
			snake.body = append([]Coordinates{newHead}, snake.body[:len(snake.body)-1]...)
		}

		// Check for collision with walls or self
		if snake.CollidesWithWalls() || snake.CollidesWithSelf() {
			GameOver()
		}
	}
}

var food *Food
var game *tl.Game
var scoreText *tl.Text

func showMenu() string {
	// color.Yellow(logo)

	// orange := color.RGB(255, 128, 0).Println("foreground orange")
	color.RGB(255, 128, 0).AddBgRGB(0, 0, 0).Println(banner)

	// color.orange(banner)

	// Create a huh form for the menu

	var choice string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose an option").
				Options(
					huh.NewOption("Start Game", "start"),
					huh.NewOption("Exit", "exit"),
				).
				Value(&choice),
		),
	)

	// Run the form
	err := form.Run()
	if err != nil {
		fmt.Println("Error running form:", err)
		os.Exit(1)
	}

	// Handle the user's choice
	if choice == "exit" {
		fmt.Println("Goodbye!")
		os.Exit(0)
	}

	// Create a huh form for the player's name
	var playerName string
	nameForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter your name").
				Value(&playerName),
		),
	)

	// Run the name form
	err = nameForm.Run()
	if err != nil {
		fmt.Println("Error running name form:", err)
		os.Exit(1)
	}

	fmt.Printf("Hello, %s! Get ready to play!\n", playerName)
	return playerName
}

func main() {
	playerName := showMenu()
	startGame(playerName)

	game = tl.NewGame()
	game.Screen().SetFps(30)

	level := tl.NewBaseLevel(tl.Cell{
		Bg: tl.ColorBlack,
		Fg: tl.ColorWhite,
		Ch: ' ',
	})

	snake := NewSnake(20, 20)
	food = NewFood(playerName)

	// Place the first food
	foodX := rand.Intn(LevelWidth-4) + 2
	foodY := rand.Intn(LevelHeight-4) + 2
	food.SetPosition(foodX, foodY)
	food.placed = true

	level.AddEntity(snake)
	level.AddEntity(food)

	// Display Score
	scoreText = tl.NewText(1, 0, "Score: 0", tl.ColorWhite, tl.ColorBlack)
	level.AddEntity(scoreText)

	game.Screen().SetLevel(level)
	game.Start()
}

func sendNotificationToHomerun(playerName string, score int) {
	dt := time.Now()

	messageBody := homerun.Message{
		Title:           fmt.Sprintf("%s scored %d points!", playerName, score),
		Message:         fmt.Sprintf("%s scored %d points!", playerName, score),
		Severity:        severityPreFix + getSeverityLevel(score),
		Author:          playerName,
		Timestamp:       dt.Format("01-02-2006 15:04:05"),
		System:          "sthings-snake",
		Tags:            "sthings-snake,score,chaos",
		AssigneeAddress: "",
		AssigneeName:    "",
		Artifacts:       "",
		Url:             "",
	}

	rendered := homerun.RenderBody(homerun.HomeRunBodyData, messageBody)
	homerun.SendToHomerun(homerunAddr, homerunToken, []byte(rendered), insecure)

	// Log to file
	logFilePath := logPath
	logEntry := fmt.Sprintf("%s - Player %s scored %d points\n", dt.Format("2006-01-02 15:04:05"), playerName, score)

	// Append log entry to file
	f, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Error opening log file: %v", err)
		return
	}
	defer f.Close()

	if _, err := f.WriteString(logEntry); err != nil {
		log.Printf("Error writing to log file: %v", err)
	}
}

func getSeverityLevel(score int) string {
	var severity string

	switch {
	case score >= 0 && score <= 5:
		severity = "1"
	case score >= 6 && score <= 10:
		severity = "2"
	case score >= 11 && score <= 15:
		severity = "3"
	case score >= 16:
		severity = "4"
	default:
		severity = "Unknown"
	}

	return severity
}
