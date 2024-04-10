# Dagger GitHub Container Registry (ghcr.io) demo

A demo of how to push an image built by dagger into a container registry, in this case github ghcr.io.

Not intended to be used as a dagger module because it's so simple, more of an example to copy from.

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

This is how you acquire "packages" write permission, then call the dagger-for-github module with DOCKER_REGISTRY, DOCKER_IMAGE_NAME, DOCKER_USERNAME and DOCKER_PASSWORD set appropriately.

Nothing in the dagger module is specific to GHCR, it's just:

* Build the container in the usual way with Dagger. Here we make a cowsay
  container that says "How now, brown cow", showing how to pass build context
  from the filesystem
* Use `ctr.WithRegistryAuth(registry, username, password).Publish(ctx, registry+"/"+imageName)`
  to publish the container to the ghcr registry.

Note: when specifying the args for the dagger module, you have to use `--registry=$DOCKER_REGISTRY` syntax for the non-secrets (regular strings), and `--password=env:DOCKER_PASSWORD` syntax for the secrets. This is because dagger only supports the `env:` syntax for secrets.