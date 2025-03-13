FROM redgate/flyway
COPY flyway.conf /flyway/conf
COPY /migration /flyway/sql
CMD ["migrate"]