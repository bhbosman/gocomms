package websocket

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"testing"
)

func TestName(t *testing.T) {
	t.Run("fff", func(t *testing.T) {
		h := hmac.New(sha512.New, []byte("4961b74efac86b25cce8fbe4c9811c4c7a787b7a5996660afcc2e287ad864363"))
		h.Write([]byte("1558014486185"))
		h.Write([]byte("GET"))
		h.Write([]byte("/v1/account/balances"))
		h.Write([]byte(""))

		println(hex.EncodeToString(h.Sum(nil)))

	})

	t.Run("fff", func(t *testing.T) {
		h := hmac.New(sha512.New, []byte("4961b74efac86b25cce8fbe4c9811c4c7a787b7a5996660afcc2e287ad864363"))
		h.Write([]byte("1558017528946"))
		h.Write([]byte("POST"))
		h.Write([]byte("/v1/orders/market"))
		h.Write([]byte(`{"customerOrderId":"ORDER-000001","pair":"BTCZAR","side":"BUY","quoteAmount":"80000"}`))

		println(hex.EncodeToString(h.Sum(nil)))

	})

}
