package main
import (
   "net/http"
   "encoding/json"
   "errors"

  "github.com/auth0/go-jwt-middleware"
  "github.com/dgrijalva/jwt-go"
  "github.com/gorilla/mux"
  "github.com/rs/cors"

   "login"
   "callback"
   "api"
)

type Response struct {
   Message string
}

type Jwks struct {
   Keys []JSONWebKeys
}

type JSONWebKeys struct {
   Kty string
   Kid string
   Use string
   N   string
   E   string
   X5c []string
}

type Profile struct{
	Name string
	Uuid int
}


var HomePage = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    var data Profile
    var emergencyContact EmergencyContact
    // data.Name = "Pinyarat"
    // data.Uuid = 109877189

    data.First_name = "Phophon"
    data.Last_name = "Insee"
    data.Program = "Electrical and Computer Engineering"
    data.Cmkl_email = "phophon.i@cmkl.ac.th"
    data.UUID = 109877190
    data.Photo = "https://lh3.googleusercontent.com/-QSuWiIaU_2o/AAAAAAAAAAI/AAAAAAAAAAA/AMZuuclw9aen0EfsX0OW-YHk3S0x6r8J8w/photo.jpg"
    data.Contact.Phone_number = "0965436527"
    data.Contact.Personnal_email = "phophon.i@cmkl.ac.th"
    data.Contact.Second_email = "phophon.i@cmkl.ac.th"
    data.Address.Addressstatus = "235/111 M.6 pruksa village"
    data.Address.City = "samutprakarn"
    data.Address.State = "bangmeaung"
    data.Address.Zip = "10270"
    data.Address.Country = "Thailand"

    emergencyContact.Firstname = ""
    emergencyContact.Lastname = ""
    emergencyContact.Relationship = ""
    emergencyContact.Phone = ""
    emergencyContact.Email = ""
    data.Emergency = append(data.Emergency, emergencyContact)

    emergencyContact.Firstname = "montean"
    emergencyContact.Lastname = "puengkeaw"
    emergencyContact.Relationship = "Dad"
    emergencyContact.Phone = "0952262928"
    emergencyContact.Email = "montean226@hotmail.com"
    data.Emergency = append(data.Emergency, emergencyContact)


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

 func getPemCert(token *jwt.Token) (string, error) {
   cert := ""
   resp, err := http.Get("https://cmkl-omega.auth0.com/.well-known/jwks.json")

   if err != nil {
       return cert, err
   }
   defer resp.Body.Close()

   var jwks = Jwks{}
   err = json.NewDecoder(resp.Body).Decode(&jwks)

   if err != nil {
       return cert, err
   }

   for k, _ := range jwks.Keys {
       if token.Header["kid"] == jwks.Keys[k].Kid {
           cert = "-----BEGIN CERTIFICATE-----\n" + jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
       }
   }

   if cert == "" {
       err := errors.New("Unable to find appropriate key.")
       return cert, err
   }

   return cert, nil
 }

func StartServer() {
   jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
      ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
         // Verify 'aud' claim
         aud := "https://omega-next.cmkl.ac.th/"
         checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)
         if !checkAud {
             return token, errors.New("Invalid audience.")
         }
         // Verify 'iss' claim
         iss := "https://cmkl-omega.auth0.com/"
         checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
         if !checkIss {
             return token, errors.New("Invalid issuer.")
         }
   
         cert, err := getPemCert(token)
         if err != nil {
             panic(err.Error())
         }
   
         result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
         return result, nil
       },
       SigningMethod: jwt.SigningMethodRS256,
})
   r := mux.NewRouter()

   r.HandleFunc("/api/v1/home", HomePage).Methods("POST")
   r.HandleFunc("/api/v1/callback", callback.CallbackHandler).Methods("GET")
   r.HandleFunc("/api/v1/login", login.LoginHandler).Methods("GET")
   r.Handle("/api/v1/profile", jwtMiddleware.Handler(api.UpdateProfileHandler)).Methods("POST")
   r.Handle("/api/v1/profile", jwtMiddleware.Handler(api.ProfileApiHandler)).Methods("GET")
   r.Handle("/api/v1/enroll", jwtMiddleware.Handler(api.EnrollmentApiHandler)).Methods("GET")
   r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

   corsWrapper := cors.New(cors.Options{
      AllowedMethods: []string{"GET", "POST"},
      AllowedHeaders: []string{"Content-Type", "Origin", "Accept", "*"},
   })

   http.ListenAndServe("0.0.0.0:8910", corsWrapper.Handler(r))
}