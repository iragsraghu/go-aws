package utility

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func ListObjects(bucketName string) ([]*s3.Object, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("REGION")),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %v", err)
	}

	// Create an S3 service client
	svc := s3.New(sess)

	// Create the ListObjectsV2Input
	listParams := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	}

	// Call the ListObjectsV2 API
	resp, err := svc.ListObjectsV2WithContext(context.TODO(), listParams)
	if err != nil {
		return nil, fmt.Errorf("failed to list objects: %v", err)
	}
	return resp.Contents, nil
}

func UploadObject(bucketName, objectKey, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %v", filePath, err)
	}
	defer file.Close()

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("REGION")),
	})
	if err != nil {
		return fmt.Errorf("failed to create session: %v", err)
	}

	// Create an S3 service client
	svc := s3.New(sess)

	// Ensure bucketName is not empty
	if bucketName == "" {
		return fmt.Errorf("bucketName is empty")
	}

	// Create the PutObjectInput
	uploadParams := &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(objectKey),
		Body:        file,
		ContentType: aws.String("application/pdf"),
	}

	// Call the PutObject API
	_, err = svc.PutObjectWithContext(context.TODO(), uploadParams)
	if err != nil {
		return fmt.Errorf("failed to upload object: %v", err)
	}

	fmt.Printf("Successfully uploaded object %s to bucket %s\n", objectKey, bucketName)
	return nil
}

func DeleteObjects(bucketName string, objectKeys []*s3.Object) error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("REGION")),
	})
	if err != nil {
		return fmt.Errorf("failed to create session: %v", err)
	}

	// Create an S3 service client
	svc := s3.New(sess)

	// Create a slice of DeleteObjectsRequestEntries
	deleteEntries := make([]*s3.ObjectIdentifier, len(objectKeys))
	for i, obj := range objectKeys {
		deleteEntries[i] = &s3.ObjectIdentifier{
			Key: obj.Key,
		}
	}

	// Create the DeleteObjectsInput
	deleteParams := &s3.DeleteObjectsInput{
		Bucket: aws.String(bucketName),
		Delete: &s3.Delete{
			Objects: deleteEntries,
			Quiet:   aws.Bool(true),
		},
	}

	// Call the DeleteObjects API
	_, err = svc.DeleteObjectsWithContext(context.TODO(), deleteParams)
	if err != nil {
		return fmt.Errorf("failed to delete objects: %v", err)
	}

	fmt.Printf("Successfully deleted %d objects from bucket %s\n", len(deleteEntries), bucketName)
	return nil
}
