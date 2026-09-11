//go:build windows

package main

import "syscall"

func setSysProcAttr() *syscall.SysProcAttr {
	return nil
}
