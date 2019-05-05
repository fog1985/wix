package main

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"
        "./sub"
)

func handleConnection(c net.Conn) {
	fmt.Printf("Serving %s\n", c.RemoteAddr().String())

        result := morse.EncodeITU(strings.Split(c.RemoteAddr().String(), ":")[0]) + "\n"
        c.Write([]byte(string(result)))
        c.Close()
}

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide a port number!")
		return
	}

	PORT := ":" + arguments[1]
	l, err := net.Listen("tcp4", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()
	rand.Seed(time.Now().Unix())

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go handleConnection(c)
	}
}
