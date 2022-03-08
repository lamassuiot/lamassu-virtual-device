package deviceview

import (
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io/ioutil"

	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/log"
	"github.com/lamassuiot/lamassu-est/pkg/client"
	"github.com/lamassuiot/lamassu-virtual-device/pkg/observer"
	"github.com/rivo/tview"
)

func GetReenrollItem(logger log.Logger, deviceData *observer.DeviceState, app *tview.Application) tview.Primitive {
	dmsEndpoint := deviceData.Config.Dms.DmsEstServer

	form := tview.NewForm().
		AddInputField("DMS EST server endpoint", dmsEndpoint, 50, nil, func(text string) {
			dmsEndpoint = text
		}).
		AddButton("Re Enroll", func() {
			level.Info(logger).Log("msg", "Re-Enrolling... ", "with", dmsEndpoint)

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

			priv, err := ioutil.ReadFile("certificates/device.key")
			if err == nil {
				block, _ := pem.Decode([]byte(priv))
				privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)

				if err == nil {
					csrBytes, err := x509.CreateCertificateRequest(rand.Reader, &template, privateKey)
					if err != nil {
						level.Error(logger).Log("err", err)
					}

					csrPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrBytes})
					ioutil.WriteFile("certificates/device.csr", csrPEM, 0777)

					clienteEst, err := client.NewLamassuEstClient(
						dmsEndpoint,
						// "/home/ikerlan/lamassu/lamassu-dms/nogit/crypto-material/dms-ca.crt",
						"/home/ikerlan/lamassu-compose-v2/tls-certificates/downstream/tls.crt",
						"certificates/device.crt",
						"certificates/device.key",
						logger,
					)
					if err != nil {
						level.Error(logger).Log("err", err)
					}

					x509Csr, err := x509.ParseCertificateRequest(csrBytes)
					if err != nil {
						level.Error(logger).Log("err", err)
					}

					level.Info(logger).Log("csr", csrPEM)

					crt, err := clienteEst.Reenroll(x509Csr)
					if err != nil {
						level.Error(logger).Log("err", err)
					}

					certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: crt.Raw})
					ioutil.WriteFile("certificates/device.crt", certPEM, 0777)

					deviceData.SetCertificate(crt, logger)
					level.Info(logger).Log("msg", "Serial", fmt.Sprintf("%x", crt.SerialNumber), "Certificate content: "+crt.Subject.String()+" Issuer: "+crt.Issuer.String())
				} else {
					if err != nil {
						level.Error(logger).Log("err", err)
					}
				}
			} else {
				if err != nil {
					level.Error(logger).Log("err", err)
				}
			}

		})

	return form
}
