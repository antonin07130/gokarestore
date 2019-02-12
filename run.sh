/bin/sh

docker build . -t gokarestore

docker run --rm gokarestore emitter

docker run --rm gokarestore processor

docker run --rm gokarestore emitter

# this one should fail if state is correctly restored
docker run --rm gokarestore processor

