package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
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

func SetupConfig(ctx context.Context) error {

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

	if _, err := io.Copy(os.Stdout, bytes.NewReader(byt)); err != nil {
		return err
	}

	return nil
}
