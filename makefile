#
#
#

SUBDIRS := \
						$(CURDIR)/cmd/anyclient/cmdline \
					  $(CURDIR)/cmd/anyserver/cmdline \
						$(CURDIR)/cmd/gspktservice/cmdline 

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

#.PHONY: $(SUBDIRS)
#$(SUBDIRS):
#	$(MAKE) -C $@
#
#.PHONY: all
#all: $(SUBDIRS)

#executable: library


# EOF
