package deviceview

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"

	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/log"

	"github.com/lamassuiot/lamassu-virtual-device/pkg/observer"
	"github.com/lamassuiot/lamassu-virtual-device/pkg/utils"
	"github.com/lamassuiot/lamassuiot/pkg/est/client"
	"github.com/rivo/tview"
)

func GetReenrollItem(logger log.Logger, deviceData *observer.DeviceState, app *tview.Application) tview.Primitive {
	devManager := deviceData.Config.Devmanager.EstServer

	statusTextView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetTextAlign(tview.AlignCenter).
		SetText(" ")
	level.Info(logger).Log("msg", devManager)
	form := tview.NewForm().
		AddInputField("Device Manager EST server endpoint", devManager, 50, nil, func(text string) {
			devManager = text
		}).
		AddButton("Re Enroll", func() {
			level.Info(logger).Log("msg", "Re-Enrolling... ", "with", devManager)
			statusTextView.SetText("Re-Enrolling...")
			serverCert, err := utils.ReadCertPool(deviceData.Config.Devmanager.Cert)
			if err != nil {
				level.Error(logger).Log("err", err)
			}
			deviceCert, err := utils.ReadCert(deviceData.PathCertificate)
			if err != nil {
				level.Error(logger).Log("err", err)
			}
			deviceKey, err := utils.ReadKey(deviceData.PathKey)
			if err != nil {
				level.Error(logger).Log("err", err)
			}
			clienteEst, err := client.NewLamassuEstClient(
				devManager,
				serverCert,
				deviceCert,
				deviceKey,
				logger,
			)
			if err != nil {
				level.Error(logger).Log("err", err)
				app.Stop()
			}
			csrContent, err := ioutil.ReadFile(deviceData.CSR)
			if err != nil {
				level.Error(logger).Log("err", err)
			}

			cpb, _ := pem.Decode(csrContent)
			x509Csr, err := x509.ParseCertificateRequest(cpb.Bytes)
			level.Info(logger).Log("msg", x509Csr)
			if err != nil {
				level.Error(logger).Log("err", err)
				statusTextView.SetText(err.Error())
			}
			var ctx context.Context
			crt, err := clienteEst.Reenroll(ctx, x509Csr)
			if err != nil {
				level.Error(logger).Log("err", err)
				statusTextView.SetText(err.Error())
			}
			if err == nil {
				sn := utils.InsertNth(utils.ToHexInt(crt.SerialNumber), 2)
				certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: crt.Raw})
				ioutil.WriteFile(deviceData.PathDevicesCertFolder+sn+".crt", certPEM, 0777)
				deviceData.SetSN(sn, logger)
				deviceData.SetCertificate(crt, logger)
				app.ForceDraw()
				level.Info(logger).Log("msg", "Serial", fmt.Sprintf("%x", crt.SerialNumber), "Certificate content: "+crt.Subject.String()+" Issuer: "+crt.Issuer.String())

				statusTextView.SetText("Reenroll success!")

			}
		})

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(form, 7, 1, false).
		AddItem(statusTextView, 7, 1, false)

	return flex
}
