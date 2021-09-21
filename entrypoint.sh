#!/bin/bash
set -eo pipefail

export STEPPATH=$(step path)

# List of env vars required for step ca bootstrap
declare -ra REQUIRED_INIT_VARS=(DOMAIN STEP_CA_URL STEP_CA_FINGERPRINT)

# Ensure all env vars required to bootstrap the environment.
function bootstrap_if_possible () {
    local missing_vars=0
    for var in "${REQUIRED_INIT_VARS[@]}"; do
        if [ -z "${!var}" ]; then
            missing_vars=1
        fi
    done
    if [ ${missing_vars} = 1 ]; then
		>&2 echo "Please provide required environment variables using the --env or or --env-file flags."
        >&2 echo "Required variables are ${REQUIRED_INIT_VARS[*]}."
    else
       step ca bootstrap --ca-url $STEP_CA_URL --fingerprint $STEP_CA_FINGERPRINT
    fi
}

if [ ! -f "${STEPPATH}/certs/root_ca.crt" ]; then
	bootstrap_if_possible
fi

exec "${@}"
