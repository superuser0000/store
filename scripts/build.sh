#!/usr/bin/env bash

go build

cp settings.json.example settings.json

# TODO: sync models
#./tochka-free-market sync-models

# TODO: sync views
#./tochka-free-market sync-views

# TODO: import cities
# TODO: import countries
# TODO: import Moscow metro
# TODO: import SPB metro
# TODO: create admin user
# TODO: setup settings.json