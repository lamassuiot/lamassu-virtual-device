package decodedcertview

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/lamassuiot/lamassu-virtual-device/pkg/observer"
	"github.com/rivo/tview"
)

type DeviceForm struct {
	DeviceID    string
	DeviceAlias string
}

var (
	grid   *tview.Grid
	logger *log.Logger
)

type concreteObserver struct {
}

func (s *concreteObserver) Update(t interface{}) {
	// do something
	level.Info(*logger).Log("Observer updated", t)
	grid.Clear()
	grid = DrawGrid(t.(*observer.DeviceState).Cert)
}

func DrawGrid(cert *x509.Certificate) *tview.Grid {
	_, _, keyBits, _ := GetPublicKeyInfo(*cert)

	grid.AddItem(tview.NewTextView().SetText("Serial Number").SetTextColor(tcell.ColorYellow), 0, 0, 1, 1, 0, 0, false)
	grid.AddItem(tview.NewTextView().SetText(fmt.Sprintf("%x", cert.SerialNumber)), 0, 1, 1, 1, 0, 0, false)

	grid.AddItem(tview.NewTextView().SetText("Certificate Subject").SetTextColor(tcell.ColorYellow), 1, 0, 1, 1, 0, 0, false)
	grid.AddItem(tview.NewTextView().SetText(cert.Subject.String()), 1, 1, 1, 1, 0, 0, false)

	grid.AddItem(tview.NewTextView().SetText("Issuer Subject").SetTextColor(tcell.ColorYellow), 2, 0, 1, 1, 0, 0, false)
	grid.AddItem(tview.NewTextView().SetText(cert.Issuer.String()), 2, 1, 1, 1, 0, 0, false)

	grid.AddItem(tview.NewTextView().SetText("Key Bits").SetTextColor(tcell.ColorYellow), 3, 0, 1, 1, 0, 0, false)
	grid.AddItem(tview.NewTextView().SetText(strconv.Itoa(keyBits)), 3, 1, 1, 1, 0, 0, false)

	return grid
}

func GetDecodedCertItem(inlog log.Logger, deviceData *observer.DeviceState) *tview.Grid {
	logger = &inlog
	concreteObserver := &concreteObserver{}
	deviceData.Attach(concreteObserver)

	grid = tview.NewGrid().
		SetRows(2, 2, 2, 2).
		SetColumns(30, 0).
		SetBorders(true)
	return DrawGrid(deviceData.Cert)
}

func GetPublicKeyInfo(cert x509.Certificate) (string, string, int, string) {
	if bytes.Compare(cert.Raw, []byte{}) == 0 {
		return "", "", 0, ""
	}

	key := cert.PublicKeyAlgorithm.String()
	var keyBits int
	switch key {
	case "RSA":
		keyBits = cert.PublicKey.(*rsa.PublicKey).N.BitLen()
	case "ECDSA":
		keyBits = cert.PublicKey.(*ecdsa.PublicKey).Params().BitSize
	}
	publicKeyDer, _ := x509.MarshalPKIXPublicKey(cert.PublicKey)
	publicKeyBlock := pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyDer,
	}
	publicKeyPem := string(pem.EncodeToMemory(&publicKeyBlock))

	var keyStrength string = "unknown"
	switch key {
	case "RSA":
		if keyBits < 2048 {
			keyStrength = "low"
		} else if keyBits >= 2048 && keyBits < 3072 {
			keyStrength = "medium"
		} else {
			keyStrength = "high"
		}
	case "ECDSA":
		if keyBits <= 128 {
			keyStrength = "low"
		} else if keyBits > 128 && keyBits < 256 {
			keyStrength = "medium"
		} else {
			keyStrength = "high"
		}
	}

	return publicKeyPem, key, keyBits, keyStrength
}
