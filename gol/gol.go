package gol

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
	b.n = size
	b.b = make([][]cell, b.n)
	for i := range b.b {
		b.b[i] = make([]cell, b.n)
	}
	b.left = -1
	b.top = size
	b.right = size
	b.bottom = -1
}
func (b *board) withinLimits(pos position) bool {
	return pos.x > b.left && pos.x < b.right && pos.y < b.top && pos.y > b.bottom
}
func getCell(board [][]cell, position2 position) cell {

}
func run() {

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

}
