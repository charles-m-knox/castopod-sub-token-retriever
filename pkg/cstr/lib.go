package cstr

import (
	"database/sql"
	"fmt"

	"git.cmcode.dev/cmcode/castopod-sub-token-retriever/pkg/smtp"
)

type Config struct {
	// Connection string for the Castopod mariadb database.
	SQLConnectionString string `json:"sqlConnectionString"`

	// SMTP configuration for sending emails.
	SMTP smtp.SMTPConfig `json:"smtp"`

	// For example: https://podcast.example.com
	CastopodBaseURL string `json:"castopodBaseUrl"`

	// Email address to send from.
	EmailFrom string `json:"emailFrom"`

	// The path to use for routing. All routes will be appended to it. For
	// example, /foo.
	BaseURLRoute string `json:"baseUrlRoute"`

	// This is the publicly accessible URL to use in emails and links for
	// interacting with this application. Include the full path, including
	// your [Config.BaseURLPath] value, as well as /token.
	//
	// A few examples (it really just depends on your reverse proxy):
	//
	// 	`https://example.com/podcast-reset/token`
	// 	`https://castopod-reset.example.com/token`
	URL string `json:"url"`

	// When the reset workflow finishes, this text will be shown to the user
	// and they will be instructed to click this URL to back to the home
	// page of your organization, for example.
	// RedirectURL string `json:"redirectUrl"`
}

// CastopodSubscription is a struct that (mostly) mirrors the SQL database's
// definition.
type CastopodSubscription struct {
	PodcastID      uint
	SubscriptionID uint
	Email          string // 255 chars max
	Token          string // hashed token
	Status         string // can only be active or suspended

	Secret string // un-hashed value of Token; not stored in the database
}

// Usage depends on your SQL backend - parameterize your queries:
//
// PostgreSQL, SQLite:
//
//	fmt.Sprintf("%v$1, $2, $3", CASTOPOD_SUBSCRIPTIONS_QUERY)
//
// MySQL, MariaDB, Oracle:
//
//	fmt.Sprintf("%v?, ?, ?", CASTOPOD_SUBSCRIPTIONS_QUERY)
//
// SQL Server:
//
//	fmt.Sprintf("%v@name, @name, @name", CASTOPOD_SUBSCRIPTIONS_QUERY)
//
// Then, use the value in the ... below:
//
//	stmt, err := db.Prepare(...)
//	defer stmt.Close()
//	rows, err := stmt.Query(val1, val2, val3)

// Use this as part of a parameterized db statement/query.
const CASTOPOD_SUBSCRIPTIONS_QUERY = "SELECT podcast.id as podcast_id, sub.id as sub_idl, podcast.handle, sub.email, sub.status FROM cp_subscriptions AS sub INNER JOIN cp_podcasts AS podcast ON podcast.id = sub.podcast_id" // WHERE sub.status = '" + castopod.CastopodStatusActive + "' AND sub.email = "

// GetNewTokenSQL returns the SQL statement that will update the database with
// a new token, provided you've already generated it.
func GetNewTokenSQL(podcastID uint, newToken string) string {
	return fmt.Sprintf("UPDATE cp_subscriptions SET token = '%v' WHERE id = %v;", newToken, podcastID)
}

// GetNewTokenURL returns the user's new podcast feed with the secret token. The
// handle must not have the @ character at the beginning.
func GetNewTokenURL(baseURL string, handle, secret string) string {
	return fmt.Sprintf("%v/@%v/feed.xml?token=%v", baseURL, handle, secret)
}

// GetNewTokenEmail returns the body content of an email that includes the
// user's new podcast feed with the secret token. Use the output of
// [GetNewTokenURL] for the first argument.
func GetNewTokenEmail(secretURL string) string {
	return fmt.Sprintf("Your new secret RSS feed is available at: %v \r\n\r\nDo not share this with anyone.\r\n\r\nAny of your old secret RSS feed URLs will no longer work.", secretURL)
}

// GetCastopodSubscription processes a row from a valid SQL query.
func (c *Config) GetCastopodSubscription(rows *sql.Rows) (CastopodSubscription, error) {
	var cs CastopodSubscription

	err := rows.Scan(&cs.PodcastID, &cs.SubscriptionID, &cs.Email, &cs.Token, &cs.Status)
	if err != nil {
		return cs, fmt.Errorf("failed to marshal row into interface: %v", err.Error())
	}

	return cs, nil
}
