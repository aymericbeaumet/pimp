package semver

func Major(input interface{}) (*Version, error) {
	v, err := Latest(input)
	if err != nil {
		return nil, err
	}
	if err := v.version.IncrementMajor(); err != nil {
		return nil, err
	}
	return v, nil
}
