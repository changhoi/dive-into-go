package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func main() {
	ctx := context.Background()
	parent, parentCancel := context.WithCancel(ctx)
	child, childCancel := context.WithCancel(parent)

	if len(os.Args) != 3 {
		fmt.Printf("usage: %s <parent delay> <child delay>\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}

	parentDelay, _ := strconv.Atoi(os.Args[1])
	childDelay, _ := strconv.Atoi(os.Args[2])

	go func() {
		time.Sleep(time.Duration(parentDelay) * time.Second)
		parentCancel()
	}()

	go func() {
		time.Sleep(time.Duration(childDelay) * time.Second)
		childCancel()
	}()

	select {
	case <-child.Done():
		fmt.Println("Child closed!:", child.Err())
		fmt.Println("Parent not closed!:", parent.Err())
	case <-parent.Done():
		fmt.Println("Parent closed!:", parent.Err())
		fmt.Println("Child closed!:", child.Err())
	}
}
