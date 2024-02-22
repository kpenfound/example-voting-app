import dagger
from dagger import dag, function

@function
def build(directory: dagger.Directory) -> dagger.Container:
    return base(directory)

@function
def serve(directory: dagger.Directory, redis: dagger.Service) -> dagger.Service:
    return (
        base(directory)
        .with_exposed_port(80)
        .with_service_binding("redis", redis)
        .with_exec(["python", "app.py"])
        .as_service()
    )

def base(directory: dagger.Directory) -> dagger.Container:
    return (
        dag.container().from_("python:3.11-slim")
        .with_exec(["apt-get", "update"])
        .with_exec(["apt-get", "install", "-y", "--no-install-recommends", "curl"])
        .with_exec(["rm", "-rf", "/var/lib/apt/lists/*"])
        .with_workdir("/usr/local/app")
        .with_file("/usr/local/app/requirements.txt", directory.file("requirements.txt"))
        .with_exec(["pip", "install", "--no-cache-dir", "-r", "requirements.txt"])
        .with_directory("/usr/local/app", directory)
    )
