package main

import (
	"crypto/rand"
	"encoding/hex"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func Test__copy(t *testing.T) {
	f, err := ioutil.TempFile("", "test")
	if err != nil {
		t.Error(err)
	}

	defer os.Remove(f.Name())

	f.WriteString("abc")
	if err := f.Close(); err != nil {
		t.Error(err)
	}

	b := make([]byte, 10)
	if _, err := rand.Read(b); err != nil {
		t.Error(err)
	}

	newpath := filepath.Join(os.TempDir(), hex.EncodeToString(b))

	if _, err := os.Stat(f.Name()); err != nil {
		t.Error("should exist file")
	}

	if _, err := os.Stat(newpath); err == nil {
		t.Error("should not exist file")
	}

	if err := _copy(f.Name(), newpath); err != nil {
		t.Error(err)
	}

	defer os.Remove(newpath)

	if _, err := os.Stat(f.Name()); err == nil {
		t.Error("should not exist file")
	}

	if _, err := os.Stat(newpath); err != nil {
		t.Error("should exist file")
	}

	s, err := ioutil.ReadFile(newpath)
	if err != nil {
		t.Error(err)
	}

	if string(s) != "abc" {
		t.Error("should abc but got", string(s))
	}
}
