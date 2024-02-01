package main

import (
	"context"
)

type Aws struct{}

// example usage: "dagger call get-secret --aws-credentials ~/.aws/credentials"
func (m *Aws) GetSecret(ctx context.Context, awsCredentials *File) (string, error) {
	ctr, err := m.WithAwsSecret(ctx, dag.Container().From("ubuntu:latest"), awsCredentials)
	if err != nil {
		return "", err
	}
	return ctr.
		WithExec([]string{"bash", "-c", "cat /root/.aws/credentials |base64"}).
		Stdout(ctx)
}

func (m *Aws) WithAwsSecret(ctx context.Context, ctr *Container, awsCredentials *File) (*Container, error) {
	credsFile, err := awsCredentials.Contents(ctx)
	if err != nil {
		return nil, err
	}
	secret := dag.SetSecret("aws-credential", credsFile)
	return ctr.WithMountedSecret("/root/.aws/credentials", secret), nil
}

// example usage: "dagger call list --aws-credentials ~/.aws/credentials"
func (m *Aws) List(ctx context.Context, awsCredentials *File) (string, error) {
	ctr := dag.Container().
		From("public.ecr.aws/aws-cli/aws-cli:latest")
	ctr, err := m.WithAwsSecret(ctx, ctr, awsCredentials)
	if err != nil {
		return "", err
	}
	return ctr.
		WithExec([]string{"s3", "ls"}).
		Stdout(ctx)
}
