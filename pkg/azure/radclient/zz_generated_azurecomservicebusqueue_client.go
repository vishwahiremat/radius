//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package radclient

import (
	"context"
	"errors"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	armruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"net/url"
	"strings"
)

// AzureComServiceBusQueueClient contains the methods for the AzureComServiceBusQueue group.
// Don't use this type directly, use NewAzureComServiceBusQueueClient() instead.
type AzureComServiceBusQueueClient struct {
	ep string
	pl runtime.Pipeline
	subscriptionID string
}

// NewAzureComServiceBusQueueClient creates a new instance of AzureComServiceBusQueueClient with the specified values.
func NewAzureComServiceBusQueueClient(con *arm.Connection, subscriptionID string) *AzureComServiceBusQueueClient {
	return &AzureComServiceBusQueueClient{ep: con.Endpoint(), pl: con.NewPipeline(module, version), subscriptionID: subscriptionID}
}

// BeginCreateOrUpdate - Creates or updates a azure.com.ServiceBusQueue resource.
// If the operation fails it returns the *ErrorResponse error type.
func (client *AzureComServiceBusQueueClient) BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, applicationName string, azureServiceBusName string, parameters AzureServiceBusResource, options *AzureComServiceBusQueueBeginCreateOrUpdateOptions) (AzureComServiceBusQueueCreateOrUpdatePollerResponse, error) {
	resp, err := client.createOrUpdate(ctx, resourceGroupName, applicationName, azureServiceBusName, parameters, options)
	if err != nil {
		return AzureComServiceBusQueueCreateOrUpdatePollerResponse{}, err
	}
	result := AzureComServiceBusQueueCreateOrUpdatePollerResponse{
		RawResponse: resp,
	}
	pt, err := armruntime.NewPoller("AzureComServiceBusQueueClient.CreateOrUpdate", "location", resp, 	client.pl, client.createOrUpdateHandleError)
	if err != nil {
		return AzureComServiceBusQueueCreateOrUpdatePollerResponse{}, err
	}
	result.Poller = &AzureComServiceBusQueueCreateOrUpdatePoller {
		pt: pt,
	}
	return result, nil
}

