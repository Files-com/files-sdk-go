package lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClassifyS3ErrorCode(t *testing.T) {
	tests := []struct {
		code      string
		class     string
		message   string
		retryable bool
	}{
		{"RequestTimeout", "request_timeout", "Upload to backend storage timed out.", true},
		{"InternalError", "service_unavailable", "Error returned by the remote service", true},
		{"TemporarilyUnavailable", "service_unavailable", "Error returned by the remote service", true},
		{"ServiceUnavailable", "service_unavailable", "Error returned by the remote service", true},
		{"SlowDown", "rate_limit", "Rate limited by remote service", true},
		{"RequestLimitExceeded", "rate_limit", "Rate limited by remote service", true},
		{"NoSuchKey", "not_found", "Resource not found on the remote service", false},
		{"NoSuchBucket", "connect_validation", "The specified bucket does not exist", false},
		{"InvalidAccessKeyId", "connect_validation", "The provided access key is invalid", false},
		{"InvalidRequest", "upload_error", "Invalid Request. The resource name or option is not supported by the remote service", false},
		{"AccessDenied", "permissions", "Access denied by the remote service", false},
		{"Forbidden", "permissions", "Access denied by the remote service", false},
		{"Unauthorized", "permissions", "Permission denied by the remote service", false},
		{"AuthorizationHeaderMalformed", "permissions", "Permission denied by the remote service", false},
		{"IllegalLocationConstraintException", "connect_validation", "The specified region does not match the region of the bucket", false},
		{"SignatureDoesNotMatch", "authentication", "Signature does not match", false},
		{"PermanentRedirect", "connect", "Permanent redirect: the bucket you are attempting to access must be addressed using the specified endpoint", false},
		{"PreconditionFailed", "if_match_condition_not_met", "The file has changed since comparison or no longer exists. This can happen if the file was deleted or modified by another process.", false},
		{"InvalidObjectState", "archived_object", "The object is in an archived storage class (such as Glacier) and must be restored before it can be accessed", false},
		{"InvalidKey", "invalid_path", "The specified path contains invalid characters", false},
		{"AccountProblem", "account_problem", "There is a problem with your Amazon Web Services account that prevents the action from completing successfully", false},
		{"PoorAccountStanding", "account_problem", "There is a problem with your Amazon Web Services account that prevents the action from completing successfully", false},
		{"SecurityPolicyViolated", "permissions", "Request violates VPC Service Controls", false},
		{"DatabaseTimeout", "database_timeout", "Remote service timeout", true},
		{"AccessDeniedException", "authentication", "Access denied when assuming role. Verify the trust policy allows this principal.", false},
		{"ExpiredTokenException", "authentication", "Session token has expired", false},
		{"SomethingWeNeverHeardOf", "unknown", "Error returned by the remote service", false},
	}

	for _, test := range tests {
		t.Run(test.code, func(t *testing.T) {
			classified := ClassifyS3ErrorCode(test.code)

			assert.Equal(t, test.class, classified.Class)
			assert.Equal(t, test.message, classified.Message)
			assert.Equal(t, test.retryable, classified.Retryable)
		})
	}
}

func TestClassifyS3Error(t *testing.T) {
	classified, ok := ClassifyS3Error(S3Error{Code: "RequestTimeout", Message: "raw provider message"})

	assert.True(t, ok)
	assert.Equal(t, "request_timeout", classified.Class)
	assert.True(t, classified.Retryable)
	assert.NotContains(t, classified.Message, "raw provider message")

	_, ok = ClassifyS3Error(assert.AnError)
	assert.False(t, ok)
}
