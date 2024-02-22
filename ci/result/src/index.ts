import { dag, Container, Directory, object, func } from "@dagger.io/dagger"

@object()
// eslint-disable-next-line @typescript-eslint/no-unused-vars
class Result {

  @func()
  build(directory: Directory): Container {
    return dag
      .container()
      .from("node:18-slim")
      .withExec(["apt-get", "update"])
      .withExec(["apt-get", "install", "-y", "--no-install-recommends", "curl", "tini"])
      .withExec(["rm", "-rf", "/var/lib/apt/lists/*"])
      .withWorkdir("/usr/local/app")
      .withFile("/usr/local/app/package.json", directory.file("package.json"))
      .withFile("/usr/local/app/package-lock.json", directory.file("package-lock.json"))
      .withExec(["npm", "ci"])
      .withExec(["npm", "cache", "clean", "--force"])
      .withExec(["mv", "/usr/local/app/node_modules", "/node_modules"])
      .withDirectory("/usr/local/app", directory)
  }

  @func()
  serve(directory: Directory, db: Service): Service {
    return this.build(directory)
      .withExposedPort(80)
      .withEnvVariable("PORT", "80")
      .withServiceBinding("db", db)
      .withEntrypoint(["/usr/bin/tini", "--"])
      .withExec(["node", "server.js"])
      .asService()
  }
}

