package workspace

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
)

const ProjectDefinitionFileName = ".hackt"

type Project struct {
	// fs is the file system that holds all project files. Use this to reload
	// this project's properties.
	fs afero.Fs
	// dotfile holds all the properties of this project. It is loaded from and
	// stored to a file within the project folder.
	dotfile *DotFile
}

func LoadProject(fs afero.Fs) (*Project, error) {
	definitionFile, err := fs.Open(ProjectDefinitionFileName)
	if err == os.ErrNotExist {
		return nil, ErrNoProjectDefinitionFile
	} else if err != nil {
		return nil, fmt.Errorf("open %s: %w", ProjectDefinitionFileName, err)
	}
	defer func() { _ = definitionFile.Close() }()

	dotfile, err := LoadDotFile(definitionFile)
	if err != nil {
		return nil, fmt.Errorf("load dot file: %w", err)
	}

	return &Project{
		fs:      fs,
		dotfile: dotfile,
	}, nil
}

func (p Project) Name() string {
	return p.dotfile.Name
}

func (p Project) Fs() afero.Fs {
	return p.fs
}

type DotFile struct {
	Name string `yaml:"name"`
}

func LoadDotFile(rd io.Reader) (*DotFile, error) {
	var dotfile DotFile
	if err := yaml.NewDecoder(rd).Decode(&dotfile); err != nil {
		return nil, err
	}

	return &dotfile, nil
}
