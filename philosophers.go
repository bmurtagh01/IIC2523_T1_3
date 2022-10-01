package main

import (
	"fmt"
	"sync"
)

type Chopstick struct {
	id int
	sync.Mutex
}
type Philosopher struct {
	id                  int
	leftFork, rightFork *Chopstick
}
type Host struct {
	philEating int
}

func Eat(wg *sync.WaitGroup) {
	contador = contador + 1
	wg.Done()
}

func AllowEating(wg *sync.WaitGroup) {
	contador = contador - 1
	wg.Done()
}

func main() {
	const n = 5
	const simultaneo = 2
	const numEating = 2

	palillos := make([]*Chopstick, n)
	for i := 0; i < n; i++ {
		palillos[i] = &Chopstick{id: i}
	}

	pedir := make(chan []int, n)

	var wg sync.WaitGroup
	wg.Add(n)

	for i := 0; i < n; i++ {
		go Adder1(&wg)
		go Adder2(&wg)
	}

	wg.Wait()
	fmt.Println("Contador deberia ser 0 pero es:", contador)
}
