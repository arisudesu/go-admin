package web

import (
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
)

func LoadTemplates(template *template.Template, directory string, patterns ...string) error {
	return filepath.WalkDir(directory, func(path string, d fs.DirEntry, err error) error {
		// Skip reading errors
		if err != nil {
			return err
		}

		// Skip directory itself
		if d.IsDir() {
			return nil
		}

		relpath, err := filepath.Rel(directory, path)
		if err != nil {
			return err
		}

		// Normalize to *nix slashes
		relpath = filepath.ToSlash(relpath)

		for _, pattern := range patterns {
			match, err := filepath.Match(pattern, relpath)
			if err != nil {
				return err
			}

			if match {
				src, err := os.ReadFile(path)
				if err != nil {
					return err
				}

				_, err = template.New(relpath).Parse(string(src))
				return err
			}
		}
		return nil
	})
}
