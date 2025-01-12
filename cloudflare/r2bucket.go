package cloudflare

import (
	"github.com/syumai/workers/cloudflare/r2"
)

// R2Bucket represents interface of Cloudflare Worker's R2 Bucket instance.
// Deprecated: use r2.Bucket instead.
type R2Bucket = r2.Bucket

// NewR2Bucket returns R2Bucket for given variable name.
// Deprecated: use r2.NewBucket instead.
func NewR2Bucket(varName string) (*R2Bucket, error) {
	return r2.NewBucket(varName)
}

// R2PutOptions represents Cloudflare R2 put options.
// Deprecated: use r2.PutOptions instead.
type R2PutOptions = r2.PutOptions
