package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/hashicorp/consul/api"
)

const (
	ttl     = time.Second * 8
	checkId = "check_health"
)

type Service struct {
	consulClient *api.Client
}

func NewService() *Service {
	client, err := api.NewClient(&api.Config{})

	if err != nil {
		log.Fatal(err)
	}
	return &Service{
		consulClient: client,
	}
}

func (s *Service) Start() {
	ln, err := net.Listen("tcp", ":3000")

	if err != nil {
		log.Fatal(err)
	}
	s.registerService()
	fmt.Println("Service Registered")
	go s.updateHealthCheck()
	s.acceptLoop(ln)

}

func (s *Service) acceptLoop(ln net.Listener) {
	fmt.Println("Accepting Connections")
	for {
		_, err := ln.Accept()

		if err != nil {
			log.Fatal(err)
		}
	}
}

func (s *Service) updateHealthCheck() {
	ticker := time.NewTicker(time.Second * 5)

	for {
		err := s.consulClient.Agent().UpdateTTL(checkId, "online", api.HealthPassing)
		if err != nil {
			log.Fatal(err)
		}
		<-ticker.C
	}
}

func (s *Service) registerService() {
	check := &api.AgentServiceCheck{
		DeregisterCriticalServiceAfter: ttl.String(),
		TLSSkipVerify:                  true,
		TTL:                            ttl.String(),
		CheckID:                        checkId,
	}

	register := &api.AgentServiceRegistration{
		ID:      "login_Service",
		Name:    "my-cluster",
		Tags:    []string{"login"},
		Address: "127.0.0.1",
		Port:    3000,
		Check:   check,
	}

	// query := map[string]any{
	// 	"type":       "service",
	// 	"service":    "mycluster",
	// 	"passinonly": true,
	// }

	// plan, err := watch.Parse(query)

	// if err != nil {
	// 	log.Fatal(err)
	// }

	err := s.consulClient.Agent().ServiceRegister(register)

	if err != nil {
		log.Fatal(err)
	}

}

func main() {

	s := NewService()
	s.Start()
}
