package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	"code.google.com/p/go.net/context"
	"github.com/robfig/cron"
)

type Radicast struct {
	reloadChan chan struct{}
	saveChan   chan *RadikoResult
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
}

func NewRadicast(path string, host string, port string, title string, output string, bitrate string, buffer int64) *Radicast {
	ctx, cancel := context.WithCancel(context.Background())

	r := &Radicast{
		reloadChan: make(chan struct{}),
		saveChan:   make(chan *RadikoResult),
		configPath: path,
		ctx:        ctx,
		cancel:     cancel,
		host:       host,
		port:       port,
		title:      title,
		output:     output,
		bitrate:    bitrate,
		buffer:     buffer,
	}
	return r
}

func (r *Radicast) Run() error {
	if err := r.ReloadConfig(); err != nil {
		return err
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
		case ret := <-r.saveChan:
			if err := ret.Save(r.output); err != nil {
				r.Log(err)
				if err := os.Remove(ret.Mp3Path); err != nil {
					r.Log(err)
				}
			}
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

					radiko := &Radiko{
						Station: station,
						Bitrate: r.bitrate,
						Buffer:  r.buffer,
					}

					ret, err := radiko.Run(r.ctx)
					if err != nil {
						r.Log(err)
						return
					}

					if ret != nil {
						r.saveChan <- ret
					}
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
