package semver

func SemverMinor(input interface{}) (*Version, error) {
	v, err := getLatestVersion(input)
	if err != nil {
		return nil, err
	}
	if err := v.version.IncrementMinor(); err != nil {
		return nil, err
	}
	return v, nil
}
