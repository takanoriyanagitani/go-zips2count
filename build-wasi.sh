#!/bin/sh

tinygo \
	build \
	-o ./zips2count.wasm \
	-target=wasip1 \
	-opt=z \
	-no-debug \
	./zips2count.go
