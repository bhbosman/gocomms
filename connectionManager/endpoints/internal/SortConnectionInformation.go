package internal

import (
	"github.com/bhbosman/gocomms/connectionManager"
	"strings"
)

type SortConnectionInformation struct {
	Connections []*connectionManager.ConnectionInformation
}

func (self *SortConnectionInformation) Len() int {
	return len(self.Connections)
}

func (self *SortConnectionInformation) Less(i, j int) bool {
	n := strings.Compare(self.Connections[i].Name, self.Connections[j].Name)
	if n != 0 {
		return n < 0
	}
	n = strings.Compare(self.Connections[i].Id, self.Connections[j].Id)
	return n < 0

}

func (self *SortConnectionInformation) Swap(i, j int) {
	self.Connections[i], self.Connections[j] = self.Connections[j], self.Connections[i]
}
