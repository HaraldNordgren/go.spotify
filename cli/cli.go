package cli

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/codegangsta/cli"
	"github.com/pblaszczyk/sscc"
)

type impl struct {
	*cli.App
	ctrl sscc.Controller
}

// NewApp returns initialized instance of ssc struct.
func NewApp() (app *impl) {
	ctrl := sscc.NewControl(&sscc.Context{})
	app = &impl{cli.NewApp(), ctrl}
	app.App.Name = "sscc"
	app.App.Version = "0.0.1"
	app.App.Usage = "commandline controller of Spotify desktop app."
	app.App.Commands = []cli.Command{
		{Name: "run", Usage: "Run Spotify desktop app.", Action: app.Start},
		{Name: "kill", Usage: "Kill Spotify desktop app.", Action: app.Kill},
		{Name: "raise", Usage: "Raise Spotify desktop app.", Action: app.Raise},
		{Name: "next", Usage: "Play next track.", Action: app.Next},
		{Name: "prev", Usage: "Play prev track.", Action: app.Prev},
		{Name: "open", Usage: "Play music identified by uri.", Action: app.Open},
		{Name: "seek", Usage: "Goto.", Action: app.Goto},
		{Name: "play", Usage: "Play current track/uri/pos.", Action: app.Play},
		{Name: "stop", Usage: "Stop.", Action: app.Stop},
		{Name: "toggle", Usage: "Play/Pause.", Action: app.Toggle},
		{Name: "status", Usage: "Status.", Action: app.Status},
		{Name: "track", Usage: "Current track.", Action: app.CurTrack},
		{Name: "setpos", Usage: "Sets position.", Action: app.SetPos},
		{Name: "length", Usage: "Length of current track.", Action: app.Length},
		{Name: "pos", Usage: "Current position.", Action: app.Pos},
		{Name: "canplay", Usage: "Can play.", Action: app.CanPlay},
		{Name: "cannext", Usage: "Can next.", Action: app.CanNext},
		{Name: "canprev", Usage: "Can prev.", Action: app.CanPrev},
		{Name: "canctrl", Usage: "Can control.", Action: app.CanControl},
		{Name: "search", Usage: "Search for artist/album/track.",
			Subcommands: []cli.Command{
				{Name: "artist", Usage: "Search for artist.", Action: app.Artist,
					Flags: []cli.Flag{
						cli.BoolFlag{Name: "i", Usage: "Quasi-interactive mode"}},
				},
				{Name: "album", Usage: "Search for album.", Action: app.Album,
					Flags: []cli.Flag{
						cli.BoolFlag{Name: "i", Usage: "Quasi-interactive mode"}},
				},
				{Name: "track", Usage: "Search for track.", Action: app.Track,
					Flags: []cli.Flag{
						cli.BoolFlag{Name: "i", Usage: "Quasi-interactive mode"}},
				},
			},
		},
	}
	return
}

var handleErr = func(err error) {
	switch {
	case sscc.IsEOF(err) || err == nil:
		return
	default:
		fmt.Printf("sscc: %q\n", err)
		os.Exit(1)
	}
}

// Start starts spotify app.
func (s *impl) Start(ctx *cli.Context) {
	handleErr(s.ctrl.Run())
}

// Raise raises spotify app's window.
func (s *impl) Raise(ctx *cli.Context) {
	handleErr(s.ctrl.Raise())
}

// Kill stops spotify app.
func (s *impl) Kill(ctx *cli.Context) {
	handleErr(s.ctrl.Kill())
}

// Next starts playing next track.
func (s *impl) Next(ctx *cli.Context) {
	handleErr(s.ctrl.Next())
}

// Prev starts playing prev track.
func (s *impl) Prev(ctx *cli.Context) {
	handleErr(s.ctrl.Prev())
}

// Open starts playing specified uri.
func (s *impl) Open(ctx *cli.Context) {
	handleErr(validateSingle(ctx.Args()))
	handleErr(s.ctrl.Open(sscc.URI(ctx.Args().First())))
}

// Open starts playing.
func (s *impl) Play(ctx *cli.Context) {
	handleErr(s.ctrl.Play())
}

// Goto pos.
func (s *impl) Goto(ctx *cli.Context) {
	handleErr(validateSingle(ctx.Args()))
	d, err := time.ParseDuration(ctx.Args().First())
	handleErr(err)
	handleErr(s.ctrl.Goto(d))
}

// Stop playing current track.
func (s *impl) Stop(ctx *cli.Context) {
	handleErr(s.ctrl.Stop())
}

