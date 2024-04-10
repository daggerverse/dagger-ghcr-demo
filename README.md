# Dagger GitHub Container Registry (ghcr.io) demo

A demo of how to push an image built by dagger into a container registry, in this case GitHub's `ghcr.io`.

This repo is not intended to be used as a dagger module because it's so simple, more of an example to copy from.

Example .github/workflows/docker-publish.yaml:

```
name: 'build-and-push'

on:
  push:
    branches:
    - main

jobs:
  dagger:
    runs-on: ubuntu-latest
    permissions:
      contents: read
      packages: write
    steps:
      - name: Checkout
        uses: actions/checkout@v4

      - name: Dagger Build & Push
        uses: dagger/dagger-for-github@v5
        with:
          version: "0.11.0"
          verb: call
          args: build-and-push --registry=$DOCKER_REGISTRY --image-name=$DOCKER_IMAGE_NAME --username=$DOCKER_USERNAME --password=env:DOCKER_PASSWORD --build-context .
        env:
          DOCKER_REGISTRY: ghcr.io
          DOCKER_IMAGE_NAME: ${{ github.repository }}
          DOCKER_USERNAME: ${{ github.actor }}
          DOCKER_PASSWORD: ${{ secrets.GITHUB_TOKEN }}
```

This is how you acquire `packages: write` permission, then call the `dagger-for-github` module with `DOCKER_REGISTRY`, `DOCKER_IMAGE_NAME`, `DOCKER_USERNAME` and `DOCKER_PASSWORD` set appropriately.

Copy and edit appropriately :-)

Note: when specifying the args for the dagger module, you have to use `--registry=$DOCKER_REGISTRY` syntax for the non-secrets (regular strings), and `--password=env:DOCKER_PASSWORD` syntax for the secrets. This is because dagger only supports the `env:` syntax for secrets.

Nothing in the dagger module is specific to GHCR, it's just:

* Build the container in the usual way with Dagger. Here we make a cowsay
  container that says "How now, brown cow", showing how to pass build context
  from the filesystem
* Use `ctr.WithRegistryAuth(registry, username, password).Publish(ctx, registry+"/"+imageName)`
  to publish the container to the ghcr registry.

Here's the complete [`dagger/main.go`](dagger/main.go):
```go
package main

import (
	"context"
)

type Cowsay struct{}

// Application specific build logic
func (m *Cowsay) Build(ctx context.Context, buildContext *Directory) *Container {
	return dag.Container().
		From("ubuntu:latest").
		WithFile("/cow.txt", buildContext.File("cow.txt")).
		WithExec([]string{"apt", "update"}).
		WithExec([]string{"apt", "install", "-y", "cowsay"}).
		WithEntrypoint([]string{"/usr/bin/cowsay", "/cow.txt"})
}

// Take the built container and push it
func (m *Cowsay) BuildAndPush(ctx context.Context, registry, imageName, username string, password *Secret, buildContext *Directory) error {
	_, err := m.Build(ctx, buildContext).
		WithRegistryAuth(registry, username, password).
		Publish(ctx, registry+"/"+imageName)
	return err
}
```

And you can see the resulting published container images: [https://github.com/lukemarsden/dagger-ghcr-demo/pkgs/container/dagger-ghcr-demo](https://github.com/lukemarsden/dagger-ghcr-demo/pkgs/container/dagger-ghcr-demo)