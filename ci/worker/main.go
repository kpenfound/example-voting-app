package main

import (
	"runtime"
)

type Worker struct {}


func (w *Worker) Serve(dir *Directory, redis *Service, db *Service) *Service {
	csproj := dir.File("Worker.csproj")

	builder := dag.Container().
		From("mcr.microsoft.com/dotnet/sdk:7.0").
		WithWorkdir("/source").
		WithFile("/source/Worker.csproj", csproj).
		WithExec([]string{"dotnet", "restore", "-a", runtime.GOARCH}).
		WithDirectory("/source", dir).
		WithExec([]string{"dotnet", "publish", "-c", "release", "-o", "/app", "-a", runtime.GOARCH, "--self-contained", "false", "--no-restore"})

	return dag.Container().
		From("mcr.microsoft.com/dotnet/runtime:7.0").
		WithWorkdir("/app").
		WithDirectory("/app", builder.Directory("/app")).
		WithServiceBinding("redis", redis).
		WithServiceBinding("db", db).
		WithExec([]string{"dotnet", "Worker.dll"}).
		AsService()
}
