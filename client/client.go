package main

import (
	"fmt"
	"log"
	"os"

	"github.com/skorobogatov/input"
	"golang.org/x/crypto/ssh"
)

func main() {
	username := "admin"
	password := "admin"
	hostname := "127.0.0.1"
	port := "3000"
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := ssh.Dial("tcp", hostname + ":" + port, config)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	sess, err := client.NewSession()
	if err != nil {
		log.Fatal("Fail:", err)
	}
	defer sess.Close()
	stdin, err := sess.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	if len(os.Args) == 1 {
		sess.Stdout = os.Stdout
		sess.Stderr = os.Stderr
		err = sess.Shell()
		if err != nil {
			log.Fatal(err)
		}
		for {
			cmd := input.Gets()
			_, err = fmt.Fprintf(stdin, "%s\n", cmd)
			if err != nil {
				log.Fatal(err)
			}
			if cmd == "exit" {
				sess.Close()
				break
			}
		}
	} else {
		sess.Stdout = os.Stdout
		sess.Stderr = os.Stderr
		cmd := ""
		for i := 1; i < len(os.Args); i++ {
			cmd += os.Args[i]
			if i != len(os.Args)-1 {
				cmd += " "
			}
		}
		err = sess.Run(cmd)
		if err != nil {
			fmt.Println("hi")
			log.Fatal(err)
		}
	}
}
