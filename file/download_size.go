package file

import (
	"errors"
	"strconv"
	"strings"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
)

// contentRange is a parsed Content-Range header. start and end are -1 for an
// unsatisfied range ("bytes */total"); totalKnown is false for "*/".
type contentRange struct {
	start, end int64
	total      int64
	totalKnown bool
}

// parseContentRange parses "bytes <start>-<end>/<total>". The unit is optional
// because some servers omit it, but any other unit is not a byte range. The
// total may be "*".
func parseContentRange(value string) (contentRange, bool) {
	value = strings.TrimSpace(value)
	if unit, rest, found := strings.Cut(value, " "); found {
		if !strings.EqualFold(unit, "bytes") {
			return contentRange{}, false
		}
		value = strings.TrimSpace(rest)
	}
	rangeSpec, totalSpec, found := strings.Cut(value, "/")
	if !found {
		return contentRange{}, false
	}
	var parsed contentRange
	if totalSpec != "*" {
		total, err := strconv.ParseInt(totalSpec, 10, 64)
		if err != nil || total < 0 {
			return contentRange{}, false
		}
		parsed.total, parsed.totalKnown = total, true
	}
	if rangeSpec == "*" {
		parsed.start, parsed.end = -1, -1
		return parsed, true
	}
	startSpec, endSpec, found := strings.Cut(rangeSpec, "-")
	if !found {
		return contentRange{}, false
	}
	start, startErr := strconv.ParseInt(startSpec, 10, 64)
	end, endErr := strconv.ParseInt(endSpec, 10, 64)
	if startErr != nil || endErr != nil || start < 0 || end < start {
		return contentRange{}, false
	}
	parsed.start, parsed.end = start, end
	return parsed, true
}

// downloadSourceChanged reports the typed download_source_changed contract: the
// download URL was issued for a source that has since changed size, so no range
// of that URL can be retried. Other 409 responses keep their ordinary handling.
func downloadSourceChanged(err error) bool {
	return errors.Is(err, files_sdk.ErrDownloadSourceChanged)
}
