package app

import (
	"io"
	"net/http"
	"testing"

	"github.com/Masterminds/semver"

	"github.com/jsnfwlr/keyper-cli/internal/app/filters"
)

func GetCurrentVersion() string {
	return buildVersion
}

type versionCache struct {
	version         string
	updateAvailable bool
}

var cacheVersion *versionCache

func GetNewVersionDetails() (updateAvailable bool, version string, fault error) {
	if cacheVersion != nil {
		return cacheVersion.updateAvailable, cacheVersion.version, nil
	}

	c := http.Client{Timeout: timeoutDuration}
	resp, err := c.Get(ReleaseAPI)
	if err != nil {
		return false, "", err
	}

	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, "", err
	}

	availableVersion, err := filters.JQ(b, ReleaseFilter)
	if err != nil {
		return false, "", err
	}

	av, err := semver.NewVersion(availableVersion)
	if err != nil {
		return false, "", err
	}
	bv, err := semver.NewVersion(buildVersion)
	if err != nil {
		return false, "", err
	}

	cacheVersion = &versionCache{
		version:         availableVersion,
		updateAvailable: av.GreaterThan(bv),
	}

	return av.GreaterThan(bv), availableVersion, nil
}

func Alpha(t *testing.T) {
	t.Helper()

	buildVersion = "v1.0.0-alpha.0"
}

func Beta(t *testing.T) {
	t.Helper()

	buildVersion = "v1.0.0-beta.0"
}

func ReleaseCandidate(t *testing.T) {
	t.Helper()

	buildVersion = "v1.0.0-rc.0"
}

func GeneralAvailability(t *testing.T) {
	t.Helper()

	buildVersion = "v1.0.0"
}
