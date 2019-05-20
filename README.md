Clone this repository

We need a valid AWS Account and a IAM user with a minimal set of permissions to create a S3 Bucket and create a SQS Queue.
Upon creating the account, get the values needed and export into the environment
```
export AWS_REGION="us-east-1" \
export AWS_DEFAULT_REGION="us-east-1" \
export AWS_ACCESS_KEY_ID="" \
export AWS_SECRET_ACCESS_KEY=""

```

Run terraform on ./terraform/src/s3-bucket folder to create the s3 bucket.
```
cd terraform/src/s3-bucket
terraform init && terraform apply
```
When asked, approve the creation of the bucket under your AWS Account


Run terraform on ./terraform/src/queue folder to create the queue with a state file on the S3 Bucket created above.
```
cd terraform/src/queue
terraform init && terraform apply
```
Again when asked for confirmation, approve it and it will create the SQS


Build the images
```
docker-compose build
```

Run docker-compose
```
docker-compose up
```

I hardcoded the queue name as `mm-sqs` both at the terraform and Go Source code

the Go API and Worker will use the AWS credentials provided as environment variables to discover the queue using the AWS SDK for Go

The API will listen on port 8123 at the path `/api` using only `POST`
The expected JSON data will be simillar to the following:
```
{
	"_id": "1",
	"suggestion": "testing first",
	"person": {
		"firstname":"Giovanni",
		"lastname": "Coutinho"
	}
}
```
It will then trigger a function to send the whole JSON to SQS
The worker will be responsible for decoding the JSON so it will be possible to get all the information we need.


The worker still has some work to be done
 - Didnt find time to better understand and work with concurrency in Go
 - Didnt focus on getting MongoDB logic working as I wanted to work on concurrency first
 - It will search for messages, print it and then delete it once.
