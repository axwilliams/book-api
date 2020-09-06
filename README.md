# Book API

## Docker

Launches the app and database containers. The database and tables will be created and seeded with example data, including an admin with the username `admin` and the password `Adminl#1`. These admin credentials should be passed as the Basic Auth Header to the `/users/token` endpoint in order to retrieve the login token. The login token should then be used as the Bearer Token for all other endpoints.

Make sure the `.env` file is loaded using `godotenv.Load()`.

The following commands must be run from inside the project root:

Run
```
docker-compose up
```

Run in background
```
docker-compose up -d
```

Stop 
```
docker-compose down
```

If you need to make changes to the app (e.g. to the `.env` file) you can stop the process using
```
docker-compose down --remove-orphans --volumes
```

... then rebuild using
```
docker-compose up --build
```

## Local

To run locally (without Docker):

1. Create a local Postgres database
2. Change the settings in the `.env` file to match the local database
3. Make sure the `.env` file is loaded as `godotenv.Load("../../.env")`
4. From insise the `cmd/books-api` directory run: `go run main.go`

The database tables will created and seeded with example data, including an admin with the username `admin` and the password `Adminl#1`. These admin credentials should be passed as the Basic Auth Header to the `/users/token` endpoint in order to retrieve the login token. The login token should then be used as the Bearer Token for all other endpoints.

## Endpoints

##### POST https://<i></i>localhost:8080/api/v1/users/token

Response:
```
HTTP/1.1 200 OK

{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6Ik..."
}
```

##### POST http://<i></i>localhost:8080/api/v1/books

Request:
```
{
    "isbn": "978-1234567891",
    "title": "Some Title",
    "author": "Some Author",
    "category": "Some Category"
}
```

Response:
```
HTTP/1.1 201 Created

{
  "id": "0296bc0e-75e4-43e5-9815-2933024d4aa7"
}
```

##### PATCH http://<i></i>localhost:8080/api/v1/books/{id}

Request:
```
{
    "isbn": "978-1234567892",
    "title": "Modified Title",
    "author": "Modified Author",
    "category": "Modified Category "
}
```

Response:
```
HTTP/1.1 200 OK
```

##### DELETE http://<i></i>localhost:8080/api/v1/books/{id}

Response:
```
HTTP/1.1 200 OK
```

##### GET http://<i></i>localhost:8080/api/v1/books/{id}

Response:
```
HTTP/1.1 200 OK

{
    "id": "0296bc0e-75e4-43e5-9815-2933024d4aa7",
    "isbn": "978-1234567891",
    "title": "Some Title",
    "author": "Some Author",
    "category": "Some Category",
}
```

##### GET http://<i></i>localhost:8080/api/v1/books

Response:
```
HTTP/1.1 200 OK

[
  {
      "id": "0296bc0e-75e4-43e5-9815-2933024d4aa7",
      "isbn": "978-1234567891",
      "title": "Some Title",
      "author": "Some Author",
      "category": "Some Category",
  },
  ...
]
```

##### GET http://<i></i>localhost:8080/api/v1/search/books

Parameters: 

`isbn` (`string`), `title` (`string`), `author` (`string`), `category` (`string`), `sort` (`string`), `order` (`string`), `limit`(`int`), `offset` (`int`).

Response:
```
HTTP/1.1 200 OK

[
  {
      "id": "0296bc0e-75e4-43e5-9815-2933024d4aa7",
      "isbn": "978-1234567891",
      "title": "Some Title",
      "author": "Some Author",
      "category": "Some Category",
  },
  ...
]
```

##### POST http://<i></i>localhost:8080/api/v1/users

Request:
```
{
    "username": "author",
    "email": "author@example.com",
    "password": "Author#1",
    "roles": ["AUTHOR"]
}
```

Response:
```
HTTP/1.1 201 Created

{
  "id": "2ee940ca-49a6-4281-972a-4da78ce7ce29"
}
```


##### PATCH http://<i></i>localhost:8080/api/v1/users/{id}

Request:
```
{
    "username": "author-modified",
    "email": "author-modified@example.com",
    "password": "Author#2",
    "password": ["AUTHOR", "ADMIN"]
}
```

Response:
```
HTTP/1.1 200 OK
```

##### DELETE http://<i></i>localhost:8080/api/v1/users/{id}

Response:
```
HTTP/1.1 200 OK
```

## Errors

Basic format:
```
HTTP/1.1 400 Bad Request

{
  "message": "ID is not in the correct form" 
}
```

Validation:
```
HTTP/1.1 422 Unprocessable Entity

{
  "message": "Validation failed",
  "errors": [
    {
      "username is a required filed",
      "email must be a valid email address"
    }
  ]
}
```