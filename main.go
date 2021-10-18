package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"emarcey/data-vault/dependencies"
	"emarcey/data-vault/dependencies/secrets"
)

func main() {
	fmt.Printf("Hello\n")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	opts, err := dependencies.ReadOpts("./server_conf.yml")
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

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

	dummySecret := secrets.NewSecret("tableName", "rowId", "columnName", "idHash", "key", "iv")
	fmt.Printf("Secret: %v\n", dummySecret)
	oSecret, err := deps.SecretsManager.GetOrPutSecret(ctx, dummySecret)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	dummySecret2 := secrets.NewSecret("tableName", "rowId", "columnName", "idHash", "key", "ivvvvv")
	oSecret2, err := deps.SecretsManager.GetOrPutSecret(ctx, dummySecret2)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	fmt.Printf("Secret: %v\n", oSecret)
	fmt.Printf("Secret2: %v\n", oSecret2)

	time.Sleep(8 * time.Second)
}
