package historicalcertview

import (
	"bytes"
	"crypto/x509"
	"fmt"
	"strings"

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
	view                     *tview.TextView
	logger                   *log.Logger
	historical_serials_certs []string
)

func (s *concreteObserver) Update(t interface{}) {
	// do something
	level.Info(*logger).Log("Observer updated", t)
	view = DrawView(t.(*observer.DeviceState).Cert)
}

func DrawView(cert *x509.Certificate) *tview.TextView {
	if bytes.Compare(cert.Raw, []byte{}) != 0 {
		historical_serials_certs = append(historical_serials_certs, fmt.Sprintf("%x", cert.SerialNumber))
	}

	inverted_order := historical_serials_certs
	for i, j := 0, len(inverted_order)-1; i < j; i, j = i+1, j-1 {
		inverted_order[i], inverted_order[j] = inverted_order[j], inverted_order[i]
	}

	view.SetText(strings.Join(inverted_order, "\n"))
	view.SetBorder(true).SetBorderPadding(0, 0, 1, 1)

	return view
}

func GetHistoricalCertItem(inlogger log.Logger, deviceData *observer.DeviceState) tview.Primitive {
	logger = &inlogger

	view = tview.NewTextView()
	view.SetTitle("Historical Serial Numbers")

	concreteObserver := &concreteObserver{}
	deviceData.Attach(concreteObserver)

	return DrawView(deviceData.Cert)
}
