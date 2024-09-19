package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/iragsraghu/go-aws/utility"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	bucketName := os.Getenv("S3_BUCKET_NAME")

	// Ask whether to upload a file
	fmt.Print("Do you want to upload a file? (yes/no): ")
	uploadResponse, _ := reader.ReadString('\n')
	uploadResponse = strings.TrimSpace(strings.ToLower(uploadResponse))

	if uploadResponse == "yes" {
		// Ask for the file path and object key
		fmt.Print("Enter the full file path (e.g., /path/to/file.pdf): ")
		filePath, _ := reader.ReadString('\n')
		filePath = strings.TrimSpace(filePath)

		fmt.Print("Enter the object key (S3 file name): ")
		objectKey, _ := reader.ReadString('\n')
		objectKey = strings.TrimSpace(objectKey)

		// Upload the file
		if err := utility.UploadObject(bucketName, objectKey, filePath); err != nil {
			log.Fatalf("Failed to upload object: %v", err)
		}

		fmt.Printf("File uploaded successfully to %s/%s\n", bucketName, objectKey)
	}

	// Ask whether to list objects in the bucket
	fmt.Print("Do you want to list objects in the bucket? (yes/no): ")
	listResponse, _ := reader.ReadString('\n')
	listResponse = strings.TrimSpace(strings.ToLower(listResponse))

	if listResponse == "yes" {
		// List objects in the bucket
		objects, err := utility.ListObjects(bucketName)
		if err != nil {
			log.Fatalf("Failed to list objects: %v", err)
		}

		if len(objects) == 0 {
			fmt.Println("No objects found in the bucket.")
		} else {
			fmt.Println("Objects in the bucket:")
			for _, obj := range objects {
				fmt.Println(obj)
			}
		}
	}

	// Ask whether to delete all objects in the bucket
	fmt.Print("Do you want to delete all objects in the bucket? (yes/no): ")
	deleteResponse, _ := reader.ReadString('\n')
	deleteResponse = strings.TrimSpace(strings.ToLower(deleteResponse))

	if deleteResponse == "yes" {
		// List and delete objects
		objects, err := utility.ListObjects(bucketName)
		if err != nil {
			log.Fatalf("Failed to list objects: %v", err)
		}
		if len(objects) != 0 {
			if err := utility.DeleteObjects(bucketName, objects); err != nil {
				log.Fatalf("Failed to delete objects: %v", err)
			}

			fmt.Printf("All objects deleted successfully from %s\n", bucketName)
		} else {
			fmt.Println("No objects found to delete")
		}
	}

	fmt.Println("Program completed.")
}
