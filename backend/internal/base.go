package internal

import (
	"agro-bot/internal/mav"
	"agro-bot/internal/mqttclient"
	"database/sql"
)

type App struct {
	DB            *sql.DB
	MavLinkClient *mav.Client
	MqttClient    *mqttclient.MqttClient
}
