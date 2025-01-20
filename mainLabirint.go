package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type Point struct {
	x, y int
}

type Cell struct {
	x, y    int
	visited bool
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	// Задать размеры лабиринта
	dimensions, _ := reader.ReadString('\n')
	dimensions = strings.TrimSpace(dimensions)
	dimParts := strings.Split(dimensions, " ")
	n, _ := strconv.Atoi(dimParts[0])
	m, _ := strconv.Atoi(dimParts[1])

	// Генерация лабиринта
	maze := generateMaze(n, m)

	// Чтение координат старта и финиша
	startEnd, _ := reader.ReadString('\n')
	startEnd = strings.TrimSpace(startEnd)
	coords := strings.Split(startEnd, " ")
	startX, _ := strconv.Atoi(coords[0])
	startY, _ := strconv.Atoi(coords[1])
	endX, _ := strconv.Atoi(coords[2])
	endY, _ := strconv.Atoi(coords[3])

	// Поиск пути
	path, err := bfs(maze, Point{startX, startY}, Point{endX, endY})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Вывод результата
	for _, p := range path {
		fmt.Printf("%d %d\n", p.x, p.y)
	}
	fmt.Println(".")
}

func generateMaze(n, m int) [][]int {
	// Инициализация
	maze := make([][]int, n)
	for i := range maze {
		maze[i] = make([]int, m)
		for j := range maze[i] {
			maze[i][j] = 0 // Сначала все стены
		}
	}

	// Алгоритм генерации лабиринта
	cells := make([]Cell, 0)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			cells = append(cells, Cell{x: i, y: j})
		}
	}

	// Открываем первую клетку
	rand.Seed(time.Now().UnixNano())
	startCell := cells[rand.Intn(len(cells))]
	maze[startCell.x][startCell.y] = 1
	cells = removeCell(cells, startCell)

	// Основной цикл генерации
	for len(cells) > 0 {
		cell := cells[rand.Intn(len(cells))]
		neighbors := getNeighbors(cell, maze)

		if len(neighbors) > 0 {
			// Выбираем случайного соседа
			neighbor := neighbors[rand.Intn(len(neighbors))]
			// Убираем стену между соседями
			maze[(cell.x+neighbor.x)/2][(cell.y+neighbor.y)/2] = 1
			maze[neighbor.x][neighbor.y] = 1
			cells = removeCell(cells, neighbor)
		}
	}

	return maze
}

func getNeighbors(cell Cell, maze [][]int) []Point {
	var neighbors []Point
	directions := []Point{{2, 0}, {-2, 0}, {0, 2}, {0, -2}}

	for _, dir := range directions {
		newX, newY := cell.x+dir.x, cell.y+dir.y
		if newX >= 0 && newX < len(maze) && newY >= 0 && newY < len(maze[0]) {
			if maze[newX][newY] == 0 {
				neighbors = append(neighbors, Point{newX, newY})
			}
		}
	}
	return neighbors
}

func removeCell(cells []Cell, cell Cell) []Cell {
	for i, c := range cells {
		if c.x == cell.x && c.y == cell.y {
			return append(cells[:i], cells[i+1:]...)
		}
	}
	return cells
}

func bfs(maze [][]int, start, end Point) ([]Point, error) {
	n := len(maze)
	m := len(maze[0])
	if maze[start.x][start.y] == 0 || maze[end.x][end.y] == 0 {
		return nil, fmt.Errorf("старт или финиш на стене")
	}

	directions := []Point{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}
	visited := make([][]bool, n)
	for i := range visited {
		visited[i] = make([]bool, m)
	}
	visited[start.x][start.y] = true

	queue := []Point{start}
	parent := make(map[Point]Point)

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if current == end {
			return reconstructPath(parent, start, end), nil
		}

		for _, dir := range directions {
			next := Point{current.x + dir.x, current.y + dir.y}
			if isValid(maze, next, visited) {
				visited[next.x][next.y] = true
				queue = append(queue, next)
				parent[next] = current
			}
		}
	}

	return nil, fmt.Errorf("путь не найден")
}

func isValid(maze [][]int, point Point, visited [][]bool) bool {
	n := len(maze)
	m := len(maze[0])
	return point.x >= 0 && point.x < n && point.y >= 0 && point.y < m && maze[point.x][point.y] != 0 && !visited[point.x][point.y]
}

func reconstructPath(parent map[Point]Point, start, end Point) []Point {
	var path []Point
	for at := end; at != start; at = parent[at] {
		path = append(path, at)
	}
	path = append(path, start)
	reverse(path)
	return path
}

func reverse(path []Point) {
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
}
