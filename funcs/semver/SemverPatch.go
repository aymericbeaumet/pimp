package semver

func SemverPatch(input interface{}) (*Version, error) {
	v, err := getLatestVersion(input)
	if err != nil {
		return nil, err
	}
	if err := v.IncrementPatch(); err != nil {
		return nil, err
	}
	return v, nil
}
