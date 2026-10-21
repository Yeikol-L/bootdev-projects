package main

import (
	"bytes"
	"fmt"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	type testCase struct {
		key       []byte
		plaintext []byte
		expected  string
	}

	runCases := []testCase{
		{[]byte("12344321"), []byte("Today I met my crush, what a hunk"), "Today I met my crush, what a hunk"},
		{[]byte("p@$$w0rd"), []byte("I hope my boyfriend never finds out about this"), "I hope my boyfriend never finds out about this"},
		{[]byte("secretky"), []byte("The best secret ever!"), "The best secret ever!"},
		{[]byte("longpass"), []byte("Just testing DES encryption with padding"), "Just testing DES encryption with padding"},
	}

	submitCases := append(runCases, []testCase{
		{[]byte("secretky"), []byte("The best secret ever!"), "The best secret ever!"},
		{[]byte("longpass"), []byte("Just testing DES encryption with padding"), "Just testing DES encryption with padding"},
	}...)

	testCases := runCases
	if withSubmit {
		testCases = submitCases
	}
	skipped := len(submitCases) - len(testCases)

	passed, failed := 0, 0

	for _, test := range testCases {
		ciphertext, err := encrypt(test.key, test.plaintext)
		if err != nil {
			t.Errorf(`---------------------------------
Encryption failed for key: %v
plaintext: %v
error: %v
Fail
`, string(test.key), string(test.plaintext), err)
			failed++
			continue
		}

		decryptedText, err := decrypt(test.key, ciphertext)
		if err != nil {
			t.Errorf(`---------------------------------
Decryption failed for key: %v
ciphertext: %v
error: %v
Fail
`, string(test.key), ciphertext, err)
			failed++
			continue
		}
		decryptedText = bytes.Trim(decryptedText, "\x00")

		if string(decryptedText) != test.expected {
			t.Errorf(`---------------------------------
Inputs:      key: %v, plaintext: %v
Expecting:   decrypted: %v
Actual:      decrypted: %v
Fail
`, string(test.key), string(test.plaintext), test.expected, string(decryptedText))
			failed++
		} else {
			fmt.Printf(`---------------------------------
Inputs:      key: %v, plaintext: %v
Expecting:   decrypted: %v
Actual:      decrypted: %v
Pass
`, string(test.key), string(test.plaintext), test.expected, string(decryptedText))
			passed++
		}
	}

	fmt.Println("---------------------------------")
	if skipped > 0 {
		fmt.Printf("%d passed, %d failed, %d skipped\n", passed, failed, skipped)
	} else {
		fmt.Printf("%d passed, %d failed\n", passed, failed)
	}
}

// withSubmit is set at compile time depending
// on which button is used to run the tests
var withSubmit = true
