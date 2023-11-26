package libdnsnextdns

import (
	"context"
	"fmt"
	"time"

	"github.com/amalucelli/nextdns-go/nextdns"
	"github.com/libdns/libdns"
)

type Opt struct {
	ApiKey             string        `json:"api_key,omitempty"`
	MaxRetries         int           `json:"max_retries,omitempty"`
	MaxWaitDur         time.Duration `json:"max_wait_dur,omitempty"`
	WaitForPropogation bool          `json:"wait_for_propogation,omitempty"`
}

// Provider implements the libdns interfaces for nextdns
type Provider struct {
	client *nextdns.Client
	opt    Opt
}

func NewProvider(ctx context.Context, opt Opt) (*Provider, error) {
	if opt.MaxRetries == 0 {
		opt.MaxRetries = 5
	}
	if opt.MaxWaitDur == 0 && opt.WaitForPropogation {
		opt.MaxWaitDur = time.Minute * 1
	}

	client, err := nextdns.New(nextdns.WithAPIKey(opt.ApiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate nextdns client, %w", err)
	}

	return &Provider{
		client: client,
		opt:    opt,
	}, nil
}

// GetRecords fetches all rewrites from NextDNS profile
func (p *Provider) GetRecords(ctx context.Context, profileID string) ([]libdns.Record, error) {
	return p.getRewrites(ctx, profileID)
}

// AppendRecords creates rewrites in the NextDNS profile
func (p *Provider) AppendRecords(ctx context.Context, profileID string, records []libdns.Record) ([]libdns.Record, error) {
	createdRecords := []libdns.Record{}

	for _, record := range records {
		newRecord, err := p.createRewrite(ctx, profileID, record)
		if err != nil {
			return nil, err
		}
		createdRecords = append(createdRecords, newRecord)
	}

	return createdRecords, nil
}

// DeleteRecords deletes rewrites from NextDNS profile
func (p *Provider) DeleteRecords(ctx context.Context, profileID string, records []libdns.Record) ([]libdns.Record, error) {
	deletedRecords := []libdns.Record{}

	for _, record := range records {
		deletedRecord, err := p.deleteRewrite(ctx, profileID, record)
		if err != nil {
			return nil, err
		}
		deletedRecords = append(deletedRecords, deletedRecord)
	}

	return deletedRecords, nil
}

// SetRecords updates rewrites in NextDNS profile
func (p *Provider) SetRecords(ctx context.Context, profileID string, records []libdns.Record) ([]libdns.Record, error) {
	updatedRecords := []libdns.Record{}

	for _, record := range records {
		updatedRecord, err := p.updateRewrite(ctx, profileID, record)
		if err != nil {
			return nil, err
		}
		updatedRecords = append(updatedRecords, updatedRecord)
	}

	return updatedRecords, nil
}

// Interface guards
var (
	_ libdns.RecordGetter   = (*Provider)(nil)
	_ libdns.RecordAppender = (*Provider)(nil)
	_ libdns.RecordSetter   = (*Provider)(nil)
	_ libdns.RecordDeleter  = (*Provider)(nil)
)
