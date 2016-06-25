#!/bin/sh
#
# dockerbuild.sh
#
# Build the software inside a Docker container
#
# @author      Nicola Asuni <info@tecnick.com>
# @copyright   2015-2016 Nicola Asuni - Tecnick.com LTD
# ------------------------------------------------------------------------------

# NOTES:
#
# This script requires docker

# EXAMPLE USAGE:
# ./dockerbuild.sh

# build the environment
docker build -t tecnickcom/crossdev ./resources/DockerDev/

# generate a docker file on the fly
cat > Dockerfile <<- EOM
FROM tecnickcom/crossdev
MAINTAINER info@tecnick.com
RUN mkdir -p /root/GO/src/github.com/tecnickcom/rndpwd
ADD ./ /root/GO/src/github.com/tecnickcom/rndpwd
WORKDIR /root/GO/src/github.com/tecnickcom/rndpwd
RUN make deps && make qa && make rpm && make deb && make bz2 && make crossbuild
EOM

# docker image name
DOCKER_IMAGE_NAME="local/build"

# build the docker container and build the project
docker build --no-cache -t ${DOCKER_IMAGE_NAME} .

# start a container using the newly created docker image
CONTAINER_ID=$(docker run -d ${DOCKER_IMAGE_NAME})

# copy the artifact back to the host
docker cp ${CONTAINER_ID}:"/root/GO/src/github.com/tecnickcom/rndpwd/target" ./

# remove the container and image
docker rm -f ${CONTAINER_ID} || true
docker rmi -f ${DOCKER_IMAGE_NAME} || true
