// Package version records versioning information about this module.

package version

import (
	"fmt"
)

// These constants determine the current version of this module.
// TODO: Prepare before release

const (
	Major      = 0
	Minor      = 0
	Patch      = 0
	Prerelease = ""
)

// String formats the version string for this module in semver format.
//
// Examples:
// 		0.1.0
//		1.4.1

func String() string {
	v := fmt.Sprintf("%d.%d.%d", Major, Minor, Patch)

	if Prerelease != "" {
		v += "-" + Prerelease
	}

	return v
}
