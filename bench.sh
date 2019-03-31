#!/bin/sh -eu
hyperfine \
	--prepare 'rm -rf .git/refs/gitstore' \
	--warmup 5 \
	'./store write foo $(uuidgen)' \
	'./gitstore write foo $(uuidgen)'
