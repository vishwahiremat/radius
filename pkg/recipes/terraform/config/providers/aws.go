/*
Copyright 2023 The Radius Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package providers

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	ucp_datamodel "github.com/radius-project/radius/pkg/ucp/datamodel"

	"github.com/radius-project/radius/pkg/azure/tokencredentials"
	"github.com/radius-project/radius/pkg/components/secret"
	"github.com/radius-project/radius/pkg/components/secret/secretprovider"
	"github.com/radius-project/radius/pkg/corerp/datamodel"
	"github.com/radius-project/radius/pkg/recipes"
	"github.com/radius-project/radius/pkg/sdk"
	"github.com/radius-project/radius/pkg/ucp/credentials"
	"github.com/radius-project/radius/pkg/ucp/resources"
	resources_aws "github.com/radius-project/radius/pkg/ucp/resources/aws"
	"github.com/radius-project/radius/pkg/ucp/ucplog"
)

// Provider's config parameters need to match the values expected by Terraform
// https://registry.terraform.io/providers/hashicorp/aws/latest/docs
const (
	AWSProviderName = "aws"

	awsRegionParam    = "region"
	awsAccessKeyParam = "access_key"
	awsSecretKeyParam = "secret_key"

	// configs for AWS IRSA
	// Ref: https://registry.terraform.io/providers/hashicorp/aws/latest/docs#assuming-an-iam-role-using-a-web-identity
	awsIRSAProvider = "assume_role_with_web_identity"
	awsRoleARN      = "role_arn"
	sessionName     = "session_name"
	tokenFile       = "web_identity_token_file"
	// The path used in Amazon Elastic Kubernetes Service (EKS) to store the service account token for a Kubernetes pod.
	// Ref: https://docs.aws.amazon.com/eks/latest/userguide/pod-configuration.html
	tokenFilePath = "/var/run/secrets/eks.amazonaws.com/serviceaccount/token"
	sessionPrefix = "radius-terraform-"
)

var _ Provider = (*awsProvider)(nil)

type awsProvider struct {
	ucpConn        sdk.Connection
	secretProvider *secretprovider.SecretProvider
}

// NewAWSProvider creates a new AWSProvider instance.
func NewAWSProvider(ucpConn sdk.Connection, secretProvider *secretprovider.SecretProvider) Provider {
	return &awsProvider{ucpConn: ucpConn, secretProvider: secretProvider}
}

// BuildConfig generates the Terraform provider configuration for AWS provider. It checks if the AWS provider/scope
// is configured on the Environment and if so, parses the scope to get the region and returns a map containing the region.
// If the scope is invalid, an error is returned.
// https://registry.terraform.io/providers/hashicorp/aws/latest/docs
func (p *awsProvider) BuildConfig(ctx context.Context, envConfig *recipes.Configuration) (map[string]any, error) {
	region, err := p.parseScope(ctx, envConfig)
	if err != nil {
		return nil, err
	}

	credentialsProvider, err := p.getCredentialsProvider()
	if err != nil {
		return nil, err
	}

	credentials, err := fetchAWSCredentials(ctx, credentialsProvider)
	if err != nil {
		return nil, err
	}

	return p.generateProviderConfigMap(credentials, region), nil
}

// parseScope parses an AWS provider scope and returns the associated region
// Example scope: /planes/aws/aws/accounts/123456789/regions/us-east-1
func (p *awsProvider) parseScope(ctx context.Context, envConfig *recipes.Configuration) (string, error) {
	logger := ucplog.FromContextOrDiscard(ctx)
	if (envConfig == nil) || (envConfig.Providers == datamodel.Providers{}) || (envConfig.Providers.AWS == datamodel.ProvidersAWS{}) || envConfig.Providers.AWS.Scope == "" {
		logger.Info("AWS provider/scope is not configured on the Environment, skipping AWS region configuration.")
		return "", nil
	}

	scope := envConfig.Providers.AWS.Scope
	parsedScope, err := resources.Parse(scope)
	if err != nil {
		return "", fmt.Errorf("invalid AWS provider scope %q is configured on the Environment, error parsing: %s", scope, err.Error())
	}

	region := parsedScope.FindScope(resources_aws.ScopeRegions)
	if region == "" {
		return "", fmt.Errorf("invalid AWS provider scope %q is configured on the Environment, region is required in the scope", scope)
	}

	return region, nil
}

func (p *awsProvider) getCredentialsProvider() (*credentials.AWSCredentialProvider, error) {
	return credentials.NewAWSCredentialProvider(p.secretProvider, p.ucpConn, &tokencredentials.AnonymousCredential{})
}

// fetchAWSCredentials fetches AWS credentials from UCP. Returns nil if credentials not found error is received or the credentials are empty.
func fetchAWSCredentials(ctx context.Context, awsCredentialsProvider credentials.CredentialProvider[credentials.AWSCredential]) (*credentials.AWSCredential, error) {
	logger := ucplog.FromContextOrDiscard(ctx)
	credentials, err := awsCredentialsProvider.Fetch(ctx, credentials.AWSPublic, "default")
	if err != nil {
		if errors.Is(err, &secret.ErrNotFound{}) {
			logger.Info("AWS credentials are not registered, skipping credentials configuration.")
			return nil, nil
		}

		return nil, err
	}

	switch credentials.Kind {
	case ucp_datamodel.AWSAccessKeyCredentialKind:
		if credentials.AccessKeyCredential == nil || credentials.AccessKeyCredential.AccessKeyID == "" || credentials.AccessKeyCredential.SecretAccessKey == "" {
			logger.Info("AWS AccessKey credentials are not registered, skipping credentials configuration.")
			return nil, nil
		}
	case ucp_datamodel.AWSIRSACredentialKind:
		if credentials.IRSACredential == nil || credentials.IRSACredential.RoleARN == "" {
			logger.Info("AWS IRSA credentials are not registered, skipping credentials configuration.")
			return nil, nil
		}
	}

	return credentials, nil
}

func (p *awsProvider) generateProviderConfigMap(credentials *credentials.AWSCredential, region string) map[string]any {
	config := make(map[string]any)
	if region != "" {
		config[awsRegionParam] = region
	}

	if credentials != nil {
		switch credentials.Kind {
		case ucp_datamodel.AWSAccessKeyCredentialKind:
			if credentials.AccessKeyCredential != nil &&
				credentials.AccessKeyCredential.AccessKeyID != "" && credentials.AccessKeyCredential.SecretAccessKey != "" {
				config[awsAccessKeyParam] = credentials.AccessKeyCredential.AccessKeyID
				config[awsSecretKeyParam] = credentials.AccessKeyCredential.SecretAccessKey
			}

		case ucp_datamodel.AWSIRSACredentialKind:
			if credentials.IRSACredential != nil && credentials.IRSACredential.RoleARN != "" {
				config[awsIRSAProvider] = map[string]any{
					awsRoleARN:  credentials.IRSACredential.RoleARN,
					sessionName: sessionPrefix + uuid.New().String(),
					tokenFile:   tokenFilePath,
				}
			}
		}
	}

	return config
}
