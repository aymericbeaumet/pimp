package semver

func Patch(input interface{}) (*Version, error) {
	v, err := Latest(input)
	if err != nil {
		return nil, err
	}
	if err := v.version.IncrementPatch(); err != nil {
		return nil, err
	}
	return v, nil
}
