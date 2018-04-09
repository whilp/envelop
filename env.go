package main

import (
	"fmt"
	"strings"
)

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
