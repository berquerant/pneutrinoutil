ARG MYSQL_VERSION=9.5.0

FROM mysql:$MYSQL_VERSION

COPY etc/my.cnf /etc/my.cnf
RUN mkdir /var/log/mysql
RUN chown mysql:mysql /var/log/mysql

CMD ["mysqld"]
