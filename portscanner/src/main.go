package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"sort"
	"strings"
	"time"
)

const URL = "192.168.1.9"

func main() {

	port := flag.Int("p", 80, "send port to scan")
	scanAll := flag.Bool("all", false, "scan all ports")
	addr := flag.String("a", "", "Provide address")
	flag.Parse()
	// addr := os.Args[len(os.Args)-1]

	fmt.Println(*addr)

	if *addr == "" {
		fmt.Println("Please provide address")
		return
	}

	if *scanAll {
		scanAllPorts(65535, *addr)
	} else {
		if *port < 1 {
			fmt.Println("not a valid port")
			return
		}
		scanPort(*port, *addr)
	}

}

func scanPort(port int, addr string) {

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		fmt.Println("port is closed or filterd")
		return
	}
	conn.Close()
	fmt.Printf("\n")

	fmt.Printf("%d is open\n", port)
}

func scanAllPorts(portsToScan int, addr string) {
	ports := make(chan int, 1000)
	results := make(chan int)
	var openports []int

	for i := 0; i < cap(ports); i++ {
		go worker(ports, results, addr)
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
	fmt.Printf("\n")

	for _, port := range openports {
		fmt.Printf("%d open\n", port)
	}
}

func worker(ports chan int, results chan int, addr string) {
	count := 1
	for p := range ports {
		address := fmt.Sprintf("%s:%d", addr, p)
		count += 1

		loading := "."

		loading = strings.Repeat(".", count%50)
		fmt.Printf("\r %s", loading)

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
