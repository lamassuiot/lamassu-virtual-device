package deviceview

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io/ioutil"

	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/log"
	"github.com/lamassuiot/lamassu-est/pkg/client"
	"github.com/lamassuiot/lamassu-virtual-device/pkg/observer"
	"github.com/rivo/tview"
)

func GetEnrollItem(logger log.Logger, deviceData *observer.DeviceState, app *tview.Application) tview.Primitive {
	aps := "Lamassu-Root-CA1-RSA4096"
	dmsEndpoint := deviceData.Config.Dms.DmsEstServer

	form := tview.NewForm().
		AddInputField("DMS EST server endpoint", dmsEndpoint, 50, nil, func(text string) {
			dmsEndpoint = text
		}).
		AddDropDown("CAs", []string{"Lamassu-Root-CA1-RSA4096", "Lamassu-Root-CA3-ECC384", "Lamassu-Root-CA2-RSA2048", "Lamassu-Root-CA4-ECC256", "test1"}, 0, func(option string, optionIndex int) {
			aps = option
		}).
		AddButton("Enroll", func() {
			level.Info(logger).Log("msg", "Enrolling... ", "with", aps+","+dmsEndpoint)

			privKey, _ := rsa.GenerateKey(rand.Reader, 4096)
			privkeyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privKey)})

			// generate private key
			subj := pkix.Name{
				CommonName:         deviceData.Config.Device.Id,
				Country:            []string{deviceData.Config.Device.EnrollingProperties.Subject.Country},
				Province:           []string{deviceData.Config.Device.EnrollingProperties.Subject.State},
				Locality:           []string{deviceData.Config.Device.EnrollingProperties.Subject.Locality},
				Organization:       []string{deviceData.Config.Device.EnrollingProperties.Subject.Organization},
				OrganizationalUnit: []string{deviceData.Config.Device.EnrollingProperties.Subject.OrganizationalUnit},
			}

			template := x509.CertificateRequest{
				Subject:            subj,
				SignatureAlgorithm: x509.SHA256WithRSA,
			}

			csrBytes, _ := x509.CreateCertificateRequest(rand.Reader, &template, privKey)
			csrPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrBytes})
			ioutil.WriteFile("certificates/device.csr", csrPEM, 0777)
			ioutil.WriteFile("certificates/device.key", privkeyPEM, 0777)

			clienteEst, _ := client.NewLamassuEstClient(
				dmsEndpoint,
				// "/home/ikerlan/lamassu/lamassu-dms/nogit/crypto-material/dms-ca.crt",
				"/home/ikerlan/lamassu-compose-v2/tls-certificates/downstream/tls.crt",
				"certificates/bootstrap.crt",
				"certificates/bootstrap.key",
				logger,
			)

			x509Csr, err := x509.ParseCertificateRequest(csrBytes)
			if err != nil {
				level.Error(logger).Log("err", err)
			}

			crt, err := clienteEst.Enroll(aps, x509Csr)
			if err != nil {
				level.Error(logger).Log("err", err)
			}

			certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: crt.Raw})
			ioutil.WriteFile("certificates/device.crt", certPEM, 0777)

			deviceData.SetCertificate(crt, logger)

			level.Info(logger).Log("msg", "Certificate content: "+crt.Subject.String()+" Issuer: "+crt.Issuer.String())
		})

	// gauge := components.NewActivityModeGauge()
	// gauge.SetPgBgColor(tcell.ColorOrange)
	// gauge.SetBorder(true)

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(form, 7, 0, false)
	// AddItem(gauge, 3, 0, false)

	// update := func() {
	// 	tick := time.NewTicker(20 * time.Millisecond)
	// 	level.Info(logger).Log("msg", "tic")
	// 	for {
	// 		select {
	// 		case <-tick.C:
	// 			gauge.Pulse()
	// 			app.Draw()
	// 		}
	// 	}
	// }
	// go update()

	return flex
}
