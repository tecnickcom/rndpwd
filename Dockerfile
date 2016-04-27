FROM tecnickcom/crossdev
MAINTAINER info@tecnick.com
RUN mkdir -p /root/GO/src/github.com/tecnickcom/rndpwd
ADD ./ /root/GO/src/github.com/tecnickcom/rndpwd
WORKDIR /root/GO/src/github.com/tecnickcom/rndpwd
RUN make deps && make qa && make rpm && make deb && make bz2 && make crossbuild
