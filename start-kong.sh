#!/bin/bash

export PATH=/tmp/kong/bin:$PATH

eval $(luarocks path --bin)
kong migrations bootstrap
kong start


