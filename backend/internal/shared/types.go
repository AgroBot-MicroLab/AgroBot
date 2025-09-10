package shared

type MissionStatus struct {
    Status bool `json:"status"`
}

type Pos struct {
    Lat float64 `json:"lat"`
    Lon float64 `json:"lon"`
}
