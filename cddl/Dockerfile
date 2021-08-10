FROM ruby

RUN bundle config --global frozen 1

WORKDIR /usr/src/app

RUN gem install cddlc

ADD cddl-files /cddl-files

ADD processing_cddl.sh ./processing_cddl.sh

CMD ./processing_cddl.sh