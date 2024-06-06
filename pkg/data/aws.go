package data

import (
	"context"
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/gopaytech/unused-exporter/pkg/model"
	"github.com/gopaytech/unused-exporter/pkg/settings"
)

var (
	ErrMissingIAMKeySecret = errors.New("aws IAM key and/or secret is missing")
	ErrMissingRegion       = errors.New("aws region is missing")
)

type AWSData struct {
	assumeRolesRegions []string
	assumeRoleDuration int
	stsClient          *sts.Client
}

func (g *AWSData) GetUnusedIP() ([]model.IPAddress, error) {
	var IPs []model.IPAddress

	for _, assumeRoleRegion := range g.assumeRolesRegions {
		assume := strings.Split(assumeRoleRegion, ",")
		role := assume[0]
		region := assume[1]

		assumeRoleOutput, err := g.stsClient.AssumeRole(context.TODO(), &sts.AssumeRoleInput{
			RoleArn:         &role,
			DurationSeconds: aws.Int32(int32(g.assumeRoleDuration)),
		})
		if err != nil {
			return nil, err
		}

		credentialProvider := credentials.NewStaticCredentialsProvider(
			*assumeRoleOutput.Credentials.AccessKeyId,
			*assumeRoleOutput.Credentials.SecretAccessKey,
			*assumeRoleOutput.Credentials.SessionToken,
		)
		cfg := aws.Config{
			Credentials: credentialProvider,
			Region:      region,
		}

		ec2Client := ec2.NewFromConfig(cfg)

		input := &ec2.DescribeAddressesInput{}
		output, err := ec2Client.DescribeAddresses(context.TODO(), input)
		if err != nil {
			return nil, err
		}

		for _, address := range output.Addresses {
			if address.AllocationId == nil && address.PublicIp != nil {
				IPs = append(IPs, model.IPAddress{
					Cloud:    "AWS",
					Region:   region,
					Value:    *address.PublicIp,
					Type:     "Public",
					Used:     false,
					Identity: "Elastic IP",
				})
			}
		}

	}
	return IPs, nil
}

func NewAWSData(settings settings.Settings) (*AWSData, error) {
	awsData := &AWSData{
		assumeRolesRegions: settings.AWSAssumeRolesRegions,
		assumeRoleDuration: settings.AWSAssumeRoleDuration,
	}

	if settings.AWSIAMKey == "" || settings.AWSIAMSecret == "" {
		return nil, ErrMissingIAMKeySecret
	}

	if settings.AWSRegion == "" {
		return nil, ErrMissingRegion
	}

	creddentialProvider := credentials.NewStaticCredentialsProvider(settings.AWSIAMKey, settings.AWSIAMSecret, "")
	stsClient := sts.New(sts.Options{
		Credentials: creddentialProvider,
		Region:      settings.AWSRegion,
	})

	awsData.stsClient = stsClient

	return awsData, nil
}
