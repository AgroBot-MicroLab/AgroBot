package shared

type MissionEventType string

const (
    EventWaypointPassed  MissionEventType = "waypoint_passed"
    EventPhotoReceived MissionEventType = "photo_received"
)

type MissionEvent struct {
    Type  MissionEventType `json:"type"`
    Data  any              `json:"data"`
}

type WaypointPassed struct {
	WaypointId uint16 `json:"waypoint_id"`
	IsLast     bool   `json:"is_last"`
}

type PhotoReceived struct {
    MissionID   int64     `json:"mission_id"`
    ImageID     int64     `json:"image_id"`
    Path  		string 	  `json:"path"`
}

type Pos struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
	Yaw float64 `json:"yaw"`
}
