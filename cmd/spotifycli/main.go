package main

import (
	"fmt"
	"os"
	"reflect"

	"github.com/pblaszczyk/go.spotify"
)

func usage() {
	fmt.Printf(`spotifycli - commandline controller for Spotify desktop app.

Usage:
  spotifycli [commands] [args...]

Commands:
  search             - Search for artist/album/track.
       artist <name> - Search for artist.
       album  <name> - Search for album.
       track  <name> - Search for track.
`)
	platfusage()
	os.Exit(1)
}

func handlerr(err error) {
	if err != nil && !spotify.IsEOF(err) {
		fmt.Fprintf(os.Stderr, "[spotifycli]: failed: %q", err)
		os.Exit(1)
	}
}

func newApp() *spotify.App {
	app, err := spotify.NewApp("")
	handlerr(err)
	return app
}

func searchArtist() {
	s := spotify.NewSearch()
	res, err := make(chan []spotify.Artist), make(chan error)
	s.Artist(os.Args[3], res, err)
	b := true
LOOP:
	for {
		select {
		case res := <-res:
			disp(res, b)
			b = false
		case err := <-err:
			handlerr(err)
			break LOOP
		}
	}
}

func searchAlbum() {
	s := spotify.NewSearch()
	res, err := make(chan []spotify.Album), make(chan error)
	s.Album(os.Args[3], res, err)
	b := true
LOOP:
	for {
		select {
		case res := <-res:
			disp(res, b)
			b = false
		case err := <-err:
			handlerr(err)
			break LOOP
		}
	}
}

func searchTrack() {
	s := spotify.NewSearch()
	res, err := make(chan []spotify.Track), make(chan error)
	s.Track(os.Args[3], res, err)
	b := true
LOOP:
	for {
		select {
		case res := <-res:
			disp(res, b)
			b = false
		case err := <-err:
			handlerr(err)
			break LOOP
		}
	}
}

func search() {
	switch os.Args[2] {
	case "artist":
		searchArtist()
	case "album":
		searchAlbum()
	case "track":
		searchTrack()
	}
	fmt.Println("")
}

func main() {
	if len(os.Args) == 1 {
		usage()
	}
	switch os.Args[1] {
	case "run":
		handlerr(newApp().Start())
	case "kill":
		handlerr(newApp().Kill())
	case "process":
		if err := newApp().Ping(); err != nil {
			fmt.Println("Not running")
			os.Exit(0)
		}
		fmt.Println("Running")
	case "search":
		if len(os.Args) != 4 {
			usage()
		}
		search()
	default:
		platform()
	}
}

func disp(r interface{}, b bool) {
	for i := reflect.ValueOf(r).Len() - 1; i >= 0; i-- {
		if !b {
			fmt.Println("")
			fmt.Println("")
		}
		b = false
		for j, l := 0, reflect.ValueOf(r).Index(i).NumField(); j < l; j++ {
			f := reflect.ValueOf(r).Index(i).Field(j)
			if f.Kind() == reflect.Slice && f.Len() > 0 {
				fmt.Printf("%q\n", reflect.ValueOf(r).Index(i).Type().Field(j).Name)
				disp(f.Interface(), true)
			} else {
				fmt.Printf("%q: %q",
					reflect.ValueOf(r).Index(i).Type().Field(j).Name, f.String())
			}
			if j < l-1 {
				fmt.Println("")
			}
		}
	}
}
