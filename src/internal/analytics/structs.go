package analytics

type Analytics struct {
	Results               []Analytic `db:"results" json:"results"`
	TotalVisitationsCount int        `db:"total_visitations_count" json:"totalVisitationsCount"`
	MostPopularCountry    string     `db:"most_popular_country" json:"mostPopularCountry"`
	MostPopularDevice     string     `db:"most_popular_device" json:"mostPopularDevice"`
}

type Analytic struct {
	Date          string `db:"date_time" json:"date"`
	DeviceType    string `db:"device_type" json:"deviceType"`
	OriginCountry string `db:"origin_country" json:"originCountry"`
	IpAddress     string `db:"ip_address" json:"ipAddress"`
}

type UserDevice struct {
	DeviceType string `db:"device_type" json:"deviceType" valid:"required"`
}
