package main

import (
	"fmt"
	"runtime"
)

func main() {
	fmt.Println("Classifieds Watcher is starting...")
	fmt.Printf("Running on: %s/%s\n", runtime.GOOS, runtime.GOARCH)
}
