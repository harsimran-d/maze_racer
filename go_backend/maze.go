package main

import (
	"fmt"
	"log"
	"math/rand/v2"
	"time"
)

type Maze struct {
	MazeMap   [][]int
	Size      int
	Start     Position
	Finish    Position
	Finishers map[string]time.Time
}

func NewMaze(size int) *Maze {

	mazeMap := initializeMazeMap(size, size)
	carvePassage(1, 1, mazeMap)
	printMaze(mazeMap)
	return &Maze{
		MazeMap:   mazeMap,
		Size:      size,
		Finish:    Position{size - 2, size - 2},
		Start:     Position{1, 1},
		Finishers: make(map[string]time.Time),
	}
}
func printMaze(maze [][]int) {
	for _, row := range maze {
		for _, cell := range row {
			if cell == wall {
				fmt.Print("#")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
}
func initializeMazeMap(width, height int) [][]int {
	maze := make([][]int, height)
	for i := range maze {
		maze[i] = make([]int, width)
		for j := range maze[i] {
			maze[i][j] = wall
		}
	}
	return maze
}

const (
	wall  = 1
	space = 0
)

var directions = []struct {
	dx, dy int
}{
	{0, -1}, // Up
	{0, 1},  // Down
	{-1, 0}, // Left
	{1, 0},  // Right
}

func shuffleDirections() {
	rand.Shuffle(len(directions), func(i, j int) {
		directions[i], directions[j] = directions[j], directions[i]
	})
}

// CarvePassage recursively carves out the maze
func carvePassage(x, y int, maze [][]int) {
	maze[y][x] = space // Mark the current cell as space
	shuffleDirections()

	for _, dir := range directions {
		nx, ny := x+dir.dx*2, y+dir.dy*2

		if ny >= 0 && ny < len(maze) && nx >= 0 && nx < len(maze[0]) && maze[ny][nx] == wall {
			maze[y+dir.dy][x+dir.dx] = space // Knock down the wall
			carvePassage(nx, ny, maze)
		}
	}
}
func (m *Maze) NewPosition(old Position, move string) (Position, bool) {
	log.Println(old, move)
	switch move {
	case "UP":
		if old.Column-1 >= 0 && m.MazeMap[old.Row][old.Column-1] == 0 {
			return Position{Row: old.Row, Column: old.Column - 1}, true
		}
	case "DOWN":
		if old.Column+1 < m.Size && m.MazeMap[old.Row][old.Column+1] == 0 {
			return Position{Row: old.Row, Column: old.Column + 1}, true
		}
	case "LEFT":
		if old.Row-1 >= 0 && m.MazeMap[old.Row-1][old.Column] == 0 {
			return Position{Row: old.Row - 1, Column: old.Column}, true
		}
	case "RIGHT":
		if old.Row+1 < m.Size && m.MazeMap[old.Row+1][old.Column] == 0 {
			return Position{Row: old.Row + 1, Column: old.Column}, true
		}
	default:
		log.Println("invalid move: ", move)
	}
	return old, false
}
