package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestKey(t *testing.T) {
	type test struct {
		in    string
		name  string
		id    string
		field string
		ok    bool
	}

	cases := []test{
		{"NAME=id.field", "NAME", "id", "field", true},
		{"NA.ME=i=d.fie=ld", "NA.ME", "i=d", "fie=ld", true},
		{"NAME=i.d.field", "NAME", "i", "d.field", true},
		{"NAME=id", "NAME", "id", "", false},
		{"NAME=.field", "NAME", "", "field", false},
		{"NAME=.", "NAME", "", "", false},
		{"=.", "", "", "", false},
		{"=", "", "", "", false},
		{".", "", "", "", false},
	}

	for _, tt := range cases {
		k := &Key{tt.in}
		if k.Name() != tt.name {
			t.Errorf("Key{%s}.Name got '%s'", tt.in, k.Name())
		}
		if k.ID() != tt.id {
			t.Errorf("Key{%s}.ID got '%s'", tt.in, k.ID())
		}
		if k.Field() != tt.field {
			t.Errorf("Key{%s}.Field got '%s'", tt.in, k.Field())
		}
		_, err := NewKey(tt.in)
		if tt.ok && err != nil {
			t.Errorf("Key{%s} got error '%s'", tt.in, err)
		} else if !tt.ok && err == nil {
			t.Errorf("Key{%s} ok but expected error", tt.in)
		}
	}
}

func TestSplit(t *testing.T) {
	type test struct {
		in   []string
		keys []string
		cmd  []string
	}

	cases := []test{
		{
			[]string{"NAME=id.field", "OTHER=id.field", "--", "command", "arg"},
			[]string{"NAME=id.field", "OTHER=id.field"},
			[]string{"command", "arg"},
		},
		{
			[]string{"NAME=id.field", "OTHER=id.field", "command", "arg"},
			[]string{},
			[]string{"NAME=id.field", "OTHER=id.field", "command", "arg"},
		},
		{
			[]string{"NAME=id.field", "--", "command", "--", "arg"},
			[]string{"NAME=id.field"},
			[]string{"command", "--", "arg"},
		},
	}

	for _, tt := range cases {
		keys, cmd := split(tt.in)
		if !reflect.DeepEqual(keys, tt.keys) {
			t.Errorf("split(%v) keys got %v", tt.in, keys)
		}
		if !reflect.DeepEqual(cmd, tt.cmd) {
			t.Errorf("split(%v) cmd got %v", tt.in, cmd)
		}
	}
}

type MockStore struct {
	fail bool
}

func (s *MockStore) Get(name, key string) (string, error) {
	if s.fail {
		return "", fmt.Errorf("womp")
	}
	return "value", nil
}

func TestEnvLoad(t *testing.T) {
	store := &MockStore{}
	env := &Env{}
	keys := []string{"A=b.c", "D=e.f"}

	err := env.Load(store, keys)
	if err != nil {
		t.Errorf("Env.Load(store, %v) returned err: %s", keys, err)
	}

	got := env.Environ()
	want := []string{"A=value", "D=value"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("after Env.Load(store, %v), Environ() returned %v", keys, got)
	}

	store.fail = true
	if err := env.Load(store, keys); err == nil {
		t.Errorf("with MockStore.fail=true, expected Env.Load() to fail")
	}
}
