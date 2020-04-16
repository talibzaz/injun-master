FROM iron/base

ENV ES_HOST 139.59.85.55
ENV ES_PORT 9200
ENV ES_INDEX events
ENV ES_INDEX_TYPE event
ENV HTTP_HOST :80
ENV ARANGO_HOST http://139.59.85.55:8529
ENV ARANGO_DB eventackle
ENV ARANGO_USERNAME root
ENV ARANGO_PASSWORD qF3mKQcu7zyzBYly
ENV GRPC_HOST :34567

ENV MAIL_HOST  smtp.gmail.com
ENV MAIL_PORT 587
ENV MAIL_USERNAME eventackleinfo@gmail.com
ENV MAIL_PASSWORD Matz.each.1

ENV EVENTACKLE https://eventackle.com
ENV SENDER  eventackleinfo@gmail.com
ENV EVENT_CREATION_GUIDE https://support.eventackle.com/email/eventcreationguide.pdf
ENV EVENT_PROMOTION_GUIDE http://support.eventackle.com/wp-content/uploads/2018/11/EventPromotionGuide.pdf
ENV PASSWORD_URL https://eventackle.com/reset-password/
ENV ORGANIZER_PROFILE_URL https://eventackle.com/organisation-profile-list
ENV ACCOUNT_SETTINGS_URL https://eventackle.com/account/profile
ENV MINIO https://minio.eventackle.com

ENV MYSQL_DSN root:root@tcp(159.65.153.232:3306)/eventackle

ENV ADDRESS1 Clarence Centre, 6 St. George Circus
ENV ADDRESS2 London
ENV ADDRESS3 United Kingdom SE1 6FE
ENV PHONE   +44 20 8242 6566
ENV EMAIL info@eventackle.com


EXPOSE 80
EXPOSE 34567
ADD injun-linux-amd64 /
CMD ["./injun-linux-amd64" ,"start"]