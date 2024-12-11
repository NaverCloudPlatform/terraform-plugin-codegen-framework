// Parse returns a Specification from the JSON document contents, or any validation errors.
func NcloudParse(ctx context.Context, document []byte) (util.NcloudSpecification, error) {
	if err := spec.Validate(ctx, document); err != nil {
		return util.NcloudSpecification{}, err
	}

	var spec util.NcloudSpecification

	if err := json.Unmarshal(document, &spec); err != nil {
		return spec, err
	}

	if err := spec.Validate(ctx); err != nil {
		return spec, err
	}

	return spec, nil
}
