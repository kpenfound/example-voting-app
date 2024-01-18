package main

import (
	"context"
)

type Ci struct {}

func (c *Ci) Serve(ctx context.Context, dir *Directory) (*Service, error) {
	red := redis().AsService()
	psql := postgres().AsService()

	// Vote needs redis
	vote := dag.Vote().Build(dir).WithServiceBinding("redis", red)
	// Result needs postgres
	result := dag.Result().Build(dir).WithServiceBinding("db", psql)
	// Worker needs redis and postgres
	worker := dag.Worker().Build(dir).
		WithServiceBinding("redis", red).
		WithServiceBinding("db", psql)

	// Seed initial data
	_, err := dag.Seed().Run(ctx, dir, psql)
	if err != nil {
		return nil, err
	}

	// Proxy services
	voteSvc := vote.AsService()
	resultSvc := result.AsService()
	workerSvc := worker.AsService()
	proxy := dag.Proxy().
		WithService(voteSvc, "vote", 5000, 80).
		WithService(resultSvc, "result", 5001, 80).
		WithService(workerSvc, "worker", 9999, 9999).
		Service()

	return proxy, nil
}

// A redis container
func redis() *Container {
	return dag.Container().From("redis:alpine").WithExposedPort(6379)
}

// A postgres container
func postgres() *Container {
	return dag.Container().From("postgres:15-alpine").
		WithEnvVariable("POSTGRES_USER", "postgres").
		WithEnvVariable("POSTGRES_PASSWORD", "postgres").
		WithExposedPort(5432)
}
