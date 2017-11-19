package main

import (
	"fmt"
	"sync"
)

type worker struct {
	source chan interface{}
	quit chan struct{}
}

type threadSafeSlice struct {
	sync.Mutex
	workers []*worker
}

func (w *worker) Start() {
	w.source = make(chan interface{}, 10) // some buffer size to avoid blocking
	go func() {
		for {
			select {
			case msg := <-w.source:
				fmt.Println(msg)
			case <-w.quit:
				return
			}
		}
	}()
}

func (slice *threadSafeSlice) Push(w *worker) {
	slice.Lock()
	defer slice.Unlock()

	slice.workers = append(slice.workers, w)
}

func (slice *threadSafeSlice) Iter(routine func(*worker)) {
	slice.Lock()
	defer slice.Unlock()

	for _, worker := range slice.workers {
		routine(worker)
	}
}

func main() {
	sourceData := make(chan int)

	wr := &worker{}
	sl := &threadSafeSlice{}
	wr.Start()
	sl.Push(wr)
	go func() {
		for {
			msg := <- sourceData
			sl.Iter(func(w *worker) { w.source <- msg })
		}
	}()

	for i := 0; i < 100; i++ {
		sourceData <- i
	}
}