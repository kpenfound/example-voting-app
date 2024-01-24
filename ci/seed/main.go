package main

type Seed struct {}


func (s *Seed) Run(dir *Directory, vote *Service) *Container {
	return dag.Container().From("python:3.9-slim").
		WithExec([]string{"apt-get", "update"}).
		WithExec([]string{"apt-get", "install", "-y", "--no-install-recommends", "apache2-utils"}).
		WithExec([]string{"rm", "-rf", "/var/lib/apt/lists/*"}).
		WithWorkdir("/seed").
		WithMountedDirectory("/seed", dir).
		WithServiceBinding("vote", vote).
		WithExec([]string{"python", "make-data.py"}).
		WithExec([]string{"/seed/generate-votes.sh"})
}
