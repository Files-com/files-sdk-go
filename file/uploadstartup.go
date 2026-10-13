package file

import "fmt"

// UploadStartup selects how an adaptive (V2) upload starts its part requests to S3: UploaderParams.AdaptiveUploadStartup
// for a job, UploadWithStartup for one upload. Either way the upload adapts throughout, within every cap, and the
// files of a job share what it learns.
//
// An explicit start in UploadV2Tuning (InitialTarget or S3InitialTarget) takes precedence: the upload starts there,
// as it always has, whichever startup is selected. Uploads to other destinations, and downloads, start as before.
type UploadStartup string

const (
	// UploadStartupDefault starts with one part request and adds requests only as their measured throughput earns
	// them without queueing (lib.AdaptiveConcurrencyConfig.WarmStart). A measured fast link (96 MiB/s) goes straight
	// to the policy's start of 150 requests, or the cap if lower; so does a probe reaching it. From there it adapts
	// as UploadStartupAggressive does and can reach the same throughput and targets; probing delays startup,
	// not the concurrency or performance ceiling. The zero value selects it.
	UploadStartupDefault UploadStartup = "default"
	// UploadStartupAggressive starts at the policy's start, 150 part requests, or the cap if lower, for a link known
	// to carry them, skipping the initial concurrency probes on known high-throughput networks.
	// It is not a fixed concurrency: failures, throttling and probes that do not pay lower it.
	UploadStartupAggressive UploadStartup = "aggressive"
)

func (s UploadStartup) validate() error {
	switch s {
	case "", UploadStartupDefault, UploadStartupAggressive:
		return nil
	default:
		return fmt.Errorf("invalid upload startup %q: expected %q or %q", string(s), UploadStartupDefault, UploadStartupAggressive)
	}
}

// UploadWithStartup selects how this upload starts its part requests when it uses upload V2 (UploadStartup).
func UploadWithStartup(startup UploadStartup) UploadOption {
	return func(params uploadIO) (uploadIO, error) {
		if err := startup.validate(); err != nil {
			return params, err
		}
		params.uploadV2Startup = startup
		return params, nil
	}
}

// uploadV2WarmStarts reports whether an upload's adaptive manager warm-starts: an S3 upload under the default startup
// whose tuning does not set its start.
func uploadV2WarmStarts(plan uploadV2PartPlan, tuning UploadV2Tuning, startup UploadStartup) bool {
	return plan.target == uploadV2TargetS3 &&
		startup != UploadStartupAggressive &&
		tuning.InitialTarget == 0 &&
		tuning.S3InitialTarget == 0
}
