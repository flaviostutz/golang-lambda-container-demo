#!/bin/sh

echo "Create ECR repository using Cloudformation..."
aws cloudformation deploy --template-file cf-ecr.yml --stack-name golang-lambda-demo-ecr --capabilities CAPABILITY_NAMED_IAM

repoUri=$(aws cloudformation describe-stacks --stack-name golang-lambda-demo-ecr --query Stacks[].[Outputs[].OutputValue] --output text)

echo "ECR repo is $repoUri"

echo "Building container image and tagging with $repoUri..."
docker build -t $repoUri .

echo "Logging docker cli to ECR repo"
repoDomain=$(echo $repoUri | cut -d/ -f1)
aws ecr --no-verify-ssl get-login-password | docker login --username AWS --password-stdin $repoDomain

echo "Pushing container image to repo..."
docker push $repoUri

echo "Deploying AWS API Gateway and Lambda Functions with Cloudformation..."
aws cloudformation deploy \
    --template-file cf-api-lambda.yml \
    --stack-name golang-lambda-demo-service \
    --parameter-overrides LambdaContainerImageUri=$repoUri \
    --capabilities CAPABILITY_NAMED_IAM

out=$(aws cloudformation describe-stacks --stack-name golang-lambda-demo-service --query Stacks[].[Outputs[].OutputValue] --output text)
apiUri=$(echo $out | cut -f1)

echo "Golang API is running at $apiUri"
curl $apiUri/repo

