package observer

import (
	"crypto/x509"
	"errors"
	"strconv"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/lamassuiot/lamassu-virtual-device/pkg/config"
)

type DeviceState struct {
	// internal state
	Cert   *x509.Certificate
	Config config.Config

	observers []Observer
}

func (s *DeviceState) Attach(o Observer) (bool, error) {

	for _, observer := range s.observers {
		if observer == o {
			return false, errors.New("Observer already exists")
		}
	}
	s.observers = append(s.observers, o)
	return true, nil
}

func (s *DeviceState) Detach(o Observer) (bool, error) {

	for i, observer := range s.observers {
		if observer == o {
			s.observers = append(s.observers[:i], s.observers[i+1:]...)
			return true, nil
		}
	}
	return false, errors.New("Observer not found")
}

func (s *DeviceState) Notify(logger log.Logger) (bool, error) {
	level.Info(logger).Log("msg", "Obserer notify... "+strconv.Itoa(len(s.observers)))
	for _, observer := range s.observers {
		observer.Update(s)
	}
	return true, nil
}

func (s *DeviceState) SetCertificate(cert *x509.Certificate, logger log.Logger) {
	s.Cert = cert
	s.Notify(logger)
}
