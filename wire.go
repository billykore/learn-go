//go:build wireinject
// +build wireinject

package learning

import "github.com/google/wire"

func initApp() *app {
	wire.Build(
		newApp,
		newMemRepo, wire.Bind(new(repo), new(*memRepo)),
		newJsonRepo, wire.Bind(new(repo), new(*jsonRepo)),
	)
	return &app{}
}
