package downloader

import (
	"archive/tar"
	"bytes"
	"io"
	"io/fs"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

func createTar(tr *tar.Writer, filesystem fs.FS, vars any) error {
	return fs.WalkDir(filesystem, ".", func(path string, entry fs.DirEntry, err error) error {
		if path == "." {
			return nil
		}
		fi, err := entry.Info()
		if err != nil {
			return err
		}
		if entry.IsDir() {
			t := time.Now()
			err = tr.WriteHeader(&tar.Header{
				Name:       path,
				Size:       0,
				Mode:       int64(fi.Mode()),
				ModTime:    t,
				AccessTime: t,
				ChangeTime: t,
			})
			if err != nil {
				return err
			}
		} else if strings.HasSuffix(path, ".tmpl") {
			tmpl, err := template.New(filepath.Base(path)).ParseFS(filesystem, path)
			if err != nil {
				return err
			}
			outputBuf := bytes.NewBuffer(nil)
			err = tmpl.Execute(outputBuf, vars)
			if err != nil {
				return err
			}
			actualPath := path[0 : len(path)-5]
			contentBytes := outputBuf.Bytes()
			t := time.Now()
			err = tr.WriteHeader(&tar.Header{
				Name:       actualPath,
				Size:       int64(len(contentBytes)),
				Mode:       int64(fi.Mode()),
				ModTime:    t,
				AccessTime: t,
				ChangeTime: t,
			})
			if err != nil {
				return err
			}
			_, err = tr.Write(contentBytes)
			if err != nil {
				return err
			}
		} else {
			f, err := filesystem.Open(path)
			if err != nil {
				return err
			}
			defer f.Close()
			contentBytes, err := io.ReadAll(f)
			t := time.Now()
			err = tr.WriteHeader(&tar.Header{
				Name:       path,
				Size:       int64(len(contentBytes)),
				Mode:       int64(fi.Mode()),
				ModTime:    t,
				AccessTime: t,
				ChangeTime: t,
			})
			if err != nil {
				return err
			}
			_, err = tr.Write(contentBytes)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
