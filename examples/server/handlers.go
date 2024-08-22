package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	cstr "git.cmcode.dev/cmcode/castopod-sub-token-retriever/pkg/cstr"
	castopod "git.cmcode.dev/cmcode/go-castopod/pkg/lib"
)

type resetCode struct {
	Code    string
	Expires time.Time
}

// Stores access codes that prove the user owns the email. Make sure to access
// this with a sync mutex lock. Email is the key, code is the value.
var codes map[string]resetCode = map[string]resetCode{}

// Make sure to lock this upon application startup to avoid panics!
var codesLock sync.Mutex = sync.Mutex{}

// Receives and routes requests.
//
// Wrap router like this so that it can be used with other middleware:
//
//	http.HandlerFunc(router)
func router(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case fmt.Sprintf("%v/token", conf.BaseURLRoute):
		getNewToken(w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

// clearCode removes the code from the codes map.
func clearCode(email string) {
	codesLock.Lock()
	defer codesLock.Unlock()
	delete(codes, email)
}

// addCode adds the code to the codes map.
func addCode(email, code string) {
	codesLock.Lock()
	defer codesLock.Unlock()
	codes[email] = resetCode{Code: code, Expires: time.Now().Add(15 * time.Minute)}
}

// getCode gets the code from the codes map.
func getCode(email string) string {
	codesLock.Lock()
	defer codesLock.Unlock()
	if time.Now().After(codes[email].Expires) {
		delete(codes, email)
		return ""
	}
	return codes[email].Code
}

// Periodically clears out expired codes once per minute to save memory. Run
// this as a goroutine.
func clearCodes() {
	for {
		codesLock.Lock()
		toDelete := []string{}
		for email, code := range codes {
			if time.Now().After(code.Expires) {
				toDelete = append(toDelete, email)
			}
		}

		for _, email := range toDelete {
			delete(codes, email)
		}

		codesLock.Unlock()

		l := len(toDelete)
		if l > 0 {
			log.Printf("cleared %v codes", l)
		}

		time.Sleep(1 * time.Minute)
	}
}

const genericResponse = "If this email address is valid, an email containing a new podcast URL will be sent shortly."

var genericResponseB = []byte(genericResponse)

// getNewToken queries the database to find the given email, resets the token
// stored in the database, then emails the new token to the user.
func getNewToken(w http.ResponseWriter, r *http.Request) {
	if DB == nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte("Unable to reach backend"))
		if err != nil {
			log.Printf("failed to write nil db error: %v", err.Error())
		}
		return
	}

	q := r.URL.Query()
	email := q.Get("email")
	if email == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("Email must not be empty"))
		if err != nil {
			log.Printf("failed to write http bad req error to user, error: %v", err.Error())
		}
		return
	}

	handle := q.Get("handle")
	if handle == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("handle must not be empty"))
		if err != nil {
			log.Printf("failed to write http bad req error to user, error: %v", err.Error())
		}
		return
	}

	st, err := DB.Prepare(fmt.Sprintf("%v WHERE sub.status = ? AND sub.email = ? AND podcast.handle = ?", cstr.CASTOPOD_SUBSCRIPTIONS_QUERY))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte("Unable to prepare backend"))
		if err != nil {
			log.Printf("failed to prepare db in request, error: %v", err.Error())
		}
		return
	}
	defer st.Close()

	rows, err := st.Query(castopod.CastopodStatusActive, email, handle)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte("Unable to query backend"))
		if err != nil {
			log.Printf("failed to query db in request, error: %v", err.Error())
		}
		return
	}

	cs := []cstr.CastopodSubscription{}
	for rows.Next() {
		css, err := conf.GetCastopodSubscription(rows)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err := w.Write([]byte("Unable to process data from backend"))
			if err != nil {
				log.Printf("failed to process rows from db query in request, error: %v", err.Error())
			}
			return
		}

		cs = append(cs, css)
	}

	lcs := len(cs)

	if lcs == 0 {
		log.Printf("email=%v, handle=%v returned no results", email, handle)

		w.WriteHeader(http.StatusOK)
		_, err := w.Write(genericResponseB)
		if err != nil {
			log.Printf("failed to send zero-data back to user, error: %v", err.Error())
		}

		return
	}

	if lcs != 1 {
		log.Printf("got multiple results for %v, this shouldn't happen. Picking the first result.", email)
	}

	log.Printf("email=%v, handle=%v returned %v results", email, handle, lcs)

	code := q.Get("code")

	storedCode := getCode(email)
	if code == "" || storedCode == "" {
		ncb, err := castopod.RandomAlphanumeric(16)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err := w.Write([]byte("Server is unable to process the request at this time. Please try later or contact support."))
			if err != nil {
				log.Printf("failed to get random alphanumeric code for %v, error: %v", email, err.Error())
			}

			return
		}

		newCode := string(ncb)

		addCode(email, newCode)

		log.Printf("sending reset request email to %v", email)

		// send an email to ask the user to verify the code
		err = conf.SMTP.SendEmail(
			[]string{email},
			"Podcast Feed Reset Request",
			fmt.Sprintf(
				"Please use the following link to reset your podcast feed.\r\n\r\n%v?email=%v&handle=%v&code=%v\r\n\r\nOnce you visit the above URL, an email will be sent shortly after with your new feed URL.",
				conf.URL,
				email,
				handle,
				newCode,
			),
			conf.EmailFrom,
		)
		if err != nil {
			log.Printf("failed to send email: %v", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			_, err := w.Write([]byte("Failed to send an email to the provided address. Please try later or contact support."))
			if err != nil {
				log.Printf("failed to send email with reset URL to %v, error: %v", email, err.Error())
			}

			return
		}

		log.Printf("sent reset request email to %v", email)

		w.WriteHeader(http.StatusOK)
		_, err = w.Write(genericResponseB)
		if err != nil {
			log.Printf("failed to send reset request confirmation data back to user, error: %v", err.Error())
		}

		return
	}

	if storedCode != code && storedCode != "" {
		log.Printf("email %v entered invalid code", email)
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, err := w.Write([]byte("Code is invalid"))
		if err != nil {
			log.Printf("failed to write http invalid code error to user %v, error: %v", email, err.Error())
		}
		return
	}

	defer clearCode(email) // remove the code, since it was a successful use

	// generate a new token, write it to the db, and send it to the user via email

	newSecret, newToken := castopod.NewToken()

	newQuery := "UPDATE cp_subscriptions SET token = ? WHERE id = ?"

	if !flagTest {
		log.Printf("updating db row for email=%v, handle=%v, id=%v...", email, handle, cs[0].SubscriptionID)

		ust, err := DB.Prepare(newQuery)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err := w.Write([]byte("Unable to prepare backend for update"))
			if err != nil {
				log.Printf("failed to prepare db for update in request, error: %v", err.Error())
			}
			return
		}
		defer ust.Close()

		_, err = ust.Query(newToken, cs[0].SubscriptionID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err := w.Write([]byte("Unable to update backend"))
			if err != nil {
				log.Printf("failed to update db in request, error: %v", err.Error())
			}
			return
		}
	}

	newFeedURL := cstr.GetNewTokenURL(conf.CastopodBaseURL, handle, newSecret)
	newFeedEmail := cstr.GetNewTokenEmail(newFeedURL)

	log.Printf("email=%v, handle=%v success: %v (token=%v, id=%v)", email, handle, newQuery, newToken, cs[0].SubscriptionID)

	err = conf.SMTP.SendEmail(
		[]string{email},
		"Podcast Feed Reset Successfully",
		newFeedEmail,
		conf.EmailFrom,
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte("Failed to send an email to the provided address. Please try later or contact support."))
		if err != nil {
			log.Printf("failed to send email with reset URL to %v, error: %v", email, err.Error())
		}

		return
	}

	_, err = w.Write([]byte("An email will be sent to you shortly containing your new podcast URL."))
	if err != nil {
		log.Printf("getNewToken failed to reply to user: %v", err.Error())
	}
}
