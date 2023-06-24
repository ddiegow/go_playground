package gol

import (
	"fmt"
	"time"
)

// GAME OF LIFE - SERIAL IMPLEMENTATION

type position struct {
	x int
	y int
}
type cell struct {
	current bool
	next    bool
	pos     position
}

func (pos *position) init(x int, y int) {
	pos.x = x
	pos.y = y
}
func (pos *position) add(pos2 position) position {
	return position{pos.x + pos2.x, pos.y + pos2.y}
}

type board struct {
	b      [][]cell
	n      int
	left   int
	top    int
	right  int
	bottom int
}

func (b *board) init(size int) {
	start := time.Now()
	b.n = size
	b.b = make([][]cell, b.n)
	for i := range b.b {
		b.b[i] = make([]cell, b.n)
		for j := range b.b[i] {
			b.b[i][j].pos.init(j, i)
		}
	}
	b.left = -1
	b.top = size
	b.right = size
	b.bottom = -1
	fmt.Printf("Init duration: %d\n", time.Since(start))
}
func (b *board) activate(p position) {
	c := b.getCell(p)
	if c == nil {
		return
	}
	c.current = true
}

/*
Check if position is withing the limits of the board
*/
func (b *board) withinLimits(pos position) bool {
	return pos.x > b.left && pos.x < b.right && pos.y < b.top && pos.y > b.bottom
}

/*
Return pointer to cell based on position. Will return nil if position is not within the board limits.
*/
func (b *board) getCell(p position) *cell {
	if !b.withinLimits(p) {
		return nil
	}
	return &b.b[p.y][p.x]
}

func (b *board) countAround(c cell) int {
	totalCount := 0
	around := []position{
		{-1, 0},
		{-1, 1},
		{0, 1},
		{1, 1},
		{1, 0},
		{1, -1},
		{0, -1},
		{-1, -1},
	}
	for i := range around {
		a := b.getCell(c.pos.add(around[i]))
		if a == nil {
			continue
		}
		if a.current == true {
			totalCount++
		}
	}
	return totalCount
}

func (b *board) check() {
	start := time.Now()
	for i := range b.b {
		for j := range b.b[i] {
			count := b.countAround(b.b[i][j])
			// the all-important three rules!
			if b.b[i][j].current && (count < 2 || count > 3) { // if cell is alive
				b.b[i][j].next = false
			} else if !b.b[i][j].current && count == 3 {
				b.b[i][j].next = true
			} else {
				b.b[i][j].next = b.b[i][j].current
			}
		}
	}
	fmt.Printf("Check duration: %d\n", time.Since(start))
}
func (b *board) update() {
	start := time.Now()
	for i := range b.b {
		for j := range b.b[i] {
			b.b[i][j].current = b.b[i][j].next
		}
	}
	fmt.Printf("Update duration: %d\n", time.Since(start))
}

func (b *board) draw() {
	for i := b.n - 1; i >= 0; i-- {
		fmt.Printf("|")
		for j := range b.b[i] {
			if b.b[i][j].current {
				fmt.Printf(" x |")
			} else {
				fmt.Printf("   |")
			}

		}
		fmt.Printf("\n")

	}
	fmt.Printf("\n")
}
func Run(n int) {
	var b board
	b.init(n)
	b.activate(position{0, 1})
	b.activate(position{1, 1})
	b.activate(position{2, 1})
	b.activate(position{2, 2})
	b.activate(position{3, 2})
	b.activate(position{4, 2})
	b.activate(position{4, 4})
	count := 0
	average := int64(0)
	for i := 0; i < 5; i++ {
		//b.draw()
		start := time.Now()
		b.check()
		b.update()
		timeElapsed := time.Since(start)
		count++
		average += timeElapsed.Nanoseconds()
		time.Sleep(time.Millisecond * 25)
	}
	fmt.Printf("Serial average: %f\n", float64(average/int64(count)))
}
