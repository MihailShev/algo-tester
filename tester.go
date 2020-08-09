package tester

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type ITask interface {
	Run(data []string) string
}

type Tester struct {
	task ITask
	path string
}

func NewTester(task ITask, path string) Tester {
	t := Tester{
		task: task,
		path: path,
	}

	return t
}

func (t *Tester) RunTest() {
	nr := 0
	finish := false

	for !finish {
		inFilePath := filepath.FromSlash(fmt.Sprintf("%s/test.%d.in", t.path, nr))
		outFilePath := filepath.FromSlash(fmt.Sprintf("%s/test.%d.out", t.path, nr))

		fmt.Printf("Start test %d\n", nr)

		if isFileExist(inFilePath) && isFileExist(outFilePath) {
			isSuccess, err := execute(t.task, inFilePath, outFilePath)

			if err != nil {
				fmt.Printf("Test %d returned an error: %s\n", nr, err)
			}

			if isSuccess {
				fmt.Printf("Test %d is successful\n", nr)
			} else {
				fmt.Printf("Test %d failed\n", nr)
			}
			fmt.Println("=========================")
		} else {
			finish = true
		}

		nr += 1
	}

}

func execute(task ITask, in string, out string) (bool, error) {
	inData, err := readAllLine(in)

	if err != nil {
		return false, err
	}

	expect, err := readExpect(out)

	fmt.Printf("Expect %s\n", expect)

	if err != nil {
		return false, err
	}

	result := task.Run(inData)

	fmt.Printf("Got %s\n", result)

	return expect == result, nil
}

func readAllLine(path string) ([]string, error) {
	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	lines := make([]string, 0)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, nil
}

func readExpect(path string) (string, error) {
	bytes, err := ioutil.ReadFile(path)

	if err != nil {
		return "", err
	}

	expect := string(bytes)

	expect = strings.Trim(expect, "\n")
	expect = strings.Trim(expect, "\r")
	return expect, nil
}

func isFileExist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println("File:", path, "not created")
		return false
	}

	return true
}
