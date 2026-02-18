package dto

type FileRef struct {
	Path string `json:"path"`
	Mime string `json:"mime"`
	Size int64  `json:"size"`
}
