package files

type ImageDeleteRequestBody struct {
	ImageURL string `json:"imageURL" valid:"required"`
}
