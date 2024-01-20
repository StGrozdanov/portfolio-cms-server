package analytics

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/mssola/user_agent"
	"github.com/nleeper/goment"
	"github.com/oschwald/geoip2-golang"
	"net"
	"portfolio-cms-server/database"
	"strconv"
	"time"
)

var analyticTypes = map[string]func() (Analytics, error){
	"today":      getAnalyticsForTheDay,
	"yesterday":  getAnalyticsFromYesterday,
	"last7days":  getAnalyticsForTheLast7Days,
	"last30days": getAnalyticsForTheLast30Days,
	"last90days": getAnalyticsForTheLast90Days,
	"lastYear":   getAnalyticsForTheLastYear,
}

var allowedQuarterFormats = map[int]interface{}{
	1: nil,
	2: nil,
	3: nil,
	4: nil,
}

// Get gets analytics depending on the provided query parameters
func Get(parameter gin.Param) (Analytics, error) {
	if mapFunction, parameterIsFound := analyticTypes[parameter.Key]; parameterIsFound {
		results, err := mapFunction()
		return normaliseOutput(results, err)
	} else if parameter.Key == "quarter" {
		quarterFormatError := errors.New("the provided quarter param should be in format quarter=number where number is from 1 to 4")

		paramAsANumber, err := strconv.Atoi(parameter.Value)
		if err != nil {
			err = quarterFormatError
			return Analytics{}, err
		}

		if _, quarterValueIsValid := allowedQuarterFormats[paramAsANumber]; !quarterValueIsValid {
			err = quarterFormatError
			return Analytics{}, err
		}

		results, err := getAnalyticsForTheQuarter(paramAsANumber)
		return normaliseOutput(results, err)
	}
	return Analytics{}, errors.New("no param was found matching your criteria")
}

// Count counts the visitations for the day
func Count() (count int, err error) {
	dateNow, err := goment.New(time.Now())
	if err != nil {
		return
	}

	date := dateNow.UTC().Format("YYYY-MM-DD")

	err = database.GetSingleRecordNamedQuery(
		&count,
		`SELECT COALESCE(COUNT(id), 0) FROM analytics WHERE DATE(date_time) = :date;`,
		map[string]interface{}{"date": date},
	)
	return
}

// Track retrieves request information such as client ip, referer, browser, device and country and stores it in the
// database if not already written for the same day with the same ip (the visitations are counted as unique)
func Track(db *geoip2.Reader, ctx *gin.Context, deviceType string) (err error) {
	var (
		country     string
		countryCode string
		clientIP    = ctx.ClientIP()
		referer     = ctx.Request.Referer()
		userAgent   = user_agent.New(ctx.Request.UserAgent())
		browser, _  = userAgent.Browser()
	)

	record, err := db.Country(net.ParseIP(clientIP))

	if err == nil {
		country = record.Country.Names["en"]
		countryCode = record.Country.IsoCode
	} else {
		country = "unknown"
		countryCode = "unknown"
	}

	dateNow, err := goment.New(time.Now())
	if err != nil {
		return
	}

	date := dateNow.UTC().Format("YYYY-MM-DD HH:mm:ss")

	var exists bool

	err = database.GetSingleRecordNamedQuery(
		&exists,
		`SELECT EXISTS (SELECT id FROM analytics WHERE ip_address = :ip)`,
		map[string]interface{}{"date": date, "ip": clientIP},
	)

	if !exists {
		_, err = database.ExecuteNamedQuery(
			`INSERT INTO analytics (date_time, device_type, origin_country, ip_address, referer, browser, country_code)
				VALUES (:date, :device, :country, :ip, :referer, :browser, :country_code)`,
			map[string]interface{}{
				"date":         date,
				"device":       deviceType,
				"country":      country,
				"ip":           clientIP,
				"referer":      referer,
				"browser":      browser,
				"country_code": countryCode,
			})
	}
	return
}

// GetByCountry gets all analytics and groups them by country
func GetByCountry() (analyticsByCountry []ByCountry, err error) {
	err = database.GetMultipleRecords(
		&analyticsByCountry,
		`SELECT origin_country AS country,
					   country_code   AS code,
					   COUNT(id)      AS count
				FROM analytics
				GROUP BY origin_country, country_code
				ORDER BY count DESC;`,
	)
	return
}

