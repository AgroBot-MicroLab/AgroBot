package mqttclient

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Client struct {
	c mqtt.Client
}

var defaultClient *Client

// MustInitFromEnv инициализирует глобальный MQTT-клиент из ENV и паникует при ошибке.
func subscribeWithRetry(c mqtt.Client, topic string) {
	backoff := 500 * time.Millisecond
	for i := 1; i <= 5; i++ { // до 5 попыток
		token := c.Subscribe(topic, 1, func(_ mqtt.Client, m mqtt.Message) {
			log.Printf("[MQTT] msg topic=%s payload=%s\n", m.Topic(), string(m.Payload()))
		})

		// ждём завершения, но с таймаутом
		if token.WaitTimeout(5*time.Second) && token.Error() == nil {
			log.Println("[MQTT] subscribed:", topic)
			return
		}

		if err := token.Error(); err != nil {
			log.Printf("[MQTT] subscribe error for %s (try %d/5): %v", topic, i, err)
		} else {
			log.Printf("[MQTT] subscribe timeout for %s (try %d/5)", topic, i)
		}

		time.Sleep(backoff)
		backoff *= 2
	}

	log.Printf("[MQTT] subscribe failed for %s after retries", topic)
}

func MustInitFromEnv() *Client {
	broker := getenvDefault("MQTT_BROKER", "tcp://broker.emqx.io:1883") // пример
	clientID := getenvDefault("MQTT_CLIENT_ID", "agro-api")
	user := os.Getenv("MQTT_USERNAME") // можно пусто
	pass := os.Getenv("MQTT_PASSWORD") // можно пусто
	useTLS := strings.HasPrefix(broker, "ssl://") || strings.HasPrefix(broker, "tls://")

	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(clientID)
	if user != "" {
		opts.SetUsername(user)
		opts.SetPassword(pass)
	}
	// Реконнект и устойчивость
	opts.SetAutoReconnect(true)
	opts.SetConnectRetry(true)
	opts.SetConnectRetryInterval(2 * time.Second)
	opts.SetOrderMatters(false) // повысить пропускную способность обработчиков
	opts.SetResumeSubs(true)
	// TLS при необходимости
	if useTLS {
		opts.SetTLSConfig(&tls.Config{InsecureSkipVerify: false})
	}

	opts.OnConnect = func(c mqtt.Client) {
		log.Println("[MQTT] connected to", broker)

		if topics := os.Getenv("MQTT_SUB_TOPICS"); topics != "" {
			for _, t := range strings.Split(topics, ",") {
				t = strings.TrimSpace(t)
				if t == "" {
					continue
				}
				// НЕ блокируем OnConnect — пусть подписки идут параллельно с ретраями
				go subscribeWithRetry(c, t)
			}
		}
	}
	opts.OnConnectionLost = func(_ mqtt.Client, err error) {
		log.Println("[MQTT] connection lost:", err)
	}

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(fmt.Errorf("mqtt connect error: %w", token.Error()))
	}

	defaultClient = &Client{c: c}
	return defaultClient
}

// Publish публикует сообщение (QoS 1, без ретейна).
func Publish(topic string, payload []byte) error {
	if defaultClient == nil {
		return fmt.Errorf("mqtt not initialized")
	}
	t := defaultClient.c.Publish(topic, 1, false, payload)
	t.Wait()
	return t.Error()
}

func Close() {
	if defaultClient != nil && defaultClient.c.IsConnected() {
		defaultClient.c.Disconnect(250)
	}
}

func getenvDefault(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
