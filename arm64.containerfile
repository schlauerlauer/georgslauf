#
# BUILD
#

FROM --platform=$BUILDOS/$BUILDARCH docker.io/library/golang:1.24.1-bookworm AS build
RUN apt update && apt install -y gcc-aarch64-linux-gnu gcc-arm-linux-gnueabihf g++-arm-linux-gnueabihf ca-certificates

ARG USER=gl
ARG UID=6812

RUN adduser \
	--disabled-password \
	--gecos "" \
	--shell "/sbin/nologin" \
	--uid "${UID}" \
	"${USER}"

RUN mkdir /app && chown ${USER}:${USER} /app
WORKDIR /app

USER gl
COPY --chown=${USER}:${USER} . .
RUN ./build_arm.sh

#
# SCRATCH
#

FROM scratch

ARG UID=6812

LABEL maintainer="Janek Lauer <security@georgslauf.de>"

COPY --from=build /etc/ssl/certs /etc/ssl/certs
COPY --from=build --chown=${UID}:${UID} --chmod="700" /app/build/georgslauf /app/georgslauf

ENTRYPOINT [ "/app/georgslauf" ]
