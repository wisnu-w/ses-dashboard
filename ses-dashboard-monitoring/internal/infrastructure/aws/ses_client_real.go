package aws

import (
	"context"
	"fmt"
	"time"

	"ses-monitoring/internal/domain/settings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

type SESClient struct {
	config    *settings.AWSConfig
	lastCall  time.Time
	minDelay  time.Duration
}

func NewSESClient(config *settings.AWSConfig) *SESClient {
	return &SESClient{
		config:   config,
		minDelay: 200 * time.Millisecond, // Minimum delay between API calls
	}
}

type SuppressionStatus struct {
	Email      string `json:"email"`
	Suppressed bool   `json:"suppressed"`
	Reason     string `json:"reason"`
	LastUpdate string `json:"last_update"`
}

// rateLimitedCall ensures minimum delay between AWS API calls
func (c *SESClient) rateLimitedCall() {
	since := time.Since(c.lastCall)
	if since < c.minDelay {
		time.Sleep(c.minDelay - since)
	}
	c.lastCall = time.Now()
}

// CheckSuppressionStatus checks if email is suppressed in AWS SES
func (c *SESClient) CheckSuppressionStatus(ctx context.Context, email string) (*SuppressionStatus, error) {
	if !c.config.Enabled {
		return nil, fmt.Errorf("AWS integration is disabled")
	}
	
	c.rateLimitedCall()
	
	cfg, err := c.getAWSConfig(ctx)
	if err != nil {
		return nil, err
	}
	
	sesClient := sesv2.NewFromConfig(cfg)
	
	result, err := sesClient.GetSuppressedDestination(ctx, &sesv2.GetSuppressedDestinationInput{
		EmailAddress: aws.String(email),
	})
	
	if err != nil {
		// If not found, email is not suppressed
		return &SuppressionStatus{
			Email:      email,
			Suppressed: false,
			Reason:     "Not suppressed",
			LastUpdate: "",
		}, nil
	}
	
	return &SuppressionStatus{
		Email:      email,
		Suppressed: true,
		Reason:     string(result.SuppressedDestination.Reason),
		LastUpdate: result.SuppressedDestination.LastUpdateTime.Format("2006-01-02T15:04:05Z"),
	}, nil
}

// RemoveFromSuppression removes email from AWS SES suppression list
func (c *SESClient) RemoveFromSuppression(ctx context.Context, email string) error {
	if !c.config.Enabled {
		return fmt.Errorf("AWS integration is disabled")
	}
	
	c.rateLimitedCall()
	
	cfg, err := c.getAWSConfig(ctx)
	if err != nil {
		return err
	}
	
	sesClient := sesv2.NewFromConfig(cfg)
	
	_, err = sesClient.DeleteSuppressedDestination(ctx, &sesv2.DeleteSuppressedDestinationInput{
		EmailAddress: aws.String(email),
	})
	
	return err
}

// AddToSuppression adds email to AWS SES suppression list
func (c *SESClient) AddToSuppression(ctx context.Context, email, reason string) error {
	if !c.config.Enabled {
		return fmt.Errorf("AWS integration is disabled")
	}
	
	c.rateLimitedCall()
	
	cfg, err := c.getAWSConfig(ctx)
	if err != nil {
		return err
	}
	
	sesClient := sesv2.NewFromConfig(cfg)
	
	_, err = sesClient.PutSuppressedDestination(ctx, &sesv2.PutSuppressedDestinationInput{
		EmailAddress: aws.String(email),
		Reason:       types.SuppressionListReasonComplaint, // Default to complaint
	})
	
	return err
}

// TestConnection tests AWS SES connection
func (c *SESClient) TestConnection(ctx context.Context) error {
	if !c.config.Enabled {
		return fmt.Errorf("AWS integration is disabled")
	}
	
	if c.config.AccessKey == "" || c.config.SecretKey == "" {
		return fmt.Errorf("AWS credentials are required")
	}
	
	c.rateLimitedCall()
	
	cfg, err := c.getAWSConfig(ctx)
	if err != nil {
		return err
	}
	
	sesClient := sesv2.NewFromConfig(cfg)
	
	// Test connection by getting account sending enabled status
	_, err = sesClient.GetAccount(ctx, &sesv2.GetAccountInput{})
	return err
}

// GetSuppressionList gets all suppressed emails from AWS SES
func (c *SESClient) GetSuppressionList(ctx context.Context) ([]*SuppressionStatus, error) {
	if !c.config.Enabled {
		return nil, fmt.Errorf("AWS integration is disabled")
	}

	if c.config.AccessKey == "" || c.config.SecretKey == "" {
		return nil, fmt.Errorf("AWS credentials not configured")
	}
	
	c.rateLimitedCall()
	
	// Create context with longer timeout
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	
	cfg, err := c.getAWSConfig(ctxWithTimeout)
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS config: %w", err)
	}
	
	sesClient := sesv2.NewFromConfig(cfg)
	
	var allSuppressions []*SuppressionStatus
	var nextToken *string
	
	for {
		result, err := sesClient.ListSuppressedDestinations(ctxWithTimeout, &sesv2.ListSuppressedDestinationsInput{
			NextToken: nextToken,
			PageSize:  aws.Int32(100), // Increase page size
		})
		
		if err != nil {
			return nil, fmt.Errorf("failed to list suppressed destinations: %w", err)
		}
		
		for _, dest := range result.SuppressedDestinationSummaries {
			allSuppressions = append(allSuppressions, &SuppressionStatus{
				Email:      *dest.EmailAddress,
				Suppressed: true,
				Reason:     string(dest.Reason),
				LastUpdate: dest.LastUpdateTime.Format("2006-01-02T15:04:05Z"),
			})
		}
		
		nextToken = result.NextToken
		if nextToken == nil {
			break
		}
		
		// Rate limiting - increase delay for safety
		time.Sleep(1 * time.Second)
	}
	
	return allSuppressions, nil
}

// getAWSConfig creates AWS config with credentials and retry configuration
func (c *SESClient) getAWSConfig(ctx context.Context) (aws.Config, error) {
	return config.LoadDefaultConfig(ctx,
		config.WithRegion(c.config.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			c.config.AccessKey,
			c.config.SecretKey,
			"",
		)),
		config.WithRetryer(func() aws.Retryer {
			return retry.AddWithMaxAttempts(retry.NewStandard(), 5)
		}),
	)
}