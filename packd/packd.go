package packd

import (
    "bytes"
    "fmt"
    "os"
    "path/filepath"
)

func Packd(path string, to bytes.Buffer) error {
    return nil
}

func WalkDir(wpath string) error {
	subDirToSkip := "skip"

    err := filepath.Walk(wpath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if info.IsDir() && info.Name() == subDirToSkip {
			fmt.Printf("skipping a dir without errors: %+v \n", info.Name())
			return filepath.SkipDir
		}
		fmt.Printf("visited file or dir: %q\n", path)
		return nil
	})
	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", wpath, err)
        return err
	}
    return nil
}

