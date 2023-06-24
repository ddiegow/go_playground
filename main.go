package main

import "playground/gol"

func main() {
	for i := 0; i < 10; i++ {
		gol.Run(100)
	}
	for i := 0; i < 10; i++ {
		gol.RunCon(100)
	}

}
