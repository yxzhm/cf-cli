package codegen

import "testing"

func TestReadFile(t *testing.T) {
	ReadXmlDef("./xmldef.txt")
}
