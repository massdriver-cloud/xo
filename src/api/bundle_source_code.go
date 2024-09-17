package api

import (
	"context"

	"github.com/Khan/genqlient/graphql"
)

func GetBundleSourceCode(client graphql.Client, bundleId string, organizationId string) ([]byte, error) {
	response, err := bundleSourceCode(context.Background(), client, bundleId, organizationId)
	if err != nil {
		return nil, err
	}

	return []byte(response.BundleSourceCode.Source), nil
}
