// Push a container image into Amazon Elastic Container Registry (ECR)
//
// This module lets you push a container into ECR, automating the tedious manual steps of configuring your local docker daemon with the credentials
//
// For more info and sample usage, check the readme: https://github.com/lukemarsden/dagger-aws

package main

import (
	"context"
	"fmt"
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

func (m *Aws) AwsCli(ctx context.Context, awsCredentials *File) (*Container, error) {
	ctr := dag.Container().
		From("public.ecr.aws/aws-cli/aws-cli:latest")
	ctr, err := m.WithAwsSecret(ctx, ctr, awsCredentials)
	if err != nil {
		return nil, err
	}
	return ctr, nil
}

// Not called S3List because for some reason dagger fails to translate that to s3-list
// example usage: "dagger call list --aws-credentials ~/.aws/credentials"
func (m *Aws) List(ctx context.Context, awsCredentials *File) (string, error) {
	ctr, err := m.AwsCli(ctx, awsCredentials)
	if err != nil {
		return "", err
	}
	return ctr.
		WithExec([]string{"s3", "ls"}).
		Stdout(ctx)
}

// example usage: "dagger call ecr-get-login-password --region us-east-1 --aws-credentials ~/.aws/credentials"
func (m *Aws) EcrGetLoginPassword(ctx context.Context, awsCredentials *File, region string) (string, error) {
	ctr, err := m.AwsCli(ctx, awsCredentials)
	if err != nil {
		return "", err
	}
	return ctr.
		WithExec([]string{"--region", region, "ecr", "get-login-password"}).
		Stdout(ctx)
}

// Push ubuntu:latest to ECR under given repo 'test' (repo must be created first)
// example usage: "dagger call ecr-push-example --region us-east-1 --aws-credentials ~/.aws/credentials --aws-account-id 12345 --repo test"
func (m *Aws) EcrPushExample(ctx context.Context, awsCredentials *File, region, awsAccountId, repo string) (string, error) {
	ctr := dag.Container().From("ubuntu:latest")
	return m.EcrPush(ctx, awsCredentials, region, awsAccountId, repo, ctr)
}

func (m *Aws) EcrPush(ctx context.Context, awsCredentials *File, region, awsAccountId, repo string, pushCtr *Container) (string, error) {
	// Get the ECR login password so we can authenticate with Publish WithRegistryAuth
	ctr, err := m.AwsCli(ctx, awsCredentials)
	if err != nil {
		return "", err
	}
	regCred, err := ctr.
		WithExec([]string{"--region", region, "ecr", "get-login-password"}).
		Stdout(ctx)
	if err != nil {
		return "", err
	}
	secret := dag.SetSecret("aws-reg-cred", regCred)
	ecrHost := fmt.Sprintf("%s.dkr.ecr.%s.amazonaws.com", awsAccountId, region)
	ecrWithRepo := fmt.Sprintf("%s/%s", ecrHost, repo)

	return pushCtr.WithRegistryAuth(ecrHost, "AWS", secret).Publish(ctx, ecrWithRepo)
}
