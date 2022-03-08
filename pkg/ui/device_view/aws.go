package deviceview

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/log"
	"github.com/lamassuiot/lamassu-virtual-device/pkg/observer"

	"github.com/rivo/tview"
)

func GetAwsItem(logger log.Logger, deviceData *observer.DeviceState, app *tview.Application) tview.Primitive {
	awsEndpoint := deviceData.Config.Aws.IotCoreEndpoint
	deviceId := deviceData.Config.Device.Id
	lambdaText := "Im " + deviceId + " device :)"
	form := tview.NewForm().
		AddInputField("Iot Core Endpoint", awsEndpoint, 50, nil, func(text string) {
			awsEndpoint = text
		}).
		AddInputField("Lambda Text", lambdaText, 50, nil, func(text string) {
			deviceData.Config.Device.Id = text
		}).
		AddButton("Connect to IoT Core", func() {
			tlsconfig := NewTLSConfig()

			opts := MQTT.NewClientOptions()
			opts.AddBroker("tls://" + awsEndpoint + ":8883")
			opts.SetClientID(deviceId).SetTLSConfig(tlsconfig)
			opts.SetDefaultPublishHandler(f)
			level.Info(logger).Log("AWS endpoint", awsEndpoint)

			// Start the connection
			c := MQTT.NewClient(opts)
			if token := c.Connect(); token.Wait() && token.Error() != nil {
				level.Info(logger).Log("Token error", token.Error())

				panic(token.Error())
			}

			level.Info(logger).Log("IsConnected", c.IsConnected())

			c.Publish(deviceId+"/hello", 0, false, lambdaText)
			level.Info(logger).Log("Message published", true)

			c.Disconnect(250)

		})

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(form, 7, 1, false).
		AddItem(tview.NewBox().SetBorder(true).SetTitle("Bottom (5 rows)"), 0, 1, false)

	return flex
}

func NewTLSConfig() *tls.Config {
	// Import trusted certificates from CAfile.pem.
	// Alternatively, manually add CA certificates to
	// default openssl CA bundle.
	certpool := x509.NewCertPool()
	pemCerts, err := ioutil.ReadFile("certificates/aws-ca.crt")
	if err == nil {
		certpool.AppendCertsFromPEM(pemCerts)
	}

	// Import client certificate/key pair
	cert, err := tls.LoadX509KeyPair("certificates/device.crt", "certificates/device.key")
	if err != nil {
		panic(err)
	}

	// Just to print out the client certificate..
	cert.Leaf, err = x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		panic(err)
	}

	// Create tls.Config with desired tls properties
	return &tls.Config{
		// RootCAs = certs used to verify server cert.
		RootCAs: certpool,
		// ClientAuth = whether to request cert from server.
		// Since the server is set up for SSL, this happens
		// anyways.
		ClientAuth: tls.NoClientCert,
		// ClientCAs = certs used to validate client cert.
		ClientCAs: nil,
		// InsecureSkipVerify = verify that cert contents
		// match server. IP matches what is in cert etc.
		InsecureSkipVerify: true,
		// Certificates = list of certs client sends to server.
		Certificates: []tls.Certificate{cert},
	}
}

var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}
