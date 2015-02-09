package main

import (
	"bytes"
	"encoding/json"
	"io"
	"os"

	"golang.org/x/net/context"
)

type Config map[string][]string

func LoadConfig(path string) (Config, error) {
	f, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	var c Config
	if err := json.NewDecoder(f).Decode(&c); err != nil {
		return nil, err
	}

	return c, nil
}

func SetupConfig(ctx context.Context, path string) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666)

	if err != nil {
		return err
	}

	r := &Radiko{}
	stations, err := r.StationList(ctx)
	if err != nil {
		return err
	}

	c := Config{}

	for _, station := range stations {
		c[station] = []string{}
	}

	byt, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	if _, err := io.Copy(f, bytes.NewReader(byt)); err != nil {
		return err
	}

	return nil
}
