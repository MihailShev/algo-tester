package tester

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type InFilePath = string
type OutFilePath = string

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
	testNumber := 0
	finish := false

	for !finish {
		inFilePath := filepath.FromSlash(fmt.Sprintf("%s/test.%d.in", t.path, testNumber))
		outFilePath := filepath.FromSlash(fmt.Sprintf("%s/test.%d.out", t.path, testNumber))

		fmt.Printf("Start test %d\n", testNumber)

		if isFileExist(inFilePath) && isFileExist(outFilePath) {
			t.runTest(inFilePath, outFilePath, testNumber)
		} else {
			finish = true
		}

		testNumber += 1
	}

}

func (t *Tester) RunTestWithCount(count int) {
	for i := 0; i < count; i++ {
		inFile, outFile := t.getTestFiles(i)

		fmt.Printf("Start test %d\n", i)

		if isFileExist(inFile) && isFileExist(outFile) {
			t.runTest(inFile, outFile, i)
		}
	}
}

func (t *Tester) getTestFiles(testNumber int) (InFilePath, OutFilePath) {
	inFilePath := filepath.FromSlash(fmt.Sprintf("%s/test.%d.in", t.path, testNumber))
	outFilePath := filepath.FromSlash(fmt.Sprintf("%s/test.%d.out", t.path, testNumber))

	return inFilePath, outFilePath
}

func (t *Tester) runTest(in InFilePath, out OutFilePath, testNumber int) {
	fmt.Printf("Start test %d", testNumber)

	isSuccess, err := t.execute(in, out)

	if err != nil {
		fmt.Printf("Test %d returned an error: %s\n", testNumber, err)
	}

	if isSuccess {
		fmt.Printf("Test %d is successful\n", testNumber)
	} else {
		fmt.Printf("Test %d failed\n", testNumber)
	}
	fmt.Println("=========================")
}

func (t *Tester) execute(in string, out string) (bool, error) {
	inData, err := readAllLine(in)

	if err != nil {
		return false, err
	}

	expect, err := readExpect(out)

	fmt.Printf("Expect %s\n", expect)

	if err != nil {
		return false, err
	}

	startTime := time.Now()
	result := t.task.Run(inData)
	endTime := time.Since(startTime)

	fmt.Printf("Got %s\n", result)

	fmt.Printf("Execution time %v\n", endTime)
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
