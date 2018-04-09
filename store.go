package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
)

type Store struct {
	name   string
	stdin  io.Reader
	stderr io.Writer
	env    []string
}

func NewStore() *Store {
	op := &Store{}
	op.name = "op"
	op.stdin = os.Stdin
	op.stderr = os.Stderr
	op.env = os.Environ()
	return op
}

func (s *Store) Get(name, key string) (string, error) {
	item, err := s.GetItem(name, key)
	if err != nil {
		return "", err
	}
	return item.Find(key)
}

func (s *Store) GetItem(name, key string) (*Item, error) {
	item := &Item{}
	buf := &bytes.Buffer{}
	cmd := exec.Command(s.name, "get", "item", name)
	cmd.Stdin = s.stdin
	cmd.Stdout = buf
	cmd.Stderr = s.stderr
	cmd.Env = s.env
	if err := cmd.Run(); err != nil {
		return item, err
	}

	data, err := ioutil.ReadAll(buf)
	if err != nil {
		return item, err
	}
	if err := json.Unmarshal(data, item); err != nil {
		return item, err
	}
	return item, nil
}
