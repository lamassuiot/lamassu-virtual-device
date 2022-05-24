// Demo code for the Flex primitive.
package main

import (
	//"fmt"

	//"encoding/base64"

	"crypto/x509"
	"os"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/lamassuiot/lamassu-virtual-device/pkg/config"
	"github.com/lamassuiot/lamassu-virtual-device/pkg/observer"
	decodedcertview "github.com/lamassuiot/lamassu-virtual-device/pkg/ui/decoded_cert_view"
	deviceview "github.com/lamassuiot/lamassu-virtual-device/pkg/ui/device_view"
	historicalcertview "github.com/lamassuiot/lamassu-virtual-device/pkg/ui/historical_cert_view"
	"github.com/rivo/tview"

	rawcertview "github.com/lamassuiot/lamassu-virtual-device/pkg/ui/raw_cert_view"
)

func main() {

	var logger log.Logger
	f, _ := os.OpenFile("./vdevice.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)

	{
		logger = log.NewLogfmtLogger(log.NewSyncWriter(f))
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = level.NewFilter(logger, level.AllowDebug())
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	cfg, _ := config.NewConfig()

	obs := observer.DeviceState{
		Cert:                  &x509.Certificate{},
		Config:                cfg,
		PathCertificate:       "",
		PathKey:               "",
		CSR:                   "",
		PathDevicesFolder:     cfg.CertificatesDirectory,
		AwsEndpoint:           cfg.Aws.IotCoreEndpoint,
		AwsEndpointCAFile:     cfg.Aws.IotCoreEndpointCACertFile,
		PathDevicesCertFolder: "",
	}

	app := tview.NewApplication()

	flex := tview.NewFlex().
		AddItem(deviceview.GetItem(logger, &obs, app), 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(decodedcertview.GetDecodedCertItem(logger, &obs), 15, 0, false).
			AddItem(rawcertview.GetRawCertItem(logger, &obs), 0, 1, false), 0, 1, false).
		AddItem(historicalcertview.GetHistoricalCertItem(logger, &obs, app), 65, 0, false)

	if err := app.SetRoot(flex, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
