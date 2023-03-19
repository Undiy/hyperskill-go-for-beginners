package main

import (
	"fmt"
	"math/rand"
	"time"
)

const CELL_ALIVE = 'O'
const CELL_DEAD = ' '

func initWorld(size int, r rand.Rand) [][]rune {
	world := make([][]rune, size)

	for i := 0; i < size; i++ {
		world[i] = make([]rune, size)
		for j := 0; j < size; j++ {
			if r.Intn(2) == 1 {
				world[i][j] = CELL_ALIVE
			} else {
				world[i][j] = CELL_DEAD
			}
		}
	}

	return world
}

func printWorld(world [][]rune, step int) {
	fmt.Printf("Generation #%d\n", step)
	fmt.Printf("Alive: %d\n", countLiveCells(world))
	for _, line := range world {
		fmt.Println(string(line))
	}
}

func countLiveNeighbors(world [][]rune, y, x int) int {
	size := len(world)
	count := 0
	for _, cell := range []rune{
		world[(y-1+size)%size][(x)%size],        // N
		world[(y-1+size)%size][(x+1)%size],      // NE
		world[y][(x+1)%size],                    // E
		world[(y+1)%size][(x+1)%size],           // SE
		world[(y+1)%size][x],                    // S
		world[(y+1)%size][(x-1+size)%size],      // SW
		world[y][(x-1+size)%size],               // W
		world[(y-1+size)%size][(x-1+size)%size], // NW
	} {
		if cell == CELL_ALIVE {
			count++
		}
	}

	return count
}

func nextGeneration(world [][]rune) [][]rune {
	size := len(world)
	newWorld := make([][]rune, size)
	for i := 0; i < size; i++ {
		newWorld[i] = make([]rune, size)
		for j := 0; j < size; j++ {
			liveNeighbors := countLiveNeighbors(world, i, j)
			if liveNeighbors == 3 || (liveNeighbors == 2 && world[i][j] == CELL_ALIVE) {
				newWorld[i][j] = CELL_ALIVE
			} else {
				newWorld[i][j] = CELL_DEAD
			}
		}
	}
	return newWorld
}

func countLiveCells(world [][]rune) int {
	count := 0
	for _, line := range world {
		for _, cell := range line {
			if cell == CELL_ALIVE {
				count++
			}
		}
	}
	return count
}

func main() {
	var size int
	fmt.Scan(&size)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	step := 0

	world := initWorld(size, *r)
	for step < 10 {
		printWorld(world, step)
		step++
		world = nextGeneration(world)
	}

	printWorld(world, step)
}
