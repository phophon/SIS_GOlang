package callback

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	_ "github.com/lib/pq"
)

const (
	host     = "omega-postgresql-sgp1-08776-do-user-4090996-0.db.ondigitalocean.com"
	port     = 25061
	user     = "omega_rew"
	password = "c6eqgnwwv09cxlzo"
	dbname   = "TestPool"
	sslmode  = "require"
)

// var (
// 	googleOauthConfig = &oauth2.Config{
// 		RedirectURL:  "https://omega-next.cmkl.ac.th/api/v1/GoogleCallback",
// 		ClientID:     "346969593881-rhso1lgkgg6n5fgmqm05odobpemtsjae.apps.googleusercontent.com",
// 		ClientSecret: "0QKs98ImeI4FyX_3_VURaQXu",
// 		Scopes: []string{"https://www.googleapis.com/auth/userinfo.profile",
// 			"https://www.googleapis.com/auth/userinfo.email"},
// 		Endpoint: google.Endpoint,
// 	}
// 	// Some random string, random for each request
// 	oauthStateString = "random"
// )

type customClaims struct {
	FristName string
	LastName  string
	CmklMail  string
	jwt.StandardClaims
}

type Token struct {
	Data struct {
		AccessToken string
	}
}

type JWT struct {
	Key string
}

func CallbackHandler(w http.ResponseWriter, r *http.Request) {

	// b, err := ioutil.ReadAll(r.Body)
	// defer r.Body.Close()
	// if err != nil {
	// 	http.Error(w, err.Error(), 500)
	// 	return
	// }
	// fmt.Println("RequstBodyFormat ===== ", string(b))

	var t Token
	json.NewDecoder(r.Body).Decode(&t)

	fmt.Println("RequstBody ===== ", t.Data.AccessToken)

	var cmkl_email string
	var first_name string
	var last_name string
	var profile map[string]interface{}
	var token JWT

	// state := r.FormValue("state")
	// if state != oauthStateString {
	// 	fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
	// 	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	// 	return
	// }

	// code := r.FormValue("code")
	// token, err := googleOauthConfig.Exchange(oauth2.NoContext, code)
	// if err != nil {
	// 	fmt.Println("Code exchange failed with '%s'\n", err)
	// 	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	// 	return
	// }

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + t.Data.AccessToken)
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := json.Unmarshal([]byte(contents), &profile); err != nil {
		fmt.Println("ugh: ", err)
	}
	fmt.Println("profile: ", profile)
	fmt.Println("picture: ", profile["picture"].(string))
	mail := profile["email"].(string)

	// Qurry
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbname, sslmode)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	result, err := db.Query(`SELECT first_name, last_name, cmkl_email FROM student WHERE cmkl_email = $1;`, mail)
	if err != nil {
		panic(err)
		log.Fatal(err)
	}

	for result.Next() {
		if err := result.Scan(&first_name, &last_name, &cmkl_email); err != nil {
			log.Fatal(err)
		}
	}
	sqlStatement := `UPDATE student SET photo = $1 WHERE cmkl_email = $2`
	_, err = db.Exec(sqlStatement, profile["picture"].(string), mail)
	if err != nil {
		panic(err)
	}

	fmt.Println(cmkl_email)
	fmt.Println("passed Callback")

	claims := customClaims{
		FristName: first_name,
		LastName:  last_name,
		CmklMail:  cmkl_email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: 15000,
			Issuer:    "nameOfWebsiteHere",
		},
	}
	fmt.Println(claims)

	jwttoken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := jwttoken.SignedString([]byte("secureSecretText"))
	if err != nil {
		fmt.Println(err)
		return
	}

	bearer := "Bearer " + signedToken
	token.Key = signedToken
	fmt.Println(signedToken)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(token)
	w.Header().Add("Authorization", bearer)
	// contents, err := ioutil.ReadAll(response.Body)
	// fmt.Fprintf(w, "%s", contents)
	// http.Redirect(w, r, "/account", http.StatusSeeOther)
}
