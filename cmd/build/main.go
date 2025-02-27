package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"matins-blog/internal/markdown"
)

func main() {
	// Create dist directory
	if err := os.MkdirAll("dist", 0755); err != nil {
		log.Fatal(err)
	}

	// Copy static assets only if the directory exists
	if _, err := os.Stat("static"); !os.IsNotExist(err) {
		if err := copyDir("static", "dist/static"); err != nil {
			log.Fatal(err)
		}
	}

	// Process markdown files
	err := filepath.Walk("content", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !strings.HasSuffix(info.Name(), ".md") {
			return nil
		}

		content, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		html, err := markdown.ProcessMarkdown(content)
		if err != nil {
			return err
		}

		// Create output path
		relPath, err := filepath.Rel("content", path)
		if err != nil {
			return err
		}

		outPath := filepath.Join("dist", strings.TrimSuffix(relPath, ".md")+".html")
		if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
			return err
		}

		return ioutil.WriteFile(outPath, html, 0644)
	})

	if err != nil {
		log.Fatal(err)
	}
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		return ioutil.WriteFile(dstPath, data, info.Mode())
	})
}
