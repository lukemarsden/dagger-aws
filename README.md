# Dagger AWS module

Known to work with Dagger v0.9.5.

## S3 List

List the S3 buckets in your account:
```
dagger call -m github.com/lukemarsden/dagger-aws list --aws-credentials ~/.aws/credentials
```

## Push Image to Private ECR Repo

From CLI, to push `ubuntu:latest` to a given ECR repo, by way of example:

```
dagger call -m github.com/lukemarsden/dagger-aws \
    ecr-push-example --region us-east-1 \
    --aws-credentials ~/.aws/credentials \
    --aws-account-id 12345 --repo test
```

Check `region` and `aws-account-id` arguments and update them to match your AWS account and location of your private repo. Update `repo` to the name of your repo.

Make sure you've created the repo in your ECR account. You can do that under Amazon ECR --> Private registry --> Repositories in the AWS console.

## From Dagger Code

Call the EcrPush method on this module with the awsCredentials *File (e.g. ~/.aws/credentials) as the first argument, then the region, awsAccountId and repo as strings, then finally the container you wish to push as the final argument.

For example:

```go
func (y *YourThing) PushYourThings(ctx context.Context, awsCredentials *File) {
    ctr := dag.Container()
        .From("yourbase:image")
        .YourThings()
    // get region, awsAccountId, repo
    out, err := m.EcrPush(ctx, awsCredentials, region, awsAccountId, repo, ctr)
}
```

See `EcrPushExample` for a concrete example.