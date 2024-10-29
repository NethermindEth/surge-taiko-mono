#!/bin/bash

PRIVATE_KEY=0xbcdf20249abf0ed6d944c0288fad489e33f66b3960d9e6229c1cd214ed3bbe31 ./script/layer1/config_dcap_sgx_verifier.sh  --tcb tcb --qeid qe_identity --mrenclave dfcb4fca3073e3f3a90b05d328688c32619d56f26789c0a9797aa10e765a7807 --mrsigner ca0583a715534a8c981b914589a7f0dc5d60959d9ae79fb5353299a4231673d5 --toggle-mr-check