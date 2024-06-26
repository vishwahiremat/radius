//go:build go1.18
// +build go1.18

// Licensed under the Apache License, Version 2.0 . See LICENSE in the repository root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package v20231001preview

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
)

// PlanesClient contains the methods for the Planes group.
// Don't use this type directly, use NewPlanesClient() instead.
type PlanesClient struct {
	internal *arm.Client
}

// NewPlanesClient creates a new instance of PlanesClient with the specified values.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewPlanesClient(credential azcore.TokenCredential, options *arm.ClientOptions) (*PlanesClient, error) {
	cl, err := arm.NewClient(moduleName+".PlanesClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &PlanesClient{
	internal: cl,
	}
	return client, nil
}

// NewListPlanesPager - List all planes
//
// Generated from API version 2023-10-01-preview
//   - options - PlanesClientListPlanesOptions contains the optional parameters for the PlanesClient.NewListPlanesPager method.
func (client *PlanesClient) NewListPlanesPager(options *PlanesClientListPlanesOptions) (*runtime.Pager[PlanesClientListPlanesResponse]) {
	return runtime.NewPager(runtime.PagingHandler[PlanesClientListPlanesResponse]{
		More: func(page PlanesClientListPlanesResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *PlanesClientListPlanesResponse) (PlanesClientListPlanesResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listPlanesCreateRequest(ctx, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return PlanesClientListPlanesResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return PlanesClientListPlanesResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return PlanesClientListPlanesResponse{}, runtime.NewResponseError(resp)
			}
			return client.listPlanesHandleResponse(resp)
		},
	})
}

// listPlanesCreateRequest creates the ListPlanes request.
func (client *PlanesClient) listPlanesCreateRequest(ctx context.Context, options *PlanesClientListPlanesOptions) (*policy.Request, error) {
	urlPath := "/planes"
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-10-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listPlanesHandleResponse handles the ListPlanes response.
func (client *PlanesClient) listPlanesHandleResponse(resp *http.Response) (PlanesClientListPlanesResponse, error) {
	result := PlanesClientListPlanesResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.GenericPlaneResourceListResult); err != nil {
		return PlanesClientListPlanesResponse{}, err
	}
	return result, nil
}

