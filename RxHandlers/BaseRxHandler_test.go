package RxHandlers

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

func TestNewBaseRxHandler(t *testing.T) {
	t.Run("No Nils", func(t *testing.T) {
		_, err := NewBaseRxHandler(zap.NewNop(), func(context string, inbound bool, err error) {

		})
		assert.NoError(t, err)
	})
	t.Run("No logger", func(t *testing.T) {
		_, err := NewBaseRxHandler(nil, func(context string, inbound bool, err error) {

		})
		assert.Error(t, err)
	})
	t.Run("No callback", func(t *testing.T) {
		_, err := NewBaseRxHandler(zap.NewNop(), nil)
		assert.Error(t, err)
	})
	t.Run("No callback and logger", func(t *testing.T) {
		_, err := NewBaseRxHandler(nil, nil)
		assert.Error(t, err)
	})

}
