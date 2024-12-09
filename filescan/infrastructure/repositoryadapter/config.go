package repositoryadapter

// Tables is a struct that represents tables of computility.
type Tables struct {
	FileScan      string `json:"file_scan"                 required:"true"`
	LargeFileScan string `json:"large_file_scan"           required:"true"`
}
