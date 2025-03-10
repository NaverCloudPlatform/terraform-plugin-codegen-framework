// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cmd

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/hashicorp/cli"

	"github.com/NaverCloudPlatform/terraform-plugin-codegen-framework/internal/format"
	"github.com/NaverCloudPlatform/terraform-plugin-codegen-framework/internal/input"
	"github.com/NaverCloudPlatform/terraform-plugin-codegen-framework/internal/logging"
	"github.com/NaverCloudPlatform/terraform-plugin-codegen-framework/internal/ncloud"
	ncloud_resource "github.com/NaverCloudPlatform/terraform-plugin-codegen-framework/internal/ncloud/resource"
	"github.com/NaverCloudPlatform/terraform-plugin-codegen-framework/internal/schema"
	"github.com/NaverCloudPlatform/terraform-plugin-codegen-framework/internal/util"
	"github.com/NaverCloudPlatform/terraform-plugin-codegen-framework/internal/validate"
)

type GenerateResourcesCommand struct {
	UI              cli.Ui
	flagIRInputPath string
	flagOutputPath  string
	flagPackageName string
	flagGenRefresh  bool
}

func (cmd *GenerateResourcesCommand) Flags() *flag.FlagSet {
	fs := flag.NewFlagSet("generate resources", flag.ExitOnError)
	fs.StringVar(&cmd.flagIRInputPath, "input", "./ir.json", "path to intermediate representation (JSON)")
	fs.StringVar(&cmd.flagOutputPath, "output", "./output", "directory path to output generated code files")
	fs.StringVar(&cmd.flagPackageName, "package", "", "name of Go package for generated code files")
	fs.BoolVar(&cmd.flagGenRefresh, "gen_refresh", false, "whether render new refresh files or not")

	return fs
}

func (cmd *GenerateResourcesCommand) Help() string {
	strBuilder := &strings.Builder{}

	longestName := 0
	longestUsage := 0
	cmd.Flags().VisitAll(func(f *flag.Flag) {
		if len(f.Name) > longestName {
			longestName = len(f.Name)
		}
		if len(f.Usage) > longestUsage {
			longestUsage = len(f.Usage)
		}
	})

	strBuilder.WriteString("\nUsage: tfplugingen-framework generate resources [<args>]\n\n")
	cmd.Flags().VisitAll(func(f *flag.Flag) {
		if f.DefValue != "" {
			strBuilder.WriteString(fmt.Sprintf("    --%s <ARG> %s%s%s  (default: %q)\n",
				f.Name,
				strings.Repeat(" ", longestName-len(f.Name)+2),
				f.Usage,
				strings.Repeat(" ", longestUsage-len(f.Usage)+2),
				f.DefValue,
			))
		} else {
			strBuilder.WriteString(fmt.Sprintf("    --%s <ARG> %s%s%s\n",
				f.Name,
				strings.Repeat(" ", longestName-len(f.Name)+2),
				f.Usage,
				strings.Repeat(" ", longestUsage-len(f.Usage)+2),
			))
		}
	})
	strBuilder.WriteString("\n")

	return strBuilder.String()
}

func (a *GenerateResourcesCommand) Synopsis() string {
	return "Generate code for resources from an Intermediate Representation (IR) JSON file."
}

func (cmd *GenerateResourcesCommand) Run(args []string) int {
	ctx := context.Background()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelWarn,
	}))

	fs := cmd.Flags()
	err := fs.Parse(args)
	if err != nil {
		cmd.UI.Error(fmt.Sprintf("error parsing command flags: %s", err))
		return 1
	}

	err = cmd.runInternal(ctx, logger)
	if err != nil {
		cmd.UI.Error(fmt.Sprintf("Error executing command: %s\n", err))
		return 1
	}

	return 0
}

func (cmd *GenerateResourcesCommand) runInternal(ctx context.Context, logger *slog.Logger) error {
	// read input file
	src, err := input.Read(cmd.flagIRInputPath)
	if err != nil {
		return fmt.Errorf("error reading IR JSON: %w", err)
	}

	// validate JSON
	err = validate.JSON(src)
	if err != nil {
		return fmt.Errorf("error validating IR JSON: %w", err)
	}

	// parse and validate IR against specification
	spec, err := ncloud.NcloudParse(ctx, src)
	if err != nil {
		return fmt.Errorf("error parsing IR JSON: %w", err)
	}

	err = generateResourceCode(ctx, spec, cmd.flagOutputPath, cmd.flagPackageName, "Resource", cmd.flagGenRefresh, logger)
	if err != nil {
		return fmt.Errorf("error generating resource code: %w", err)
	}

	return nil
}

func generateResourceCode(ctx context.Context, spec util.NcloudSpecification, outputPath, packageName, generatorType string, genRefresh bool, logger *slog.Logger) error {
	ctx = logging.SetPathInContext(ctx, "resource")

	// convert IR to framework schema
	s, err := ncloud_resource.NewSchemas(spec)
	if err != nil {
		return fmt.Errorf("error converting IR to Plugin Framework schema: %w", err)
	}

	// convert framework schema to []byte
	g := schema.NewGeneratorSchemas(s)
	schemas, err := g.Schemas(packageName, generatorType)
	if err != nil {
		return fmt.Errorf("error converting Plugin Framework schema to Go code: %w", err)
	}

	// format schema code
	formattedSchemas, err := format.Format(schemas)
	if err != nil {
		return fmt.Errorf("error formatting Go code: %w", err)
	}

	// --- NCLOUD Logic ---

	// write code
	err = ncloud.WriteNcloudResources(formattedSchemas, spec, outputPath, packageName, genRefresh)
	if err != nil {
		return fmt.Errorf("error writing Go code to output: %w", err)
	}

	err = ncloud.WriteNcloudResourceTests(formattedSchemas, spec, outputPath, packageName)
	if err != nil {
		return fmt.Errorf("error writing Go code to output: %w", err)
	}

	// Render refresh file conditionally
	if genRefresh {
		err = ncloud.WriteNcloudResourceRefresh(formattedSchemas, spec, outputPath, packageName)
		if err != nil {
			return fmt.Errorf("error writing Go code to output: %w", err)
		}
	}

	return nil
}
