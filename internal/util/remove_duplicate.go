package util

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

func RemoveDuplicates(filePath string) error {
	src, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	typeMap := make(map[string]bool)
	funcMap := make(map[string]map[string]bool)
	nonReceiverFuncMap := make(map[string]bool)

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "", src, parser.AllErrors)
	if err != nil {
		return fmt.Errorf("failed to parse file: %w", err)
	}

	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("package %s\n\n", node.Name.Name))

	if len(node.Imports) > 0 {
		builder.WriteString("import (\n")
		for _, imp := range node.Imports {
			builder.WriteString(fmt.Sprintf("\t%s\n", imp.Path.Value))
		}
		builder.WriteString(")\n\n")
	}

	ast.Inspect(node, func(n ast.Node) bool {
		switch decl := n.(type) {
		case *ast.GenDecl:
			if decl.Tok == token.TYPE {
				for _, spec := range decl.Specs {
					typeSpec := spec.(*ast.TypeSpec)
					typeName := typeSpec.Name.Name
					if typeMap[typeName] {
						continue
					}
					typeMap[typeName] = true
					typePos := fset.Position(decl.Pos()).Offset
					typeEnd := fset.Position(decl.End()).Offset
					builder.WriteString(string(src[typePos:typeEnd]) + "\n\n")
				}
			}
		case *ast.FuncDecl:
			funcName := decl.Name.Name
			var receiverType string

			if decl.Recv != nil && len(decl.Recv.List) > 0 {
				recvExpr := decl.Recv.List[0].Type
				switch expr := recvExpr.(type) {
				case *ast.Ident:
					receiverType = expr.Name
				case *ast.StarExpr:
					if ident, ok := expr.X.(*ast.Ident); ok {
						receiverType = ident.Name
					}
				}

				if receiverType != "" {
					if _, exists := funcMap[receiverType]; !exists {
						funcMap[receiverType] = make(map[string]bool)
					}
					if funcMap[receiverType][funcName] {
						return false
					}
					funcMap[receiverType][funcName] = true
				}
			} else {
				if nonReceiverFuncMap[funcName] {
					return false
				}
				nonReceiverFuncMap[funcName] = true
			}

			funcPos := fset.Position(decl.Pos()).Offset
			funcEnd := fset.Position(decl.End()).Offset
			builder.WriteString(string(src[funcPos:funcEnd]) + "\n\n")
		default:
			return true
		}
		return true
	})

	if err := os.WriteFile(filePath, []byte(builder.String()), 0644); err != nil {
		return fmt.Errorf("failed to write modified file: %w", err)
	}

	return nil
}

func RemoveCustomType(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	var result []string
	scanner := bufio.NewScanner(file)
	// 전체 파일 라인을 읽음
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err = scanner.Err(); err != nil {
		return err
	}

	i := 0
	for i < len(lines) {
		// 현재 줄이 "CustomType"을 포함하는 경우,
		// 해당 줄과 그 이후 4줄(총 5줄)을 건너뛰기
		if strings.Contains(lines[i], "CustomType") {
			// 삭제할 5줄 블럭 범위
			i += 5
			continue
		}
		result = append(result, lines[i])
		i++
	}

	output := strings.Join(result, "\n")
	// io.WriteFile 대신 os.WriteFile을 사용하여 올바르게 파일을 씁니다.
	return os.WriteFile(filePath, []byte(output), 0644)
}
