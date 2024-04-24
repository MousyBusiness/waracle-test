package main

import (
	"fmt"
	"time"
)

func main() {
	for {
		fmt.Println("Running", time.Now())
		time.Sleep(time.Second)
	}
}
