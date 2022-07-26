package main

import (
	"flag"
	"fmt"
	"os/exec"
	"runtime"
)

func main() {
	os := runtime.GOOS
	println(os)

	var port int
	flag.IntVar(&port, "p", 80, "send port to scan")
	scanAll := flag.Bool("all", false, "scan all ports")
	flag.Parse()

	println(port)
	println(*scanAll)
}

func windows() {

	cmd := exec.Command("/bin/sh", "./test.sh")

	out, _ := cmd.Output()

	output := string(out)
	fmt.Println(output)
}
