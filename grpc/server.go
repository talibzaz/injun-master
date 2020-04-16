package grpc

import (
	"net"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"github.com/graphicweave/injun/proto"

	evtService "github.com/graphicweave/injun/grpc/event"
	mailService  "github.com/graphicweave/injun/grpc/mail"
	"github.com/graphicweave/injun/mail"
)

func StartGRPCServer(addr string) error {

	listener, err := net.Listen("tcp", addr)

	if err != nil {
		log.Errorln("failed to create a TCP listener for gRPC server")
		log.Fatalln("Error:", err.Error())
		return err
	}
	server := grpc.NewServer()
    mail.RegisterMailServiceServer(server, mailService.MailService{})
	event.RegisterEventServiceServer(server, evtService.EventService{})

	return server.Serve(listener)
}
