package files

import (
	"bytes"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
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

// UploadProjectImage takes a form data file, processes it and transforms it to a bytes reader, generates a key
// and uploads the image to the s3 bucket. When the file is uploaded - inserts the image into the database and
// returns an array of existing images for the given project. (along with the newly created)
func UploadProjectImage(file *multipart.FileHeader, projectTitle string) (projectImages json.RawMessage, err error) {
	randomId, _ := uuid.NewRandom()

	fileKey := fmt.Sprintf("project-%s-%s", projectTitle, randomId.String())
	contentType := file.Header.Get("Content-Type")

	fileContent, _ := file.Open()
	buffer := make([]byte, file.Size)
	_, _ = fileContent.Read(buffer)
	fileBytes := bytes.NewReader(buffer)

	err = utils.UploadToS3(fileBytes, fileKey, contentType)
	if err != nil {
		return
	}

	imageURL := utils.GetTheFullS3BucketURL() + "/" + fileKey

	err = executeUploadProjectImageQuery(imageURL, projectTitle)
	if err != nil {
		return
	}

	err = database.GetSingleRecordNamedQuery(
		&projectImages,
		`SELECT CAST(arr.object AS JSONB) -> 'imgUrl' AS project_images
				FROM users,
					 JSONB_ARRAY_ELEMENTS(projects) WITH ORDINALITY arr(object)
				WHERE CAST(arr.object AS JSONB) ->> 'title' = :project_name;`,
		map[string]interface{}{"project_name": projectTitle},
	)

	return
}

func executeUploadProjectImageQuery(imageURL, projectTitle string) error {
	_, err := database.ExecuteNamedQuery(
		`UPDATE users
				SET projects = JSONB_BUILD_ARRAY(
						(SELECT arr.object
						 FROM users,
							  JSONB_ARRAY_ELEMENTS(projects) WITH ORDINALITY arr(object)
						 WHERE CAST(arr.object AS JSONB) ->> 'title' != :project_title),
						JSONB_SET(
								(SELECT arr.object
								 FROM users,
									  JSONB_ARRAY_ELEMENTS(projects) WITH ORDINALITY arr(object)
								 WHERE CAST(arr.object AS JSONB) ->> 'title' = :project_title),
								'{imgUrl}',
								JSONB_BUILD_ARRAY(
										(SELECT REPLACE(REPLACE(REPLACE(CAST(arr.object AS JSONB) ->> 'imgUrl', '[', ''), ']', ''), '"', '')
										 FROM users,
											  JSONB_ARRAY_ELEMENTS(projects) WITH ORDINALITY arr(object)
										 WHERE CAST(arr.object AS JSONB) ->> 'title' = :project_title),
										CAST(:img_url AS text)
								    )
						)
							   )`,
		map[string]interface{}{"project_title": projectTitle, "img_url": imageURL},
	)
	return err
}
