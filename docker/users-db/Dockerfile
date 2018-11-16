FROM mongo:3.6
LABEL maintainer=<tim.curless@thinkahead.com>

ADD ./scripts /tmp/scripts

RUN mkdir -p /data/db-users \
    && echo "dbpath = /data/db-users" > /etc/mongodb.conf \
    && chown -R mongodb:mongodb /data/db-users

RUN su - mongodb && mongod --fork --logpath /var/log/mongodb.log --dbpath /data/db-users \
    && /tmp/scripts/mongo_create_insert.sh \
    && mongod --dbpath /data/db-users --shutdown \
    && chown -R mongodb /data/db-users

VOLUME /data/db-users

CMD ["mongod", "--config", "/etc/mongodb.conf", "--smallfiles"]
