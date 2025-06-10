package bundle

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"path/filepath"
	"xo/src/api"
	"xo/src/massdriver"
	"xo/src/telemetry"

	"github.com/massdriver-cloud/massdriver-sdk-go/massdriver/client"
	"github.com/massdriver-cloud/massdriver-sdk-go/massdriver/config"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"

	oras "oras.land/oras-go/v2"
	"oras.land/oras-go/v2/registry/remote"
	"oras.land/oras-go/v2/registry/remote/auth"
	"oras.land/oras-go/v2/registry/remote/retry"
)

func PullV0(ctx context.Context, client *massdriver.MassdriverClient, outFile io.Writer) error {
	_, span := otel.Tracer("xo").Start(ctx, "BundlePull")
	telemetry.SetSpanAttributes(span)
	defer span.End()

	bundleBytes, getErr := api.GetBundleSourceCode(client.GQLCLient, client.Specification.BundleID, client.Specification.OrganizationID)
	if getErr != nil {
		span.RecordError(getErr)
		span.SetStatus(codes.Error, getErr.Error())
		return fmt.Errorf("an error occurred while getting bundle source code: %w", getErr)
	}

	decodedBundleBytes, decodeErr := base64.StdEncoding.DecodeString(string(bundleBytes))
	if decodeErr != nil {
		span.RecordError(decodeErr)
		span.SetStatus(codes.Error, decodeErr.Error())
		return fmt.Errorf("an error occurred while decoding bundle source code: %w", decodeErr)
	}

	_, writeErr := outFile.Write(decodedBundleBytes)
	if writeErr != nil {
		span.RecordError(writeErr)
		span.SetStatus(codes.Error, writeErr.Error())
		return fmt.Errorf("an error occurred while writing bundle source code: %w", writeErr)
	}

	return nil
}

func PullV1(ctx context.Context, repo oras.Target, target oras.Target, tag string) (v1.Descriptor, error) {
	_, span := otel.Tracer("xo").Start(ctx, "BundlePullV1")
	telemetry.SetSpanAttributes(span)
	defer span.End()

	return oras.Copy(ctx, repo, tag, target, tag, oras.DefaultCopyOptions)
}

func GetRepo(mdClient *client.Client, organizationSlug string, bundleName string) (oras.Target, error) {
	if mdClient.Auth.Method != config.AuthAPIKey {
		return nil, fmt.Errorf("bundle publish requires API key auth")
	}
	// reg := mdClient.Auth.URL
	// repo, repoErr := remote.NewRepository(filepath.Join(reg, mdClient.Auth.AccountID, b.Name))
	reg := "2d67-47-229-209-228.ngrok-free.app"
	repo, repoErr := remote.NewRepository(filepath.Join(reg, organizationSlug, bundleName))
	if repoErr != nil {
		return nil, repoErr
	}

	repo.Client = &auth.Client{
		Client: retry.DefaultClient,
		Cache:  auth.NewCache(),
		Credential: auth.StaticCredential(reg, auth.Credential{
			Username: "myuser",
			Password: "mypass",
		}),
	}

	return repo, nil
}
