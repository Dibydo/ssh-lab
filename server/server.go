package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/unknwon/com"
	"golang.org/x/crypto/ssh"
)

func Listen() {
	config := ssh.ServerConfig{
		PasswordCallback: func(conn ssh.ConnMetadata, password []byte) (*ssh.Permissions, error) {
			if conn.User() == "admin" {
				if string(password) == "admin" {
					return &ssh.Permissions{}, nil
				}
			}
			return nil, errors.New("")
		},
		MaxAuthTries: -1,
	}
	keyPath := "host.rsa"
	if !com.IsExist(keyPath) {
		os.MkdirAll(filepath.Dir(keyPath), os.ModePerm)
		_, stderr, err := com.ExecCmd("ssh-keygen", "-f", keyPath, "-t", "rsa", "-N", "")
		if err != nil {
			panic(fmt.Sprintf("Key fail: %v - %s", err, stderr))
		}
		fmt.Printf("Nice key: %s\n", keyPath)
	}
	privateBytes, err := ioutil.ReadFile(keyPath)
	if err != nil {
		panic("Fail")
	}
	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		panic("Fail")
	}
	config.AddHostKey(private)
	listen(&config)
}
func listen(config *ssh.ServerConfig) {
	fmt.Println("Start listening")
	listener, err := net.Listen("tcp", "127.0.0.1:3000")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
		}
		fmt.Printf("Handshaking for %s\n", conn.RemoteAddr())
		sConn, chans, _, err := ssh.NewServerConn(conn, config)
		if err != nil {
			log.Print(err)
			continue
		}
		fmt.Printf("Connection from %s (%s)\n", sConn.RemoteAddr(), sConn.ClientVersion())
		handle(chans)
	}
}

func handle(chans <-chan ssh.NewChannel) {
	for newChan := range chans {
		if newChan.ChannelType() != "session" {
			newChan.Reject(ssh.UnknownChannelType, "chzh")
			continue
		}
		channel, requests, err := newChan.Accept()
		if err != nil {
			log.Fatalf("Could not accept channel: %v", err)
		}
		ex := false
		nRequests := make(chan *ssh.Request, 15)
		var com string
		for req := range requests {
			fmt.Println(req.Type)
			if req.Type == "exec" {
				ex = true
				com = string(req.Payload)
			}
			nRequests <- req
			if req.Type == "shell" || req.Type == "exec" {
				break
			}
		}
		go func(in <-chan *ssh.Request) {
			for req := range in {
				req.Reply(req.Type == "shell" || req.Type == "exec", nil)
			}
		}(nRequests)
		if ex {
			defer channel.Close()
			splitted := strings.Split(com, " ")
			fmt.Println(splitted)
			bsp := make([]byte, 0)
			for _, s := range []byte(splitted[0]) {
				if !(s < 97 || s > 122) {
					bsp = append(bsp, byte(s))
				}
			}
			splitted[0] = string(bsp)
			res, err := exec.Command(splitted[0], splitted[1:]...).Output()
			if err != nil {
				log.Print(err)
				channel.Write([]byte(err.Error() + "\n"))
			}
			channel.Write(res)
			channel.SendRequest("exit-status", false, []byte{0, 0, 0, 0})
			channel.Close()
		} else {
			a := make([]byte, 30)
			defer channel.Close()
			for {
				channel.Write([]byte(">"))
				channel.Read(a)
				cmd := string(a)
				i := strings.Index(cmd, "\n")
				cmd = cmd[:i]
				if cmd == "exit" {
					channel.Close()
					break
				}
				splitted := strings.Split(cmd, " ")
				fmt.Println(splitted)
				res, err := exec.Command(splitted[0], splitted[1:]...).Output()
				if err != nil {
					log.Print(err)
					channel.Write([]byte(err.Error() + "\n"))
				}
				channel.Write(res)
				channel.SendRequest("exit-status", false, []byte{0, 0, 0, 0})
			}
		}
	}
}
func main() {
	Listen()
}
