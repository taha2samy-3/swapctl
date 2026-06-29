package system

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func CreateSwapFile(path, size string) error {
	exec.Command("fallocate", "-l", size, path).Run()
	info, err := os.Stat(path)
	if err != nil || info.Size() < 1024 {
		num := ""
		unit := "G"
		for _, char := range size {
			if char >= '0' && char <= '9' {
				num += string(char)
			} else {
				unit = strings.ToUpper(string(char))
			}
		}
		return exec.Command("dd", "if=/dev/zero", "of="+path, "bs=1"+unit, "count="+num, "status=progress").Run()
	}
	return nil
}

func FormatSwap(path string) error {
	out, err := exec.Command("mkswap", path).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", string(out))
	}
	return nil
}

func EnableSwap(path string) error {
	out, err := exec.Command("swapon", path).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", string(out))
	}
	return nil
}

func GetActiveSwap() (string, error) {
	out, err := exec.Command("swapon", "--show").Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func SetPermissions(path string) error {
	return os.Chmod(path, 0600)
}

func DisableSwap(path string) error {
	exec.Command("swapoff", path).Run()
	return nil
}

func SetKernelParam(path string, value int) error {
	return os.WriteFile(path, []byte(strconv.Itoa(value)), 0644)
}
