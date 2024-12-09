package primitive

type FileType interface {
	Value() int
}

type fileType int

func (t fileType) Value() int {
	return int(t)
}

var (
	UnKnownFileType  fileType = 0
	ImageFileType    fileType = 1
	AudioFileType    fileType = 2
	VideoFileType    fileType = 3
	DocumentFileType fileType = 4
	MarkdownFileType fileType = 5
)
