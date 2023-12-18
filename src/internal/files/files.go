package files

import (
	"bytes"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"mime/multipart"
	"portfolio-cms-server/database"
	"portfolio-cms-server/utils"
	"strings"
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
	formattedImage := strings.ReplaceAll(imageURL, " ", "+")
	_, err := database.ExecuteNamedQuery(
		`WITH result_objects AS (SELECT arr.object AS object_result
                        FROM users,
                             JSONB_ARRAY_ELEMENTS(projects) WITH ORDINALITY arr(object)
                        WHERE CAST(arr.object AS JSONB) ->> 'title' != :project_title)
				UPDATE users
				SET projects = (
						(SELECT JSONB_AGG(object_result) FROM result_objects)
						||
						JSONB_SET(
								(SELECT arr.object
								 FROM users,
									  JSONB_ARRAY_ELEMENTS(projects) WITH ORDINALITY arr(object)
								 WHERE CAST(arr.object AS JSONB) ->> 'title' = :project_title),
								'{imgUrl}',
								JSONB_BUILD_ARRAY(
										(SELECT REPLACE(REPLACE(REPLACE(CAST(arr.object AS JSONB) ->> 'imgUrl', '[', ''), ']', ''), '"',
														'')
										 FROM users,
											  JSONB_ARRAY_ELEMENTS(projects) WITH ORDINALITY arr(object)
										 WHERE CAST(arr.object AS JSONB) ->> 'title' = :project_title),
										CAST(:img_url AS text)
								)
						)
					)`,
		map[string]interface{}{"project_title": projectTitle, "img_url": formattedImage},
	)
	return err
}

// UploadJobImage takes a form data file, processes it and transforms it to a bytes reader, generates a key
// and uploads the image to the s3 bucket. When the file is uploaded - inserts the image into the database and
// returns an array of existing images for the given job. (along with the newly created)
func UploadJobImage(file *multipart.FileHeader, company string) (jobImages json.RawMessage, err error) {
	randomId, _ := uuid.NewRandom()

	fileKey := fmt.Sprintf("job-%s-%s", company, randomId.String())
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

	err = executeUploadJobImageQuery(imageURL, company)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = database.GetSingleRecordNamedQuery(
		&jobImages,
		`SELECT CAST(arr.object AS JSONB) -> 'imgUrl' AS job_images
				FROM users,
					 JSONB_ARRAY_ELEMENTS(jobs) WITH ORDINALITY arr(object)
				WHERE CAST(arr.object AS JSONB)  ->> 'company' = :company_name;`,
		map[string]interface{}{"company_name": company},
	)

	return
}

func executeUploadJobImageQuery(imageURL, company string) error {
	formattedImage := strings.ReplaceAll(imageURL, " ", "+")

	_, err := database.ExecuteNamedQuery(
		`WITH result_objects AS (SELECT arr.object AS object_result
                        FROM users,
                             JSONB_ARRAY_ELEMENTS(jobs) WITH ORDINALITY arr(object)
                        WHERE CAST(arr.object AS JSONB) ->> 'company' != :company_name)
				UPDATE users
				SET jobs = (
						(SELECT JSONB_AGG(object_result) FROM result_objects)
						||
						JSONB_SET(
								(SELECT arr.object
								 FROM users,
									  JSONB_ARRAY_ELEMENTS(jobs) WITH ORDINALITY arr(object)
								 WHERE CAST(arr.object AS JSONB) ->> 'company' = :company_name),
								'{imgUrl}',
								JSONB_BUILD_ARRAY(
										(SELECT REPLACE(REPLACE(REPLACE(CAST(arr.object AS JSONB) ->> 'imgUrl', '[', ''), ']', ''),
														'"',
														'')
										 FROM users,
											  JSONB_ARRAY_ELEMENTS(jobs) WITH ORDINALITY arr(object)
										 WHERE CAST(arr.object AS JSONB) ->> 'company' = :company_name),
										CAST(:img_url AS TEXT))
						)
					);`,
		map[string]interface{}{"company_name": company, "img_url": formattedImage},
	)
	return err
}

// UploadPartnerImage takes a form data file, processes it and transforms it to a bytes reader, generates a key
// and uploads the image to the s3 bucket. When the file is uploaded - inserts the image into the database and
// returns an array of existing images for all partners. (along with the newly created)
func UploadPartnerImage(file *multipart.FileHeader) (partnerImages json.RawMessage, err error) {
	randomId, _ := uuid.NewRandom()

	fileKey := fmt.Sprintf("partner-%s", randomId.String())
	contentType := file.Header.Get("Content-Type")

	fileContent, _ := file.Open()
	buffer := make([]byte, file.Size)
	_, _ = fileContent.Read(buffer)
	fileBytes := bytes.NewReader(buffer)

	err = utils.UploadToS3(fileBytes, fileKey, contentType)
	if err != nil {
		return
	}

	fileBytes2 := bytes.NewReader(buffer)
	img, _, err := image.DecodeConfig(fileBytes2)
	if err != nil {
		return
	}

	imageURL := utils.GetTheFullS3BucketURL() + "/" + fileKey

	err = database.GetSingleRecordNamedQuery(
		&partnerImages,
		`UPDATE users
				SET partners = partners || (SELECT JSONB_BUILD_OBJECT(
														   'imgURL', CAST(:img_url AS TEXT),
														   'width', CAST(:width AS INT),
														   'height', CAST(:height AS INT)
												   ))
				RETURNING partners;`,
		map[string]interface{}{"img_url": imageURL, "width": img.Width, "height": img.Height},
	)

	return
}

// UploadCarouselImage takes a form data file, processes it and transforms it to a bytes reader, generates a key
// and uploads the image to the s3 bucket. When the file is uploaded - inserts the image into the database and
// returns an array of existing images for all carousels. (along with the newly created)
func UploadCarouselImage(file *multipart.FileHeader) (carouselImages []json.RawMessage, err error) {
	randomId, _ := uuid.NewRandom()

	fileKey := fmt.Sprintf("carousel-%s", randomId.String())
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

	JSONFriendlyImageURL := "\"" + imageURL + "\""

	_, err = database.ExecuteNamedQuery(
		`UPDATE users SET carousel = carousel || JSONB_BUILD_OBJECT('imgURL', CAST(:img_url AS JSONB));`,
		map[string]interface{}{"img_url": JSONFriendlyImageURL},
	)
	if err != nil {
		return
	}

	err = database.GetMultipleRecords(
		&carouselImages,
		`SELECT arr.object::JSONB -> 'imgURL' AS carousel_images
				FROM users,
					 JSONB_ARRAY_ELEMENTS(carousel) WITH ORDINALITY arr(object)
				WHERE id = 1;`,
	)

	return
}

// DeleteImage deletes an image from the s3 bucket
func DeleteImage(imageURL string) (err error) {
	err = utils.DeleteFromS3(imageURL)
	if err != nil {
		return
	}
	return
}
