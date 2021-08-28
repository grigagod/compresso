# Usage

## Docker development usage:
```
make develop
```

## Local development usage:
```
make local // runs local enviroment
make auth.run // runs auth service
make local.stop // stops local enviroment
```

## Auth service API docs:

http://localhost:5000/swagger/

## Docker-compose files:
```
docker-compose.local.yml - run postgresql, migrate containers
docker-compose.dev.yml - run auth service, postgresql, migrate containers
```
