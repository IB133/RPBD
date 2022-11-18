package scan

import (
	"fmt"
	"net"
	"sort"
	"sync"
)

// TODO: реализация в main не возвращает слайс открытых портов
// необходимо реализовать функцию Scan по по аналогии с кодом в main.go
// только её нужно дополнить, чтобы вернуть слайс открытых портов
// отсортированных по возрастанию
func worker(ports chan int, wg *sync.WaitGroup, openPorts *[]int, address string) {
	for p := range ports {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", address, p))
		if err != nil {
			// fmt.Println("port closed:", p)
			wg.Done()
			continue
		}
		conn.Close()

		fmt.Println("port opened:", p)
		*openPorts = append(*openPorts, p)
		wg.Done()
	}
}

func Scan(address string) []int {
	ports := make(chan int, 200)
	wg := sync.WaitGroup{}
	openPorts := make([]int, 0)
	for i := 0; i < cap(ports); i++ {
		go worker(ports, &wg, &openPorts, address)
	}

	for i := 1; i < 10000; i++ {
		wg.Add(1)
		ports <- i
	}

	wg.Wait()
	close(ports)
	sort.Ints(openPorts)
	return openPorts
}
