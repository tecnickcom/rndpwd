# Dockerfile
#
# Ubuntu 16.04 (Xenial) with cross-compiling development packages
# for C, C++ and GO languages
#
# @author      Nicola Asuni <info@tecnick.com>
# @copyright   Nicola Asuni - Tecnick.com LTD
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
ENV GOPATH=/root/GO
ENV PATH=$GOPATH/bin:$PATH

# add i386 architecture for wine
RUN dpkg --add-architecture i386

# add repositories and update
RUN apt-get update && apt-get -y dist-upgrade && \
    apt-get install -y apt-utils software-properties-common python-software-properties && \
    apt-add-repository multiverse && \
    apt-get update

# install development packages and debugging tools
RUN apt-get install -y wget nano vim curl rsync git screen python build-essential autotools-dev autoconf make cmake gcc g++ libboost-all-dev golang golang-golang-x-tools alien rpm dh-make dh-golang debhelper devscripts pbuilder strace lsof gdb valgrind pkg-config openssh-server openssh-client upx-ucl

RUN apt-get install -y binutils-mingw-w64 binutils-mingw-w64-i686 binutils-mingw-w64-x86-64 g++-mingw-w64 g++-mingw-w64-i686 g++-mingw-w64-x86-64 gcc-mingw-w64 gcc-mingw-w64-base gcc-mingw-w64-i686 gcc-mingw-w64-x86-64 gdb-mingw-w64 gdb-mingw-w64-target gfortran-mingw-w64 gnat-mingw-w64 gnat-mingw-w64-base gnat-mingw-w64-i686 gnat-mingw-w64-x86-64 gobjc++-mingw-w64 gobjc++-mingw-w64-i686 gobjc++-mingw-w64-x86-64 gobjc-mingw-w64 gobjc-mingw-w64-i686 gobjc-mingw-w64-x86-64 libz-mingw-w64 libz-mingw-w64-dev mingw-w64 mingw-w64-common mingw-w64-i686-dev mingw-w64-tools mingw-w64-x86-64-dev mingw-ocaml ocaml-mingw-w64 ocaml-mingw-w64-i686 ocaml-mingw-w64-x86-64 pkg-config dos2unix ubuntu-restricted-extras ubuntu-restricted-addons wine64-development-tools winetricks wine

# configure GO
RUN mkdir -p /root/GO/bin && \
    mkdir -p /root/GO/pkg && \
    mkdir -p /root/GO/src && \
    echo 'export GOPATH=/root/GO' >> /root/.profile && \
    echo 'export PATH=$GOPATH/bin:$PATH' >> /root/.profile
