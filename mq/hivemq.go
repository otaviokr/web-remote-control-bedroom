package mq

import (
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// MqttPublisher can send messages to a MQ.
type MqttPublisher struct {
	clientId string
	Client mqtt.Client
}

// NewPublisher instatiate a new MQ Publisher.
func NewPublisher(broker string, port int) (*MqttPublisher, error) {
	clientId := fmt.Sprintf("okr_bdrm_%d", time.Now().Unix())

	// Connecting to MQTT.
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID(clientId) // ClientID can be anything, just make sure it is unique
	// opts.SetUsername("username") // Not needed because we are using the public HiveMQ
	// opts.SetPassword("password") // Not needed because we are using the public HiveMQ
	// opts.SetDefaultPublishHandler(getMessageReceivedHandler(device))
	// opts.OnConnect = getOnConnectHandler()
	// opts.OnConnectionLost = getConnectionLostHandler()

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return &MqttPublisher{}, token.Error()
	}

	return &MqttPublisher {
		clientId: clientId,
		Client:   client,
	}, nil
}

// Publish will send the message *text* to topic *topic*.
func (m *MqttPublisher) Publish(topic, text string) (error) {
	if !m.Client.IsConnected() {
		if token := m.Client.Connect(); token.Wait() && token.Error() != nil {
			return token.Error()
		}
	}
	if token := m.Client.Publish(topic, 0, false, text); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}
