package main

import (
	"context"
)

type Ci struct {}

// Serve the whole stack locally
func (c *Ci) Serve(ctx context.Context, directory *Directory) (*Service, error) {
	red := redis()
	db := postgres()

	// Vote needs redis
	vote := dag.Vote().Serve(directory.Directory("vote"), red)
	// Result needs postgres
	result := dag.Result().Serve(directory.Directory("result"), db)
	// Worker needs redis and postgres
	worker := dag.Worker().Serve(directory.Directory("worker"), red, db)

	proxy := dag.Proxy().
		WithService(vote, "vote", 5000, 80).
		WithService(result, "result", 5001, 80).
		WithService(worker, "worker", 9999, 9999).
		Service()

	return proxy, nil
}

// Build the whole stack
func (c *Ci) Build(ctx context.Context, directory *Directory) error {
	vote := dag.Vote().Build(directory.Directory("vote"))
	_, err := vote.Sync(ctx)
	if err != nil {
		return err
	}

	result := dag.Result().Build(directory.Directory("result"))
	_, err = result.Sync(ctx)
	if err != nil {
		return err
	}

	worker := dag.Worker().Build(directory.Directory("worker"))
	_, err = worker.Sync(ctx)
	if err != nil {
		return err
	}

	return nil
}

// A redis container
func redis() *Service {
	return dag.Container().From("redis:alpine").
		WithExposedPort(6379).
		AsService()
}

// A postgres container
func postgres() *Service {
	return dag.Container().From("postgres:15-alpine").
		WithEnvVariable("POSTGRES_USER", "postgres").
		WithEnvVariable("POSTGRES_PASSWORD", "postgres").
		WithExposedPort(5432).
		AsService()
}
