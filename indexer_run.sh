#!/bin/bash

sleep 15 && while true; do build/./indexer --config $INDEX_CONFIG --host $SUGGEST_HOST; sleep 300; done