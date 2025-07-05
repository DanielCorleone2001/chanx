package main

import (
	"fmt"
	"sync"
	"time"

	"chanx"
)

// Message represents a custom type for demonstration
type Message struct {
	ID      int
	Content string
}

func main() {
	// Example 1: Broadcasting integers
	fmt.Println("Example 1: Broadcasting integers")
	source := make(chan int)
	dest1 := make(chan int)
	dest2 := make(chan int)

	// Setup broadcast
	if err := chanx.Broadcast(source, dest1, dest2); err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	// Start receivers
	go func() {
		defer wg.Done()
		fmt.Println("Receiver 1:")
		for v := range dest1 {
			fmt.Printf("Received: %d\n", v)
		}
	}()

	go func() {
		defer wg.Done()
		fmt.Println("Receiver 2:")
		for v := range dest2 {
			fmt.Printf("Received: %d\n", v)
		}
	}()

	// Send values
	source <- 42
	source <- 100
	close(source)
	wg.Wait()

	// Example 2: Broadcasting custom type
	fmt.Println("\nExample 2: Broadcasting custom messages")
	msgSource := make(chan Message)
	msgDest1 := make(chan Message)
	msgDest2 := make(chan Message)

	if err := chanx.Broadcast(msgSource, msgDest1, msgDest2); err != nil {
		panic(err)
	}

	wg.Add(2)

	// Start message receivers
	go func() {
		defer wg.Done()
		fmt.Println("Message Receiver 1:")
		for msg := range msgDest1 {
			fmt.Printf("Got message: ID=%d, Content=%s\n", msg.ID, msg.Content)
		}
	}()

	go func() {
		defer wg.Done()
		fmt.Println("Message Receiver 2:")
		for msg := range msgDest2 {
			fmt.Printf("Got message: ID=%d, Content=%s\n", msg.ID, msg.Content)
		}
	}()

	// Send messages
	msgSource <- Message{ID: 1, Content: "Hello"}
	time.Sleep(100 * time.Millisecond) // Small delay for better output readability
	msgSource <- Message{ID: 2, Content: "World"}
	close(msgSource)
	wg.Wait()
}
