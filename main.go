package main

import (
	"net"
	"os"
	"os/exec"
	"runtime/pprof"
	"syscall"
	"time"

	"github.com/docopt/docopt-go"
	"github.com/kovetskiy/lorg"
	"github.com/reconquest/sign-go"
	"github.com/seletskiy/mainframe/engine"
	"github.com/seletskiy/mainframe/fonts"
	"github.com/seletskiy/mainframe/server"
)

var version = "1.0"

var usage = `mainframe — terminal.

Usage:
  terminal -h | --help
  terminal [options] [-s=<socket>] listen
  terminal [options] [-s=<socket>] exec -- <command>...

Options:
  -h --help             Show this help.
  -s --socket <socket>  Path to control socket.
                         [default: /tmp/mainframe.sock]
  --profile <path>      Write CPU profile to specified file.
`

type Opts struct {
	Socket  string
	Command []string

	Listen bool
	Exec   bool

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
	addr, err := net.ResolveUnixAddr("unix", opts.Socket)
	if err != nil {
		log.Fatal(err)
	}

	connection, err := net.DialUnix("unix", nil, addr)
	if err != nil {
		log.Fatal(err)
	}

	socket, err := connection.File()
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
