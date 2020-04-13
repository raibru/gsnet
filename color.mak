#
# print colored output
RESET_COLOR    = \033[0m
make_std_color = \033[3$1m      # defined for 1 through 7
make_color     = \033[38;5;$1m  # defined for 1 through 255
OK_COLOR       = $(strip $(call make_std_color,2))
WRN_COLOR      = $(strip $(call make_std_color,3))
ERR_COLOR      = $(strip $(call make_std_color,1))
STD_COLOR      = $(strip $(call make_color,8))

COLOR_OUTPUT = 2>&1 |                                        \
    while IFS='' read -r line; do                            \
        if  [[ $$line == FAIL* ]]; then                      \
            echo -e "$(ERR_COLOR)$${line}$(RESETCOLOR)";     \
        elif [[ $$line == *:[\ ]FAIL:* ]]; then              \
            echo -e "$(ERR_COLOR)$${line}$(RESETCOLOR)";     \
        elif [[ $$line == [\-][\-][\-][\ ]FAIL:* ]]; then    \
            echo -e "$(ERR_COLOR)$${line}$(RESETCOLOR)";     \
        elif [[ $$line == WARN* ]]; then                     \
            echo -e "$(WRN_COLOR)$${line}$(RESET_COLOR)";    \
        elif [[ $$line == PASS ]]; then                       \
            echo -e "$(OK_COLOR)$${line}$(RESET_COLOR)";     \
        elif [[ $$line == [\-][\-][\-][\ ]PASS:* ]]; then    \
            echo -e "$(OK_COLOR)$${line}$(RESETCOLOR)";     \
        elif [[ $$line == ok* ]]; then                       \
            echo -e "$(OK_COLOR)$${line}$(RESET_COLOR)";     \
        else                                                 \
            echo -e "$(STD_COLOR)$${line}$(RESET_COLOR)";    \
        fi;                                                  \
    done; exit $${PIPESTATUS[0]};

