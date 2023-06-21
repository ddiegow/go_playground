package playing

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func child(i int, c chan string) {
	n := rand.Intn(5)
	time.Sleep(time.Second * time.Duration(n))
	c <- "Hello from thread number " + strconv.Itoa(i) + ". I waited for " + strconv.Itoa(n) + " seconds."
}
func PlayWithChans() {
	c := make(chan string)

	for i := 0; i < 10; i++ {
		go child(i, c)
	}
	for i := 0; i < 10; i++ {
		fmt.Println(<-c)
	}
}
