package entities

type UploadJobMessage struct {
	UploadID  string `json:"upload_id"`
	PostID    string `json:"post_id"`
	FileName  string `json:"file_name"`
	FileType  string `json:"file_type"`
	FileData  []byte `json:"file_data"`
	RequestID string `json:"request_id"`
}
