package main

import (
	"fmt"
	"playground/sieve"
	"time"
)

func main() {
	//playing.PlayWithChans()
	//fmt.Println(fibo.Fibonacci(50))
	//capitalize.Capitalize([]string{"Hello world", "How are you doing"})
	//input := []string{"apple", "bakery", "ladder", "zebra", "Akira", "left", "jail", "Mizuho", "Game", "Lost", "Addendum", "paradox", "cap", "destruction", "studio", "urine", "release", "bin", "flawed", "lose", "dish", "path", "cultivate", "charity", "abridge", "keep", "skip", "stun", "navy", "rich", "computer", "fool", "wait", "remain", "irony", "contraction", "braid", "skeleton", "sequence", "monstrous", "avenue", "recruit", "socialist", "head", "hill", "button", "cattle", "self", "pop", "miner", "bargain", "advocate", "identity", "package", "slot", "question", "extraterrestrial", "domination", "sheet", "establish", "voter"}
	//fmt.Println(mergesort.Sort(input))
	start := time.Now()
	//sieve.Run() // THIS TAKES FOREVER WITH LARGE NUMBERS
	fmt.Printf("The function took %d to run\n", time.Since(start))
	start = time.Now()
	sieve.ConcurrentSieve()
	fmt.Printf("The function took %d to run\n", time.Since(start))
}
