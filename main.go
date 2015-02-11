package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/net/context"
)

var (
	host       = flag.String("host", "0.0.0.0", "host")
	port       = flag.String("port", "3355", "port")
	buffer     = flag.Int64("buffer", 60, "buffer for recording")
	output     = flag.String("output", "output", "output")
	bitrate    = flag.String("bitrate", "64k", "bitrate")
	title      = flag.String("title", "radicast", "title")
	configPath = flag.String("config", "config.json", "path of config.json")
	setup      = flag.Bool("setup", false, "initialize json configuration")
	converter  = flag.String("converter", "", "ffmpeg or avconv. If not set this option, radicast search its automatically.")
)

func main() {
	flag.Parse()

	if *setup {
		runSetup()
		return
	}

	runRadicast()
}

func runRadicast() {
	r := NewRadicast(*configPath, *host, *port, *title, *output, *bitrate, *buffer, *converter)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill, syscall.SIGHUP)

	go func() {
		for {
			s := <-signalChan
			r.Log("got signal:", s)
			switch s {
			case syscall.SIGHUP:
				r.ReloadConfig()
			default:
				r.Stop()
			}
		}
	}()

	if err := r.Run(); err != nil {
		log.Fatal(err)
	}
}

func runSetup() {
	ctx, cancel := context.WithCancel(context.Background())

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)

	go func() {
		s := <-signalChan
		log.Println("got signal:", s)
		cancel()
	}()

	if err := SetupConfig(ctx); err != nil {
		log.Fatal(err)
	}
	return
}
