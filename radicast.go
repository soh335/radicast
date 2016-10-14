package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"sync"

	"github.com/robfig/cron"
)

type Radicast struct {
	reloadChan chan struct{}
	saveChan   chan *Radiko
	configPath string
	cron       *cron.Cron
	m          sync.Mutex
	wg         sync.WaitGroup
	ctx        context.Context
	cancel     context.CancelFunc
	host       string
	port       string
	title      string
	output     string
	bitrate    string
	buffer     int64
	converter  string
}

func NewRadicast(path string, host string, port string, title string, output string, bitrate string, buffer int64, converter string) *Radicast {
	ctx, cancel := context.WithCancel(context.Background())

	r := &Radicast{
		reloadChan: make(chan struct{}),
		saveChan:   make(chan *Radiko),
		configPath: path,
		ctx:        ctx,
		cancel:     cancel,
		host:       host,
		port:       port,
		title:      title,
		output:     output,
		bitrate:    bitrate,
		buffer:     buffer,
		converter:  converter,
	}
	return r
}

func (r *Radicast) Run() error {
	if err := r.ReloadConfig(); err != nil {
		return err
	}

	if _, err := os.Stat(r.output); err != nil {
		if err := os.MkdirAll(r.output, 0777); err != nil {
			return err
		}
	}

	go func() {

		s := &Server{
			Output: r.output,
			Title:  r.title,
			Addr:   net.JoinHostPort(r.host, r.port),
		}
		if err := s.Run(); err != nil {
			r.Log(err)
			r.Stop()
		}
	}()

	for {
		select {
		case <-r.ctx.Done():
			r.wg.Wait()
			return r.ctx.Err()
		case <-r.reloadChan:
			if err := r.ReloadConfig(); err != nil {
				r.Log(err)
			}
		// if same program is recorded, write files as parallely and may occure error. so write file as serially by channel.
		case radiko := <-r.saveChan:
			func() {
				defer os.RemoveAll(radiko.TempDir)
				if err := radiko.Result.Save(r.output); err != nil {
					r.Log(err)
				}
			}()
		}
	}
}

func (r *Radicast) Stop() {
	r.cancel()
}

func (r *Radicast) ReloadConfig() error {
	r.m.Lock()
	defer r.m.Unlock()

	if r.cron != nil {
		r.cron.Stop()
		r.Log("stop previous cron")
	}

	config, err := LoadConfig(*configPath)
	if err != nil {
		return err
	}

	c := cron.New()
	for station, specs := range config {
		for _, spec := range specs {
			func(station string, spec string) {
				r.Log("station:", station, " spec:", spec)
				c.AddFunc(spec, func() {
					r.wg.Add(1)
					defer r.wg.Done()

					dir, err := ioutil.TempDir("", "radiko")
					if err != nil {
						r.Log(err)
						return
					}

					radiko := &Radiko{
						Station:   station,
						Bitrate:   r.bitrate,
						Buffer:    r.buffer,
						Converter: r.converter,
						TempDir:   dir,
					}

					if err := radiko.Run(r.ctx); err != nil {
						os.RemoveAll(radiko.TempDir)
						r.Log(err)
						return
					}

					r.saveChan <- radiko
				})
			}(station, spec)
		}
	}
	c.Start()

	r.cron = c

	r.Log("start new cron")

	return nil
}

func (r *Radicast) Log(v ...interface{}) {
	log.Println("[radicast]", fmt.Sprint(v...))
}
