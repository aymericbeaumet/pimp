package semver

func Minor(input interface{}) (*Version, error) {
	v, err := Latest(input)
	if err != nil {
		return nil, err
	}
	if err := v.version.IncrementMinor(); err != nil {
		return nil, err
	}
	return v, nil
}
