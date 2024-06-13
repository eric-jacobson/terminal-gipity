package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/nsf/termbox-go"
)

const SPACE = " "
const BACKSPACE = "\x08"

func main() {
	godotenv.Load()

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("no OPENAI_API_KEY environment variable found")
	}

	err := termbox.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer termbox.Close()

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
