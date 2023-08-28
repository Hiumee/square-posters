go get github.com/aws/aws-lambda-go/lambda
GOARCH=arm64 GOOS=linux go build -tags="lambda.norpc,awslambda" -o bootstrap .
zip function.zip bootstrap default.png