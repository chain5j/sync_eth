FROM alpine:edge

ENV TZ=Asia/Shanghai
RUN rm -rf /etc/localtime &&\
    ln -sv /usr/share/zoneinfo/$TZ /etc/localtime &&\
    echo '$TZ' >/etc/timezone

ENV APPDIR=/data
ENV TEMPDIR=/temp
RUN mkdir -p $APPDIR $APPDIR/conf $TEMPDIR

COPY . $TEMPDIR

RUN cp $TEMPDIR/sync_eth /usr/bin \
    && cp $TEMPDIR/conf/config.yaml $APPDIR/conf/

RUN rm -rf $TEMPDIR

WORKDIR /data
ENTRYPOINT ["sync_eth"]
