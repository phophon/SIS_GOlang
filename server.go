package main
import (
   "fmt"
   "net/http"
   "encoding/json"
   "io/ioutil"
   "errors"
//    "strings"

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


func homePage(w http.ResponseWriter, r *http.Request) {
    url := "http://localhost:3000//api/profile"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("authorization", "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6IlFqVTNSa1kwTVVWR1F6WkNRVU0xTXpFMFJFUkZORUpCTWtaRlJqRTRNRGc0TWtRNU1VVTFNdyJ9.eyJpc3MiOiJodHRwczovL2Nta2wtb21lZ2EuYXV0aDAuY29tLyIsInN1YiI6IjNUbTZJQWFoeENiYm5ydlNWN1E0OFdld3JnY3FzWVJPQGNsaWVudHMiLCJhdWQiOiJodHRwczovL29tZWdhLW5leHQuY21rbC5hYy50aC8iLCJpYXQiOjE1OTc2NjIzMDMsImV4cCI6MTU5Nzc0ODcwMywiYXpwIjoiM1RtNklBYWh4Q2JibnJ2U1Y3UTQ4V2V3cmdjcXNZUk8iLCJndHkiOiJjbGllbnQtY3JlZGVudGlhbHMifQ.g_C7jEzhY3dSBppAOxU3a9niBuXJ18UJbyxODqQDNovTsI4IIQVHx_Ph2sTdnj-HPHAcWPigS-mTXVFtZI336SOBnCU39etisidPju9ommD-_ntvvPzfsE7DL_vEQjSBJbhb9WgesOFBou38k_qosru1MBG9GiTp2Vm2SIhomyZeCuYUtDtxcQMiSal1xzPctK8Foj5ipd1Hr74O3gkDrkvIWXVhaitWyezMduuoM7n1ikxffAGqqUVOSec-JJW4fkDmMFI6kB1rNpO6Dwk4FxRqh_kN116TLiX33pKx8UVPOMX0-CyKT1EDPHf4LNLylqpGLoJjapOaapMFlXZElg")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
  
    fmt.Println(res)
    fmt.Println(string(body))
    fmt.Fprint(w, string(body))
 }

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

   r.HandleFunc("/home", homePage).Methods("GET")
   r.HandleFunc("/callback", callback.CallbackHandler).Methods("GET")
   r.HandleFunc("/login", login.LoginHandler).Methods("GET")
   r.Handle("/api/v1/profile", api.ProfileApiHandler).Methods("GET")
   r.Handle("/api/enroll", jwtMiddleware.Handler(api.EnrollmentApiHandler)).Methods("GET")
   r.Handle("/api/v1/profile", jwtMiddleware.Handler(api.UpdateProfileHandler)).Methods("POST")
   r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

   corsWrapper := cors.New(cors.Options{
      AllowedMethods: []string{"GET", "POST"},
      AllowedHeaders: []string{"Content-Type", "Origin", "Accept", "*"},
   })

   http.ListenAndServe("0.0.0.0:8910", corsWrapper.Handler(r))
}