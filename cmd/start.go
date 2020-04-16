package cmd

import (
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
	"time"
	"github.com/spf13/viper"
	log "github.com/sirupsen/logrus"
	"github.com/graphicweave/injun/http"
	"github.com/graphicweave/injun/grpc"
	"github.com/graphicweave/injun/elastic"
	"context"
	"net"
)

var startCmd = &cobra.Command{
	Use:  "start",
	Long: "starts the injun search server",
	PreRun: func(cmd *cobra.Command, args []string) {

		// Set logrus formatter for GELF logging
		log.SetFormatter(&log.JSONFormatter{})

		// Parse env variables
		viper.AutomaticEnv()

		es, err := elastic.NewElasticSearch(context.Background())
		if err != nil {
			log.Info("failed to create elastic search client")
			log.Infoln("Error :", err)
			return
		}

		err = es.SetupIndex()
		if err != nil {
			log.Info("failed to create elastic search index")
			log.Infoln("Error :", err)
		}
		log.Infoln("determining local non-loopback IP")
		log.Infoln("non-loopback IP is: ", getIp())
	},
	Run: func(cmd *cobra.Command, args []string) {

		var gracefulStop = make(chan os.Signal)

		httpCloseChan := make(chan bool)
		grpcCloseChan := make(chan bool)

		signal.Notify(gracefulStop, syscall.SIGTERM)
		signal.Notify(gracefulStop, syscall.SIGINT)
		signal.Notify(gracefulStop, syscall.SIGKILL)

		go func() {
			sig := <-gracefulStop
			log.Infof("caught signal: %+v\n", sig)
			// TODO: refactor time
			log.Infoln("Waiting for 1 second to finish processing")
			time.Sleep(1 * time.Second)

			httpCloseChan <- true
			grpcCloseChan <- true

			os.Exit(0)
		}()

		httpAddr := viper.GetString("HTTP_HOST")
		grpcAddr := viper.GetString("GRPC_HOST")

		go func() {
			log.Infoln("started HTTP server on " + httpAddr)

			if err := http.StartServer(httpAddr); err != nil {
				log.Infoln("error while starting injun HTTP server")
				log.Fatalln(err.Error())
				httpCloseChan <- true
			}
		}()

		go func() {
			log.Infoln("started gRPC server on " + grpcAddr)

			if err := grpc.StartGRPCServer(grpcAddr); err != nil {
				log.Infoln("error while starting injun gRPC server")
				log.Fatalln(err.Error())
				grpcCloseChan <- true
			}
		}()

		<-httpCloseChan
		<-grpcCloseChan
	},
}

func getIp() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatalln(err)
		return ""
	}
	for _, v := range addrs {
		if ip, ok := v.(*net.IPNet); ok && !ip.IP.IsLoopback() {
			if ip.IP.To4() != nil {
				return ip.IP.String()
			}
		}
	}
	log.Fatalln("failed to determine local non-loopback IP")
	return ""
}
