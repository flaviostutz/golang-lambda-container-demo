# golang-lambda-container-demo

Demo on how to run AWS Lambda functions in Golang built with a custom container and routing API calls using Gin adapter for AWS API Gateway Proxy requests.

Using this structure you can make your existing Golang Gin backend applications to run as AWS Lambda services with minimal effort. In this example, only the main.go file had to be changed.

During development, run this container with `docker-compose up`. It will start a regular Gin HTTP server with no AWS Lambda bindings.

For Lambda deployments, the same container image can be deployed as a AWS Lambda function without modifications.

Check complete script with Cloudformation templates for provisioning all AWS resources below.

## Usage

### For local or plain HTTP deployments

```sh
docker-compose up
curl localhost:3000/repo
curl -X POST localhost:3000/repo/key3 --data 'value3'
curl localhost:3000/repo
```

### For Lambda Function deployments

* Configure aws cli environment on your machine with `aws configure`

* run `build-deploy-aws-lambda-demo.sh`

```sh
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
version=$(docker inspect --format='{{index .Id}}' $repoUri | cut -d\: -f2)
repoUriTag=$repoUri\:$version
docker tag $repoUri $repoUriTag
docker push $repoUriTag

echo "Deploying AWS API Gateway and Lambda Functions with Cloudformation..."
aws cloudformation deploy \
    --template-file cf-api-lambda.yml \
    --stack-name golang-lambda-demo-service \
    --parameter-overrides LambdaContainerImageUri=$repoUriTag \
    --capabilities CAPABILITY_NAMED_IAM

sleep 3
out=$(aws cloudformation describe-stacks --stack-name golang-lambda-demo-service --query Stacks[].[Outputs[].OutputValue] --output text)
apiUri=$(echo $out | cut -f1)

echo "Golang API is running at $apiUri"
set +x
curl $apiUri/repo
```

* It will:

  * Build and push this image to a AWS ECR repository
  * Deploy the Lambda function to AWS refering to the container image
  * Create AWS API Gateway with Proxy to the function
  * Call the API endpoint for testing a request

## Steps to add support to a Golang application to work as Lambda Proxy

* Dockerfile
  * create ENV ENV ENDPOINT 'http'
  * http will launch regular Gin http server. 'lambda' will launch Lambda Gin bridge

* Add "ENDPOINT" flag to startup.sh and main.go

* run `go get github.com/awslabs/aws-lambda-go-api-proxy/gin` to add bridge as dependency

* Add ginLambda = ginadapter.New(httpServer.Router) to main.go init()

* Add to main.go main()

```golang
	if h.Opt.Endpoint == "http" {
    //starting regular gin server (without Lambda proxy)
		err := h.Start()
    ...
    return
  }

  //start Lambda Proxy to Gin bridge
	lambda.Start(func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		return ginLambda.ProxyWithContext(ctx, req)
	})

```



## Details on how to deploy API with Lambda Proxy Integration

https://docs.aws.amazon.com/apigateway/latest/developerguide/set-up-lambda-proxy-integrations.html

