package internal

import (
	"fmt"
	"go.uber.org/fx"
)

func ProvideStackName(stackName string) fx.Option {
	return fx.Provide(
		fx.Annotated{
			Name: "StackName",
			Target: func() (string, error) {
				if stackName == "" {
					return "", fmt.Errorf("stackname not defined. Please resolve")
				}
				return stackName, nil
			},
		},
	)
}
