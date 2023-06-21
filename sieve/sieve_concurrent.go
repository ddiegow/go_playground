package sieve

import (
	"fmt"
	"math"
	"sync"
)

func sieveSection(primes []int, section <-chan int, result chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for n := range section {
		isPrime := true
		for i := range primes {
			if n%primes[i] == 0 {
				isPrime = false
				break
			}
		}
		if isPrime {
			result <- n
		}
	}
}

func sieve(n int) <-chan int {
	var wg sync.WaitGroup
	initialPrimes := make([]int, 0)
	result := make(chan int)
	numbers := int(math.Sqrt(float64(n)))
	wg.Add(numbers - 1)
	firstSection := make([]bool, numbers)
	firstSection[0] = true // true indicates non-prime
	firstSection[1] = true
	for i := range firstSection {
		if i < 2 {
			continue
		}
		if !firstSection[i] {
			for j := 2 * i; j < numbers; j += i {
				firstSection[j] = true
			}
		}
	}
	for i := range firstSection {
		if firstSection[i] == false {
			initialPrimes = append(initialPrimes, i)
		}
	}
	for i := 1; i < numbers; i++ {
		section := make(chan int)
		go func(section chan int, i int) {
			go sieveSection(initialPrimes, section, result, &wg)
			for j := i * numbers; j < (i+1)*numbers; j++ {
				section <- j
			}
			close(section)
		}(section, i)

	}
	go func() {
		for i := range initialPrimes {
			result <- initialPrimes[i]
		}
	}()
	go func() {
		wg.Wait()    // wait for all threads to be done
		result <- -1 // send the termination signal
	}()
	return result
}

func ConcurrentSieve() {
	primes := sieve(1000000)
	for n := range primes {
		if n == -1 { // got the termination signal, all threads are done
			break // exit the loop
		}
		fmt.Println(n)
	}
}
