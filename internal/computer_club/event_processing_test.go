package computer_club

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"log"
	"os"
	"testing"
)

func TestEventProcessing(t *testing.T) {
	var testFiles = []struct {
		input, expectedOutput, actualOutput string
	}{
		{"../../testdata/input1.txt",
			"../../testdata/expected_output1.txt",
			"../../testdata/actual_output1.txt"},
	}

	for _, tf := range testFiles {
		input, err1 := os.Open(tf.input)
		expectedOutput, err2 := os.Open(tf.expectedOutput)
		actualOutput, err3 := os.Create(tf.actualOutput)
		if err := errors.Join(err1, err2, err3); err != nil {
			log.Fatalln(err)
		}

		ProcessComputerClubDayEvents(input, actualOutput)
		if err := actualOutput.Close(); err != nil {
			log.Fatalln(err)
		}

		actualOutput, err := os.Open(tf.actualOutput)
		if err != nil {
			log.Fatalln(err)
		}

		equal, lineN, expectedLine, actualLine := fileCmp(expectedOutput, actualOutput)

		if err := errors.Join(input.Close(), expectedOutput.Close(), actualOutput.Close(), os.Remove(tf.actualOutput)); err != nil {
			log.Fatalln(err)
		}

		if !equal {
			t.Errorf("Process events of %s: expected %s, got %s at line %d",
				tf.input, expectedLine, actualLine, lineN)
		} else {
			t.Logf("--- OK: Process events of %s\n", tf.input)
		}
	}
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
