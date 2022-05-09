package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
)

var ZoneId = os.Getenv("ROUTE53_ZONE_ID")
var RecordTTL = Must(time.ParseDuration(os.Getenv("ROUTE53_RECORD_TTL")))
var r53Config = Must(config.LoadDefaultConfig(context.Background(), config.WithRegion(os.Getenv("AWS_REGION"))))
var r53 = route53.NewFromConfig(r53Config)

func StoreToken(ctx context.Context, base, t, url string) error {
	fqdn := t + "." + base
	qurl := fmt.Sprintf("%q", url)
	ttlsec := int64(RecordTTL.Seconds())
	if ttlsec < 60 {
		// assume any value less than 1 minute is invalid, default to 5 minutes
		ttlsec = 300
	}
	_, err := r53.ChangeResourceRecordSets(ctx, &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &types.ChangeBatch{
			Changes: []types.Change{{
				Action: types.ChangeActionCreate,
				ResourceRecordSet: &types.ResourceRecordSet{
					Name:            &fqdn,
					Type:            types.RRTypeTxt,
					ResourceRecords: []types.ResourceRecord{{Value: &qurl}},
					TTL:             &ttlsec,
				},
			}},
		},
		HostedZoneId: &ZoneId,
	})
	if err != nil {
		return fmt.Errorf("Cannot create TXT record for %q: %w", fqdn, err)
	}
	return nil
}
