# Main Makefile for ssllabs
#
# Copyright 2018 © by Ollivier Robert <roberto@keltia.net>
#

.PATH=	cmd/ssllabs:.
GOBIN=	${GOPATH}/bin

GO=		go
GSRCS=	cmd/ssllabs/main.go
SRCS=	ssllabs.go ssllabs_subr.go types.go utils.go

BIN=	ssllabs
EXE=	${BIN}.exe

OPTS=	-ldflags="-s -w" -v

all: ${BIN}

${BIN}: ${GSRCS} ${SRCS} ${USRCS}
	${GO} build ${OPTS} ./cmd/...

${EXE}: ${GSRCS} ${SRCS} ${USRCS}
	GOOS=windows ${GO} build ${OPTS} ./cmd/...

build: ${SRCS} ${USRCS}
	${GO} build ${OPTS} ./cmd/...

test: build
	${GO} test ./...

windows: ${EXE}
	GOOS=windows ${GO} build ${OPTS} ./cmd/...

install:
	${GO} install ${OPTS} ./cmd/...

lint:
	gometalinter .

clean:
	${GO} clean .
	${GO} clean ./cmd/...
	-rm -f ${BIN} ${EXE}

push:
	git push --all
	git push --tags
