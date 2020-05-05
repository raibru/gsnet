#!/bin/bash
#
#


session="work"

tmux new-session -d -s $session -n shell \; \
  split-window -h \; \
  new-window -n 'develop' \; \
    split-window -h \; \
    split-window -v \; \
    select-pane -t 0 \; \
  new-window -n 'runtime' \; \
    send-keys 'cd build/develop/serverservice && clear' C-m \; \
    split-window -v -p 70 \; \
    send-keys 'cd build/develop/gspktservice && clear' C-m \; \
    split-window -v -p 50 \; \
    send-keys 'cd build/develop/clientservice && clear' C-m \; \
    select-pane -t 0 \; \
    split-window -h \; \
    send-keys 'cd build/develop/serverservice && clear && tail -f anyserver.log' C-m \; \
    select-pane -t 2 \; \
    split-window -h \; \
    send-keys 'cd build/develop/gspktservice && clear' C-m \; \
    select-pane -t 4 \; \
    split-window -h \; \
    send-keys 'cd build/develop/clientservice && clear && tail -f anyclient.log' C-m \; \
    select-pane -t 0 \;

tmux select-window -t $session:1
tmux attach -t $session

