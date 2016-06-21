// +build windows

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

func (c *ctx) isRunning(name string) (bool, error) {
	filename := filepath.Join(c.s.Destination, name)
	f, err := os.OpenFile(filename+".exe", os.O_RDWR, 0600)
	if err != nil {
		x, ok := err.(*os.PathError)
		if !ok {
			return false, fmt.Errorf("invalid type")
		}
		e, ok := x.Err.(syscall.Errno)
		if !ok {
			return false, fmt.Errorf("invalid error type")
		}
		if e == syscall.FILE_MAP_EXECUTE {
			return true, nil
		}
		return false, err
	}
	defer f.Close()

	return false, nil
}
