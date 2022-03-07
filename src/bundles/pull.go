package bundles

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3API interface {
	GetObject(ctx context.Context,
		params *s3.GetObjectInput,
		optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)

	ListObjectsV2(ctx context.Context,
		params *s3.ListObjectsV2Input,
		optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error)
}

func Pull(bucket string, key string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client := s3.NewFromConfig(cfg)

	// bucket := "xo-prod-bundlebucket-0000"
	// prefix := ""

	// input := &s3.ListObjectsV2Input{
	// 	Bucket: &bucket,
	// 	Prefix: &prefix,
	// }

	// resp, err := client.ListObjectsV2(context.TODO(), input)

	// if err != nil {
	// 	fmt.Println("Got error retrieving list of objects:")
	// 	return err
	// }

	// fmt.Println("Objects in " + bucket + ":")

	// for _, item := range resp.Contents {
	// 	fmt.Println("Name:          ", *item.Key)
	// 	fmt.Println("Last modified: ", *item.LastModified)
	// 	fmt.Println("Size:          ", item.Size)
	// 	fmt.Println("Storage class: ", item.StorageClass)
	// 	fmt.Println("")

	getInput := &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	}

	obj, err := client.GetObject(context.TODO(), getInput)
	if err != nil {
		return err
	}
	outFile, err := os.Create(filepath.Base(key))
	if err != nil {
		return err
	}
	// handle err
	defer outFile.Close()
	_, err = io.Copy(outFile, obj.Body)
	if err != nil {
		return err
	}
	// handle err

	// }

	// fmt.Println("Found", len(resp.Contents), "items in bucket", bucket)
	// fmt.Println("")

	return nil
}
