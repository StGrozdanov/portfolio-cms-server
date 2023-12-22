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

type ByCountry struct {
	Country     string `db:"country" json:"country"`
	CountryCode string `db:"code" json:"code"`
	Count       int    `db:"count" json:"count"`
}

type ByBrowser struct {
	Browser string `db:"browser" json:"browser"`
	Count   int    `db:"count" json:"count"`
}

type ByDevice struct {
	Device string `db:"device" json:"device"`
	Count  int    `db:"count" json:"count"`
}
