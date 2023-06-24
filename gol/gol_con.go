package gol

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

// GAME OF LIFE - SIMPLE CONCURRENT IMPLEMENTATION
// This implementation is about 4 times faster than the serial one (on a moderately powerful computer)
// TODO: compare implementing the cells around as pointers idea vs the old calculation one. Could be that the current double access to memory required by the pointers makes it somewhat slower than the original

type positionCon struct {
	x int
	y int
}
type cellCon struct {
	current bool
	next    bool
	pos     positionCon
	around  []*cellCon
}

func (pos *positionCon) initCon(x int, y int) {
	pos.x = x
	pos.y = y
}
func (pos *positionCon) addCon(pos2 positionCon) positionCon {
	return positionCon{pos.x + pos2.x, pos.y + pos2.y}
}

type boardCon struct {
	b      [][]cellCon
	n      int
	left   int
	top    int
	right  int
	bottom int
}

func (b *boardCon) initCon(size int) {
	start := time.Now()
	b.n = size
	b.b = make([][]cellCon, b.n)
	b.left = -1
	b.top = -1
	b.right = b.n
	b.bottom = b.n
	around := []positionCon{
		{-1, 0},
		{-1, 1},
		{0, 1},
		{1, 1},
		{1, 0},
		{1, -1},
		{0, -1},
		{-1, -1},
	}

	for i := range b.b {
		b.b[i] = make([]cellCon, b.n)
		for j := range b.b[i] {
			b.b[i][j].pos.initCon(j, i)
			for k := range around {
				a := b.getCellCon(b.b[i][j].pos.addCon(around[k]))
				if a == nil {
					continue
				}
				b.b[i][j].around = append(b.b[i][j].around, a)
			}
		}
	}
	fmt.Printf("Init duration: %d\n", time.Since(start))
}
func (b *boardCon) activateCon(p positionCon) {
	c := b.getCellCon(p)
	if c == nil {
		return
	}
	c.current = true
}

/*
Check if position is within the limits of the board
*/
func (b *boardCon) withinLimitsCon(pos positionCon) bool {
	return pos.x > b.left && pos.x < b.right && pos.y < b.top && pos.y > b.bottom
}

/*
Return pointer to cell based on position. Will return nil if position is not within the board limits.
*/
func (b *boardCon) getCellCon(p positionCon) *cellCon {
	if !b.withinLimitsCon(p) {
		return nil
	}
	return &b.b[p.y][p.x]
}

func (b *boardCon) countAround(c cellCon) int {
	totalCount := 0
	for i := range c.around {
		if c.around[i].current == true {
			totalCount++
		}
	}
	return totalCount
}
func (b *boardCon) checkRow(i int) {
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
func (b *boardCon) checkCon() {
	start := time.Now()
	sem := make(chan int, runtime.NumCPU()) // basic semaphore to use as many threads are there are cores
	var wg sync.WaitGroup
	for i := range b.b {
		sem <- 1 // wait for the semaphore
		wg.Add(1)
		go func(rowNum int, w *sync.WaitGroup) {
			b.checkRow(rowNum)
			w.Done()
			<-sem // signal the semaphore
		}(i, &wg)
	}
	wg.Wait() // wait for all threads to finish
	fmt.Printf("Check duration: %d\n", time.Since(start))
}
func (b *boardCon) updateCon() {
	sem := make(chan int, runtime.NumCPU())
	var wg sync.WaitGroup
	start := time.Now()
	for i := range b.b {
		sem <- 1 // wait for the semaphore
		wg.Add(1)
		go func(col int, w *sync.WaitGroup) {
			for j := range b.b[i] {
				b.b[i][j].current = b.b[i][j].next
			}
			<-sem
			w.Done()
		}(i, &wg)
	}
	wg.Wait()
	fmt.Printf("Update duration: %d\n", time.Since(start))
}

func (b *boardCon) drawCon() {
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
func RunCon(n int) {
	var b boardCon
	b.initCon(n)
	b.activateCon(positionCon{0, 1})
	b.activateCon(positionCon{1, 1})
	b.activateCon(positionCon{2, 1})
	b.activateCon(positionCon{2, 2})
	b.activateCon(positionCon{3, 2})
	b.activateCon(positionCon{4, 2})
	b.activateCon(positionCon{4, 4})
	average := int64(0)
	count := 0
	for i := 0; i < 5; i++ {
		//fmt.Print("\033[H\033[2J")
		//b.drawCon()

		start := time.Now()
		b.checkCon()
		b.updateCon()
		timeElapsed := time.Since(start)
		count++
		average += timeElapsed.Nanoseconds()
		time.Sleep(time.Millisecond * 25)
	}
	fmt.Printf("Concurrent average: %f\n", float64(average/int64(count)))
}
