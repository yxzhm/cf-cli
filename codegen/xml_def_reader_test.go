package codegen

import "testing"

func TestReadFile(t *testing.T) {
	ReadXmlDef("./xmldef.txt")
	ReadSingleXmlDef("./xmldef.txt", "cips.112.001.02")
}
