package main

import (
	"fmt"
	"github.com/gliderlabs/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const helpMessage = `Commands:
help - display this message
ls - list all files and directories
mkdir <dirname> - create directory <dirname>
deldir <dirname> - remove directory <dirname>
over - exit program`

func tPrintLn(t *terminal.Terminal, text string) {
	t.Write([]byte(text + "\n"))
}

func handleUserInput(t *terminal.Terminal) {
	for {
		line, err := t.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal("Terminal reading error ", err)
			break
		}

		input := strings.Split(line, " ")
		command := input[0]
		args := input[1:]

		switch command {
		case "help":
			tPrintLn(t, helpMessage)
		case "ls":
			files, _ := ioutil.ReadDir("./")
			for _, f := range files {
				tPrintLn(t, f.Name())
			}
		case "mkdir":
			if len(args) != 1 {
				tPrintLn(t, "Invalid number of arguments.")
				continue
			}

			err := os.Mkdir(args[0], 0777)
			if err != nil {
				tPrintLn(t, "Failed to create dir.")
			} else {
				tPrintLn(t, "Dir created successfully.")
			}
		case "deldir":
			if len(args) != 1 {
				tPrintLn(t, "Invalid number of arguments.")
				continue
			}

			err := os.RemoveAll(args[0])
			if err != nil {
				tPrintLn(t, "Failed to remove dir.")
			} else {
				tPrintLn(t, "Dir removed successfully.")
			}
		case "over":
			os.Exit(1)
		default:
			tPrintLn(t, "Try another command, type help for instructions.")
		}
	}

	log.Println("Terminal closed")
}

func handleSession(s ssh.Session) {
	t := terminal.NewTerminal(s, "Enter command: ")
	fmt.Fprintf(t, "%s\n",
		helpMessage,
	)
	handleUserInput(t)
}

func main() {
	ssh.Handle(handleSession)

	log.Println("starting ssh server on port 2222...")
	log.Fatal(ssh.ListenAndServe(":2222", nil))
}

