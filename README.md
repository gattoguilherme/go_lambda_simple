# 🚀 Deploy a Go "Hello World" API to AWS Lambda + API Gateway (REST)

This guide shows how to build, package, deploy, and test a simple **Go REST API** running on **AWS Lambda** integrated with **API Gateway**.

---

## 🧱 1. Prerequisites

| Tool | Purpose | Check |
|------|----------|--------|
| **Go** (>= 1.22) | Build and compile the Lambda binary | `go version` |
| **AWS CLI** | Deploy Lambda & API Gateway | `aws --version` |
| **IAM Role** | Lambda execution role | Must include `AWSLambdaBasicExecutionRole` |
| **ZIP Utility** | To package the function | Built-in on Windows/macOS/Linux |

---

## ⚙️ 2. Create Project Structure

```bash
mkdir go-hello-lambda
cd go-hello-lambda
go mod init go-hello-lambda
go mod tidy
```

---

## 🧩 3. Implement Go Code

Create a file named `main.go`:

```go
package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	name := request.QueryStringParameters["name"]
	if name == "" {
		name = "World"
	}
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       fmt.Sprintf("Hello, %s!", name),
	}, nil
}

func main() {
	lambda.Start(handler)
}
```

Install Lambda dependencies:

```bash
go get github.com/aws/aws-lambda-go/events
go get github.com/aws/aws-lambda-go/lambda
```

---

## 🏗️ 4. Build for AWS Lambda Environment

AWS Lambda runs on **Amazon Linux 2**.  
You must cross-compile your binary for **Linux x86_64** (or ARM if selected).

### On Windows PowerShell:
```powershell
$env:GOOS = "linux"
$env:GOARCH = "amd64"
$env:CGO_ENABLED = "0"
go build -o bootstrap main.go
```

### On macOS/Linux:
```bash
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bootstrap main.go
```

✅ Output: `bootstrap` (no `.exe` extension)

---

## 📦 5. Package the Function

### On Windows PowerShell:
```powershell
Compress-Archive -Path bootstrap -DestinationPath function.zip -Force
```

### On macOS/Linux:
```bash
zip function.zip bootstrap
```

Ensure your ZIP contains only:
```
function.zip
└── bootstrap
```

---

## ☁️ 6. Create the Lambda Function

You can create via **AWS CLI** or **AWS Console**.

### ▶️ Using AWS CLI
```bash
aws lambda create-function   --function-name go-hello-lambda   --runtime provided.al2   --handler bootstrap   --zip-file fileb://function.zip   --role arn:aws:iam::<YOUR_ACCOUNT_ID>:role/<YOUR_LAMBDA_ROLE>   --architectures x86_64   --region sa-east-1
```

> The IAM role must have the policy `AWSLambdaBasicExecutionRole`.

### 🖥️ Using AWS Console
1. Go to **AWS Lambda → Create function**
2. Choose **Author from scratch**
3. Function name: `go-hello-lambda`
4. Runtime: **Amazon Linux 2**
5. Architecture: **x86_64**
6. Choose or create an execution role with **AWSLambdaBasicExecutionRole**
7. After creation:
   - Under **Code → Upload from → .zip file**
   - Upload your `function.zip`
   - Set **Handler = bootstrap**
   - Click **Deploy**

---

## 🌐 7. Create REST API Gateway

### 🪄 Steps (Console)
1. Open **API Gateway → Create API**
2. Select **REST API → Build**
3. Choose **New API**
4. Name: `go-hello-api`
5. Click **Create API**

#### ➕ Create Resource
1. In the left panel, click **Actions → Create Resource**
2. Resource Name: `hello`
3. Resource Path: `/hello`
4. Save

#### ⚙️ Create Method
1. With `/hello` selected → **Actions → Create Method → GET**
2. Integration type: **Lambda Function**
3. Select Region
4. Function name: `go-hello-lambda`
5. Click **Save**
6. Confirm with **“Add permissions”**

#### 🚀 Deploy API
1. Click **Actions → Deploy API**
2. Deployment Stage: **New Stage → prod**
3. Click **Deploy**

You’ll now get an **Invoke URL**, such as:
```
https://abcd1234.execute-api.sa-east-1.amazonaws.com/prod
```

---

## 🧪 8. Test the Endpoint

### Default:
```bash
curl "https://abcd1234.execute-api.sa-east-1.amazonaws.com/prod/hello"
```
**Response:**
```
Hello, World!
```

### With a query parameter:
```bash
curl "https://abcd1234.execute-api.sa-east-1.amazonaws.com/prod/hello?name=Guilherme"
```
**Response:**
```
Hello, Guilherme!
```

---

## 🔁 9. Update the Lambda (after code changes)

Each time you change your Go code:

```bash
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bootstrap main.go
zip function.zip bootstrap
aws lambda update-function-code   --function-name go-hello-lambda   --zip-file fileb://function.zip   --region sa-east-1
```

---

## 🧠 10. Architecture Overview

```
Client
  ↓
API Gateway (REST)
  ↓
Lambda (Go binary on Amazon Linux 2)
  ↓
CloudWatch Logs (optional)
```

---

## ✅ Summary

| Step | Description |
|------|--------------|
| 1 | Create Go project |
| 2 | Write handler using `aws-lambda-go` |
| 3 | Cross-compile for Linux |
| 4 | Package ZIP file |
| 5 | Create Lambda function |
| 6 | Create REST API Gateway integration |
| 7 | Deploy & test |
| 8 | Update easily after changes |

---

**🎉 Congratulations!**  
You now have a fully deployed **Go REST API** running on **AWS Lambda**, exposed through **API Gateway**, and callable via HTTPS.
