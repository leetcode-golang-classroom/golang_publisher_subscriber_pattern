package main

import (
	"fmt"
	"sync"
)

type PubSub[T any] struct {
	subscribers []chan T
	mu          sync.RWMutex
	closed      bool
}

func NewPubSub[T any]() *PubSub[T] {
	return &PubSub[T]{
		mu: sync.RWMutex{},
	}
}

func (s *PubSub[T]) Subscribe() <-chan T {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return nil
	}

	r := make(chan T)
	s.subscribers = append(s.subscribers, r)
	return r
}

func (s *PubSub[T]) Publish(value T) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.closed {
		return
	}

	for _, ch := range s.subscribers {
		ch <- value
	}
}
func (s *PubSub[T]) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return
	}
	for _, ch := range s.subscribers {
		close(ch)
	}
	s.closed = true
}
func main() {
	pubSub := NewPubSub[string]()

	wg := sync.WaitGroup{}
	// sub1
	s1 := pubSub.Subscribe()
	wg.Add(1)
	go func() {
		for {
			select {
			case val, ok := <-s1:
				if !ok {
					fmt.Println("sub 1, exiting")
					wg.Done()
					return
				}
				fmt.Println("sub1, value ", val)
			}
		}
	}()
	// sub2
	s2 := pubSub.Subscribe()
	wg.Add(1)
	go func() {
		for val := range s2 {
			fmt.Println("sub 2, value", val)
		}
		wg.Done()
		fmt.Println("sub 2, exiting")
	}()

	// wait
	pubSub.Publish("one")
	pubSub.Publish("two")
	pubSub.Publish("three")

	pubSub.Close()

	wg.Wait()

	fmt.Println("complete")
}
