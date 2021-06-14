package main

import (
	"context"
	"flag"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/flaviostutz/golang-lambda-container-demo/handlers"
	"github.com/sirupsen/logrus"
)

var repo map[string]string
var httpServer *handlers.HTTPServer
var ginLambda *ginadapter.GinLambda

func init() {
	//flags parse
	logLevel := "info"
	readonly := false
	flag.StringVar(&logLevel, "loglevel", "info", "Log level")
	flag.BoolVar(&readonly, "readonly", false, "Repo is readonly")
	flag.Parse()

	l, err := logrus.ParseLevel(logLevel)
	if err != nil {
		panic("Invalid loglevel")
	}
	logrus.SetLevel(l)

	//key value repo
	repo = make(map[string]string)
	repo["key1"] = "value1"
	repo["key2"] = "value2"

	//setup gin routes
	opt := handlers.Options{
		Readonly: readonly,
	}
	httpServer, err = handlers.NewHTTPServer(opt, repo)
	if err != nil {
		panic(err)
	}
	ginLambda = ginadapter.New(httpServer.Router)
}

func main() {
	//event handling for a single function
	// lambda.Start(SingleFunction)

	endpoint := os.Getenv("ENDPOINT")
	if endpoint == "http" {
		logrus.Infof("Starting HTTP server with Gin at :3000 (no AWS Lambda)")
		httpServer.Server.ListenAndServe()
	}

	logrus.Infof("Starting Lambda function (Gin wrapper)")
	lambda.Start(GinProxyHandler)
}

func GinProxyHandler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

// type Evt struct {
// 	Name string `json:"name"`
// 	Age  int    `json:"age"`
// }
// type Resp struct {
// 	Message string `json:"message"`
// }

// func SingleFunction(ctx context.Context, event Evt) (Resp, error) {
// 	return Resp{Message: fmt.Sprintf("%s is %d years old!", event.Name, event.Age)}, nil
// }
