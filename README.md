# DT_echo_training
DeepTrack's Echo training sample.

## Install
Preparing Sources
```
cd $WORKDIR
git@github.com:SeiyaNakamura/DT_echo_training.git
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
