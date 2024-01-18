import { dag, Container, Directory, object, func } from "@dagger.io/dagger"

@object
// eslint-disable-next-line @typescript-eslint/no-unused-vars
class Result {

  @func
  build(dir: Directory): Container {
    return dag
      .container()
      .from("node:18-slim")
      .withExec(["apt-get", "update"])
      .withExec(["apt-get", "install", "-y", "--no-install-recommends", "curl", "tini"])
      .withExec(["rm", "-rf", "/var/lib/apt/lists/*"])
      .withWorkdir("/usr/local/app")
      .withFile("/usr/local/app/package.json", dir.file("package.json"))
      .withFile("/usr/local/app/package-lock.json", dir.file("package-lock.json"))
      .withExec(["npm", "ci"])
      .withExec(["npm", "cache", "clean", "--force"])
      .withExec(["mv", "/usr/local/app/node_modules", "/node_modules"])
      .withDirectory("/usr/local/app", dir)
      .withExposedPort(80)
      .withEnvVariable("PORT", "80")
      .withEntrypoint(["/usr/bin/tini", "--"])
      .withExec(["node", "server.js"])
  }
}
