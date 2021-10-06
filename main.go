package main

import (
	"context"
	"fmt"
	"os"

	"emarcey/data-vault/dependencies"
)

func main() {
	fmt.Printf("Hello\n")
	ctx := context.Background()
	opts := dependencies.DependenciesInitOpts{LoggerType: "text", Env: "local"}
	deps, err := dependencies.MakeDependencies(ctx, opts)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	deps.Logger.Info("Heyo")
	deps.Logger.Debug("Heyo")
	tracer := deps.Tracer(ctx, "dummy")
	defer tracer.Close()
	tracer.AddBreadcrumb(map[string]interface{}{"hi": "there"})

}
