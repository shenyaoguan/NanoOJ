NAME=nano-oj
BINDIR=build
VERSION=1.0.0
BUILDTIME=$(shell date -u)
GOBUILD=cd judge && go mod tidy && go build -ldflags '-s -w -X "main.version=$(VERSION)" -X "main.buildTime=$(BUILDTIME)"'
FRONTBUILD=cd web && pnpm i && vite build --outDir=../$(BINDIR)/static --emptyOutDir

.PHONY: web

all: linux-amd64 windows-amd64 web

web:
	$(FRONTBUILD)

web-dev:
	cd web && pnpm i && vite dev

linux-amd64: 
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o ../$(BINDIR)/$(NAME)-$@-$(VERSION)

windows-amd64: 
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o ../$(BINDIR)/$(NAME)-$@-$(VERSION).exe

clean:
	rm -rf $(BINDIR)/*
