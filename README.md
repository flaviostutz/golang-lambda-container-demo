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

* It will:

  * Build and push this image to a AWS ECR repository
  * Deploy the Lambda function to AWS refering to the container image
  * Create AWS API Gateway with Proxy to the function
  * Call the API endpoint for testing a request

## Details on how to deploy API with Lambda Proxy Integration

https://docs.aws.amazon.com/apigateway/latest/developerguide/set-up-lambda-proxy-integrations.html

