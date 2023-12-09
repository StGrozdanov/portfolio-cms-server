package files

import (
	"bytes"
	"mime/multipart"
	"portfolio-cms-server/database"
	"portfolio-cms-server/utils"
)

// UploadCV takes a form data file, processes it and transforms it to a bytes reader with a file key and
// content type of application/pdf and uploads it to the s3 bucket. When the file is uploaded it uses a
// database function to determine if the database cv link should be updated or not and updates it if needed.
func UploadCV(file *multipart.FileHeader) (fileURL string, err error) {
	fileKey := "cv"
	contentType := "application/pdf"

	fileContent, _ := file.Open()
	buffer := make([]byte, file.Size)
	_, _ = fileContent.Read(buffer)
	fileBytes := bytes.NewReader(buffer)

	err = utils.UploadToS3(fileBytes, fileKey, contentType)
	if err != nil {
		return
	}

	fileURL = utils.GetTheFullS3BucketURL() + "/" + fileKey

	_, err = database.ExecuteNamedQuery(
		`SELECT set_user_cv_link_if_not_already_exist( :user_id, :URL )`,
		map[string]interface{}{"URL": fileURL, "user_id": 1},
	)

	if err != nil {
		return
	}
	return
}
