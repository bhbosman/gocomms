package common

import (
	"github.com/reactivex/rxgo/v2"
	"golang.org/x/net/context"
	"testing"
)

func TestName(t *testing.T) {
	rxgo.Defer(
		[]rxgo.Producer{
			func(ctx context.Context, next chan<- rxgo.Item) {

			},
		},
	).
		Map(
			func(ctx context.Context, i interface{}) (interface{}, error) {
				return i, nil
			},
		).
		Map(
			func(ctx context.Context, i interface{}) (interface{}, error) {
				return i, nil
			},
		).
		Map(
			func(ctx context.Context, i interface{}) (interface{}, error) {
				return i, nil
			},
		).
		Map(
			func(ctx context.Context, i interface{}) (interface{}, error) {
				return i, nil
			},
		)
}
