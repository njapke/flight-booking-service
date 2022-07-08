#!/usr/bin/env bash

set -euo pipefail

rm -rf chi
wget https://github.com/go-chi/chi/archive/refs/tags/v5.0.7.zip
unzip v5.0.7.zip && rm v5.0.7.zip
mv chi-5.0.7 chi
