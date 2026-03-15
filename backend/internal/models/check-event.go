package models

type CheckEvent struct {
	MonitorID    string
	Url          string
	Status_ok    bool
	ResponseTime float64
}
