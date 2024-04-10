// A demo of how to
//
// This module has been generated via dagger init and serves as a reference to
// basic module structure as you get started with Dagger.
//
// Two functions have been pre-created. You can modify, delete, or add to them,
// as needed. They demonstrate usage of arguments and return types using simple
// echo and grep commands. The functions can be called from the dagger CLI or
// from one of the SDKs.
//
// The first line in this comment block is a short description line and the
// rest is a long description with more detail on the module's purpose or usage,
// if appropriate. All modules should have a short description.

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
		WithExec([]string{"apt", "install", "-y", "cowsay"}).
		WithEntrypoint([]string{"/usr/bin/cowsay", "/cow.txt"})
}

// Take the built container and push it
func (m *Cowsay) BuildAndPush(ctx context.Context, registry, imageName, username string, password *Secret, buildContext *Directory) error {
	_, err := m.Build(ctx, buildContext).
		WithRegistryAuth(registry, username, password).
		Publish(ctx, imageName)
	return err
}
