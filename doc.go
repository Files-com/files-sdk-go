// Package files_sdk provides Files.com API models, request parameters, configuration,
// pagination, and error handling. Resource-specific clients are in subpackages such
// as folder and user.
//
// Initialize a Config and embed it in a resource client to authenticate with an API
// key or session. Package-level resource functions use the default configuration,
// including the FILES_API_KEY environment variable.
//
// Listing methods return lazy iterators. Range over listing.All(), checking each
// error before using its resource. All walks the remaining results without
// restarting the listing. Next and Err remain available for existing callers.
// WithContext supplies cancellation and deadlines for requests. ResponseError
// supports errors.As for details and errors.Is for API error types and groups.
package files_sdk
