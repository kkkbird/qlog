FROM elastic/filebeat:7.3.2

ADD --chown=root:filebeat . ./

USER root
RUN chmod go-w ./filebeat.yml ./modules.d/qlog.yml ./module/qlog/filelog/manifest.yml ./module/qlog/udplog/manifest.yml
USER filebeat

EXPOSE 6060/udp

CMD [ "filebeat", "-e"]