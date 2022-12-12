go get -u github.com/aws/aws-lambda-go/cmd/build-lambda-zip
go install github.com/aws/aws-lambda-go/cmd/build-lambda-zip
$env:GOOS = "linux"
$env:CGO_ENABLED = "0"
$env:GOARCH = "amd64"
go build -tags awslambda -o main
~\Go\Bin\build-lambda-zip.exe -output main.zip main default.png