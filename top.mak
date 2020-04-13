#
#
#

_OUTTOP ?= $(_ROOT)/build/bin

.PHONY: all
all:

_MAKEFILES := $(filter %/makefile,$(MAKEFILE_LIST))

_INCLUDED_FROM := $(patsubst $(_ROOT)/%,%,$(if $(_MAKEFILES),$(patsubst %/makefile,%,$(word $(words $(_MAKEFILES)),$(_MAKEFILES)))))

ifeq ($(_INCLUDED_FROM),)
	_MODULE := $(patsubst $(_ROOT)/%,%,$(CURDIR))
else
	_MODULE := $(_INCLUDED_FROM)
endif

_MODULE_PATH := $(_ROOT)/$(_MODULE)
_MODULE_NAME := $(subst /,_,$(_MODULE))
$(_MODULE_NAME)_OUTPUT := $(_OUTTOP)/$(_MODULE)

_GOEXT := .go
_SERVICE :=

# EOF
