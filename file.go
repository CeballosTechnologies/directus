package directus

import "time"

type File struct {
	Charset          interface{} `json:"charset,omitempty"`
	Description      string      `json:"description,omitempty"`
	Duration         interface{} `json:"duration,omitempty"`
	Embed            interface{} `json:"embed,omitempty"`
	FilenameDisk     string      `json:"filename_disk,omitempty"`
	FilenameDownload string      `json:"filename_download,omitempty"`
	Filesize         string      `json:"filesize,omitempty"`
	Folder           string      `json:"folder,omitempty"`
	Height           string      `json:"height,omitempty"`
	Id               string      `json:"id,omitempty"`
	Location         string      `json:"location,omitempty"`
	Metadata         string      `json:"metadata,omitempty"`
	ModifiedBy       interface{} `json:"modified_by,omitempty"`
	ModifiedOn       time.Time   `json:"modified_on,omitempty"`
	Storage          string      `json:"storage,omitempty"`
	Tags             interface{} `json:"tags,omitempty"`
	Title            string      `json:"title,omitempty"`
	Type             string      `json:"type,omitempty"`
	UploadedBy       string      `json:"uploaded_by,omitempty"`
	UploadedOn       time.Time   `json:"uploaded_on,omitempty"`
	Width            interface{} `json:"width,omitempty"`
}
