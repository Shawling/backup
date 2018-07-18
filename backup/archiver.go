package backup

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

type Archiver interface {
	DestExt() string
	Archiver(src, dst string) error
}

type zipper struct {
}

func (z *zipper) Archiver(src, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), os.ModePerm); err != nil {
		return err
	}
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	w := zip.NewWriter(out)
	defer w.Close()
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil // skip
		}
		if err != nil {
			return err
		}
		in, err := os.Open(path)
		if err != nil {
			return err
		}
		defer in.Close()
		f, err := w.Create(path)
		if err != nil {
			return err
		}
		_, err = io.Copy(f, in)
		if err != nil {
			return err
		}
		return nil
	})
}

func (z *zipper) DestExt() string {
	return ".zip"
}

// ZIP 建立不同的 Archive 对象，对外暴露只有名字不同，类型一样。但是不同名字的 Archieve 却可以再同名函数下完成不同的功能。
var ZIP Archiver = (*zipper)(nil)
