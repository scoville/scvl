scvl
===

scvl is a simple URL shortener written in go.

## Required

* Redis
* MySQL

## Setup

Place ``.env`` to root directory like below:

```
SESSION_SECRET=SUPER_SECRET_TOKEN
GOOGLE_CLIENT_ID=hogehoge.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=SUPER_SECRET
GOOGLE_REDIRECT_URL="http://127.0.0.1:8080/oauth/google/callback"
DB_URL="root:@/scvl_development?charset=utf8&parseTime=True&loc=Local"
ALLOWED_DOMAINS="sc0ville.com,en-courage.com,en-courage.net"
```

Website: [scvl.site](http://scvl.site)
