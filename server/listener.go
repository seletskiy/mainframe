package server

import (
	"net"

	"github.com/apex/log"
	"github.com/kovetskiy/lorg"
	"github.com/reconquest/karma-go"
	"github.com/seletskiy/mainframe/engine"
)

type Listener struct {
	listener net.Listener
	closed   bool
	log      lorg.Logger
	engine   *engine.Engine
}

func Listen(
	path string,
	engine *engine.Engine,
	log lorg.Logger,
) (*Listener, error) {
	var err error

	listener := Listener{
		engine: engine,
		log:    log,
	}

	listener.listener, err = net.Listen("unix", path)
	if err != nil {
		return nil, karma.Format(
			err,
			"{protocol} unable to listen",
		)
	}

	go listener.accept()

	return &listener, nil
}

func (listener *Listener) Close() error {
	listener.closed = true

	err := listener.listener.Close()
	if err != nil {
		return err
	}

	return nil
}

func (listener *Listener) accept() {
	for {
		connection, err := listener.listener.Accept()
		if err != nil {
			if listener.closed {
				return
			}

			log.Error(
				karma.Format(err, "unable to accept connection").Error(),
			)

			continue
		}

		context, err := listener.engine.CreateWindow(640, 480, "mainframe")
		if err != nil {
			log.Error(
				karma.Format(err, "unable to create window").Error(),
			)

			continue
		}

		client := engine.Client{
			Connection: connection,
			Context:    context,
		}

		go client.Serve()
	}
}
