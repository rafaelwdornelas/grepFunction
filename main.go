package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <directory>")
		fmt.Println("Example: go run main.go .\\mob_009_multiplepinning\\")
		return
	}
	directory := os.Args[1]
	result, err := grepFunction("java/lang/String;Ljava/lang/String;\\)", "java/lang/Object;Ljava/lang/String;\\)", directory)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(result)
}
func grepFunction(pattern1, pattern2, directory string) (string, error) {
	var result strings.Builder

	regex1, err := regexp.Compile(pattern1)
	if err != nil {
		return "", err
	}

	regex2, err := regexp.Compile(pattern2)
	if err != nil {
		return "", err
	}

	err = filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		var buffer []string

		for scanner.Scan() {
			line := scanner.Text()
			if regex1.MatchString(line) {
				buffer = append(buffer, line)
				for i := 0; i < 10 && scanner.Scan(); i++ {
					buffer = append(buffer, scanner.Text())
				}
			}
		}

		for _, line := range buffer {
			if regex2.MatchString(line) {
				result.WriteString("File: ")
				result.WriteString(path)
				result.WriteString("\nLine: ")
				result.WriteString(line)
				result.WriteString("\n")
			}
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	return result.String(), nil
}
