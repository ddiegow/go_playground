package fibo

import "sync"

func fibo(n int, c chan int, prevfibos *map[int]int, mu *sync.Mutex) {
	// the lock/unlock is a bit messy but it works really well
	mu.Lock()
	// c is the "return"
	if n <= 1 {
		mu.Unlock()
		c <- n
	} else if _, ok := (*prevfibos)[n]; ok {
		// can't put the lock here because another process might get access to the variable and change it before we can send it back
		c <- (*prevfibos)[n]
		mu.Unlock()
	} else {
		mu.Unlock()
		myret := make(chan int)
		go fibo(n-1, myret, prevfibos, mu)
		go fibo(n-2, myret, prevfibos, mu)
		res1 := <-myret
		res2 := <-myret
		mu.Lock()
		(*prevfibos)[n] = res1 + res2
		mu.Unlock()
		c <- res1 + res2
	}
}
func Fibonacci(n int) int {
	var mu sync.Mutex // need a lock to manipulate the previous fibonacci numbers storage buffer
	result := make(chan int)
	prevfibos := make(map[int]int) // this map will store previously calculated fibonacci numbers to avoid repeating calculations. It speeds up the calculation A LOT for higher numbers
	go fibo(n, result, &prevfibos, &mu)
	return <-result
}
