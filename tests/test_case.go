package tests

import (
	"github.com/goravel/framework/testing"

	"players/bootstrap"
)

func init() {
	bootstrap.Boot()
}

type TestCase struct {
	testing.TestCase
}
