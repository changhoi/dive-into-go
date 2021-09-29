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
	if len(os.Args) != 2 {
		fmt.Printf("usage %s <deadline>\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}

	t, _ := strconv.Atoi(os.Args[1])
	deadline := time.Now().Add(time.Duration(t) * time.Second)

	ctx := context.Background()
	ctx, cancel := context.WithDeadline(ctx, deadline)
	go func() {
		select {
		case <-time.After(time.Duration(3) * time.Second):
			fmt.Println("TOO LONG DEADLINE!")
			cancel()
		}
	}()

	select {
	case <-ctx.Done():
		fmt.Println("closed!:", ctx.Err())
	}
}
