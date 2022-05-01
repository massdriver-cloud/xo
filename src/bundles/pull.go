package bundles

import (
	"context"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

type S3API interface {
	GetObject(ctx context.Context,
		params *s3.GetObjectInput,
		optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)

	ListObjectsV2(ctx context.Context,
		params *s3.ListObjectsV2Input,
		optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error)
}

func Pull(ctx context.Context, bundleBucket string, organizationId string, bundleId string) error {
	_, span := otel.Tracer("xo").Start(ctx, "BundlePull")
	defer span.End()

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client := s3.NewFromConfig(cfg)

	key := path.Join("bundles", organizationId, bundleId, "bundle.tar.gz")

	log.Info().Msg("attempting to pull s3://" + bundleBucket + "/" + key)

	getInput := &s3.GetObjectInput{
		Bucket: &bundleBucket,
		Key:    &key,
	}

	obj, err := client.GetObject(context.TODO(), getInput)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	outFile, err := os.Create(filepath.Base(key))
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	defer outFile.Close()
	_, err = io.Copy(outFile, obj.Body)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	return nil
}
