// A demo of how to push to a container registry, in this case github ghcr.io
//
// Not intended to be used as a dagger module because it's so simple, more of an
// example to copy from.
//
// For GitHub Actions yaml and dagger go code to copy, see the README at
// https://github.com/lukemarsden/dagger-ghcr-demo

package main

import (
	"context"
)

type DaggerGhcrDemo struct{}

// Application specific build logic
func (m *DaggerGhcrDemo) Build(ctx context.Context, buildContext *Directory) *Container {
	return dag.Container().
		From("ubuntu:latest").
		WithFile("/cow.txt", buildContext.File("cow.txt")).
		WithExec([]string{"apt", "update"}).
		WithExec([]string{"apt", "install", "-y", "cowsay"}).
		WithEntrypoint([]string{"bash", "-c", "/usr/games/cowsay < /cow.txt"})
}

// Take the built container and push it
func (m *DaggerGhcrDemo) BuildAndPush(ctx context.Context, registry, imageName, username string, password *Secret, buildContext *Directory) error {
	_, err := m.Build(ctx, buildContext).
		WithRegistryAuth(registry, username, password).
		Publish(ctx, registry+"/"+imageName)
	return err
}
