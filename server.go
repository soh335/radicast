package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"

	"github.com/gorilla/mux"
)

type Server struct {
	Output string
	Title  string
	Addr   string
}

func (s *Server) errorHandler(f func(http.ResponseWriter, *http.Request) error) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			s.Log(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

func (s *Server) Run() error {

	s.Log("start ", s.Addr)

	router := mux.NewRouter()

	router.HandleFunc("/podcast/{program}.mp3", s.errorHandler(func(w http.ResponseWriter, r *http.Request) error {
		dir := mux.Vars(r)["program"]

		mp3Path, mp3Stat, err := s.mp3Path(dir)

		if _, err := os.Stat(mp3Path); err != nil {
			http.NotFound(w, r)
			return nil
		}

		xmlPath, _, err := s.xmlPath(dir)

		if _, err := os.Stat(xmlPath); err != nil {
			http.NotFound(w, r)
			return nil
		}

		f, err := os.Open(mp3Path)

		if err != nil {
			return err
		}

		defer f.Close()

		http.ServeContent(w, r, mp3Stat.Name(), mp3Stat.ModTime(), f)
		return nil
	}))

	router.HandleFunc("/rss", s.errorHandler(func(w http.ResponseWriter, r *http.Request) error {

		baseUrl, err := url.Parse("http://" + r.Host)

		if err != nil {
			return err
		}

		rss, err := s.rss(baseUrl)

		if err != nil {
			return err
		}

		var b bytes.Buffer

		b.WriteString(xml.Header)

		enc := xml.NewEncoder(&b)
		enc.Indent("", "    ")
		if err := enc.Encode(rss); err != nil {
			return err
		}

		if _, err := io.Copy(w, &b); err != nil {
			return err
		}

		return nil
	}))

	return http.ListenAndServe(s.Addr, router)
}

func (s *Server) rss(baseUrl *url.URL) (*PodcastRss, error) {

	dirs, err := ioutil.ReadDir(s.Output)

	if err != nil {
		return nil, err
	}

	items := PodcastItems{}

	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}

		item, err := s.itemByDir(dir.Name(), baseUrl)

		if err != nil {
			s.Log(err)
			continue
		}

		items = append(items, *item)
	}

	sort.Sort(sort.Reverse(items))

	rss := NewPodcastRss()

	channel := PodcastChannel{}
	channel.Title = s.Title
	channel.Items = items

	rss.Channel = channel

	return rss, nil
}

func (s *Server) itemByDir(dir string, baseUrl *url.URL) (*PodcastItem, error) {

	_, mp3Stat, err := s.mp3Path(dir)

	if err != nil {
		return nil, err
	}

	xmlPath, _, err := s.xmlPath(dir)

	if err != nil {
		return nil, err
	}

	xmlFile, err := os.Open(xmlPath)

	if err != nil {
		return nil, err
	}

	defer xmlFile.Close()

	dec := xml.NewDecoder(xmlFile)

	var prog RadikoProg
	if err := dec.Decode(&prog); err != nil {
		return nil, err
	}

	u, err := url.Parse("/podcast/" + dir + ".mp3")

	if err != nil {
		return nil, err
	}

	ft, _ := prog.FtTime()

	var item PodcastItem

	item.Title = fmt.Sprintf("%s (%s)", prog.Title, ft)
	item.ITunesAuthor = prog.Pfm
	item.ITunesSummary = prog.Info

	item.Enclosure.Url = baseUrl.ResolveReference(u).String()
	item.Enclosure.Type = "audio/mpeg"
	item.Enclosure.Length = int(mp3Stat.Size())
	item.PubDate = PubDate{mp3Stat.ModTime()}

	return &item, nil
}

func (s *Server) mp3Path(dir string) (string, os.FileInfo, error) {
	return s.pathStat(dir, "podcast.mp3")
}

func (s *Server) xmlPath(dir string) (string, os.FileInfo, error) {
	return s.pathStat(dir, "podcast.xml")
}

func (s *Server) pathStat(dir string, name string) (string, os.FileInfo, error) {
	p := filepath.Join(s.Output, dir, name)
	stat, err := os.Stat(p)

	if err != nil {
		return "", nil, err
	}

	return p, stat, nil
}

func (s *Server) Log(v ...interface{}) {
	log.Println("[server]", fmt.Sprint(v...))
}
