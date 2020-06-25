#!/bin/bash
#
#


session="gsnet"

tmux new-session -d -s $session -n shell \; \
  split-window -h \; \
  new-window -n 'develop' \; \
    split-window -h \; \
    split-window -v \; \
    select-pane -t 0 \; \
  new-window -n 'runtime-1' \; \
    send-keys 'cd build/develop/server_rx && clear' C-m \; \
    split-window -v -p 70 \; \
    send-keys 'cd build/develop/pktservice && clear' C-m \; \
    split-window -v -p 50 \; \
    send-keys 'cd build/develop/client_tx && clear' C-m \; \
    select-pane -t 0 \; \
    split-window -h \; \
    send-keys 'cd build/develop/server_rx && clear && tail -f anyserver.log' C-m \; \
    select-pane -t 2 \; \
    split-window -h \; \
    send-keys 'cd build/develop/pktservice && clear && tail -f pktservice.log' C-m \; \
    select-pane -t 4 \; \
    split-window -h \; \
    send-keys 'cd build/develop/client_tx && clear && tail -f anyclient.log' C-m \; \
    select-pane -t 0 \; \
  new-window -n 'runtime-2' \; \
    send-keys 'cd build/develop/server_tx && clear' C-m \; \
    split-window -v -p 70 \; \
    send-keys 'cd build/develop/pktservice && clear' C-m \; \
    split-window -v -p 50 \; \
    send-keys 'cd build/develop/client_rx && clear' C-m \; \
    select-pane -t 0 \; \
    split-window -h \; \
    send-keys 'cd build/develop/server_tx && clear && tail -f anyserver.log' C-m \; \
    select-pane -t 2 \; \
    split-window -h \; \
    send-keys 'cd build/develop/pktservice && clear && tail -f pktservice.log' C-m \; \
    select-pane -t 4 \; \
    split-window -h \; \
    send-keys 'cd build/develop/client_rx && clear && tail -f anyclient.log' C-m \; \
    select-pane -t 0 \;

tmux select-window -t $session:1
tmux attach -t $session

