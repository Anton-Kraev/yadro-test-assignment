package computer_club

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEventProcessing(t *testing.T) {
	testFiles := getTestFiles("../../testdata")

	for _, tf := range testFiles {
		input, err1 := os.Open(tf.input)
		expectedOutput, err2 := os.Open(tf.expectedOutput)
		actualOutput, err3 := os.Create(tf.actualOutput)
		checkFatal(errors.Join(err1, err2, err3))

		ProcessComputerClubDayEvents(input, actualOutput)

		checkFatal(actualOutput.Close())
		actualOutput, err := os.Open(tf.actualOutput)
		checkFatal(err)

		equal, lineN, expectedLine, actualLine := fileCmp(expectedOutput, actualOutput)

		checkFatal(errors.Join(input.Close(), expectedOutput.Close(), actualOutput.Close(), os.Remove(tf.actualOutput)))

		if !equal {
			t.Errorf(
				"Process events of %s: expected %s, got %s at line %d",
				tf.input, expectedLine, actualLine, lineN,
			)
		} else {
			t.Logf("--- OK: Process events of %s\n", tf.input)
		}
	}
}

type testFile struct {
	input, expectedOutput, actualOutput string
}

func getTestFiles(testdataDir string) []testFile {
	var testFiles []testFile

	err := filepath.Walk(testdataDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			base := filepath.Base(path)
			if strings.HasPrefix(base, "input") && strings.HasSuffix(base, ".txt") {
				id := strings.TrimSuffix(strings.TrimPrefix(base, "input"), ".txt")
				testFiles = append(testFiles, testFile{
					input:          path,
					expectedOutput: filepath.Join(testdataDir, fmt.Sprintf("expected_output%s.txt", id)),
					actualOutput:   filepath.Join(testdataDir, fmt.Sprintf("actual_output%s.txt", id)),
				})
			}
		}
		return nil
	})

	checkFatal(err)
	return testFiles
}

func fileCmp(fileExpected, fileActual io.Reader) (equal bool, lineN int, expected, actual string) {
	expectedScanner := bufio.NewScanner(fileExpected)
	actualScanner := bufio.NewScanner(fileActual)

	for {
		lineN++
		hasExpectedLine := expectedScanner.Scan()
		hasActualLine := actualScanner.Scan()

		if !hasExpectedLine && !hasActualLine {
			break
		} else if hasExpectedLine && !hasActualLine {
			return false, lineN, expectedScanner.Text(), ""
		} else if !hasExpectedLine {
			return false, lineN, "", actualScanner.Text()
		}
		if !bytes.Equal(expectedScanner.Bytes(), actualScanner.Bytes()) {
			return false, lineN, expectedScanner.Text(), actualScanner.Text()
		}
	}

	return true, 0, "", ""
}
