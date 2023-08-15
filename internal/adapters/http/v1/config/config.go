package config

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 15 August 2023
 */
import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Conf struct {
	Port string `yaml:"port"` 
}

func Parse(fpath string) (Conf, error) {

	conf := Conf{}
	if fpath == `` {
		return conf, nil
	}

	filename, err := filepath.Abs(fpath)
	if err != nil {
		return conf, fmt.Errorf("can`t open confd %v, %v", fpath, err)
	}

	r, err := os.ReadFile(filename)
	if err != nil {
		return conf, fmt.Errorf("can`t read conf %v, %v", filename, err)
	}

	err = yaml.Unmarshal(r, &conf)
	if err != nil {
		return conf, fmt.Errorf("can`t unmarshal conf %v", err)
	}
	return conf, nil
}
