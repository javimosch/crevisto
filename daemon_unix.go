//go:build !windows

package main

import "syscall"

func setSysProcAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{Setsid: true}
}
