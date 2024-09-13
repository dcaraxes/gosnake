package main

import (
	"fmt"
	"math/rand"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type direction int

const (
	up direction = iota
	down
	left
	right
)

const (
	boardWidth  = 20
	boardHeight = 10
)

type point struct {
	x, y int
}

type model struct {
	snake    []point
	dir      direction
	food     point
	gameOver bool
}

func initialModel() model {
	s := []point{{5, 5}, {4, 5}, {3, 5}}
	return model{
		snake:    s,
		dir:      right,
		food:     randomFood(s),
		gameOver: false,
	}
}

func randomFood(snake []point) point {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	var p point
	for {
		p = point{rand.Intn(boardWidth), rand.Intn(boardHeight)}
		collides := false
		for _, sp := range snake {
			if sp == p {
				collides = true
				break
			}
		}
		if !collides {
			break
		}
	}
	return p
}

func (m model) Init() tea.Cmd {
	return tea.Tick(time.Millisecond*200, func(t time.Time) tea.Msg {
		return t
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.gameOver {
		return m, tea.Quit
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "up":
			if m.dir != down {
				m.dir = up
			}
		case "down":
			if m.dir != up {
				m.dir = down
			}
		case "left":
			if m.dir != right {
				m.dir = left
			}
		case "right":
			if m.dir != left {
				m.dir = right
			}
		}
	case time.Time:
		m = m.moveSnake()
		return m, tea.Tick(time.Millisecond*200, func(t time.Time) tea.Msg {
			return t
		})
	}
	return m, nil
}

func (m model) moveSnake() model {
	head := m.snake[0]
	newHead := head

	switch m.dir {
	case up:
		newHead.y--
	case down:
		newHead.y++
	case left:
		newHead.x--
	case right:
		newHead.x++
	}

	if newHead.x < 0 || newHead.y < 0 || newHead.x >= boardWidth || newHead.y >= boardHeight {
		m.gameOver = true
		return m
	}

	for _, p := range m.snake {
		if p == newHead {
			m.gameOver = true
			return m
		}
	}

	m.snake = append([]point{newHead}, m.snake...)

	if newHead == m.food {
		m.food = randomFood(m.snake)
	} else {
		m.snake = m.snake[:len(m.snake)-1]
	}

	return m
}

func (m model) View() string {
	if m.gameOver {
		return "Game Over! Press q to quit.\n"
	}

	board := make([][]rune, boardHeight)
	for i := range board {
		board[i] = make([]rune, boardWidth)
		for j := range board[i] {
			board[i][j] = ' '
		}
	}

	for _, p := range m.snake {
		board[p.y][p.x] = 'o'
	}

	board[m.food.y][m.food.x] = 'x'

	var result string
	topBottomBorder := "*" + repeatString("*", boardWidth) + "*\n"
	result += topBottomBorder

	for _, row := range board {
		result += "*"
		for _, cell := range row {
			result += string(cell)
		}
		result += "*\n"
	}

	result += topBottomBorder
	return result + "Use arrow keys to move. Press q to quit.\n"
}

func repeatString(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
	}
}
