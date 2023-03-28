package services

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	log "github.com/sirupsen/logrus"
)

func DownloadS3Object(bucket string, object string, path string, region string) error {
	contextLogger := log.WithFields(log.Fields{
		"bucket": bucket,
		"object": object,
		"path":   path,
	})

	file, err := os.Create(path)
	if err != nil {
		contextLogger.Error("Unable to create file locally")
		return err
	}

	defer file.Close()

	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)

	downloader := s3manager.NewDownloader(sess)

	_, err = downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(object),
		})

	if err != nil {
		contextLogger.Error("Unable to download item from S3")
		return err
	}

	return nil
}
