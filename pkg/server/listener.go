package server

import (
	"net"

	"github.com/reconquest/karma-go"
	"github.com/seletskiy/mainframe/pkg/engine"
	"github.com/seletskiy/mainframe/pkg/log"
)

type Listener struct {
	listener net.Listener
	closed   bool
	engine   *engine.Engine
}

func Listen(
	path string,
	engine *engine.Engine,
) (*Listener, error) {
	var err error

	listener := Listener{
		engine: engine,
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

		client := engine.Client{
			Connection: connection,
			Engine:     listener.engine,
		}

		go client.Serve()
	}
}
