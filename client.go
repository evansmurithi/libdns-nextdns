package libdnsnextdns

import (
	"context"
	"fmt"

	"github.com/amalucelli/nextdns-go/nextdns"
	"github.com/libdns/libdns"
)

func (p *Provider) getRewrites(ctx context.Context, profileID string) ([]libdns.Record, error) {
	rewrites, err := p.client.Rewrites.List(ctx, &nextdns.ListRewritesRequest{
		ProfileID: profileID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch NextDNS rewrites, %w", err)
	}

	records := []libdns.Record{}

	for _, rewrite := range rewrites {
		record := libdns.Record{
			ID:    rewrite.ID,
			Name:  rewrite.Name,
			Value: rewrite.Content,
			Type:  rewrite.Type,
		}

		records = append(records, record)
	}

	return records, nil
}

func (p *Provider) createRewrite(ctx context.Context, profileID string, record libdns.Record) (libdns.Record, error) {
	rewriteID, err := p.client.Rewrites.Create(ctx, &nextdns.CreateRewritesRequest{
		ProfileID: profileID,
		Rewrites: &nextdns.Rewrites{
			Name:    record.Name,
			Content: record.Value,
		},
	})
	if err != nil {
		return record, fmt.Errorf("failed to create NextDNS rewrite, %w", err)
	}

	record.ID = rewriteID
	return record, nil
}

func (p *Provider) deleteRewrite(ctx context.Context, profileID string, record libdns.Record) (libdns.Record, error) {
	err := p.client.Rewrites.Delete(ctx, &nextdns.DeleteRewritesRequest{
		ProfileID: profileID,
		ID:        record.ID,
	})
	if err != nil {
		return record, fmt.Errorf("failed to delete NextDNS rewrite, %w", err)
	}

	return record, nil
}

func (p *Provider) updateRewrite(ctx context.Context, profileID string, record libdns.Record) (libdns.Record, error) {
	_, err := p.deleteRewrite(ctx, profileID, record)
	if err != nil {
		return record, err
	}

	return p.createRewrite(ctx, profileID, record)
}
