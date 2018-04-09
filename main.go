package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
)

func main() {
	flag.Parse()
	args := flag.Args()

	e := &Env{os.Environ()}
	s := NewStore()
	c := NewCache(*s)

	status := 0
	keys, command := split(args)

	if err := e.Load(c, keys); err != nil {
		log.Print(err)
		status = 1
	}

	cmd := exec.Command(command[0], command[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Env = e.Environ()
	if err := cmd.Run(); err != nil {
		log.Print(err)
		status = 1
	}
	os.Exit(status)
}

func split(args []string) ([]string, []string) {
	pos := 0
	width := 0
	for i, arg := range args {
		if arg == "--" {
			pos = i
			width = 1
			break
		}
	}
	return args[:pos], args[pos+width:]
}

type env interface {
	Set(string, string)
	Environ() []string
}

type store interface {
	Get(string, string) (string, error)
}
