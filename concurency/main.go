package main

import (
	"fmt"

	"github.com/IB133/RPBD/concurency/scan"
)

func main() {
	openPorts := scan.Scan("127.0.0.1")
	for _, i := range openPorts {
		fmt.Println(i)
	}
}
