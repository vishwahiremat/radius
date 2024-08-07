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

package radinit

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2_types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/radius-project/radius/pkg/cli/aws"
	"github.com/radius-project/radius/pkg/cli/output"
	"github.com/radius-project/radius/pkg/cli/prompt"
	"github.com/radius-project/radius/pkg/to"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func Test_enterAWSCloudProvider(t *testing.T) {
	ctrl := gomock.NewController(t)
	prompter := prompt.NewMockInterface(ctrl)
	client := aws.NewMockClient(ctrl)
	outputSink := output.MockOutput{}
	runner := Runner{Prompter: prompter, awsClient: client, Output: &outputSink}
	ec2Regions := []ec2_types.Region{
		{RegionName: to.Ptr("region")},
		{RegionName: to.Ptr("region2")},
	}
	regions := []string{"region", "region2"}

	setAWSAccessKeyIDPrompt(prompter, "access-key-id")
	setAWSSecretAccessKeyPrompt(prompter, "secret-access-key")
	setAWSCallerIdentity(client, &sts.GetCallerIdentityOutput{Account: to.Ptr("account-id")})
	setAWSAccountIDConfirmPrompt(prompter, "account-id", prompt.ConfirmYes)
	setAWSListRegions(client, &ec2.DescribeRegionsOutput{Regions: ec2Regions})
	setAWSRegionPrompt(prompter, regions, "region")

	provider, err := runner.enterAWSCloudProvider(context.Background())
	require.NoError(t, err)

	expected := &aws.Provider{
		AccessKeyID:     "access-key-id",
		SecretAccessKey: "secret-access-key",
		Region:          "region",
		AccountID:       "account-id",
	}
	require.Equal(t, expected, provider)
	require.Equal(t, []any{output.LogOutput{Format: awsAccessKeysCreateInstructionFmt}}, outputSink.Writes)
}
