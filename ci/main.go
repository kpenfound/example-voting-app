package main

import (
	"context"
)

type Ci struct {}

func (c *Ci) Serve(ctx context.Context, dir *Directory) (*Service, error) {
	red := redis()
	db := postgres()

	// Vote needs redis
	vote := dag.Vote().Serve(dir.Directory("vote"), red)
	// Result needs postgres
	result := dag.Result().Serve(dir.Directory("result"), db)
	// Worker needs redis and postgres
	worker := dag.Worker().Serve(dir.Directory("worker"), red, db)

	// Seed initial data
// This causes a failure: https://github.com/dagger/dagger/issues/6493
	_, err := dag.Seed().Run(dir.Directory("seed-data"), vote).Sync(ctx)
	if err != nil {
		return nil, err
	}

	proxy := dag.Proxy().
		WithService(vote, "vote", 5000, 80).
		WithService(result, "result", 5001, 80).
		WithService(worker, "worker", 9999, 9999).
		Service()

	return proxy, nil
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