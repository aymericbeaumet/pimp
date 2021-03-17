package semver

func SemverMajor(input interface{}) (*Version, error) {
	v, err := SemverLatest(input)
	if err != nil {
		return nil, err
	}
	if err := v.version.IncrementMajor(); err != nil {
		return nil, err
	}
	return v, nil
}
