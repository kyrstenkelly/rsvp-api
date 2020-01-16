# Go RSVP API
Endpoint for creating RSVPs, guests and events written in Golang.

## Development

Create a `.env` file and fill in the missing values:
```bash
DB_HOST=
DB_PORT=5432
DB_USER=
DB_PASS=
DB_NAME=rsvps
DB_SSL_ENABLED=disable
AUTH_CLIENT_AUDIENCE=rsvps-api
AUTH_CLIENT_DOMAIN=https://jamesandkyrsten.auth0.com/
AUTH_CLIENT_SECRET=
```

Run with:
```
$ ./rsvp-api
```

## Documentation

### [API Docs](./docs/api.md)
### [Database Structure](./docs/database.md)

## Testing

TODO!

## Future Plans

* Learn about testing in Go :P
* Add tests!
* Clean up repetitive error handling code

   




