package config

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/kkyr/fig"
	"github.com/nenecchuu/arcana/env"
)

var (
	allowedFileExt map[string]bool = map[string]bool{".yaml": true, ".json": true}
)

func ReadConfig(cfg interface{}, path string, module string) error {
	var (
		fileName string
		fileExt  string
		environ  string
		wd       string
		err      error
	)

	fileExt, err = getFileExt(path)

	if err != nil {
		return err
	}

	environ = env.GetEnvironmentName()
	v, ok := allowedFileExt[fileExt]

	if !v || !ok {
		return errors.New("config file doesn't exist or invalid file format")
	}

	fileName = module + "." + environ + fileExt
	wd, err = os.Getwd()

	if err != nil {
		return err
	}

	return fig.Load(cfg,
		fig.File(fileName),
		fig.UseEnv(strings.ToUpper(module)),
		fig.Dirs(wd, path),
	)
}

func ReadConfigFromSpecifiedFile(cfg interface{}, configFileName string) error {
	var (
		wd  string
		err error
	)

	wd, err = os.Getwd()
	if err != nil {
		return err
	}

	return fig.Load(cfg,
		fig.File(configFileName),
		fig.UseEnv(""),
		fig.Dirs(wd, "config"),
	)
}

func getFileExt(root string) (string, error) {
	var (
		fileExt string
		err     error
	)

	err = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		ext := filepath.Ext(d.Name())
		v, ok := allowedFileExt[ext]

		if v && ok {
			fileExt = ext
		}

		return nil
	})

	return fileExt, err
}
