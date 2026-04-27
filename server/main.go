package main

import (
	"fmt"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":10001")
	if err != nil {
		fmt.Printf("Все пошло по пизде\n")
		panic(err)
	}

	defer listener.Close()
	fmt.Println("сервер висит на 10k1")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Ошибка accept %d\n", err)
			continue
		}
		go handleClient(conn)
	}
}
