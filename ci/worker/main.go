package main

import (
	"runtime"
)

type Worker struct {}

func (w *Worker) Build(directory *Directory) *Container {
	builder := dotnetPublish(directory)

	return dag.Container().
		From("mcr.microsoft.com/dotnet/runtime:7.0").
		WithWorkdir("/app").
		WithDirectory("/app", builder.Directory("/app"))
}

func (w *Worker) Serve(directory *Directory, redis *Service, db *Service) *Service {
	builder := dotnetPublish(directory)

	return dag.Container().
		From("mcr.microsoft.com/dotnet/runtime:7.0").
		WithWorkdir("/app").
		WithDirectory("/app", builder.Directory("/app")).
		WithServiceBinding("redis", redis).
		WithServiceBinding("db", db).
		WithExec([]string{"dotnet", "Worker.dll"}).
		AsService()
}

func dotnetPublish(directory *Directory) *Container {
	csproj := directory.File("Worker.csproj")
	return dag.Container().
		From("mcr.microsoft.com/dotnet/sdk:7.0").
		WithWorkdir("/source").
		WithFile("/source/Worker.csproj", csproj).
		WithExec([]string{"dotnet", "restore", "-a", runtime.GOARCH}).
		WithDirectory("/source", directory).
		WithExec([]string{"dotnet", "publish", "-c", "release", "-o", "/app", "-a", runtime.GOARCH, "--self-contained", "false", "--no-restore"})

}
