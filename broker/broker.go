package broker

import (
	"net/http"
	"os"

	"code.cloudfoundry.org/lager"
	"github.com/henrytk/aws-service-broker/provider"
	usb "github.com/henrytk/universal-service-broker/broker"
)

func NewAWSServiceBroker(config usb.Config, awsProvider *provider.AWSProvider) http.Handler {
	logger := lager.NewLogger("aws-service-broker")
	logger.RegisterSink(lager.NewWriterSink(os.Stdout, config.API.LagerLogLevel))

	serviceBroker := usb.New(config, awsProvider, logger)
	return usb.NewAPI(serviceBroker, logger, config)
}
