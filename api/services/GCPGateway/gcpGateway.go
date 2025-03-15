package gcpgateway

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"webapp/config"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

func UploadFile(bucketName string, objectPath string, fileName string, file multipart.File) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithQuotaProject(config.GetEnvConfig().GOOGLE_PROJECT_ID))
	if err != nil {
		fmt.Println(err)
		return errors.New("Failed to create Cloud Storage client")
	}
	defer client.Close()
	// 3. get the bucker objects
	bucket := client.Bucket(bucketName)
	object := bucket.Object(objectPath + fileName)
	writer := object.NewWriter(ctx)

	// 4. Copy the file contents from the HTTP request to GCS
	if _, err := io.Copy(writer, file); err != nil {
		fmt.Println(err)
		return errors.New("Failed to upload file to Cloud Storage")

	}
	if err := writer.Close(); err != nil {
		fmt.Println(err)
		return errors.New("Failed to finalize the upload")
	}
	return nil
}
