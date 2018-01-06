package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"code.cloudfoundry.org/lager"
	"github.com/henrytk/aws-service-broker/provider"
	"github.com/henrytk/universal-service-broker/broker"
)

var configFilePath string

func main() {
	flag.StringVar(&configFilePath, "config", "", "Location of the config file")
	flag.Parse()

	file, err := os.Open(configFilePath)
	if err != nil {
		log.Fatalf("Error opening config file %s: %s\n", configFilePath, err)
	}
	defer file.Close()

	config, err := broker.NewConfig(file)
	if err != nil {
		log.Fatalf("Error validating config file: %v\n", err)
	}

	awsProvider, err := provider.NewAWSProvider(config.Provider)
	if err != nil {
		log.Fatalf("Error creating AWS Provider: %v\n", err)
	}

	logger := lager.NewLogger("aws-service-broker")
	logger.RegisterSink(lager.NewWriterSink(os.Stdout, config.API.LagerLogLevel()))

	serviceBroker := broker.New(config, awsProvider, logger)
	server := broker.NewAPI(serviceBroker, logger, config)

	listener, err := net.Listen("tcp", ":"+config.API.Port)
	if err != nil {
		log.Fatalf("Error listening to port %s: %s", config.API.Port, err)
	}
	fmt.Println("AWS Service Broker started on port " + config.API.Port + "...")
	http.Serve(listener, server)
}
