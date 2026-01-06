package aws

import (
	"context"
	"fmt"
	"log"
	"net/http"
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
		minDelay: 1 * time.Second, // Increased minimum delay between API calls
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
		log.Printf("AWS config error for %s: %v", email, err)
		return err
	}
	
	sesClient := sesv2.NewFromConfig(cfg)
	
	log.Printf("Attempting to remove %s from AWS SES suppression list", email)
	_, err = sesClient.DeleteSuppressedDestination(ctx, &sesv2.DeleteSuppressedDestinationInput{
		EmailAddress: aws.String(email),
	})
	
	if err != nil {
		log.Printf("AWS API error removing %s: %v", email, err)
		return fmt.Errorf("failed to remove %s from AWS SES: %w", email, err)
	}
	
	log.Printf("Successfully removed %s from AWS SES suppression list", email)
	return nil
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

// GetSuppressionList gets all suppressed emails from AWS SES using manual pagination
func (c *SESClient) GetSuppressionList(ctx context.Context) ([]*SuppressionStatus, error) {
	if !c.config.Enabled {
		return nil, fmt.Errorf("AWS integration is disabled")
	}

	if c.config.AccessKey == "" || c.config.SecretKey == "" {
		return nil, fmt.Errorf("AWS credentials not configured")
	}
	
	c.rateLimitedCall()
	
	// Create context with timeout for large suppression lists
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 20*time.Minute)
	defer cancel()
	
	cfg, err := c.getAWSConfig(ctxWithTimeout)
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS config: %w", err)
	}
	
	sesClient := sesv2.NewFromConfig(cfg)
	
	log.Printf("Starting to fetch AWS suppression list...")
	
	var allSuppressions []*SuppressionStatus
	nextToken := ""  // Start with empty string like Python
	pageCount := 0
	
	// Loop while nextToken is not nil (like Python: while next_token is not None)
	for nextToken != "END" {
		select {
		case <-ctxWithTimeout.Done():
			log.Printf("Operation timed out after processing %d pages (%d records)", pageCount, len(allSuppressions))
			return allSuppressions, nil // Return partial results
		default:
		}
		
		// Prepare input - include NextToken only if it's not empty
		input := &sesv2.ListSuppressedDestinationsInput{
			PageSize: aws.Int32(1000), // Maximum page size like Python
		}
		if nextToken != "" {
			input.NextToken = aws.String(nextToken)
		}
		
		result, err := sesClient.ListSuppressedDestinations(ctxWithTimeout, input)
		if err != nil {
			log.Printf("Failed to fetch page %d: %v", pageCount+1, err)
			return allSuppressions, fmt.Errorf("failed to get page %d: %w", pageCount+1, err)
		}
		
		// Process current batch
		for _, dest := range result.SuppressedDestinationSummaries {
			allSuppressions = append(allSuppressions, &SuppressionStatus{
				Email:      *dest.EmailAddress,
				Suppressed: true,
				Reason:     string(dest.Reason),
				LastUpdate: dest.LastUpdateTime.Format("2006-01-02T15:04:05Z"),
			})
		}
		
		pageCount++
		log.Printf("Fetched page %d: %d records (total: %d)", pageCount, len(result.SuppressedDestinationSummaries), len(allSuppressions))
		
		// Get next token - set to "END" if nil (like Python: next_token becomes None)
		if result.NextToken != nil {
			nextToken = *result.NextToken
		} else {
			nextToken = "END" // Signal to end loop
			log.Printf("Completed fetching all pages. Total: %d records", len(allSuppressions))
		}
		
		// Exponential backoff delay to avoid rate limiting
		baseDelay := time.Duration(pageCount) * 200 * time.Millisecond
		if baseDelay < 1*time.Second {
			baseDelay = 1 * time.Second // Minimum 1 second
		}
		if baseDelay > 5*time.Second {
			baseDelay = 5 * time.Second // Maximum 5 seconds
		}
		log.Printf("Waiting %v before next page to avoid rate limiting...", baseDelay)
		time.Sleep(baseDelay)
	}
	
	log.Printf("Successfully fetched all %d suppressed destinations from AWS", len(allSuppressions))
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
			return retry.AddWithMaxAttempts(
				retry.AddWithMaxBackoffDelay(
					retry.NewStandard(), 30*time.Second), 5)
		}),
		config.WithHTTPClient(&http.Client{
			Timeout: 30 * time.Second,
		}),
	)
}