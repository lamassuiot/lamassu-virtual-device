package deviceview

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	//"github.com/go-kit/kit/log/level"
	"github.com/gdamore/tcell/v2"
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/log"
	"github.com/lamassuiot/lamassu-virtual-device/pkg/observer"
	"github.com/rivo/tview"
)

var (
	app           *tview.Application
	view          *tview.Grid
	logger        *log.Logger
	mqttLogsCount int
	mqttLogsGrid  *tview.Grid
	deviceID      string
)

type concreteObserver struct {
}

func (s *concreteObserver) Update(t interface{}) {
	// do something
	level.Info(*logger).Log("AWS Observer updated", t)
	view.Clear()
	view = DrawAWS(t.(*observer.DeviceState))
}

func DrawAWS(obs *observer.DeviceState) *tview.Grid {
	level.Info(*logger).Log("msg", obs.PathCertificate)
	mqttLogsCount = 0
	mqttLogsGrid.Clear()

	if _, err := os.Stat(obs.PathCertificate); err != nil {
		level.Info(*logger).Log("msg", obs.PathCertificate, "err", err)
		return view
	}

	awsEndpoint := obs.Config.Aws.IotCoreEndpoint

	r, err := ioutil.ReadFile(obs.PathCertificate)
	if err != nil {
		level.Info(*logger).Log("err", err)
	}

	block, _ := pem.Decode(r)
	cert, err := x509.ParseCertificate(block.Bytes)
	deviceID = cert.Subject.CommonName

	view.AddItem(tview.NewTextView().SetText("Iot Core Endpoint").SetTextColor(tcell.ColorYellow), 0, 0, 1, 1, 0, 0, false)
	view.AddItem(tview.NewTextView().SetText(awsEndpoint), 0, 1, 1, 1, 0, 0, false)

	view.AddItem(tview.NewTextView().SetText("MQTT Publish Topic").SetTextColor(tcell.ColorYellow), 1, 0, 1, 1, 0, 0, false)
	view.AddItem(tview.NewTextView().SetText(deviceID+"/helloworld"), 1, 1, 1, 1, 0, 0, false)

	view.AddItem(tview.NewTextView().SetText("MQTT Subscription Topic").SetTextColor(tcell.ColorYellow), 2, 0, 1, 1, 0, 0, false)
	view.AddItem(tview.NewTextView().SetText(deviceID+"/helloworld"), 2, 1, 1, 1, 0, 0, false)

	return view
}

func addMqttLog(msg string, color tcell.Color) {
	mqttLogsGrid.AddItem(tview.NewTextView().SetText(msg).SetTextColor(color), mqttLogsCount, 0, 1, 1, 0, 0, false)
	mqttLogsCount = mqttLogsCount + 1
	app.ForceDraw()
}

func GetAWSItem(inlog log.Logger, deviceData *observer.DeviceState, inapp *tview.Application) tview.Primitive {
	logger = &inlog
	concreteObserver := &concreteObserver{}
	deviceData.Attach(concreteObserver)

	mqttLogsGrid = tview.NewGrid().
		SetColumns(0).
		SetBorders(false)

	view = tview.NewGrid().
		SetRows(2, 2).
		SetColumns(30, 0)

	app = inapp
	view = DrawAWS(deviceData)
	form := ConnectButton(deviceData)
	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(view, 7, 1, false).
		AddItem(form, 7, 1, false).
		AddItem(mqttLogsGrid, 7, 1, false)
	return flex
}

func ConnectButton(obs *observer.DeviceState) *tview.Form {
	awsEndpoint := obs.Config.Aws.IotCoreEndpoint
	form := tview.NewForm().
		AddButton("Connect to IoT Core", func() {
			level.Info(*logger).Log("msg", obs.PathCertificate)
			lambdaText := "Im " + deviceID + " device :)"
			tlsconfig, err := NewTLSConfig(*logger, obs.PathCertificate, obs.PathKey, obs.AwsEndpointCAFile)
			if err == nil {
				opts := MQTT.NewClientOptions()
				opts.AddBroker("tls://" + awsEndpoint + ":8883")
				opts.SetClientID(deviceID).SetTLSConfig(tlsconfig)
				opts.SetDefaultPublishHandler(f)

				// Start the connection
				c := MQTT.NewClient(opts)
				addMqttLog("Connecting to AWS IoT Core ...", tcell.ColorYellow)
				succesfulConnection := false
				if token := c.Connect(); token.Wait() && token.Error() != nil {
					level.Info(*logger).Log("Token error", token.Error())
					addMqttLog("Error while connecting to IoT Core. Retrying ...", tcell.ColorYellow)
					time.Sleep(5 * time.Second)
					if token := c.Connect(); token.Wait() && token.Error() != nil {
						level.Info(*logger).Log("Token error", token.Error())
						addMqttLog("Error while connecting to IoT Core. Desisting ...", tcell.ColorRed)
					} else {
						succesfulConnection = true
					}
				} else {
					succesfulConnection = true
				}

				if succesfulConnection {
					addMqttLog(fmt.Sprintf("IsConnected? %t", c.IsConnected()), tcell.ColorYellow)
					level.Info(*logger).Log("IsConnected", c.IsConnected())

					addMqttLog("Publishing to AWS with topic: "+deviceID+"/hello-ping", tcell.ColorYellow)

					msg := struct {
						Message  string `json:"msg"`
						DeviceID string `json:"device_id"`
					}{
						Message:  lambdaText,
						DeviceID: deviceID,
					}
					msgbytes, err := json.Marshal(msg)
					if err != nil {
						level.Info(*logger).Log("err", err)
					}
					level.Info(*logger).Log("msg", string(msgbytes))
					addMqttLog("Message : "+string(msgbytes), tcell.ColorGreen)
					level.Info(*logger).Log("msg", "publishing...")

					tokenResponse := c.Publish(deviceID+"/hello-ping", 0, false, string(msgbytes))
					level.Info(*logger).Log("msg", "waiting...")
					tokenResponse.Wait()
					level.Info(*logger).Log("er2", tokenResponse.Error())
				}

			} else {
				level.Info(*logger).Log("err", err)
			}

		})
	return form
}
func NewTLSConfig(logger log.Logger, certFile string, keyFile string, awsIotEndpointCAFile string) (*tls.Config, error) {
	certpool := x509.NewCertPool()
	pemCerts, err := ioutil.ReadFile(awsIotEndpointCAFile)
	if err != nil {
		level.Info(logger).Log("err", err)
		return nil, err
	}
	certpool.AppendCertsFromPEM(pemCerts)
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		level.Info(logger).Log("err", err)
		return nil, err
	}

	// Just to print out the client certificate..
	cert.Leaf, err = x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		level.Info(logger).Log("err", err)
		return nil, err
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
	}, nil
}

var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}
