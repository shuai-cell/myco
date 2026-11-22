package dshcli

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/sys/windows/registry"
)

func Lookup() (binary, entry string) {
	cli, err := exec.LookPath("dsh")
	if err == nil && strings.EqualFold(filepath.Ext(cli), ".exe") {
		return cli, ""
	}
	node, err := exec.LookPath("node.exe")
	if err != nil {
		return "", ""
	}
	var roots []string
	if cli != "" && filepath.IsAbs(cli) {
		roots = append(roots, filepath.Dir(cli))
	}
	for _, hive := range []registry.Key{registry.CURRENT_USER, registry.LOCAL_MACHINE} {
		key, err := registry.OpenKey(hive, `Software\Microsoft\Windows\CurrentVersion\Uninstall\1bf39983-50d0-5fe0-9ef4-cece76f67c5e`, registry.QUERY_VALUE)
		if err != nil {
			continue
		}
		dir, _, err := key.GetStringValue("InstallLocation")
		key.Close()
		if err == nil && filepath.IsAbs(dir) {
			for _, resources := range []string{"app", "app.asar.unpacked"} {
				roots = append(roots, filepath.Join(dir, "resources", resources, "dsh"))
			}
		}
	}
	for _, root := range roots {
		entry := filepath.Join(root, "node_modules", "@deepseek-ai", "dsh", "lib", "bin.js")
		if info, err := os.Stat(entry); err == nil && info.Mode().IsRegular() {
			return node, entry
		}
	}
	return "", ""
}
