package system

import (
	"bufio"
	"bytes"
	"os/exec"
	"strings"
)

type MountPoint struct {
	Path string
	Free string
}

func GetMountPoints() ([]MountPoint, error) {
	var list []MountPoint
	cmd := exec.Command("df", "-h")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(&out)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "/dev/") && !strings.Contains(line, "overlay") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) >= 6 {
			list = append(list, MountPoint{
				Path: fields[5],
				Free: fields[3],
			})
		}
	}
	return list, nil
}
