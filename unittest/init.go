package unittest

import (
	"os"
	"rosella/setup"
)

func init() {
	os.Chdir("../")
	setup.Setup()
}
