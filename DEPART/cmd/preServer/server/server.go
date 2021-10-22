package server

import (
	"preServer/conf"
	"preServer/rpc"
	"sync"
)

type server struct {
	config *conf.ServConfig

	httpServer *rpc.HttpServer
	stop       chan struct{}
	lock       sync.RWMutex
}

func NewServer(config *conf.ServConfig) (*server, error) {
	serv := &server{
		config:     config,
		httpServer: rpc.NewHttpServer(config.HttpCfg),
		stop:       make(chan struct{}),
	}
	return serv, nil
}

func (s *server) Start() error {
	s.httpServer.Start()

	return nil
}

func (s *server) Wait() {
	stop := s.stop
	<-stop
}

func (s *server) Stop() {
	s.lock.Lock()
	defer s.lock.Unlock()

	close(s.stop)
}
