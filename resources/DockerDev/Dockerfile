# Dockerfile
#
# Ubuntu 16.04 (Xenial) with cross-compiling development packages
# for C, C++ and GO languages
#
# @author      Nicola Asuni <info@tecnick.com>
# @copyright   2015-2016 Nicola Asuni - Tecnick.com LTD
# ------------------------------------------------------------------------------

# ------------------------------------------------------------------------------
# NOTES:
#
# This script requires Docker (https://www.docker.com/)
#
#     sudo apt-get install docker docker.io
#
# Add your user to the "docker" group:
#     sudo groupadd docker
#     sudo gpasswd -a <YOURUSER> docker
#     sudo service docker restart
#
# To create the container execute:
#     docker build --tag="tecnickcom/crossdev:latest" .
#
# To log into the newly created container:
#     docker run -t -i tecnickcom/crossdev /bin/bash
#
# To get the container ID:
#     CONTAINER_ID=`docker ps -a | grep tecnickcom/crossdev | cut -c1-12`
#
# To delete the newly created docker container:
#     docker rm -f $CONTAINER_ID
#
# To delete the docker image:
#     docker rmi -f tecnickcom/crossdev
#
# To delete all containers
#     docker rm $(docker ps -a -q)
#
# To delete all images
#     docker rmi $(docker images -q)
# ------------------------------------------------------------------------------

FROM ubuntu:xenial
MAINTAINER info@tecnick.com

ENV DEBIAN_FRONTEND noninteractive
ENV TERM linux
RUN echo 'debconf debconf/frontend select Noninteractive' | debconf-set-selections
ENV HOME /root
ENV DISPLAY :0
ENV GOPATH=/root
ENV PATH=$GOPATH/bin:$PATH

# add i386 architecture
RUN dpkg --add-architecture i386

# add repositories and update
RUN apt-get update && apt-get -y dist-upgrade && \
    apt-get install -y apt-utils software-properties-common python-software-properties && \
    apt-add-repository universe && \
    apt-add-repository multiverse && \
    apt-get update

# install development packages and debugging tools
RUN apt-get install -y wget nano vim curl rsync git screen build-essential python python-dev python-pip libffi-dev autotools-dev autoconf make cmake gcc g++ golang golang-golang-x-tools alien rpm dh-make dh-golang debhelper devscripts pbuilder strace lsof gdb valgrind pkg-config doxygen doxygen-latex lcov libc6 libc6-dev astyle pylint pyflakes flawfinder libboost-all-dev openssh-server openssh-client upx-ucl

# add mingw for cross compiling the Windows version in Linux
RUN apt-get install -y dos2unix mingw-w64 binutils-mingw-w64 mingw-w64-tools mingw-w64-i686-dev mingw-w64-x86-64-dev wine winetricks libwine-development ubuntu-restricted-extras ubuntu-restricted-addons wine64-development-tools binfmt-support

RUN pip install --upgrade pip && \
    pip install cffi && \
    pip install --upgrade autopep8

# configure GO
RUN mkdir -p /root/bin && \
    mkdir -p /root/pkg && \
    mkdir -p /root/src && \
    echo 'export GOPATH=/root' >> /root/.profile && \
echo 'export PATH=$GOPATH/bin:$PATH' >> /root/.profile