// CreateOrUpdate - Creates or updates a azure.com.ServiceBusQueue resource.
// If the operation fails it returns the *ErrorResponse error type.
func (client *AzureComServiceBusQueueClient) createOrUpdate(ctx context.Context, resourceGroupName string, applicationName string, azureServiceBusName string, parameters AzureServiceBusResource, options *AzureComServiceBusQueueBeginCreateOrUpdateOptions) (*http.Response, error) {
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, applicationName, azureServiceBusName, parameters, options)
	if err != nil {
		return nil, err
	}
	resp, err := 	client.pl.Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusCreated, http.StatusAccepted) {
		return nil, client.createOrUpdateHandleError(resp)
	}
	 return resp, nil
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *AzureComServiceBusQueueClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, applicationName string, azureServiceBusName string, parameters AzureServiceBusResource, options *AzureComServiceBusQueueBeginCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.CustomProviders/resourceProviders/radiusv3/Application/{applicationName}/azure.com.ServiceBusQueue/{azureServiceBusName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if applicationName == "" {
		return nil, errors.New("parameter applicationName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{applicationName}", url.PathEscape(applicationName))
	if azureServiceBusName == "" {
		return nil, errors.New("parameter azureServiceBusName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{azureServiceBusName}", url.PathEscape(azureServiceBusName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(	client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2018-09-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, parameters)
}

// createOrUpdateHandleError handles the CreateOrUpdate error response.
func (client *AzureComServiceBusQueueClient) createOrUpdateHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
		errType := ErrorResponse{raw: string(body)}
	if err := runtime.UnmarshalAsJSON(resp, &errType); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}

// BeginDelete - Deletes a azure.com.ServiceBusQueue resource.
// If the operation fails it returns the *ErrorResponse error type.
func (client *AzureComServiceBusQueueClient) BeginDelete(ctx context.Context, resourceGroupName string, applicationName string, azureServiceBusName string, options *AzureComServiceBusQueueBeginDeleteOptions) (AzureComServiceBusQueueDeletePollerResponse, error) {
	resp, err := client.deleteOperation(ctx, resourceGroupName, applicationName, azureServiceBusName, options)
	if err != nil {
		return AzureComServiceBusQueueDeletePollerResponse{}, err
	}
	result := AzureComServiceBusQueueDeletePollerResponse{
		RawResponse: resp,
	}
	pt, err := armruntime.NewPoller("AzureComServiceBusQueueClient.Delete", "location", resp, 	client.pl, client.deleteHandleError)
	if err != nil {
		return AzureComServiceBusQueueDeletePollerResponse{}, err
	}
	result.Poller = &AzureComServiceBusQueueDeletePoller {
		pt: pt,
	}
	return result, nil
}

// Delete - Deletes a azure.com.ServiceBusQueue resource.
// If the operation fails it returns the *ErrorResponse error type.
func (client *AzureComServiceBusQueueClient) deleteOperation(ctx context.Context, resourceGroupName string, applicationName string, azureServiceBusName string, options *AzureComServiceBusQueueBeginDeleteOptions) (*http.Response, error) {
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, applicationName, azureServiceBusName, options)
	if err != nil {
		return nil, err
	}
	resp, err := 	client.pl.Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusAccepted, http.StatusNoContent) {
		return nil, client.deleteHandleError(resp)
	}
	 return resp, nil
}

// deleteCreateRequest creates the Delete request.
func (client *AzureComServiceBusQueueClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, applicationName string, azureServiceBusName string, options *AzureComServiceBusQueueBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.CustomProviders/resourceProviders/radiusv3/Application/{applicationName}/azure.com.ServiceBusQueue/{azureServiceBusName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if applicationName == "" {
		return nil, errors.New("parameter applicationName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{applicationName}", url.PathEscape(applicationName))
	if azureServiceBusName == "" {
		return nil, errors.New("parameter azureServiceBusName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{azureServiceBusName}", url.PathEscape(azureServiceBusName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(	client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2018-09-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// deleteHandleError handles the Delete error response.
func (client *AzureComServiceBusQueueClient) deleteHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
		errType := ErrorResponse{raw: string(body)}
	if err := runtime.UnmarshalAsJSON(resp, &errType); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}

// Get - Gets a azure.com.ServiceBusQueue resource by name.
// If the operation fails it returns the *ErrorResponse error type.
func (client *AzureComServiceBusQueueClient) Get(ctx context.Context, resourceGroupName string, applicationName string, azureServiceBusName string, options *AzureComServiceBusQueueGetOptions) (AzureComServiceBusQueueGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, applicationName, azureServiceBusName, options)
	if err != nil {
		return AzureComServiceBusQueueGetResponse{}, err
	}
	resp, err := 	client.pl.Do(req)
	if err != nil {
		return AzureComServiceBusQueueGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return AzureComServiceBusQueueGetResponse{}, client.getHandleError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *AzureComServiceBusQueueClient) getCreateRequest(ctx context.Context, resourceGroupName string, applicationName string, azureServiceBusName string, options *AzureComServiceBusQueueGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.CustomProviders/resourceProviders/radiusv3/Application/{applicationName}/azure.com.ServiceBusQueue/{azureServiceBusName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if applicationName == "" {
		return nil, errors.New("parameter applicationName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{applicationName}", url.PathEscape(applicationName))
	if azureServiceBusName == "" {
		return nil, errors.New("parameter azureServiceBusName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{azureServiceBusName}", url.PathEscape(azureServiceBusName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(	client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2018-09-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *AzureComServiceBusQueueClient) getHandleResponse(resp *http.Response) (AzureComServiceBusQueueGetResponse, error) {
	result := AzureComServiceBusQueueGetResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.AzureServiceBusResource); err != nil {
		return AzureComServiceBusQueueGetResponse{}, err
	}
	return result, nil
}

// getHandleError handles the Get error response.
func (client *AzureComServiceBusQueueClient) getHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
		errType := ErrorResponse{raw: string(body)}
	if err := runtime.UnmarshalAsJSON(resp, &errType); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}

// List - List the azure.com.ServiceBusQueue resources deployed in the application.
// If the operation fails it returns the *ErrorResponse error type.
func (client *AzureComServiceBusQueueClient) List(ctx context.Context, resourceGroupName string, applicationName string, options *AzureComServiceBusQueueListOptions) (AzureComServiceBusQueueListResponse, error) {
	req, err := client.listCreateRequest(ctx, resourceGroupName, applicationName, options)
	if err != nil {
		return AzureComServiceBusQueueListResponse{}, err
	}
	resp, err := 	client.pl.Do(req)
	if err != nil {
		return AzureComServiceBusQueueListResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return AzureComServiceBusQueueListResponse{}, client.listHandleError(resp)
	}
	return client.listHandleResponse(resp)
}

// listCreateRequest creates the List request.
func (client *AzureComServiceBusQueueClient) listCreateRequest(ctx context.Context, resourceGroupName string, applicationName string, options *AzureComServiceBusQueueListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.CustomProviders/resourceProviders/radiusv3/Application/{applicationName}/azure.com.ServiceBusQueue"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if applicationName == "" {
		return nil, errors.New("parameter applicationName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{applicationName}", url.PathEscape(applicationName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(	client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2018-09-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// listHandleResponse handles the List response.
func (client *AzureComServiceBusQueueClient) listHandleResponse(resp *http.Response) (AzureComServiceBusQueueListResponse, error) {
	result := AzureComServiceBusQueueListResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.AzureServiceBusList); err != nil {
		return AzureComServiceBusQueueListResponse{}, err
	}
	return result, nil
}

// listHandleError handles the List error response.
func (client *AzureComServiceBusQueueClient) listHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
		errType := ErrorResponse{raw: string(body)}
	if err := runtime.UnmarshalAsJSON(resp, &errType); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}
