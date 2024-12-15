package main

import "log"

type Maze struct {
	MazeMap [][]int
	Size    int
}

func NewMaze(size int) *Maze {
	mazeMap := getMazeMap(size)
	return &Maze{
		MazeMap: mazeMap,
		Size:    size,
	}
}
func getMazeMap(size int) [][]int {
	maze := make([][]int, size)
	for i := range maze {
		maze[i] = make([]int, size)
	}
	return maze
}

func (m *Maze) NewPosition(old Position, move string) (Position, bool) {
	log.Println(old, move)
	switch move {
	case "up":
		if old.Column-1 >= 0 && m.MazeMap[old.Row][old.Column-1] == 0 {
			return Position{Row: old.Row, Column: old.Column - 1}, true
		}
	case "down":
		if old.Column+1 < m.Size && m.MazeMap[old.Row][old.Column+1] == 0 {
			return Position{Row: old.Row, Column: old.Column + 1}, true
		}
	case "left":
		if old.Row-1 >= 0 && m.MazeMap[old.Row-1][old.Column] == 0 {
			return Position{Row: old.Row - 1, Column: old.Column}, true
		}
	case "right":
		if old.Row+1 < m.Size && m.MazeMap[old.Row+1][old.Column] == 0 {
			return Position{Row: old.Row + 1, Column: old.Column}, true
		}
	default:
		log.Println("invalid move: ", move)
	}
	return old, false
}
