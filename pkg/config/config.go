package config

import (
	"errors"
	"io/fs"
	"os"

	"github.com/aymericbeaumet/pimp/pkg/util"
	"gopkg.in/yaml.v2"
)

type configYAML struct {
	Pimpfiles []string `yaml:"pimpfiles"`
}

func Load(path string) (*Config, error) {
	var out Config
	var allowErrNotExist bool

	if len(path) == 0 {
		allowErrNotExist = true
		path = "~/.pimprc"
	}

	expanded, err := util.NormalizePath(path)
	if err != nil {
		return nil, err
	}

	var parsed configYAML
	f, err := os.Open(expanded)
	if err != nil {
		if !(errors.Is(err, fs.ErrNotExist) && allowErrNotExist) {
			return nil, err
		}
	} else {
		defer f.Close()
		if err := yaml.NewDecoder(f).Decode(&parsed); err != nil && err.Error() != "EOF" {
			return nil, err
		}
	}

	if err := out.setPimpfiles(parsed.Pimpfiles); err != nil {
		return nil, err
	}

	return &out, nil
}

type Config struct {
	Pimpfiles []*os.File `json:"pimpfiles"`
}

func (conf *Config) Close() {
	for _, pimpfile := range conf.Pimpfiles {
		defer pimpfile.Close()
	}
}

func (conf *Config) setPimpfiles(pimpfiles []string) error {
	var skipErrNotExist bool

	if len(pimpfiles) == 0 {
		skipErrNotExist = true
		pimpfiles = []string{"~/.Pimpfile.go", "~/.Pimpfile"}
	}

	for _, name := range pimpfiles {
		normalized, err := util.NormalizePath(name)
		if err != nil {
			return err
		}
		f, err := os.Open(normalized)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) && skipErrNotExist {
				continue
			}
			return err
		}
		// no need to f.Close() here as we are taking care of it in the
		// conf.Close() method
		conf.Pimpfiles = append(conf.Pimpfiles, f)
	}

	return nil
}
