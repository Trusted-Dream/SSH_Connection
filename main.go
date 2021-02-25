package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/BurntSushi/toml"
	"golang.org/x/crypto/ssh"
)

type Conn struct {
	*ssh.Client
}

type Config struct {
	IP      string `toml:"IP"`
	Port    string `toml:"PORT"`
	User    string `toml:"USER"`
	RsaFile string `toml:"KEY"`
	Env     string `toml:"GO_ENV"`
}

var conf Config

func connection() (*Conn, error) {
	_, err := toml.DecodeFile(".env", &conf)
	if err != nil {
		log.Fatal(err)
	}

	buf, err := ioutil.ReadFile(conf.RsaFile)
	if err != nil {
		panic(err)
	}
	key, err := ssh.ParsePrivateKey(buf)
	if err != nil {
		panic(err)
	}
	config := &ssh.ClientConfig{
		User: conf.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	conn, err := ssh.Dial("tcp", conf.IP+":"+conf.Port, config)
	if err != nil {
		log.Println(err)
	}

	return &Conn{conn}, nil
}

func Command(cmd string) string {
	conn, _ := connection()
	defer conn.Close()
	session, err := conn.NewSession()
	if err != nil {
		log.Println(err)
	}
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run(cmd); err != nil {
		log.Fatal("Failed to run: " + err.Error())
	}
	msg := "COMMAND:" + cmd + "\n" + b.String()
	return msg
}

func main() {
	msg := Command("which python")
	fmt.Println(msg)
	msg = Command("ls -la")
	fmt.Println(msg)
	msg = Command("date")
	fmt.Println(msg)
}
