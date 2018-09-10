package main

import (
	"net"
	"os"
	"os/exec"
	"syscall"

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
  terminal [-s=<socket>] -L
  terminal [-s=<socket>] -E -- <command>...

Options:
  -h --help             Show this help.
  -L --listen           Listen control socket.
  -s --socket <socket>  Path to control socket.
                         [default: /tmp/mainframe.sock]
  -E --execute          Exec specified command.
`

type Opts struct {
	Socket  string
	Command []string

	Listen  bool
	Execute bool

	Separator bool `docopt:"--"`
}

var (
	log lorg.Logger
)

func main() {
	args, err := docopt.ParseArgs(usage, nil, "mainframe "+version)
	if err != nil {
		panic(err)
	}

	log = lorg.NewLog()

	var opts Opts

	err = args.Bind(&opts)
	if err != nil {
		log.Fatal(err)
	}

	switch {
	case opts.Listen:
		listen(opts)

	case opts.Execute:
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

	for engine.Running() {
		err := engine.Render()
		if err != nil {
			log.Fatal(err)
		}
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
