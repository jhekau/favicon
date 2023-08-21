package yaml

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

func Parse(fpath string, conf any) error {

	if fpath == `` {
		return nil
	}

	filename, err := filepath.Abs(fpath)
	if err != nil {
		return fmt.Errorf("can`t open confd %v, %v", fpath, err)
	}

	r, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("can`t read conf %v, %v", filename, err)
	}

	err = yaml.Unmarshal(r, &conf)
	if err != nil {
		return fmt.Errorf("can`t unmarshal conf %v", err)
	}
	return nil
}
