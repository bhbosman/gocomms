package internal

import (
	"context"
	"github.com/bhbosman/gocomms/internal"
	"github.com/bhbosman/gocomms/intf"
	"github.com/google/uuid"
	"go.uber.org/fx"
	"net"
	"net/url"
)

func CreateStackData(params struct {
	fx.In
	LifeCycle  fx.Lifecycle
	StackState []*internal.StackState
	Conn       net.Conn
	Url        *url.URL
	Ctx        context.Context
	CancelFunc internal.CancelFunc
	Cfr        intf.IConnectionReactorFactory
}) (map[uuid.UUID]interface{}, error) {
	result := make(map[uuid.UUID]interface{})
	var err error
	for _, iteration := range params.StackState {
		localItem := iteration
		var stackData interface{}
		stackData, err = localItem.Create(params.Conn, params.Url, params.Ctx, params.CancelFunc, params.Cfr)
		if err != nil {
			return nil, err
		}
		if stackData == nil {
			stackData = internal.NewNoCloser()
		}
		result[localItem.Id] = stackData
		params.LifeCycle.Append(
			fx.Hook{
				OnStop: func(ctx context.Context) error {
					if localItem.Destroy != nil {
						localData, _ := result[localItem.Id]
						return localItem.Destroy(localData)
					}
					return nil
				},
			})
	}
	return result, nil
}
