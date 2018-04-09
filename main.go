package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
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

// NAME=id.field
type Key struct {
	key string
}

func NewKey(key string) (*Key, error) {
	k := &Key{key}
	var err error
	switch "" {
	case k.Name():
		err = fmt.Errorf("key missing name: %s", key)
	case k.ID():
		err = fmt.Errorf("key missing ID: %s", key)
	case k.Field():
		err = fmt.Errorf("key missing field: %s", key)
	}
	return k, err
}

func (k *Key) Name() string {
	i := strings.Index(k.key, "=")
	if i < 0 {
		return ""
	}
	return k.key[0:i]
}

func (k *Key) ID() string {
	rest := k.key[len(k.Name())+1:]
	i := strings.Index(rest, ".")
	if i < 0 {
		i = len(rest)
	}
	return rest[:i]
}

func (k *Key) Field() string {
	rest := k.key[len(k.Name())+1:]
	i := strings.Index(rest, ".")
	if i < 0 {
		return ""
	}
	return rest[i+1:]
}

type env interface {
	Set(string, string)
	Environ() []string
}

type Env struct {
	env []string
}

func (e *Env) Set(key, value string) {
	e.env = append(e.env, fmt.Sprintf("%s=%s", key, value))
}

func (e *Env) Environ() []string {
	return e.env
}

func (e *Env) Load(s store, keys []string) error {
	for _, key := range keys {
		var val string
		k, err := NewKey(key)
		if err != nil {
			return err
		}
		if val, err = s.Get(k.ID(), k.Field()); err != nil {
			return err
		}
		e.Set(k.Name(), val)
	}
	return nil
}

type store interface {
	Get(string, string) (string, error)
}

type CacheKey struct {
	name string
	key  string
}

type Cache struct {
	Store
	data map[CacheKey]Item
}

func NewCache(s Store) *Cache {
	c := &Cache{}
	c.Store = s
	c.data = make(map[CacheKey]Item)
	return c
}

func (c *Cache) Get(name, key string) (string, error) {
	item, err := c.GetItem(name, key)
	if err != nil {
		return "", err
	}
	return item.Find(key)
}

func (c *Cache) GetItem(name, key string) (*Item, error) {
	ck := CacheKey{name, key}
	if item, ok := c.data[ck]; ok {
		return &item, nil
	}

	item, err := c.Store.GetItem(name, key)
	if err != nil {
		return &Item{}, err
	}
	c.data[ck] = *item
	return item, nil
}

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

type Item struct {
	Details *Details `json:"details"`
}

type Details struct {
	Fields   []*Field   `json:"fields"`
	Sections []*Section `json:"sections"`
}

type Section struct {
	Fields []*SectionField `json:"fields"`
}

type SectionField struct {
	Key        string `json:"k"`
	FieldName  string `json:"n"`
	Title      string `json:"t"`
	FieldValue string `json:"v"`
}

func (f *SectionField) Name() string {
	return f.Title
}

func (f *SectionField) Value() string {
	return f.FieldValue
}

type Field struct {
	Designation string `json:"designation"`
	FieldName   string `json:"name"`
	Type        string `json:"type"`
	FieldValue  string `json:"value"`
}

func (f *Field) Name() string {
	return f.FieldName
}

func (f *Field) Value() string {
	return f.FieldValue
}

type value interface {
	Name() string
	Value() string
}

func (i *Item) Find(key string) (string, error) {
	values := []value{}
	for _, field := range i.Details.Fields {
		values = append(values, field)
	}
	for _, section := range i.Details.Sections {
		for _, field := range section.Fields {
			values = append(values, field)
		}
	}

	for _, field := range values {
		if field.Name() == key {
			return field.Value(), nil
		}
	}
	return "", fmt.Errorf("could not find value for key: %s", key)
}

type run interface {
	Run([]string) error
}

type Runner struct{}

func (r *Runner) Run(e env, cmd []string) error {
	return nil
}
