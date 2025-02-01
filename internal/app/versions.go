package app

import (
	"bytes"
	"io"
	"net/http"
	"testing"
	"text/template"

	"github.com/Masterminds/semver"

	"github.com/jsnfwlr/keyper-cli/internal/app/filters"
)

func GetCurrentVersion() string {
	return buildVersion
}

func GetBuildBranch() string {
	return buildBranch
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

	// TODO @jsnfwlr: once buildBranch is implemented in the github actions,
	// this function will need to take it in to account when checking for updates
	// eg: a new version for the main branch may have been released, but if the
	// buildBranch of the installed app is nightly, the stable release may not be
	// newer than the installed version
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

func Nightly(t *testing.T) {
	t.Helper()

	buildVersion = "v1.0.0-nightly.0"
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

func OutputVersion() string {
	tmpl := VersionTemplate()

	funcMap := template.FuncMap{
		"latestVersion": func() string {
			_, ver, err := GetNewVersionDetails()
			if err != nil {
				return "unknown"
			}

			return ver
		},
		"releaseURL": func() string {
			return ReleaseURL
		},
		"canUpdate": func() bool {
			ua, _, _ := GetNewVersionDetails()
			return ua
		},
	}

	out := bytes.Buffer{}

	appInfo := struct {
		AppName string
		Short   string
		Version string
	}{
		AppName: AppName,
		Short:   RootShort,
		Version: buildVersion,
	}

	err := template.Must(template.New("version").Funcs(funcMap).Parse(tmpl)).Execute(&out, appInfo)
	if err != nil {
		panic(err)
	}

	return out.String()
}
