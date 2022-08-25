go get github.com/aws/aws-lambda-go/lambda
GOOS=linux go build -tags awslambda -o main .
zip function.zip main default.png