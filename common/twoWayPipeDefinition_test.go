package common_test

import (
	"github.com/golang/mock/gomock"
	"testing"
)





func CreateOneWayPipe[source any, destination any]() {

}

func TestTwoWayPipeDefinition(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
}
