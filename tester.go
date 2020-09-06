package tester

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type InFilePath = string
type OutFilePath = string

const MaxOutputLen = 250

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
		t.RunTestNum(i)
	}
}

func (t *Tester) RunTestNum(testNum int) {
	inFile, outFile := t.getTestFiles(testNum)

	fmt.Printf("Start test %d\n", testNum)

	if isFileExist(inFile) && isFileExist(outFile) {
		t.runTest(inFile, outFile, testNum)
	}
}

func (t *Tester) getTestFiles(testNumber int) (InFilePath, OutFilePath) {
	inFilePath := filepath.FromSlash(fmt.Sprintf("%s/test.%d.in", t.path, testNumber))
	outFilePath := filepath.FromSlash(fmt.Sprintf("%s/test.%d.out", t.path, testNumber))

	return inFilePath, outFilePath
}

func (t *Tester) runTest(in InFilePath, out OutFilePath, testNumber int) {
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

	fmt.Printf("Expect %s\n", cutStr(MaxOutputLen, expect))

	if err != nil {
		return false, err
	}

	startTime := time.Now()
	result := t.task.Run(inData)
	endTime := time.Since(startTime)

	fmt.Printf("Got %s\n", cutStr(MaxOutputLen, result))

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
	scanner.Buffer(make([]byte, 0), math.MaxInt64)

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

func cutStr(maxLen int, str string) string {
	if len(str) > maxLen {
		return fmt.Sprintf("%.*s..", maxLen, str)
	}

	return str
}
