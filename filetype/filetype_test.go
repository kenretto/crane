package filetype

import "testing"

func TestFileType(t *testing.T) {
	t.Log(FileType(FileBytes("testdata/index.html")))
	t.Log(FileType(FileBytes("testdata/file.txt")))
	t.Log(FileType(FileBytes("testdata/ini.ini")))
}
