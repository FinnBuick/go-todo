//go:build ci
// +build ci

package controller

import (
	"errors"
)

func copyToClipboard(text string) error {
	// No-op in CI environment
	return errors.New("clipboard not available in CI")
}

func readFromClipboard() (string, error) {
	return "", errors.New("clipboard not available in CI")
}

var clipboardAvailable = false
