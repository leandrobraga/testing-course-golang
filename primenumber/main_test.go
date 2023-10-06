package main

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

// For run a especific test just use the argument --run: go run --run name_test
// go run --run Teste_alpha_isPrime or go run -v --run Teste_alpha_isPrime (for verbose mode)

// alpha_ indicate is a suite test.
// Suite test can have any name
// For run a suite test: go test --run Test_alpha or go test -v --run Test_alpha for verbose mode
// Of course whether suite test is other just change the name Test_anything
func Test_alpha_isPrime(t *testing.T) {
	testPrimes := []struct {
		name     string
		testNum  int
		expected bool
		msg      string
	}{
		{"prime", 7, true, "7, is a prime number."},
		{"not prime", 8, false, "8 is not a prime number because it divisible by 2!"},
		{"0", 0, false, "0 is not prime, by definition!"},
		{"1", 1, false, "1 is not prime, by definition!"},
		{"-1", -1, false, "Negative numbers are not prime, by defition!"},
	}

	for _, e := range testPrimes {
		result, msg := isPrime(e.testNum)

		if e.expected && !result {
			t.Errorf("%s: expected true but got false", e.name)
		}

		if !e.expected && result {
			t.Errorf("%s: expected false but got true", e.name)
		}

		if msg != e.msg {
			t.Errorf("%s: expected %s but got %s", e.name, msg, e.msg)
		}
	}

}

func Test_alpha_prompt(t *testing.T) {
	// save a copy of os.Stdout
	oldOut := os.Stdout

	// create a read and write pipe
	r, w, _ := os.Pipe()

	// set os.Stdout to our write pipe
	os.Stdout = w

	prompt()

	// close our writer
	_ = w.Close()

	// reset os.Stdout to what it was before
	os.Stdout = oldOut

	// read the output of our prompt() func from our read pipe
	out, _ := io.ReadAll(r)

	// perform our test
	if string(out) != "-> " {
		t.Errorf("incorrect prompt: expect ->  but got %s", string(out))
	}
}

func Test_intro(t *testing.T) {
	// save a copy of os.Stdout
	oldOut := os.Stdout

	// create a read and write pipe
	r, w, _ := os.Pipe()

	// set os.Stdout to our write pipe
	os.Stdout = w

	intro()

	// close our writer
	_ = w.Close()

	// reset os.Stdout to what it was before
	os.Stdout = oldOut

	// read the output of our prompt() func from our read pipe
	out, _ := io.ReadAll(r)

	// perform our test
	if !strings.Contains(string(out), "Enter a whole number") {
		t.Errorf("intro text not correct; got %s", string(out))
	}
}

func Test_checkNumbers(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"empty", "", "Please enter a whole number!"},
		{"zero", "0", "0 is not prime, by definition!"},
		{"one", "1", "1 is not prime, by definition!"},
		{"two", "2", "2, is a prime number."},
		{"three", "3", "3, is a prime number."},
		{"negative", "-1", "Negative numbers are not prime, by defition!"},
		{"typed", "three", "Please enter a whole number!"},
		{"decimal", "1.1", "Please enter a whole number!"},
		{"quit", "q", ""},
		{"QUIT", "Q", ""},
	}

	for _, e := range tests {
		input := strings.NewReader(e.input)
		scanner := bufio.NewScanner(input)
		res, _ := checkNumbers(scanner)

		if !strings.EqualFold(res, e.expected) {
			t.Errorf("%s: expected %s but got %s", e.name, e.expected, res)
		}
	}

}

func Test_readUserInput(t *testing.T) {
	// to test this function, we need a channel, and a  instance of an io.Reader
	doneChan := make(chan bool)

	// create a reference to bytes.Buffer
	var stdin bytes.Buffer

	stdin.Write([]byte("1\nq\n"))

	go readUserInput(&stdin, doneChan)
	<-doneChan
	close(doneChan)
}
