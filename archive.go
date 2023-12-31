package activate_toolchain

import (
	"context"
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/yankeguo/activate-toolchain/pkg/unarchive"
)

type InstallArchiveOptions struct {
	ProvideURLs    func() (urls []string, err error)
	Name           string
	File           string
	DirectoryLevel int
	DirectoryPath  string
}

// InstallArchive installs an archive to a directory
func InstallArchive(ctx context.Context, opts InstallArchiveOptions) (dir string, err error) {
	var home string
	if home, err = os.UserHomeDir(); err != nil {
		return
	}
	base := filepath.Join(home, ".atc")
	if err = os.MkdirAll(base, 0755); err != nil {
		return
	}

	dir = filepath.Join(base, opts.Name)

	// check if already installed
	{
		var stat os.FileInfo
		if stat, err = os.Stat(dir); err == nil {
			if stat.IsDir() {
				return
			}
			err = errors.New("target directory is not a directory: " + dir)
			return
		} else {
			if !os.IsNotExist(err) {
				return
			}
			err = nil
		}
	}

	// temp directory for atomic installation
	dirTemp := filepath.Join(base, opts.Name+".tmp")
	os.RemoveAll(dirTemp)
	if err = os.MkdirAll(dirTemp, 0755); err != nil {
		return
	}
	defer os.RemoveAll(dirTemp)

	file := filepath.Join(base, opts.File)

	// ensure file exists
	{
		var stat os.FileInfo
		if stat, err = os.Stat(file); err == nil {
			if stat.IsDir() {
				err = errors.New("target file is a directory: " + file)
				return
			}
		} else {
			if !os.IsNotExist(err) {
				return
			}
			err = nil

			log.Println("downloading:", opts.File)

			var urls []string
			if opts.ProvideURLs != nil {
				urls, err = opts.ProvideURLs()
				if err != nil {
					return
				}
			} else {
				err = errors.New("no urls provided")
				return
			}

			if err = AdvancedFetchFile(ctx, urls, file); err != nil {
				return
			}
		}
	}

	log.Println("extracting:", opts.File)

	// extract file
	{
		var f *os.File
		if f, err = os.OpenFile(file, os.O_RDONLY, 0644); err != nil {
			return
		}

		if err = unarchive.Unarchive(f, dirTemp); err != nil {
			_ = f.Close()
			return
		}

		_ = f.Close()
	}

	dirSrc := dirTemp

	if opts.DirectoryLevel > 0 {
	outerLoop:
		for i := 0; i < opts.DirectoryLevel; i++ {
			var dirs []os.DirEntry
			if dirs, err = os.ReadDir(dirSrc); err != nil {
				return
			}
			for _, dir := range dirs {
				if dir.IsDir() {
					dirSrc = filepath.Join(dirSrc, dir.Name())
					continue outerLoop
				}
			}
			err = errors.New("no directory found in " + dirSrc)
			return
		}
	}

	if opts.DirectoryPath != "" {
		dirSrc = filepath.Join(dirSrc, opts.DirectoryPath)
	}

	os.RemoveAll(dir)
	if err = os.Rename(dirSrc, dir); err != nil {
		return
	}

	log.Println("installed:", opts.Name)

	os.RemoveAll(file)

	return
}
