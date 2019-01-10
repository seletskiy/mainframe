package main

import (
	"bufio"
	"io"
	"net"
	"os"
	"os/exec"
	"runtime/pprof"
	"strings"
	"syscall"

	"github.com/docopt/docopt-go"
	"github.com/reconquest/karma-go"
	"github.com/reconquest/sign-go"
	"github.com/seletskiy/mainframe/pkg/engine"
	"github.com/seletskiy/mainframe/pkg/fonts"
	"github.com/seletskiy/mainframe/pkg/log"
	"github.com/seletskiy/mainframe/pkg/protocol/messages"
	"github.com/seletskiy/mainframe/pkg/protocol/text"
	"github.com/seletskiy/mainframe/pkg/server"
)

var version = "1.0"

var usage = `mainframe — terminal.

Usage:
  terminal -h | --help
  terminal [options] [-s=<socket>] listen
  terminal [options] [-s=<socket>] open [--open-args=] -- <command>...

Options:
  -h --help              Show this help.
  -s --socket <socket>   Path to control socket.
                          [default: /tmp/mainframe.sock]
  --profile <path>       Write CPU profile to specified file.
  --open-args <options>  Parameters for new window in text protocol format.
  --font <path>          Font file to use.
  --font-size <size>     Font size to use in points. [default: 14]
  --font-dpi <dpi>       Screen DPI to render font for. [default: 72]
`

type Opts struct {
	Socket string `docopt:"--socket"`

	Listen bool `docopt:"listen"`

	Font     string  `docopt:"--font"`
	FontDPI  float64 `docopt:"--font-dpi"`
	FontSize float64 `docopt:"--font-size"`

	Profile string `docopt:"--profile"`

	Open     bool     `docopt:"open"`
	OpenArgs string   `docopt:"--open-args"`
	Command  []string `docopt:"<command>"`

	Separator bool `docopt:"--"`
}

func main() {
	args, err := docopt.ParseArgs(usage, nil, "mainframe "+version)
	if err != nil {
		panic(err)
	}

	var opts Opts

	err = args.Bind(&opts)
	if err != nil {
		log.Fatal(err)
	}

	if opts.Profile != "" {
		profile, err := os.Create("cpu.profile")
		if err != nil {
			panic(err)
		}

		pprof.StartCPUProfile(profile)
		defer pprof.StopCPUProfile()
	}

	switch {
	case opts.Listen:
		listen(opts)

	case opts.Open:
		open(opts)
	}
}

func listen(opts Opts) {
	font, err := fonts.Load(
		opts.Font,
		fonts.FontDPI(opts.FontDPI),
		fonts.FontSize(opts.FontSize),
		fonts.FontHinting(true),
	)
	if err != nil {
		panic(err)
	}

	engine := engine.New()

	listener, err := server.Listen(opts.Socket, engine)
	if err != nil {
		log.Fatal(err)
	}

	engine.SetFont(font)

	go sign.Notify(
		func(os.Signal) bool {
			listener.Close()
			engine.Stop()

			return false
		},
		syscall.SIGINT,
		syscall.SIGKILL,
		syscall.SIGTERM,
	)

	err = engine.Init()
	if err != nil {
		log.Fatal(err)
	}

	for engine.Running() {
		err := engine.Loop()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func open(opts Opts) {
	connection, err := net.Dial("unix", opts.Socket)
	if err != nil {
		log.Fatal(err)
	}

	err = openWindow(opts, connection)
	if err != nil {
		log.Fatal(err)
	}

	socket, err := connection.(*net.UnixConn).File()
	if err != nil {
		log.Fatal(err)
	}

	path, err := exec.LookPath(opts.Command[0])
	if err != nil {
		log.Fatal(err)
	}

	env := os.Environ()

	// TODO: allow to pass file descriptor in --fd option.
	syscall.Dup2(int(socket.Fd()), 3)
	syscall.Exec(path, opts.Command, env)
}

func openWindow(opts Opts, connection net.Conn) error {
	message, err := text.Parse(`open ` + opts.OpenArgs)
	if err != nil {
		return err
	}

	options := message.(*messages.Open)

	if options.Title == "" {
		options.Title = "mainframe: " + strings.Join(opts.Command, " ")
	}

	_, err = connection.Write(text.Serialize(options))
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(connection)

	if !scanner.Scan() {
		return io.EOF
	}

	response, err := text.Parse(scanner.Text())
	if err != nil {
		return err
	}

	if response.Tag() != "ok" {
		return karma.Format(
			err,
			"error while creating window",
		)
	}

	return nil
}
