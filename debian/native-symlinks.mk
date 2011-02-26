# Create architecture-independent toolchain names

# Reuse GOARCH detection logic
export GOROOT = .
include src/Make.inc

bindir = debian/golang-go/usr/bin
basenames = a c cov g l nm prof

.PHONY: all
all:
	for b in $(basenames) ; do \
	  ln -s $O$$b $(bindir)/golang-$$b; \
	done
