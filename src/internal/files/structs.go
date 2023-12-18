package files

type ImageDeleteRequestBody struct {
	ImageURL string `json:"imageURL" valid:"required"`
}

type CarouselImages struct {
	ImgURL string `json:"imgURL" db:"img_url"`
	Width  int    `json:"width" db:"width"`
	Height int    `json:"height" db:"height"`
}
