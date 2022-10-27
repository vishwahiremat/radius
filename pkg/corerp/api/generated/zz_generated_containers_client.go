//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package generated

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	armruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"net/url"
	"strings"
)

// ContainersClient contains the methods for the Containers group.
// Don't use this type directly, use NewContainersClient() instead.
type ContainersClient struct {
	host string
	rootScope string
	pl runtime.Pipeline
}

// NewContainersClient creates a new instance of ContainersClient with the specified values.
// rootScope - The scope in which the resource is present. For Azure resource this would be /subscriptions/{subscriptionID}/resourceGroup/{resourcegroupID}
// credential - used to authorize requests. Usually a credential from azidentity.
// options - pass nil to accept the default values.
func NewContainersClient(rootScope string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ContainersClient, error) {
	if options == nil {
		options = &arm.ClientOptions{}
	}
	ep := cloud.AzurePublic.Services[cloud.ResourceManager].Endpoint
	if c, ok := options.Cloud.Services[cloud.ResourceManager]; ok {
		ep = c.Endpoint
	}
	pl, err := armruntime.NewPipeline(moduleName, moduleVersion, credential, runtime.PipelineOptions{}, options)
	if err != nil {
		return nil, err
	}
	client := &ContainersClient{
		rootScope: rootScope,
		host: ep,
pl: pl,
	}
	return client, nil
}