// GetByBrowser gets all analytics and groups them by browser
func GetByBrowser() (analyticsByBrowser []ByBrowser, err error) {
	err = database.GetMultipleRecords(&analyticsByBrowser,
		`SELECT browser,
					   COUNT(id) AS count
				FROM analytics
				WHERE browser IS NOT NULL
				GROUP BY browser
				ORDER BY count DESC;`,
	)
	return
}

// GetByDevice gets all analytics and groups them by device
func GetByDevice() (analyticsByDevice []ByDevice, err error) {
	err = database.GetMultipleRecords(&analyticsByDevice,
		`SELECT device_type AS device,
					   COUNT(id)   AS count
				FROM analytics
				GROUP BY device_type
				ORDER BY count DESC;`,
	)
	return
}

func getAnalyticsForTheDay() (analyticsResponse Analytics, err error) {
	dateNow, err := goment.New(time.Now())
	if err != nil {
		return
	}
	date := dateNow.UTC().Format("YYYY-MM-DD")
	return getAnalyticsForTheDateQuery(date)
}

func getAnalyticsFromYesterday() (analyticsResponse Analytics, err error) {
	dateNow, err := goment.New(time.Now())
	if err != nil {
		return
	}
	date := dateNow.Subtract(1, "days").UTC().Format("YYYY-MM-DD")
	return getAnalyticsForTheDateQuery(date)
}

func getAnalyticsForTheLast7Days() (analyticsResponse Analytics, err error) {
	dateNow, err := goment.New(time.Now())
	if err != nil {
		return
	}

	startDate := dateNow.Subtract(7, "days").UTC().Format("YYYY-MM-DD")
	endDate := dateNow.Add(7, "days").UTC().Format("YYYY-MM-DD")

	return getAnalyticsBetweenTheDatesQuery(startDate, endDate)
}

func getAnalyticsForTheLast30Days() (analyticsResponse Analytics, err error) {
	dateNow, err := goment.New(time.Now())
	if err != nil {
		return
	}

	startDate := dateNow.Subtract(1, "month").UTC().Format("YYYY-MM-DD")
	endDate := dateNow.Add(1, "month").UTC().Format("YYYY-MM-DD")

	return getAnalyticsBetweenTheDatesQuery(startDate, endDate)
}

func getAnalyticsForTheLast90Days() (analyticsResponse Analytics, err error) {
	dateNow, err := goment.New(time.Now())
	if err != nil {
		return
	}

	startDate := dateNow.Subtract(3, "months").UTC().Format("YYYY-MM-DD")
	endDate := dateNow.Add(3, "months").UTC().Format("YYYY-MM-DD")

	return getAnalyticsBetweenTheDatesQuery(startDate, endDate)
}

func getAnalyticsForTheLastYear() (analyticsResponse Analytics, err error) {
	dateNow, err := goment.New(time.Now())
	if err != nil {
		return
	}

	startDate := dateNow.Subtract(1, "year").UTC().Format("YYYY-MM-DD")
	endDate := dateNow.Add(1, "year").UTC().Format("YYYY-MM-DD")

	return getAnalyticsBetweenTheDatesQuery(startDate, endDate)
}

func getAnalyticsForTheQuarter(quarter int) (analyticsResponse Analytics, err error) {
	dateNow, err := goment.New(time.Now())
	if err != nil {
		return
	}

	startDate := dateNow.SetQuarter(quarter).StartOf("quarter").UTC().Format("YYYY-MM-DD")
	endDate := dateNow.SetQuarter(quarter).EndOf("quarter").UTC().Format("YYYY-MM-DD")

	return getAnalyticsBetweenTheDatesQuery(startDate, endDate)
}

