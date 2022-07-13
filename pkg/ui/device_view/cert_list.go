package deviceview

import (
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"os"
	"strings"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/lamassuiot/lamassu-virtual-device/pkg/observer"

	"github.com/rivo/tview"
)

func DrawCertDropDown(app *tview.Application, obs *observer.DeviceState, logger log.Logger) tview.Primitive {
	var deviceCertificate, deviceKey, deviceCsr, deviceCertificateFolder string
	form := tview.NewForm()
	form2 := tview.NewForm()
	form3 := tview.NewForm()
	form.
		AddButton("Refresh", func() {
			form.Clear(false)
			form2.Clear(false)
			form.AddDropDown("Device ID", getDeviceIds(obs.PathDevicesFolder), 0, func(option string, optionIndex int) {
				if len(option) > 0 && optionIndex > 0 {
					obs.ClearSN(logger)
					level.Info(logger).Log("msg", option)
					name := strings.Split(option, "-")
					id := name[len(name)-5] + "-" + name[len(name)-4] + "-" + name[len(name)-3] + "-" + name[len(name)-2] + "-" + name[len(name)-1]
					deviceCertificateFolder = obs.PathDevicesFolder + option + "/certificates/"
					deviceKey = obs.PathDevicesFolder + option + "/" + id + ".key"
					deviceCsr = obs.PathDevicesFolder + option + "/" + id + ".csr"
					for _, sn := range getCertSN(deviceCertificateFolder) {
						obs.SetSN(sn, logger)
					}
					form2.Clear(false)
					form2.AddDropDown("Cert SN", obs.SN, 0, func(option string, optionIndex int) {
						if len(option) > 0 && optionIndex > 0 {

							deviceCertificate = deviceCertificateFolder + option + ".crt"

							certContent, err := ioutil.ReadFile(deviceCertificate)
							if err != nil {
								level.Info(logger).Log("err", err)
							}
							cpb, _ := pem.Decode(certContent)

							if err == nil {
								cert, err := x509.ParseCertificate(cpb.Bytes)
								if err == nil {
									level.Info(logger).Log("msg", "Set Certificate Refresh")
									obs.SetPath(deviceCertificate, deviceKey, deviceCsr, logger)
									obs.SetPathDevicesCertFolder(deviceCertificateFolder, logger)
									obs.SetCertificate(cert, logger)
									app.ForceDraw()
								}
							}

						}
					})
					app.ForceDraw()
				}

			})
			app.ForceDraw()
		}).
		AddDropDown("Device ID", getDeviceIds(obs.PathDevicesFolder), 0, func(option string, optionIndex int) {
			if len(option) > 0 && optionIndex > 0 {
				obs.ClearSN(logger)
				level.Info(logger).Log("msg", option)
				name := strings.Split(option, "-")
				id := name[len(name)-5] + "-" + name[len(name)-4] + "-" + name[len(name)-3] + "-" + name[len(name)-2] + "-" + name[len(name)-1]
				deviceCertificateFolder = obs.PathDevicesFolder + option + "/certificates/"
				deviceKey = obs.PathDevicesFolder + option + "/" + id + ".key"
				deviceCsr = obs.PathDevicesFolder + option + "/" + id + ".csr"

				for _, sn := range getCertSN(deviceCertificateFolder) {
					obs.SetSN(sn, logger)
				}
				form2.Clear(false)
				form2.AddDropDown("Cert SN", obs.SN, 0, func(option string, optionIndex int) {
					if len(option) > 0 && optionIndex > 0 {

						deviceCertificate = deviceCertificateFolder + option + ".crt"

						certContent, err := ioutil.ReadFile(deviceCertificate)
						if err != nil {
							level.Info(logger).Log("err", err)
						}
						cpb, _ := pem.Decode(certContent)

						if err == nil {
							cert, err := x509.ParseCertificate(cpb.Bytes)
							if err == nil {
								level.Info(logger).Log("msg", "Set Certificate")
								obs.SetPath(deviceCertificate, deviceKey, deviceCsr, logger)
								obs.SetPathDevicesCertFolder(deviceCertificateFolder, logger)
								obs.SetCertificate(cert, logger)
								app.ForceDraw()
							}
						}
						app.ForceDraw()
					}
				})
				app.ForceDraw()

			}
		})
	form3.
		AddButton("Update", func() {
			form2.Clear(false)
			form2.AddDropDown("Cert SN", obs.SN, 0, func(option string, optionIndex int) {
				if len(option) > 0 && optionIndex > 0 {

					deviceCertificate = deviceCertificateFolder + option + ".crt"

					certContent, err := ioutil.ReadFile(deviceCertificate)
					if err != nil {
						level.Info(logger).Log("err", err)
					}
					cpb, _ := pem.Decode(certContent)

					if err == nil {
						cert, err := x509.ParseCertificate(cpb.Bytes)
						if err == nil {
							obs.SetPath(deviceCertificate, deviceKey, deviceCsr, logger)
							obs.SetPathDevicesCertFolder(deviceCertificateFolder, logger)
							obs.SetCertificate(cert, logger)
							app.ForceDraw()
						}
					}

				}
			})
			app.ForceDraw()

		})
	app.ForceDraw()
	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(form, 5, 1, false).
		AddItem(form2, 4, 1, false).
		AddItem(form3, 4, 1, false)
	return flex
}

func DevicesListForm(app *tview.Application, obs *observer.DeviceState, logger log.Logger) tview.Primitive {
	form := DrawCertDropDown(app, obs, logger)
	return form
}

func getDeviceIds(path string) []string {
	files, _ := ioutil.ReadDir(path)
	var s []string
	s = append(s, "")

	for _, file := range files {
		fileName := file.Name()
		if len(fileName) > 0 {
			if !contains(s, fileName) {
				s = append(s, fileName)
			}
		}

	}
	return s
}
func getCertSN(path string) []string {
	files, _ := ioutil.ReadDir(path)
	var s []string
	s = append(s, "")

	for _, file := range files {
		fileName := file.Name()
		if len(fileName) > 0 {
			fileName := fileName[:len(fileName)-4]
			if !contains(s, fileName) {
				s = append(s, fileName)
			}
		}

	}
	return s
}
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
