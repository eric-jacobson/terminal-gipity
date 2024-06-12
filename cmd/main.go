package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/nsf/termbox-go"
)

func main() {
	err := termbox.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer termbox.Close()

	const SPACE = " "
	const BACKSPACE = "\x08"

	var input strings.Builder

	fmt.Print(">")

main_loop:
	for {
		switch keyPress := termbox.PollEvent(); keyPress.Type {
		case termbox.EventKey:
			switch keyPress.Key {
			case termbox.KeyEsc:
				break main_loop
			case termbox.KeyEnter:
				echo(input.String())
				input.Reset()
				break
			case termbox.KeyBackspace, termbox.KeyBackspace2:
				input.WriteString(BACKSPACE)
				fmt.Print(BACKSPACE)
			case termbox.KeySpace:
				input.WriteString(SPACE)
				fmt.Print(SPACE)
			default:
				if keyPress.Ch != 0 {
					char := fmt.Sprintf("%c", keyPress.Ch)
					input.WriteString(char)
					fmt.Print(char)
				}
			}
		}
	}
}

func echo(input string) {
	fmt.Printf("\n\tEcho: %v\n>", input)
}