// nolint: funlen
func getAnalyticsForTheDateQuery(date string) (analyticsResults Analytics, err error) {
	err = database.GetMultipleRecordsNamedQuery(
		&analyticsResults.Results,
		`SELECT  date_time,
					   device_type,
					   origin_country,
					   ip_address
				FROM analytics
					WHERE DATE(date_time) = :date;`,
		map[string]interface{}{"date": date},
	)

	_ = database.GetSingleRecordNamedQuery(
		&analyticsResults.MostPopularCountry,
		`SELECT origin_country AS most_popular_country
            	FROM analytics
                	WHERE DATE(date_time) = :date
                GROUP BY origin_country
                ORDER BY COUNT(origin_country) DESC
                LIMIT 1`,
		map[string]interface{}{"date": date},
	)

	_ = database.GetSingleRecordNamedQuery(
		&analyticsResults.MostPopularDevice,
		`SELECT device_type AS most_popular_device
                FROM analytics
                	WHERE DATE(date_time) = :date
                GROUP BY device_type
                ORDER BY COUNT(device_type) DESC
                LIMIT 1`,
		map[string]interface{}{"date": date},
	)
	analyticsResults.TotalVisitationsCount = len(analyticsResults.Results)

	_ = database.GetMultipleRecordsNamedQuery(
		&analyticsResults.VisitationsByCountry,
		`SELECT origin_country AS country,
					   country_code   AS code,
					   COUNT(id)      AS count
				FROM analytics
				WHERE DATE(date_time) = :date
				GROUP BY origin_country, country_code
				ORDER BY count DESC;`,
		map[string]interface{}{"date": date},
	)

	_ = database.GetMultipleRecordsNamedQuery(
		&analyticsResults.VisitationsByBrowser,
		`SELECT browser,
					   COUNT(id) AS count
				FROM analytics
				WHERE browser IS NOT NULL AND DATE(date_time) = :date
				GROUP BY browser
				ORDER BY count DESC;`,
		map[string]interface{}{"date": date},
	)

	_ = database.GetMultipleRecordsNamedQuery(
		&analyticsResults.VisitationsByDevice,
		`SELECT device_type AS device,
					   COUNT(id)   AS count
				FROM analytics
				WHERE DATE(date_time) = :date
				GROUP BY device_type
				ORDER BY count DESC;`,
		map[string]interface{}{"date": date},
	)

	return
}

// nolint: funlen
func getAnalyticsBetweenTheDatesQuery(startDate, endDate string) (analyticsResults Analytics, err error) {
	err = database.GetMultipleRecordsNamedQuery(
		&analyticsResults.Results,
		`SELECT date_time,
					   device_type,
					   origin_country,
					   ip_address
				FROM analytics
					WHERE DATE(date_time) BETWEEN :startDate AND :endDate;`,
		map[string]interface{}{
			"startDate": startDate,
			"endDate":   endDate,
		},
	)

	_ = database.GetSingleRecordNamedQuery(
		&analyticsResults.MostPopularCountry,
		`SELECT origin_country
				FROM analytics
					WHERE DATE(date_time) BETWEEN :startDate AND :endDate
				GROUP BY origin_country
				ORDER BY COUNT(origin_country) DESC
				LIMIT 1`,
		map[string]interface{}{
			"startDate": startDate,
			"endDate":   endDate,
		},
	)

	_ = database.GetSingleRecordNamedQuery(
		&analyticsResults.MostPopularDevice,
		`SELECT device_type
				FROM analytics
					WHERE DATE(date_time) BETWEEN :startDate AND :endDate
				GROUP BY device_type
				ORDER BY COUNT(device_type) DESC
				LIMIT 1`,
		map[string]interface{}{
			"startDate": startDate,
			"endDate":   endDate,
		},
	)
	analyticsResults.TotalVisitationsCount = len(analyticsResults.Results)

	_ = database.GetMultipleRecordsNamedQuery(
		&analyticsResults.VisitationsByCountry,
		`SELECT origin_country AS country,
					   country_code   AS code,
					   COUNT(id)      AS count
				FROM analytics
				WHERE DATE(date_time) BETWEEN :startDate AND :endDate
				GROUP BY origin_country, country_code
				ORDER BY count DESC;`,
		map[string]interface{}{
			"startDate": startDate,
			"endDate":   endDate,
		},
	)

	_ = database.GetMultipleRecordsNamedQuery(
		&analyticsResults.VisitationsByBrowser,
		`SELECT browser,
					   COUNT(id) AS count
				FROM analytics
				WHERE browser IS NOT NULL AND DATE(date_time) BETWEEN :startDate AND :endDate
				GROUP BY browser
				ORDER BY count DESC;`,
		map[string]interface{}{
			"startDate": startDate,
			"endDate":   endDate,
		},
	)

	_ = database.GetMultipleRecordsNamedQuery(
		&analyticsResults.VisitationsByDevice,
		`SELECT device_type AS device,
					   COUNT(id)   AS count
				FROM analytics
				WHERE DATE(date_time) BETWEEN :startDate AND :endDate
				GROUP BY device_type
				ORDER BY count DESC;`,
		map[string]interface{}{
			"startDate": startDate,
			"endDate":   endDate,
		},
	)
	return
}

func normaliseOutput(analytics Analytics, err error) (Analytics, error) {
	if err != nil {
		return analytics, err
	}

	if analytics.Results == nil {
		analytics.Results = []Analytic{}
	}

	return analytics, err
}
