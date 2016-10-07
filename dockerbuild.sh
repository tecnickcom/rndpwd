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
# This script requires Docker

# EXAMPLE USAGE:
# VENDOR=vendorname PROJECT=projectname ./dockerbuild.sh

# Get vendor and project name
: ${VENDOR:=vendor}
: ${PROJECT:=project}

# Name of the base development Docker image
DOCKERDEV=${VENDOR}/dev_${PROJECT}

# Build the base environment and keep it cached locally
docker build -t ${DOCKERDEV} ./resources/DockerDev/

# Define the project root path
PRJPATH=/root/src/${PROJECT}

# Generate a temporary Dockerfile to build and test the project
# NOTE: The exit status of the RUN command is stored to be returned later,
#       so in case of error we can continue without interrupting this script.
cat > Dockerfile <<- EOM
FROM ${DOCKERDEV}
RUN mkdir -p ${PRJPATH}
ADD ./ ${PRJPATH}
WORKDIR ${PRJPATH}
RUN make buildall || (echo \$? > target/buildall.exit)
EOM

# Define the temporary Docker image name
DOCKER_IMAGE_NAME=${VENDOR}/build_${PROJECT}

# Build the Docker image
docker build --no-cache -t ${DOCKER_IMAGE_NAME} .

# Start a container using the newly created Docker image
CONTAINER_ID=$(docker run -d ${DOCKER_IMAGE_NAME})

# Copy all build/test artifacts back to the host
docker cp ${CONTAINER_ID}:"${PRJPATH}/target" ./

# Remove the temporary container and image
docker rm -f ${CONTAINER_ID} || true
docker rmi -f ${DOCKER_IMAGE_NAME} || true
