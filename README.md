# DT_echo_training
DeepTrack's Echo training sample.

## Install
Preparing Sources
```
cd $WORKDIR
git clone git@github.com:SeiyaNakamura/DT_echo_training.git
```

## Launching Applications
Running with Docker
```
cd $WORKDIR/DT_echo_training
docker-compose up -d
```

## Enter Docker containers
Please check <container_name> by `docker-compose ps` etc.
```
cd $WORKDIR/DT_echo_training
docker container exec -it <container_name> sh
```

## App configuration
- /app/entites
  - Data Structure for Business Rules
- /app/infrastructure
  - DB and framework definition
- /app/interfaces
  - Adapters that convert data to internal formats used by external form `usecase` and `entities`, or convert data from internal to convenient formats for external functions
- /app/usecase
  - Appropriation's inherent business rules
- /app/template
  - HTML file location

## Technologies used
- Language: Go v1.17.7
- Web Framework: Echo v4.1.6
- RDB: MySQL
- Hot reload tool: Air @latest
- Container technology: Docker
- Architecture: MVC
