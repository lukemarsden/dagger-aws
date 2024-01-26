package main

import (
	"context"
)

type Aws struct{}

// example usage: "dagger call container-echo --string-arg yo"
func (m *Aws) ContainerEcho(stringArg string) *Container {
	return dag.Container().From("alpine:latest").WithExec([]string{"echo", stringArg})
}

// example usage: "dagger call grep-dir --directory-arg . --pattern GrepDir"
func (m *Aws) GrepDir(ctx context.Context, directoryArg *Directory, pattern string) (string, error) {
	return dag.Container().
		From("alpine:latest").
		WithMountedDirectory("/mnt", directoryArg).
		WithWorkdir("/mnt").
		WithExec([]string{"grep", "-R", pattern, "."}).
		Stdout(ctx)
}

// example usage: "dagger call get-secret --aws-credentials ~/.aws/credentials"
func (m *Aws) GetSecret(ctx context.Context, awsCredentials *File) (string, error) {
	credsFile, err := awsCredentials.Contents(ctx)
	if err != nil {
		return "", err
	}
	secret := dag.SetSecret("aws-credential", credsFile)
	return dag.Container().
		From("ubuntu:latest").
		WithMountedSecret("/secret", secret).
		WithExec([]string{"bash", "-c", "cat secret |base64"}).
		Stdout(ctx)
}
