#
#
#

SUBDIRS := \
						$(CURDIR)/cmd/anyclient/cmdline \
					  $(CURDIR)/cmd/anyserver/cmdline \
						$(CURDIR)/cmd/pktservice/cmdline 

#include $(addsuffix /makefile, $(SUBDIRS))

.PHONY: all
all:
	for dir in $(SUBDIRS); do \
		$(MAKE) -C $$dir;     \
	done

.PHONY: build
build:
	for dir in $(SUBDIRS); do \
		$(MAKE) build -C $$dir;     \
	done

.PHONY: build-win
build-win:
	for dir in $(SUBDIRS); do \
		$(MAKE) build-win -C $$dir;     \
	done

.PHONY: test
test:
	for dir in $(SUBDIRS); do \
		$(MAKE) test -C $$dir;     \
	done

.PHONY: clean
clean:
	for dir in $(SUBDIRS); do \
		$(MAKE) clean -C $$dir;     \
	done

.PHONY: check
check:
	for dir in $(SUBDIRS); do \
		$(MAKE) check -C $$dir;     \
	done

.PHONY: deploy-dev
deploy-dev:
	$(MAKE) -f deploy.mak

#.PHONY: $(SUBDIRS)
#$(SUBDIRS):
#	$(MAKE) -C $@
#
#.PHONY: all
#all: $(SUBDIRS)

#executable: library


# EOF
