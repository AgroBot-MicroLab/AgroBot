package mqttclient

import (
    "fmt"
    "log"
    "os"
    "strings"
    "time"

    paho "github.com/eclipse/paho.mqtt.golang"
)

type MqttClient struct {
    c paho.Client
}

func New() *MqttClient {
    broker := "tcp://broker.emqx.io:1883"
    cid := fmt.Sprintf("agro-bot-%d", time.Now().UnixNano())

    opts := paho.NewClientOptions().
        AddBroker(broker).
        SetClientID(cid).
        SetCleanSession(true).
        SetAutoReconnect(true).
        SetConnectTimeout(5 * time.Second).
        SetKeepAlive(30 * time.Second)

    opts.OnConnect = func(c paho.Client) {
        log.Printf("[MQTT] connected to %s as %s", broker, cid)
        if topics := os.Getenv("MQTT_SUB_TOPICS"); topics != "" {
            for t := range strings.SplitSeq(topics, ",") {
                t = strings.TrimSpace(t)
                if t == "" {
                    continue
                }
                tok := c.Subscribe(t, 1, func(_ paho.Client, m paho.Message) {
                    log.Printf("[MQTT] %s: %s", m.Topic(), string(m.Payload()))
                })
                tok.Wait()
                if err := tok.Error(); err != nil {
                    log.Printf("subscribe %q error: %v", t, err)
                }
            }
        }
    }

    opts.OnConnectionLost = func(_ paho.Client, err error) {
        log.Println("[MQTT] connection lost:", err)
    }

    c := paho.NewClient(opts)
    if t := c.Connect(); t.Wait() && t.Error() != nil {
        panic(fmt.Errorf("mqtt connect error: %w", t.Error()))
    }
    return &MqttClient{c: c}
}

func (mc *MqttClient) Publish(topic string, payload []byte) error {
    if mc == nil || mc.c == nil || !mc.c.IsConnectionOpen() {
        return fmt.Errorf("mqtt not connected")
    }
    t := mc.c.Publish(topic, 1, false, payload)
    if !t.WaitTimeout(3 * time.Second) {
        return fmt.Errorf("publish timeout")
    }
    return t.Error()
}

func (mc *MqttClient) Close() {
    if mc != nil && mc.c != nil {
        mc.c.Disconnect(250)
    }
}

