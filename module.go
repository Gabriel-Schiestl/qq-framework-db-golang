package db

import (
	"go.uber.org/fx"
)

func DBModule(serviceName string, entities ...any) fx.Option {
    return fx.Module(
        "db",
        fx.Provide(func(config IDbProvider) (*GormDatabase, error) {
            return NewDatabase(DatabaseOptions{
                ServiceName: serviceName,
                Config:      config,
            }, entities...)
        }),
    )
}
