package lib

import (
	"slices"
	"testing"
)

func TestSortVersions(t *testing.T) {
	got := sortVersions([]string{"", "  ", "not a version"})
	if len(got) != 0 {
		t.Errorf("sortVersions({\"\", \"  \", \"not a version\"}) = %v; want empty slice", got)
	}

	got = sortVersions([]string{"v1", "v0.3", "v1.2.5", "   "})
	if !slices.Equal(got, []string{"v0.3", "v1", "v1.2.5"}) {
		t.Errorf(
			"sortVersions({\"v1\", \"v0.3\", \"v1.2.5\", \"   \"}) = %v; want {\"v0.3\", \"v1\", \"v1.2.5\"}", got)
	}
}

func TestMaxVersion(t *testing.T) {
	got := maxVersion([]string{"", "  ", "not a version"})
	if got != "" {
		t.Errorf("maxVersion({\"\", \"  \", \"not a version\"}) = %v; want empty string", got)
	}

	got = maxVersion([]string{"v1", "v0.3", "v1.2.5", "   "})
	if got != "v1.2.5" {
		t.Errorf(
			"maxVersion({\"v1\", \"v0.3\", \"v1.2.5\", \"   \"}) = %v; want v1.2.5", got)
	}
}
