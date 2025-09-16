package shared

type MissionStatus struct {
	WaypointId uint16 `json:"waypoint_id"`
	IsLast     bool   `json:"is_last"`
}

type Pos struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}
