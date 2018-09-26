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
	"time"

	"github.com/docopt/docopt-go"
	"github.com/kovetskiy/lorg"
	"github.com/reconquest/karma-go"
	"github.com/reconquest/sign-go"
	"github.com/seletskiy/mainframe/engine"
	"github.com/seletskiy/mainframe/fonts"
	"github.com/seletskiy/mainframe/protocol/messages"
	"github.com/seletskiy/mainframe/protocol/text"
	"github.com/seletskiy/mainframe/server"
)

var version = "1.0"

var usage = `mainframe — terminal.

Usage:
  terminal -h | --help
  terminal [options] [-s=<socket>] listen
  terminal [options] [-s=<socket>] exec [--options=] -- <command>...

Options:
  -h --help             Show this help.
  -s --socket <socket>  Path to control socket.
                         [default: /tmp/mainframe.sock]
  --profile <path>      Write CPU profile to specified file.
  --options <options>   Parameters for new window in text protocol format.
`

type Opts struct {
	Socket  string
	Command []string

	Listen bool
	Exec   bool

	Options string

	Profile string

	Separator bool `docopt:"--"`
}

var (
	log *lorg.Log
)

func main() {
	args, err := docopt.ParseArgs(usage, nil, "mainframe "+version)
	if err != nil {
		panic(err)
	}

	log = lorg.NewLog()
	log.SetLevel(lorg.LevelTrace)

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

	case opts.Exec:
		execute(opts)
	}
}

func listen(opts Opts) {
	font, err := fonts.Load("resources/font/font.tar")
	if err != nil {
		panic(err)
	}

	engine := engine.New(log)

	listener, err := server.Listen(opts.Socket, engine, log)
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

	var (
		notch = time.Now()
		fps   = 0
	)

	for engine.Running() {
		if time.Now().Sub(notch) > time.Second {
			log.Tracef("FPS: %d", fps)
			fps = 0
			notch = time.Now()
		}

		err := engine.Render()
		if err != nil {
			log.Fatal(err)
		}

		fps++
	}
}

func execute(opts Opts) {
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

	syscall.Dup2(int(socket.Fd()), 3)
	syscall.Exec(path, opts.Command, env)
}

func openWindow(opts Opts, connection net.Conn) error {
	message, err := text.Parse(`open ` + opts.Options)
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
