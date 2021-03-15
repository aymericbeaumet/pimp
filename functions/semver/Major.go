package semver

import "github.com/blang/semver/v4"

func Major(input interface{}) (*semver.Version, error) {
	v, err := getLatestVersion(input)
	if err != nil {
		return nil, err
	}
	if err := v.IncrementMajor(); err != nil {
		return nil, err
	}
	return v, nil
}
