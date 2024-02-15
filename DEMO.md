# The Demo Flow

## Check out the app

There's an architecture diagram in README.md and some cool stuff to dig into

## Run a module

Ruff is a Python linter. It's installed in the `dagger.json` already

The `vote/` application is Python. Lint it.

`dagger -m ruff functions`

`dagger -m ruff call check --directory vote`

## Build the applications

`dagger functions`

`dagger call build --directory .`

## Run the whole stack

`dagger call serve --directory . up --ports 5000:5000,5001:5001`

Seed some intial data (in another terminal)

`dagger -m seed call run --dir seed-data/ --vote tcp://localhost:5000`

