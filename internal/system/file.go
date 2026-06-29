package system

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func AppendToFstab(path, swapPath string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if strings.Contains(string(content), swapPath) {
		return nil
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(fmt.Sprintf("%s none swap sw 0 0\n", swapPath))
	return err
}

func UpdateSysctl(path, key, value string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, key+"=") {
			lines = append(lines, line)
		}
	}
	lines = append(lines, fmt.Sprintf("%s=%s", key, value))
	return os.WriteFile(path, []byte(strings.Join(lines, "\n")+"\n"), 0644)
}
