#!/bin/sh

iz0="./sample.d/input0.zip"
iz1="./sample.d/input1.zip"

geninput(){
	echo creating input zip files...

	mkdir -p sample.d

	printf hw0 > ./sample.d/z0t0.txt
	printf hw1 > ./sample.d/z0t1.txt

	printf hw2 > ./sample.d/z1t2.txt
	printf hw3 > ./sample.d/z1t3.txt

	ls ./sample.d/z0*.txt | zip -@ -T -v -o ./sample.d/input0.zip
	ls ./sample.d/z1*.txt | zip -@ -T -v -o ./sample.d/input1.zip

}

test -f "${iz0}" || geninput
test -f "${iz1}" || geninput

ls \
	"${iz0}" \
	"${iz1}" |
	cut -d/ -f3- |
	sed 's,^,/guest-i.d/,' |
	wazero \
		run \
		-mount "./sample.d:/guest-i.d:ro" \
		./zips2count.wasm |
	jq -c
