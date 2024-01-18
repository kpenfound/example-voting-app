package main

import (
	"context"
)

type Seed struct {}


func (s *Seed) Run(ctx context.Context, dir *Directory, psql *Service) error {
	_, err := dag.Container().From("python:3.9-slim").
		WithExec([]string{"apt-get", "update"}).
		WithExec([]string{"apt-get", "install", "-y", "--no-install-recommends", "apache2-utils"}).
		WithExec([]string{"rm", "-rf", "/var/lib/apt/lists/*"}).
		WithWorkdir("/seed").
		WithMountedDirectory("/seed", dir.Directory("seed-data")).
		WithExec([]string{"python", "make-data.py"}).
		WithExec([]string{"/seed/generate-votes.sh"}).
		WithServiceBinding("db", psql).
		Sync(ctx)

	return err
}
