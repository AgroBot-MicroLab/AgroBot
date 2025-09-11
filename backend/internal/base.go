package internal

import (
	"database/sql"
	"agro-bot/internal/mav"
	"agro-bot/internal/mqttclient"
)

type App struct {
	DB            *sql.DB
    MavLinkClient *mav.Client
    MqttClient    *mqttclient.MqttClient
}

