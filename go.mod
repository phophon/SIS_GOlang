module 01-Login

go 1.12

require (
	api v0.0.0-00010101000000-000000000000
	app v0.0.0
	auth v0.0.0
	callback v0.0.0
	github.com/auth0/go-jwt-middleware v0.0.0-20200507191422-d30d7b9ece63

	github.com/codegangsta/negroni v1.0.0
	github.com/coreos/go-oidc v2.2.1+incompatible
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gorilla/handlers v1.4.2
	github.com/gorilla/mux v1.7.4
	github.com/lib/pq v1.8.0
	github.com/rs/cors v1.7.0
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
	home v0.0.0
	login v0.0.0
	logout v0.0.0
	middlewares v0.0.0
	templates v0.0.0
	user v0.0.0
)

replace app => ./app

replace auth => ./auth

replace callback => ./routes/callback

replace home => ./routes/home

replace login => ./routes/login

replace logout => ./routes/logout

replace middlewares => ./routes/middlewares

replace user => ./routes/user

replace templates => ./routes/templates

replace api => ./routes/api

replace handler => ./routes/handler