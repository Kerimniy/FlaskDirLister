package services

import (
	"crypto/rand"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/securecookie"
	"golang.org/x/crypto/bcrypt"
)

var SECRET_KEY = make([]byte, 64)
var s = securecookie.New(SECRET_KEY, nil)

func InitSecretKey() {

	file, f_err := os.Open(".SECRET_KEY")
	if f_err != nil {
		_, e := rand.Read(SECRET_KEY)
		f, err := os.Create(".SECRET_KEY")
		_, e1 := f.Write(SECRET_KEY)
		if e != nil || err != nil || e1 != nil {
			log.Fatal(e)
		}

	} else {
		_, err2 := file.Read(SECRET_KEY)
		if err2 != nil {
			_, e := rand.Read(SECRET_KEY)
			f, err := os.Create(".SECRET_KEY")
			_, e1 := f.Write(SECRET_KEY)
			if e != nil || err != nil || e1 != nil {
				log.Fatal(e)
			}

		}
	}
	s = securecookie.New(SECRET_KEY, nil)
}

func checkPassword(hashedPassword []byte, password string) bool {
	err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	return err == nil
}

func hashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

}

func setSignedCookie(w http.ResponseWriter, value string) {
	encoded, _ := s.Encode("session", value)
	cookie := &http.Cookie{
		Name:        "session",
		Value:       encoded,
		Secure:      true,
		Partitioned: true,
		HttpOnly:    true,
		Path:        "/",
		SameSite:    http.SameSiteNoneMode,
		Expires:     time.Unix(time.Now().Unix()+31_536_000, 0),
	}
	http.SetCookie(w, cookie)
}

func deleteCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:        "session",
		Value:       "",
		Secure:      true,
		Partitioned: true,
		HttpOnly:    true,
		Path:        "/",
		SameSite:    http.SameSiteNoneMode,
		Expires:     time.Unix(0, 0),
	}

	http.SetCookie(w, cookie)

}

func getSignedCookie(r *http.Request, w http.ResponseWriter) string {
	if cookie, err := r.Cookie("session"); err == nil {

		var decoded string

		if err = s.Decode("session", cookie.Value, &decoded); err == nil {

			return decoded
		} else {
			deleteCookie(w)
		}
	}
	return ""
}
