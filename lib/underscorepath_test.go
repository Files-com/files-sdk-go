package lib

import "testing"

func TestNormalizeAPIPathPreservesPathIdentity(t *testing.T) {
	got := NormalizeAPIPath("/../../remote\\path//./to/file.txt")
	want := "remote/path/to/file.txt"
	if got != want {
		t.Fatalf("NormalizeAPIPath() = %q, want %q", got, want)
	}

	got = NormalizeAPIPath("remote/../path/to/file.txt")
	want = "remote/path/to/file.txt"
	if got != want {
		t.Fatalf("NormalizeAPIPath() = %q, want %q", got, want)
	}
}

func TestUnderscoreDestinationPath(t *testing.T) {
	got := UnderscoreDestinationPath("RemoteServers", 42, "/../../remote\\path//./to/file.txt")
	want := "_/RemoteServers/42/remote/path/to/file.txt"
	if got != want {
		t.Fatalf("UnderscoreDestinationPath() = %q, want %q", got, want)
	}
}
