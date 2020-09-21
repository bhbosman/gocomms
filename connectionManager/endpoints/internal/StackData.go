package internal

import (
	"github.com/bhbosman/gocomms/connectionManager"
)

type StackData struct {
	Name     string
	Index    int
	InValue  connectionManager.StackPropertyValue
	OutValue connectionManager.StackPropertyValue
}

type StackDataSort struct {
	Data []*StackData
}

func (self *StackDataSort) Len() int {
	return len(self.Data)
}

func (self *StackDataSort) Less(i, j int) bool {
	return self.Data[i].Index-self.Data[j].Index < 0

}

func (self *StackDataSort) Swap(i, j int) {
	self.Data[i], self.Data[j] = self.Data[j], self.Data[i]
}
