package service

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/IB133/RPBD/console_game/krot"
)

func info(w io.Writer, s string) {
	fmt.Fprintf(w, s)
}

func Start() {
	k := krot.New()
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println("Select an action:")
		fmt.Print("1.Dig a hole\n2.Eat grass\n3.Go fight\n4.Sleep\n5.Check yout stats\n")
		scanner.Scan()
		input := scanner.Text()
		switch input {
		case "1":
			fmt.Print("1.Intense\n2.Lazy\n")
			scanner.Scan()
			input = scanner.Text()
			if input == "1" {
				k.Dig(true)
			} else if input == "2" {
				k.Dig(false)
			}
		case "2":
			fmt.Print("1.Green\n2.Withered\n")
			scanner.Scan()
			input = scanner.Text()
			if input == "1" {
				k.Eat(true)
			} else if input == "2" {
				k.Eat(false)
			}
		case "3":
			fmt.Print("1.Weak(30)\n2.Middle(50)\n3.Strong(70)")
			scanner.Scan()
			input = scanner.Text()
			if input == "1" {
				fmt.Println(k.Fight(30))
			} else if input == "2" {
				fmt.Println(k.Fight(50))
			} else if input == "3" {
				fmt.Println(k.Fight(70))
			}
		case "4":
			k.Sleep()
			fmt.Println("You've been asleep all day")
		case "5":
			info(os.Stdout, k.Stats())
		}
		if k.IsLoose() {
			fmt.Println("You loose")
		}
		if k.IsWin() {
			fmt.Println("You win!!!")
		}
		fmt.Println("Night is coming")
		k.Sleep()
		fmt.Println("New day")
	}
}