// CreateOrUpdate - Create or update a Container.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-03-15-privatepreview
// containerName - The name of the conatiner.
// containerResource - containers details
// options - ContainersClientCreateOrUpdateOptions contains the optional parameters for the ContainersClient.CreateOrUpdate
// method.
func (client *ContainersClient) CreateOrUpdate(ctx context.Context, containerName string, containerResource ContainerResource, options *ContainersClientCreateOrUpdateOptions) (ContainersClientCreateOrUpdateResponse, error) {
	req, err := client.createOrUpdateCreateRequest(ctx, containerName, containerResource, options)
	if err != nil {
		return ContainersClientCreateOrUpdateResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return ContainersClientCreateOrUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusCreated) {
		return ContainersClientCreateOrUpdateResponse{}, runtime.NewResponseError(resp)
	}
	return client.createOrUpdateHandleResponse(resp)
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *ContainersClient) createOrUpdateCreateRequest(ctx context.Context, containerName string, containerResource ContainerResource, options *ContainersClientCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/{rootScope}/providers/Applications.Core/containers/{containerName}"
	urlPath = strings.ReplaceAll(urlPath, "{rootScope}", client.rootScope)
	if containerName == "" {
		return nil, errors.New("parameter containerName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{containerName}", url.PathEscape(containerName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-03-15-privatepreview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, containerResource)
}

// createOrUpdateHandleResponse handles the CreateOrUpdate response.
func (client *ContainersClient) createOrUpdateHandleResponse(resp *http.Response) (ContainersClientCreateOrUpdateResponse, error) {
	result := ContainersClientCreateOrUpdateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ContainerResource); err != nil {
		return ContainersClientCreateOrUpdateResponse{}, err
	}
	return result, nil
}

// Delete - Delete a Container.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-03-15-privatepreview
// containerName - The name of the conatiner.
// options - ContainersClientDeleteOptions contains the optional parameters for the ContainersClient.Delete method.
func (client *ContainersClient) Delete(ctx context.Context, containerName string, options *ContainersClientDeleteOptions) (ContainersClientDeleteResponse, error) {
	req, err := client.deleteCreateRequest(ctx, containerName, options)
	if err != nil {
		return ContainersClientDeleteResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return ContainersClientDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusAccepted, http.StatusNoContent) {
		return ContainersClientDeleteResponse{}, runtime.NewResponseError(resp)
	}
	return ContainersClientDeleteResponse{}, nil
}

// deleteCreateRequest creates the Delete request.
func (client *ContainersClient) deleteCreateRequest(ctx context.Context, containerName string, options *ContainersClientDeleteOptions) (*policy.Request, error) {
	urlPath := "/{rootScope}/providers/Applications.Core/containers/{containerName}"
	urlPath = strings.ReplaceAll(urlPath, "{rootScope}", client.rootScope)
	if containerName == "" {
		return nil, errors.New("parameter containerName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{containerName}", url.PathEscape(containerName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-03-15-privatepreview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Gets the properties of an Container.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-03-15-privatepreview
// containerName - The name of the conatiner.
// options - ContainersClientGetOptions contains the optional parameters for the ContainersClient.Get method.
func (client *ContainersClient) Get(ctx context.Context, containerName string, options *ContainersClientGetOptions) (ContainersClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, containerName, options)
	if err != nil {
		return ContainersClientGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return ContainersClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return ContainersClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *ContainersClient) getCreateRequest(ctx context.Context, containerName string, options *ContainersClientGetOptions) (*policy.Request, error) {
	urlPath := "/{rootScope}/providers/Applications.Core/containers/{containerName}"
	urlPath = strings.ReplaceAll(urlPath, "{rootScope}", client.rootScope)
	if containerName == "" {
		return nil, errors.New("parameter containerName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{containerName}", url.PathEscape(containerName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-03-15-privatepreview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *ContainersClient) getHandleResponse(resp *http.Response) (ContainersClientGetResponse, error) {
	result := ContainersClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ContainerResource); err != nil {
		return ContainersClientGetResponse{}, err
	}
	return result, nil
}

// NewListByScopePager - List all containers in the given scope.
// Generated from API version 2022-03-15-privatepreview
// options - ContainersClientListByScopeOptions contains the optional parameters for the ContainersClient.ListByScope method.
func (client *ContainersClient) NewListByScopePager(options *ContainersClientListByScopeOptions) (*runtime.Pager[ContainersClientListByScopeResponse]) {
	return runtime.NewPager(runtime.PagingHandler[ContainersClientListByScopeResponse]{
		More: func(page ContainersClientListByScopeResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *ContainersClientListByScopeResponse) (ContainersClientListByScopeResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listByScopeCreateRequest(ctx, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return ContainersClientListByScopeResponse{}, err
			}
			resp, err := client.pl.Do(req)
			if err != nil {
				return ContainersClientListByScopeResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return ContainersClientListByScopeResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByScopeHandleResponse(resp)
		},
	})
}

// listByScopeCreateRequest creates the ListByScope request.
func (client *ContainersClient) listByScopeCreateRequest(ctx context.Context, options *ContainersClientListByScopeOptions) (*policy.Request, error) {
	urlPath := "/{rootScope}/providers/Applications.Core/containers"
	urlPath = strings.ReplaceAll(urlPath, "{rootScope}", client.rootScope)
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-03-15-privatepreview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByScopeHandleResponse handles the ListByScope response.
func (client *ContainersClient) listByScopeHandleResponse(resp *http.Response) (ContainersClientListByScopeResponse, error) {
	result := ContainersClientListByScopeResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ContainerResourceList); err != nil {
		return ContainersClientListByScopeResponse{}, err
	}
	return result, nil
}

// Update - Update the properties of an existing Container.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-03-15-privatepreview
// containerName - The name of the conatiner.
// containersResource - Container details
// options - ContainersClientUpdateOptions contains the optional parameters for the ContainersClient.Update method.
func (client *ContainersClient) Update(ctx context.Context, containerName string, containersResource ContainerResource, options *ContainersClientUpdateOptions) (ContainersClientUpdateResponse, error) {
	req, err := client.updateCreateRequest(ctx, containerName, containersResource, options)
	if err != nil {
		return ContainersClientUpdateResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return ContainersClientUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusCreated) {
		return ContainersClientUpdateResponse{}, runtime.NewResponseError(resp)
	}
	return client.updateHandleResponse(resp)
}

// updateCreateRequest creates the Update request.
func (client *ContainersClient) updateCreateRequest(ctx context.Context, containerName string, containersResource ContainerResource, options *ContainersClientUpdateOptions) (*policy.Request, error) {
	urlPath := "/{rootScope}/providers/Applications.Core/containers/{containerName}"
	urlPath = strings.ReplaceAll(urlPath, "{rootScope}", client.rootScope)
	if containerName == "" {
		return nil, errors.New("parameter containerName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{containerName}", url.PathEscape(containerName))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-03-15-privatepreview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, containersResource)
}

// updateHandleResponse handles the Update response.
func (client *ContainersClient) updateHandleResponse(resp *http.Response) (ContainersClientUpdateResponse, error) {
	result := ContainersClientUpdateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ContainerResource); err != nil {
		return ContainersClientUpdateResponse{}, err
	}
	return result, nil
}
