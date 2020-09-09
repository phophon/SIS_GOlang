module 01-Login/routes/callback

go 1.12

require (
	app v0.0.0
	auth v0.0.0
	github.com/coreos/go-oidc v2.1.0+incompatible
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/lib/pq v1.8.0
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
)

replace app => ../../app

replace auth => ../../auth