// Play/Pause current track.
func (s *impl) Toggle(ctx *cli.Context) {
	handleErr(s.ctrl.Toggle())
}

// CurTrack displays info about current track.
func (s *impl) CurTrack(ctx *cli.Context) {
	track, err := s.ctrl.CurTrack()
	handleErr(err)
	fmt.Println(track)
}

func (s *impl) SetPos(ctx *cli.Context) {
	handleErr(validateSingle(ctx.Args()))
	d, err := time.ParseDuration(ctx.Args().First())
	handleErr(err)
	handleErr(s.ctrl.SetPos(d))
}

func (s *impl) Length(ctx *cli.Context) {
	l, err := s.ctrl.Length()
	handleErr(err)
	fmt.Println(l)
}

func (s *impl) Pos(ctx *cli.Context) {
	d, err := s.ctrl.Pos()
	handleErr(err)
	fmt.Println(d)
}

func (s *impl) CanPlay(ctx *cli.Context) {
	b, err := s.ctrl.CanPlay()
	handleErr(err)
	fmt.Println(b)
}

func (s *impl) CanNext(ctx *cli.Context) {
	b, err := s.ctrl.CanNext()
	handleErr(err)
	fmt.Println(b)
}

func (s *impl) CanPrev(ctx *cli.Context) {
	b, err := s.ctrl.CanPrev()
	handleErr(err)
	fmt.Println(b)
}

func (s *impl) CanControl(ctx *cli.Context) {
	b, err := s.ctrl.CanControl()
	handleErr(err)
	fmt.Println(b)
}

// interactive runs in limited interactive mode if configured.
func (s *impl) interactive(ctx *cli.Context) {
	if ctx.Bool("i") {
		fmt.Print("Play: ")
		r := bufio.NewReader(os.Stdin)
		uri, _, err := r.ReadLine()
		handleErr(err)
		handleErr(s.ctrl.Open(sscc.URI(uri)))
	}
}

// Search for artist.
func (s *impl) Artist(ctx *cli.Context) {
	handleErr(validateSingle(ctx.Args()))
	res, err := make(chan []sscc.Artist), make(chan error)
	s.ctrl.SearchArtist(ctx.Args().First(), res, err)
LOOP:
	for {
		select {
		case res := <-res:
			disp(res)
		case err := <-err:
			handleErr(err)
			break LOOP
		}
	}
	fmt.Println("")
	s.interactive(ctx)
}

// Search for album.
func (s *impl) Album(ctx *cli.Context) {
	handleErr(validateSingle(ctx.Args()))
	res, err := make(chan []sscc.Album), make(chan error, 1)
	s.ctrl.SearchAlbum(ctx.Args().First(), res, err)
LOOP:
	for {
		select {
		case res := <-res:
			disp(res)
		case err := <-err:
			handleErr(err)
			break LOOP
		}
	}
	fmt.Println("")
	s.interactive(ctx)
}

// Search for track.
func (s *impl) Track(ctx *cli.Context) {
	handleErr(validateSingle(ctx.Args()))
	res, err := make(chan []sscc.Track), make(chan error)
	s.ctrl.SearchTrack(ctx.Args().First(), res, err)
LOOP:
	for {
		select {
		case res := <-res:
			disp(res)
		case err := <-err:
			handleErr(err)
			break LOOP
		}
	}
	s.interactive(ctx)
}

func (s *impl) Status(ctx *cli.Context) {
	status, err := s.ctrl.Status()
	handleErr(err)
	fmt.Println(status)
}

func validateSingle(args cli.Args) error {
	n := len(args)
	if n == 0 {
		return fmt.Errorf("please specify valid argument")
	}
	if n > 1 {
		return fmt.Errorf("invalid number of arguments")
	}
	return nil
}

func disp(r interface{}) {
	for i := reflect.ValueOf(r).Len() - 1; i >= 0; i-- {
		for j, l := 0, reflect.ValueOf(r).Index(i).NumField(); j < l; j++ {
			f := reflect.ValueOf(r).Index(i).Field(j)
			if f.Kind() == reflect.Slice {
				fmt.Printf("%q\n", reflect.ValueOf(r).Index(i).Type().Field(j).Name)
				disp(f.Interface())
			} else {
				fmt.Printf("%q: %q",
					reflect.ValueOf(r).Index(i).Type().Field(j).Name, f.String())
			}
			if j < l-1 {
				fmt.Println("")
			}
		}
		fmt.Println("")
	}
}
