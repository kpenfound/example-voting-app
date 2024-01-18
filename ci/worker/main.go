package main

import (
	"runtime"
)

type Worker struct {}


func (w *Worker) Build(dir *Directory) *Container {
	csproj := dag.Directory().WithDirectory("/", dir, DirectoryWithDirectoryOpts{ Include: []string{ "*.csproj" }})

	builder := dag.Container(ContainerOpts{ Platform: "linux/amd64" }).
		From("mcr.microsoft.com/dotnet/sdk:7.0").
		WithWorkdir("/source").
		WithMountedDirectory("/source", csproj).
		WithExec([]string{"dotnet", "restore", "-a", runtime.GOARCH}).
		WithMountedDirectory("/source", dir).
		WithExec([]string{"dotnet", "publish", "-c", "release", "-o", "/app", "-a", runtime.GOARCH, "--self-contained", "false", "--no-restore"})

	return dag.Container().
		From("mcr.microsoft.com/dotnet/runtime:7.0").
		WithWorkdir("/app").
		WithDirectory("/app", builder.Directory("/app")).
		WithEntrypoint([]string{"dotnet", "Worker.dll"})
}
