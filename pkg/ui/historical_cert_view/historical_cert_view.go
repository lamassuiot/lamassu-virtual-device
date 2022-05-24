package historicalcertview

import (
	"crypto/x509"
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
	view   *tview.TextView
	app    *tview.Application
	logger *log.Logger
)

func (s *concreteObserver) Update(t interface{}) {
	level.Info(*logger).Log("Observer updated", t)
	view = DrawView(t.(*observer.DeviceState).Cert, t.(*observer.DeviceState))
}

func DrawView(cert *x509.Certificate, deviceData *observer.DeviceState) *tview.TextView {

	view.SetText(strings.Join(deviceData.SN, "\n"))
	view.SetBorder(true).SetBorderPadding(0, 0, 1, 1)

	return view
}

func ClearView() *tview.TextView {
	view = tview.NewTextView()
	view.SetTitle("Historical Serial Numbers3")

	view.SetBorder(true).SetBorderPadding(0, 0, 1, 1)

	return view
}

func GetHistoricalCertItem(inlogger log.Logger, deviceData *observer.DeviceState, inapp *tview.Application) tview.Primitive {
	logger = &inlogger

	app = inapp

	view = tview.NewTextView()
	view.SetTitle("Historical Serial Numbers")

	concreteObserver := &concreteObserver{}
	deviceData.Attach(concreteObserver)

	return DrawView(deviceData.Cert, deviceData)
}
