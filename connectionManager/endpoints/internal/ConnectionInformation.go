package internal

import "time"

type ConnectionInformation struct {
	Id              string
	Name            string
	Status          string
	ConnectionTime  time.Time
	StackProperties []*StackData
}
