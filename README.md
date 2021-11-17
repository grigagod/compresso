# Compresso
This service is designed for storing and processing (downsampling, compressing and converting) video. It provides a RESTful API that allows users to upload videos, make requests(tickets) to process them, download original and processed videos.

## Local usage
Checkout ```Makefile``` for more information and usecases.
### Binary services:

```
source .env.local // configure (export environment vars)
make local.up // run local enviroment(Postgresql, RabbitMQ)
make SERVICE.run // clean + build + execute binary of SERVICE, availible options for SERVICE: auth, videoapi, videosvc
```

### Containerized services:
Services configuration must be placed in ```.env.dev``` in the project's root directory.
```
make dev.up // run local enviroment(Postgresql, RabbitMQ) and all services
make dev.down 
```

## Testing and linting
```
make unit // run unit tests
make integration // run integration tests, requires environment(AWS vars, RabbitMQ instance)
make lint // lint code, golangci-lint should be installed
```

## Endpoints
### Auth service API:
```
POST /register -- create new user
POST /login -- login with crea
```
### Video service API:
```
POST /videos - Upload new user's video
GET  /videos - Get all user's videos
GET  /videos/{id} - Get user's video with the given id
POST /tickets - Create new ticket
GET  /tickets - Get all user's tickets
GET  /tickets/{id} - Get ticket with the given id
```

## Documentation
Generated OpenAPI documentation can be found in the ```docs/``` folder.
### Swagger UI for auth service api:
http://localhost:5000/swagger/
### Swagger UI for video service api:
http://localhost:5001/swagger/
