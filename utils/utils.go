package utils

import (
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"time"
)

func AbsolutePath(dir, file string) (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	tokenCacheDir := filepath.Join(usr.HomeDir, dir)
	os.MkdirAll(tokenCacheDir, 0700)
	return filepath.Join(tokenCacheDir, url.QueryEscape(file)), err
}

func UnixMili() int64 {
	return time.Now().UnixNano() / 1000000
}
