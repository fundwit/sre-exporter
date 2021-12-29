package testinfra

import (
	"fmt"
	"os"
	"path"
)

type TempFile struct {
	file string
}

func NewFileWithContent(file, content string) (*TempFile, error) {
	realFile := path.Join(os.TempDir(), file)
	dir := path.Dir(realFile)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return nil, err
	}

	f, err := os.Create(realFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	_, err = f.WriteString(content)
	if err != nil {
		f.Close()
		os.Remove(realFile)
		return nil, err
	}

	return &TempFile{file: realFile}, nil
}

func (tf *TempFile) Clear() {
	if tf.file != "" {
		if err := os.Remove(tf.file); err != nil {
			fmt.Printf("failed to remove file %s\n, %v", tf.file, err)
		} else {
			fmt.Printf("success to remove file %s\n", tf.file)
		}
	}
}
