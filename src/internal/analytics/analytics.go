package analytics

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/nleeper/goment"
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
	return
}

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
