#
#
#

include $(_ROOT)/color.mak

GOCMD       := go
GOBUILD     := $(GOCMD) build
GOCLEAN     := $(GOCMD) clean
GOTEST      := $(GOCMD) test
GOGET       := $(GOCMD) get
BUILD_DIR   := ../../../build/bin

$(_MODULE_NAME)_BINARY := $(_OUTTOP)/$(BINARY)

ifneq ($(_NO_RULES),T)

ifneq ($($(_MODULE_NAME)_DEFINED),T)

#
# Some debugging behavior
#
.PHONY: print
print:             ## print out some defined variables
	@echo Module name       : $(_MODULE_NAME)
	@echo Binary name       : $(BINARY)
	@echo Module name binary: $($(_MODULE_NAME)_BINARY)
	@echo Module clean      : $(_CLEAN)
	@echo MOD_NAME_DEFINED  : $($(_MODULE_NAME)_DEFINED)

print-%:           ## print out a given variable named for %
	@echo $* = [$($*)] Type: $(origin $*) 

.PHONY: printvars
printvars:          ## print out all used variables
	@$(foreach V,$(sort $(.VARIABLES)),  \
		$(if $(filter-out environ% default automatic,  \
		$(origin $V)),$(info $V=$($V) ($(value $V)))))

ifdef TRACE
.PHONY: _trace _value
_trace: ; @$(MAKE) --no-print-directory TRACE= $(TRACE) ='$$(warning TRACE $(TRACE))$(shell $(MAKE) TRACE=$(TRACE) _value)'
_value: ; @echo '$(value $(TRACE))'
endif


#
# Some work
#
#all: $($(_MODULE_NAME)_BINARY)
#
##$(_MODULE_NAME): $($(_MODULE_NAME)_BINARY)
#.PHONY: $(_MODULE_NAME)
#$(_MODULE_NAME): $($(_MODULE_NAME)_GOMAIN)
#
#_IGNORE := $(shell mkdir -p $($(_MODULE_NAME)_OUTPUT))
#_CLEAN := clean-$(_MODULE_NAME)
#
#.PHONY: clean $(_CLEAN)
#clean: $(_CLEAN)          ## clean service builds
#
#$(_CLEAN):
#	rm -rf $($(patsubst clean-%,%,$@)_OUTPUT)
#
#$($(_MODULE_NAME)_BINARY): build test
#build: 
#				$(GOBUILD) -o $($(_MODULE_NAME)_BINARY) -v $(GOMAIN)
#test: 
#				$(GOTEST) -v ./...
#run:
#				$(GOBUILD) -o $($(_MODULE_NAME)_BINARY) -v $(GOMAIN)
#				$($(_MODULE_NAME)_BINARY) --version



#----------------------------------------
#all: build test
#build: 
#				$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) -v $(GOMAIN)
#test: 
#				$(GOTEST) -v ./...
#clean: 
#				$(GOCLEAN)
#				rm -f $(BUILD_DIR)/$(BINARY_NAME)
#run:
#				$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) -v $(GOMAIN)
#				$(BUILD_DIR)/$(BINARY_NAME) --version

#all: $($(_MODULE_NAME)_BINARY)
#
#.PHONY: $(_MODULE_NAME)
#$(_MODULE_NAME): $($(_MODULE_NAME)_BINARY)
#
#_IGNORE := $(shell mkdir -p $($(_MODULE_NAME)_OUTPUT))
#_CLEAN := clean-$(_MODULE_NAME)
#
#.PHONY: clean $(_CLEAN)
#clean: $(_CLEAN)
#
#$(_CLEAN):
#	rm -rf $($(patsubst clean-%,%,$@)_OUTPUT)
#
#$($(_MODULE_NAME)_OUTPUT)/%.o: $(_MODULE_PATH)/%.c
#	@$(COMPILE.c) -o '$@' '$<'
#
#$($(_MODULE_NAME)_OUTPUT)/$(BINARY).a: $($(_MODULE_NAME)_OBJS)
#	@$(AR) r '$@' $^
#	@ranlib '$@'
#
#$($(_MODULE_NAME)_OUTPUT)/$(BINARY)$(_EXEEXT): $($(_MODULE_NAME)_OBJS)
#	@$(LINK.cpp) $^ -o'$@'
#$(_MODULE_NAME)_DEFINED := T

endif

endif

# EOF
