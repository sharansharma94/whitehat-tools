package main

import (
	"context"
	"fmt"
	"net"
	"sort"
	"time"
)

const URL = "192.168.1.9"

func main() {

	portsToScan := 65535

	ports := make(chan int, 1000)
	results := make(chan int)
	var openports []int

	for i := 0; i < cap(ports); i++ {
		go worker(ports, results)
	}

	go func() {
		for i := 1; i <= portsToScan; i++ {
			ports <- i
		}
	}()

	for i := 0; i < portsToScan; i++ {
		port := <-results
		if port != 0 {
			openports = append(openports, port)
		}
	}

	close(ports)
	close(results)
	sort.Ints(openports)
	for _, port := range openports {
		fmt.Printf("%d open\n", port)
	}
}

func worker(ports, results chan int) {
	for p := range ports {
		address := fmt.Sprintf("%s:%d", URL, p)
		fmt.Println(address)

		var d net.Dialer
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		conn, err := d.DialContext(ctx, "tcp", address)
		if err != nil {
			results <- 0
			continue
		}
		conn.Close()
		results <- p
	}
}
