package main

import (
	"log"
	"net"
	"time"

	"github.com/hashicorp/consul/api"
)

const ttl = time.Second * 8

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Start() {
	ln, err := net.Listen("tcp", ":3000")

	if err != nil {
		log.Fatal(err)
	}

	s.acceptLoop(ln)

}

func (s *Service) acceptLoop(ln net.Listener) {

	for {
		_, err := ln.Accept()

		if err != nil {
			log.Fatal(err)
		}
	}
}

func (s *Service) registerService() {
	check := &api.AgentServiceCheck{
		DeregisterCriticalServiceAfter: ttl.String(),
		TLSSkipVerify:                  true,
		TTL:                            ttl.String(),
		CheckID:                        "checkalive",
	}

	register := &api.AgentServiceRegistration{
		ID:      "login_Service",
		Name:    "my-cluster",
		Tags:    []string{"login"},
		Address: "127.0.0.1",
		Port:    3000,
		Check:   check,
	}

}

func (s *Service) Init() {

}

func main() {

}
