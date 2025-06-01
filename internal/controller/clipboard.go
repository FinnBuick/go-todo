//go:build !ci
// +build !ci

package controller

import (
	"golang.design/x/clipboard"
)

func copyToClipboard(text string) error {
	err := clipboard.Init()
	if err != nil {
		return err
	}
	clipboard.Write(clipboard.FmtText, []byte(text))
	return nil
}
