package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Chopstick struct {
	id int
	mutex sync.Mutex
}
type Host struct {
	filEating []int
}
type Philosopher struct {
	id                  int
	leftFork, rightFork *Chopstick
	host *Host
	numEat int
}

const (
	n = 5
	simultaneo = 2
	numEatingTotal = 2
)

func (host *Host) Manage(requestChan <-chan *Philosopher, finishChan <-chan *Philosopher, isFinishedChan <-chan bool, respuestaChan chan<- bool) {
    for {
        select {
        case <-isFinishedChan:
            // terminamos codigo
            return
        
        // Sacamos Filosofo
        case filo := <-finishChan:
            fmt.Printf("philosopher %d signaled finish eating\n", filo.id)
            buscadorId := -1
            for id, e := range host.filEating {
                if e == filo.id {
                    buscadorId = id
                    break
                }
            }
            host.filEating = append(host.filEating[:buscadorId], host.filEating[buscadorId+1:]...)
        
        // Manejamos request
        case filo := <-requestChan:
            if len(host.filEating) >= numEatingTotal {
                respuestaChan <- false
                continue
            }
            // Reservamos palillos
            filo.leftFork.mutex.Lock()
            filo.rightFork.mutex.Lock()

            respuestaChan <- true
            host.filEating = append(host.filEating, filo.id)

        }

    }
}
func (filo *Philosopher) Eat(requestChan chan<- *Philosopher, finishChan chan<- *Philosopher, respuestaChan <-chan bool, wg *sync.WaitGroup) {
    defer wg.Done()
    for filo.numEat < numEatingTotal {

        // Pedimos comer
        response := AllowEating(filo, requestChan, respuestaChan)
        // Saltamos si no está permitido
        if !response {
            continue
        }
        filo.numEat++
        fmt.Printf("Filosofo %d comiendo\n", filo.id)
        time.Sleep(time.Millisecond * time.Duration(rand.Intn(500)))
        fmt.Printf("Filosofo %d termino de comer\n", filo.id)

        // Soltamos palillos
        filo.rightFork.mutex.Unlock()
        filo.leftFork.mutex.Unlock()
        finishChan <- filo
    }
}

// https://codereview.stackexchange.com/questions/278596/golang-implementation-of-dining-philosophers-variant/278605#278605
func intSlicesEqual(a, b []int) bool {
    if len(a) != len(b) {
        return false
    }
    for i, v := range a {
        if v != b[i] {
            return false
        }
    }
    return true
}

func AllowEating(filo *Philosopher, requestChan chan<- *Philosopher, respuestaChan <-chan bool) (bool){
    requestChan <- filo
    res := <-respuestaChan
    return res

}

func main() {
	palillos := make([]*Chopstick, n)
	for i := 0; i < n; i++ {
		palillos[i] = &Chopstick{id: i}
	}

	host := Host{ filEating: make([]int, 0, numEatingTotal),}

	philos := make([]*Philosopher, n)
    for i := 0; i < n; i++ {
        philos[i] = &Philosopher{id: i + 1, numEat: 0, rightFork: palillos[i], leftFork: palillos[(i+1)%5], host: &host}
    }

	requestChan := make(chan *Philosopher)
    finishChan := make(chan *Philosopher)

    respuestaChan := make(chan bool)
    isFinishedChan := make(chan bool)

	var wg sync.WaitGroup
	wg.Add(n)

	go host.Manage(requestChan, finishChan, isFinishedChan, respuestaChan)
	// agregamos filosofos
    for i := 0; i < n; i++ {
        go philos[i].Eat(requestChan, finishChan, respuestaChan, &wg)
    }
    wg.Wait()
    // Avisamos host que terminamos
    isFinishedChan <- true
}
