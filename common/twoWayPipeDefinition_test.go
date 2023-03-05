package common_test

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx/fxtest"
	"testing"
)

func TestTwoWayPipeDefinition(t *testing.T) {
	ctrl := gomock.NewController(t)

	app := fxtest.New(t)
	assert.NoError(t, app.Err())
	app.RequireStart()

	app.RequireStop()
	defer ctrl.Finish()
}
