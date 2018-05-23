package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
)

type Directory struct {
	Path      string   `json:"path"`
	Filepaths []string `json:"files"`
	Size      int64    `json:"size"`
}

type Scanner struct {
	*Directory
	Size int64 `json:"size"`
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("Please specify a directory to generate manifests")
	}
	directory := os.Args[1]
	os.Remove(directory + "/mainfest.json")
	s := &Scanner{&Directory{directory, nil, 0}, 0}
	if err := filepath.Walk(directory, s.visit); err != nil {
		log.Fatalf("Error: %v", err.Error())
	}

	GenerateManifest(s.Directory)
}

func (s *Scanner) visit(fp string, f os.FileInfo, err error) error {
	if path.Ext(fp) != ".json" && !f.IsDir() {
		s.Directory.Filepaths = append(s.Directory.Filepaths, f.Name())
		s.Directory.Size += f.Size()
		s.Size = s.Size + f.Size()
	}

	return nil
}

// GenerateManifest containing filepaths
func GenerateManifest(d *Directory) {
	manifestJSON, err := json.Marshal(d)
	if err != nil {
		log.Fatalf("Error cannot create manifest: %v", err.Error())
	}
	ioutil.WriteFile(d.Path+"/manifest.json", manifestJSON, 0644)
}
