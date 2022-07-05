package internal

import (
	"fmt"
	"github.com/bhbosman/gocommon/Services/interfaces"
	"go.uber.org/fx"
)

func ProvideUniqueReferenceService(UniqueSessionNumber interfaces.IUniqueReferenceService) fx.Option {
	return fx.Provide(
		fx.Annotated{
			Target: func() (interfaces.IUniqueReferenceService, error) {
				if UniqueSessionNumber == nil {
					return nil, fmt.Errorf("interfaces.IUniqueReferenceService is nil. Please resolve")
				}
				return UniqueSessionNumber, nil
			},
		},
	)
}
