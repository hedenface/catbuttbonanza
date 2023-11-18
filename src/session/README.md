# Cat Butt Bonanza > Session

This is the session manager for the UI and APIs.

Default port: 8081

```
make session
./session
```

or

```
make run
```

This will eventually have the option to use a database as storage, instead of only in-memory.


## API

### /session/create

Method `PUT`
JSON Body

Returns: Session object

### /session/read/<session id>

Method `GET`

Returns: Session object

### /session/update

Method `PATCH`
JSON Body

Returns: `{"message":"session updated"}`

### /session/delete/<session id>

Method `DELETE`

Returns: `{"message":"session deleted"}`
