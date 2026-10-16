package main

import (
	"fmt"
	"time"
)

func main() {
	bufferChan := make(chan bool, 5)
	for true {
		select {
		case bufferChan <- true:
			go func() {
				time.Sleep(time.Second)
				defer func() {
					_ = <-bufferChan
				}()
			}()
			fmt.Println("Intento")
		default:
			fmt.Println("Se lleno")
			time.Sleep(time.Second)
		}
	}
}
