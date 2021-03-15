package semver

func Minor(input interface{}) (*Version, error) {
	v, err := getLatestVersion(input)
	if err != nil {
		return nil, err
	}
	if err := v.IncrementMinor(); err != nil {
		return nil, err
	}
	return v, nil
}
