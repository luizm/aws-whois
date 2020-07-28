//+build mage

package main

import (
	"fmt"

	"github.com/magefile/mage/sh"
)

func Build() error {
	if err := sh.Run("go", "mod", "download"); err != nil {
		return err
	}
	return sh.Run("go", "install", "./...")
}

func Fmt() error {
	return sh.Run("go", "fmt", ".")
}

func Test() error {
	out, err := sh.Output("go", "test", "-v", "./...")
	if err != nil {
		return err
	}

	fmt.Println(out)
	return nil
}
