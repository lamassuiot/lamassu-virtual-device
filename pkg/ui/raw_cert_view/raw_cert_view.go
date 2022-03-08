package rawcertview

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/lamassuiot/lamassu-virtual-device/pkg/observer"
	"github.com/rivo/tview"
)

type DeviceForm struct {
	DeviceID    string
	DeviceAlias string
}

type concreteObserver struct {
}

var (
	view   *tview.TextView
	logger *log.Logger
)

func (s *concreteObserver) Update(t interface{}) {
	// do something
	level.Info(*logger).Log("Observer updated", t)
	view = DrawView(t.(*observer.DeviceState).Cert)
}

func DrawView(cert *x509.Certificate) *tview.TextView {
	txt := "This device has not been yet enrolled"
	if bytes.Compare(cert.Raw, []byte{}) != 0 {
		txt = string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw}))
	}

	view.SetText(txt)
	view.SetBorder(true).SetBorderPadding(0, 0, 1, 1)

	return view
}

func GetRawCertItem(inlogger log.Logger, deviceData *observer.DeviceState) tview.Primitive {
	logger = &inlogger

	view = tview.NewTextView()

	concreteObserver := &concreteObserver{}
	deviceData.Attach(concreteObserver)

	return DrawView(deviceData.Cert)
}
