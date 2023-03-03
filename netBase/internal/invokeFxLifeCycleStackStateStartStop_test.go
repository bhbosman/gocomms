package internal

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx/fxtest"
	"testing"
)

func TestInvokeFxLifeCycleStackStateStartStop(t *testing.T) {
	testApp := fxtest.New(t)
	assert.NoError(t, testApp.Err())
}
