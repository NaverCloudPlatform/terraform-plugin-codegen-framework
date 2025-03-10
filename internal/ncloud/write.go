package ncloud

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/NaverCloudPlatform/terraform-plugin-codegen-framework/internal/util"
	"github.com/NaverCloudPlatform/terraform-plugin-codegen-spec/spec"
)

func WriteNcloudResources(resourcesSchema map[string][]byte, spec util.NcloudSpecification, outputDir, packageName string, genRefresh bool) error {
	for k, v := range resourcesSchema {
		dirName := ""

		if packageName == "" {
			dirName = k

			err := os.MkdirAll(filepath.Join(outputDir, dirName), os.ModePerm)
			if err != nil {
				return err
			}
		}

		filename := fmt.Sprintf("%s.go", k)

		n := NewResource(spec, k, packageName)

		f, err := os.Create(filepath.Join(outputDir, dirName, filename))
		if err != nil {
			return err
		}

		_, err = f.Write(v)
		if err != nil {
			return err
		}

		_, err = f.Write(n.RenderInitial())
		if err != nil {
			return err
		}

		_, err = f.Write(n.RenderImportState())
		if err != nil {
			return err
		}

		_, err = f.Write(n.RenderCreate())
		if err != nil {
			return err
		}

		_, err = f.Write(n.RenderRead())
		if err != nil {
			return err
		}

		_, err = f.Write(n.RenderUpdate())
		if err != nil {
			return err
		}

		_, err = f.Write(n.RenderDelete())
		if err != nil {
			return err
		}

		_, err = f.Write(n.RenderModel())
		if err != nil {
			return err
		}

		filePath := f.Name()

		util.RemoveDuplicates(filePath)
		util.RemoveCustomType(filePath)
	}

	return nil
}

// WriteDataSources uses the packageName to determine whether to create a directory and package per data source.
// If packageName is an empty string, this indicates that the flag was not set, and the default behaviour is
// then to create a package and directory per data source. If packageName is set then all generated code is
// placed into the same directory and package.
func WriteNcloudDataSources(dataSourcesSchema map[string][]byte, spec util.NcloudSpecification, outputDir, packageName string) error {
	for k, v := range dataSourcesSchema {
		dirName := ""

		if packageName == "" {
			dirName = k

			err := os.MkdirAll(filepath.Join(outputDir, dirName), os.ModePerm)
			if err != nil {
				return err
			}
		}

		filename := fmt.Sprintf("%s_data_source.go", k)

		n := NewDataSources(&spec, k, packageName)

		f, err := os.Create(filepath.Join(outputDir, dirName, filename))
		if err != nil {
			return err
		}

		_, err = f.Write(v)
		if err != nil {
			return err
		}

		// --- NCLOUD Logic ---
		_, err = f.Write(n.RenderInitial())
		if err != nil {
			return err
		}

		_, err = f.Write(n.RenderRead())
		if err != nil {
			return err
		}

		_, err = f.Write(n.RenderModel())
		if err != nil {
			return err
		}

		filePath := f.Name()

		util.RemoveDuplicates(filePath)
		util.RemoveCustomType(filePath)
	}

	return nil
}

// WriteDataSources uses the packageName to determine whether to create a directory and package per data source.
// If packageName is an empty string, this indicates that the flag was not set, and the default behaviour is
// then to create a package and directory per data source. If packageName is set then all generated code is
// placed into the same directory and package.
func WriteNcloudDataSourceTests(dataSourcesSchema map[string][]byte, spec util.NcloudSpecification, outputDir, packageName string) error {
	for k := range dataSourcesSchema {
		dirName := ""

		if packageName == "" {
			dirName = fmt.Sprintf("%s_data_source", k)

			err := os.MkdirAll(filepath.Join(outputDir, dirName), os.ModePerm)
			if err != nil {
				return err
			}
		}

		filename := fmt.Sprintf("%s_data_source_test.go", k)

		n := NewDataSources(&spec, k, packageName)

		f, err := os.Create(filepath.Join(outputDir, dirName, filename))
		if err != nil {
			return err
		}

		_, err = f.Write(n.RenderTest())
		if err != nil {
			return err
		}

		filePath := f.Name()

		util.RemoveDuplicates(filePath)
		util.RemoveCustomType(filePath)
	}

	return nil
}

// WriteResources uses the packageName to determine whether to create a directory and package per resource.
// If packageName is an empty string, this indicates that the flag was not set, and the default behaviour is
// then to create a package and directory per resource. If packageName is set then all generated code is
// placed into the same directory and package.
func WriteNcloudResourceTests(resourcesSchema map[string][]byte, spec util.NcloudSpecification, outputDir, packageName string) error {
	for k := range resourcesSchema {
		dirName := ""

		if packageName == "" {
			dirName = k

			err := os.MkdirAll(filepath.Join(outputDir, dirName), os.ModePerm)
			if err != nil {
				return err
			}
		}

		filename := fmt.Sprintf("%s_test.go", k)

		n := NewResource(spec, k, packageName)

		f, err := os.Create(filepath.Join(outputDir, dirName, filename))
		if err != nil {
			return err
		}

		_, err = f.Write(n.RenderTest())
		if err != nil {
			return err
		}

		filePath := f.Name()

		util.RemoveDuplicates(filePath)
		util.RemoveCustomType(filePath)
	}

	return nil
}

func WriteNcloudResourceRefresh(resourcesSchema map[string][]byte, spec util.NcloudSpecification, outputDir, packageName string) error {
	for k := range resourcesSchema {
		dirName := ""

		if packageName == "" {
			dirName = k

			err := os.MkdirAll(filepath.Join(outputDir, dirName), os.ModePerm)
			if err != nil {
				return err
			}
		}

		filename := fmt.Sprintf("%s_refresh.go", k)

		n := NewResource(spec, k, packageName)

		f, err := os.Create(filepath.Join(outputDir, dirName, filename))
		if err != nil {
			return err
		}

		_, err = f.Write(n.RenderRefresh())
		if err != nil {
			return err
		}

		_, err = f.Write(n.RenderWait())
		if err != nil {
			return err
		}

		filePath := f.Name()

		util.RemoveDuplicates(filePath)
		util.RemoveCustomType(filePath)
	}

	return nil
}

func WriteNcloudDataSourceRefresh(resourcesSchema map[string][]byte, spec util.NcloudSpecification, outputDir, packageName string) error {
	for k := range resourcesSchema {
		dirName := ""

		if packageName == "" {
			dirName = k

			err := os.MkdirAll(filepath.Join(outputDir, dirName), os.ModePerm)
			if err != nil {
				return err
			}
		}

		filename := fmt.Sprintf("%s_refresh.go", k)

		n := NewDataSources(&spec, k, packageName)

		f, err := os.Create(filepath.Join(outputDir, dirName, filename))
		if err != nil {
			return err
		}

		_, err = f.Write(n.RenderRefresh())
		if err != nil {
			return err
		}

		// TODO - Implement RenderWait() method
		// _, err = f.Write(n.RenderWait())
		// if err != nil {
		// 	return err
		// }

		filePath := f.Name()

		util.RemoveDuplicates(filePath)
	}

	return nil
}

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
