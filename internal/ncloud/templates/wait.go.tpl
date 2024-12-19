{{ define "Wait" }}
// Template for generating Terraform provider Waiting operation code
// Needed data is as follows.
// RefreshObjectName string
// ReadMethodName string
// Endpoint string
// ReadReqBody string
// ReadReqBody

func (plan *{{.RefreshObjectName | ToPascalCase}}Model) waitResourceCreated(ctx context.Context, id string) error {
	stateConf := &retry.StateChangeConf{
		Pending: []string{"CREATING"},
		Target:  []string{"CREATED"},
		Refresh: func() (interface{}, string, error) {
			c := ncloudsdk.NewClient("{{.Endpoint}}", os.Getenv("NCLOUD_ACCESS_KEY"), os.Getenv("NCLOUD_SECRET_KEY"))
			response, err := c.{{.ReadMethodName}}_TF(&ncloudsdk.{{.ReadMethodName}}Request{
					{{.ReadReqBody}}
			})
			if err != nil {
				return response, "CREATING", nil
			}
			if response != nil {
				return response, "CREATED", nil
			}

			return response, "CREATING", nil
		},
		Timeout:    conn.DefaultTimeout,
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error occured while waiting for resource to be created: %s", err)
	}
	return nil
}

func (plan *{{.RefreshObjectName | ToPascalCase}}Model) waitResourceDeleted(ctx context.Context, id string) error {
	stateConf := &retry.StateChangeConf{
		Pending: []string{"DELETING"},
		Target:  []string{"DELETED"},
		Refresh: func() (interface{}, string, error) {
			c := ncloudsdk.NewClient("{{.Endpoint}}", os.Getenv("NCLOUD_ACCESS_KEY"), os.Getenv("NCLOUD_SECRET_KEY"))
			response, err := c.{{.ReadMethodName}}_TF(&ncloudsdk.{{.ReadMethodName}}Request{
					{{.ReadReqBody}}
			})
			if err == nil {
				return response, "DELETED", nil
			}

			return response, "DELETING", nil
		},
		Timeout:    conn.DefaultTimeout,
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error occured while waiting for resource to be deleted: %s", err)
	}
	return nil
}

{{ end }}