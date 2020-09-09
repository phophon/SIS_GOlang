package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"api"
	"callback"
)

var (
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8910/api/v1/GoogleCallback",
		ClientID:     "346969593881-rhso1lgkgg6n5fgmqm05odobpemtsjae.apps.googleusercontent.com",
		ClientSecret: "0QKs98ImeI4FyX_3_VURaQXu",
		Scopes: []string{"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint: google.Endpoint,
	}
	// Some random string, random for each request
	oauthStateString = "random"
)

const htmlIndex = `<html><body>
<a href="/api/v1/login">Log in with Google</a>
</body></html>
`

type Profile struct {
	First_name string
	Last_name  string
	Program    string
	Cmkl_email string
	UUID       string
	Photo      string
	Contact    struct {
		Phone_number    string
		Personnal_email string
		Second_email    string
	}
	Emergency   [2]EmergencyContact
	Useraddress [3]Address
}

type Address struct {
	Addressstatus string
	City          string
	State         string
	Zip           string
	Country       string
}

type EmergencyContact struct {
	Firstname    string
	Lastname     string
	Relationship string
	Phone        string
	Email        string
}

var HomePage = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var data Profile
	var emergencyContact EmergencyContact
	var address Address
	// data.Name = "Pinyarat"
	// data.Uuid = 109877189

	data.First_name = "Phophon"
	data.Last_name = "Insee"
	data.Program = "Electrical and Computer Engineering"
	data.Cmkl_email = "phophon.i@cmkl.ac.th"
	data.UUID = "59f109be-c815-48d8-abee-80bf968add40"
	data.Photo = "https://lh3.googleusercontent.com/-QSuWiIaU_2o/AAAAAAAAAAI/AAAAAAAAAAA/AMZuuclw9aen0EfsX0OW-YHk3S0x6r8J8w/photo.jpg"
	data.Contact.Phone_number = "0965436527"
	data.Contact.Personnal_email = "phophon.i@cmkl.ac.th"
	data.Contact.Second_email = "phophon.i@cmkl.ac.th"

	address.Addressstatus = "235/111 M.6 pruksa village"
	address.City = "samutprakarn"
	address.State = "bangmeaung"
	address.Zip = "10270"
	address.Country = "Thailand"
	data.Useraddress[0] = address

	address.Addressstatus = "235/111 M.6 pruksa village"
	address.City = "samutprakarn"
	address.State = "bangmeaung"
	address.Zip = "10270"
	address.Country = "Thailand"
	data.Useraddress[1] = address

	address.Addressstatus = "235/111 M.6 pruksa village"
	address.City = "samutprakarn"
	address.State = "bangmeaung"
	address.Zip = "10270"
	address.Country = "Thailand"
	data.Useraddress[2] = address

	emergencyContact.Firstname = ""
	emergencyContact.Lastname = ""
	emergencyContact.Relationship = ""
	emergencyContact.Phone = ""
	emergencyContact.Email = ""
	data.Emergency[0] = emergencyContact

	emergencyContact.Firstname = "montean"
	emergencyContact.Lastname = "puengkeaw"
	emergencyContact.Relationship = "Dad"
	emergencyContact.Phone = "0952262928"
	emergencyContact.Email = "montean226@hotmail.com"
	data.Emergency[1] = emergencyContact

	// reqBody, err := json.Marshal(map[string]string{})

	// resp, err := http.Post("http://localhost:8910/api/v1/profile",
	// 	"application/json", bytes.NewBuffer(reqBody))
	// if err != nil {
	// 	print(err)
	// }

	// defer resp.Body.Close()
	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	print(err)
	// }
	// fmt.Println(string(body))
	// fmt.Fprint(w, string(body))

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(data)
})

func handleMain(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, htmlIndex)
}

func handleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func StartServer() {
	r := mux.NewRouter()

	r.HandleFunc("/", handleMain)
	r.HandleFunc("/api/v1/login", handleGoogleLogin)
	r.HandleFunc("/api/v1/GoogleCallback", callback.CallbackHandler)
	r.HandleFunc("/api/v1/home", HomePage).Methods("POST")
	r.HandleFunc("/api/v1/callback", callback.CallbackHandler).Methods("GET")
	r.HandleFunc("/api/v1/profile", api.ProfileApiHandler).Methods("GET")
	r.HandleFunc("/api/v1/profile", api.UpdateProfileHandler).Methods("POST")
	r.HandleFunc("/api/v1/enroll", api.EnrollmentApiHandler).Methods("GET")
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	corsWrapper := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"Content-Type", "Origin", "Accept", "*"},
	})

	http.ListenAndServe("0.0.0.0:8910", corsWrapper.Handler(r))
}
