package database

import (
	"context"
	"github.com/cenkalti/backoff/v4"
	log "github.com/sirupsen/logrus"
	utils2 "portfolio-cms-server/utils"
	"time"
)

// CloseConnection use for scenario as sqlx ping function. The connection is not automatically closed after
// ping and it's a good idea to handle the case, especially because of the reconnect mech for loop - connections
// might add up.
func CloseConnection() {
	err := instance.DB.Close()
	if err != nil {
		utils2.GetLogger().WithFields(log.Fields{"error": err.Error()}).Error("Error on closing DB connection attempt")
	}
}

// Ping makes a simple ping with 3 second timeout
func Ping() error {
	var ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return instance.DB.PingContext(ctx)
}

// GetSingleRecordNamedQuery selects a single record from a named query and parses it to a destination
func GetSingleRecordNamedQuery(destination interface{}, query string, args interface{}) (err error) {
	err = backoff.Retry(func() error {
		var ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		namedStatement, err := instance.DB.PrepareNamed(query)
		if err != nil {
			return err
		}
		return namedStatement.Unsafe().GetContext(ctx, destination, args)
	}, utils2.RetryConfig())
	return
}

// GetMultipleRecords selects multiple records from the database
func GetMultipleRecords(destination interface{}, query string) (err error) {
	err = backoff.Retry(func() error {
		var ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		return instance.DB.Unsafe().SelectContext(ctx, destination, query)
	}, utils2.RetryConfig())
	return
}