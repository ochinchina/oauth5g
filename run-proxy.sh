#!/bin/bash

export GODEBUG=x509ignoreCN=0
oauth5g proxy -c proxy.yaml
