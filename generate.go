//go:generate go run . clean-graphql-generated
//go:generate go run github.com/99designs/gqlgen generate --config internal/graphql/gqlgen.yml

package main

import (
	"bufio"
	"fmt"
	"github.com/99designs/gqlgen/codegen/config"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "clean-graphql-generated" {
		fmt.Printf("Cleaning up generated graphql files...\n")
		f, err := os.OpenFile("generate.go", os.O_RDONLY, 0644)
		if err != nil {
			panic(err)
		}
		defer func(f *os.File) {
			err := f.Close()
			if err != nil {
				panic(err)
			}
		}(f)

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			ln := scanner.Text()
			if strings.HasPrefix(ln, "//go:generate") {
				parts := strings.Split(ln, " ")
				if len(parts) == 7 && parts[1] == "go" && parts[2] == "run" &&
					parts[3] == "github.com/99designs/gqlgen" &&
					parts[4] == "generate" &&
					parts[5] == "--config" &&
					(strings.HasSuffix(parts[6], ".yml") || strings.HasSuffix(parts[6], ".yaml")) {
					cfg, err := os.ReadFile(parts[6])
					if err != nil {
						panic(err)
					}
					gc := config.Config{}
					err = yaml.Unmarshal(cfg, &gc)
					if err != nil {
						panic(err)
					}
					if gc.Exec.Layout == "single-file" {
						if err = os.Remove(gc.Exec.Filename); err == nil {
							fmt.Printf("removed %s\n", gc.Exec.Filename)
						}
						if err = os.Remove(gc.Model.Filename); err == nil {
							fmt.Printf("removed %s\n", gc.Model.Filename)
						}
						for _, pat := range gc.SchemaFilename {
							entries, err := filepath.Glob(pat)
							if err != nil {
								panic(err)
							}
							for _, entry := range entries {
								resolver := gqlToResolverName(gc.Resolver.Dir(), entry, gc.Resolver.FilenameTemplate)
								if err = os.Remove(resolver); err == nil {
									fmt.Printf("removed %s\n", resolver)
								}
							}
						}
					}
				}
			}
		}
	}
}

// Check [github.com/99designs/gqlgen@v0.17.62/plugin/resolvergen/resolver.go:309]
func gqlToResolverName(base, gqlname, filenameTmpl string) string {
	gqlname = filepath.Base(gqlname)
	ext := filepath.Ext(gqlname)
	if filenameTmpl == "" {
		filenameTmpl = "{name}.resolvers.go"
	}
	filename := strings.ReplaceAll(filenameTmpl, "{name}", strings.TrimSuffix(gqlname, ext))
	return filepath.Join(base, filename)
}
