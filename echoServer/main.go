package main

import (
	"bufio"
	"io"
	"log"
	"net"
)

func main() {

	listner, err := net.Listen("tcp", ":9090")

	if err != nil {
		log.Fatalf("%e", err)
	}

	for {

		conn, err := listner.Accept()

		if err != nil {
			log.Fatalf("%e", err)
		}

		log.Println("new connection")
		go echoCopy(conn)

	}

}

// echo with  allocating buffer size but can have many problems
func echo(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 512)

	for {
		size, err := conn.Read(buffer)

		if err == io.EOF {
			log.Println("Client disconnected")
			break
		}

		if err != nil {
			log.Println("error in Reading data", err)
		}

		log.Println("size is : ", size)
		log.Println("Recieved : ", string(buffer))

		if _, err := conn.Write(buffer[0:size]); err != nil {
			log.Fatalln("unable to write data ", err)
		}
	}

}

func echoBuff(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	str, err := reader.ReadString('\n')

	if err != nil {
		log.Fatalln("error while reading string")
	}

	writer := bufio.NewWriter(conn)

	size, err := writer.WriteString(str)

	if err != nil {
		log.Fatalln("error while writing the line")
	}

	log.Printf("written %d ", size)

	writer.Flush()
}

func echoCopy(conn net.Conn) {

	defer conn.Close()

	if _, err := io.Copy(conn, conn); err != nil {
		log.Println(err)
	}

}
