//go:build !windows

package dshcli

import "os/exec"

func Lookup() (binary, entry string) {
	if cli, err := exec.LookPath("dsh"); err == nil {
		return cli, ""
	}
	return "", ""
}
