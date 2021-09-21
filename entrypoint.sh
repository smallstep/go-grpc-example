#!/bin/bash
set -eo pipefail

export STEPPATH=$(step path)

# List of env vars required for step ca bootstrap
declare -ra REQUIRED_INIT_VARS=(DOMAIN STEP_CA_URL STEP_CA_FINGERPRINT)

# List of env variables required for create and renew a certificate
declare -ra RENEW_INIT_VARS=(STEP_CA_PROVISIONER STEP_CA_PASSWORD)

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

# Ensure all env vars required to create a certificate.
function create_and_renew_certificate () {
    local missing_vars=0
    for var in "${RENEW_INIT_VARS[@]}"; do
        if [ -z "${!var}" ]; then
            missing_vars=1
        fi
    done
    if [ ${missing_vars} = 1 ]; then
		>&2 echo "Please provide required environment variables using the --env or or --env-file flags."
        >&2 echo "Required variables are ${REQUIRED_INIT_VARS[*]} and ${RENEW_INIT_VARS[*]}."
    else
        step ca certificate --provisioner "${STEP_CA_PROVISIONER}" --provisioner-password-file <(echo ${STEP_CA_PASSWORD}) ${DOMAIN} tls.crt tls.key
        step ca renew --daemon tls.crt tls.key &
    fi
}

if [ ! -f "${STEPPATH}/certs/root_ca.crt" ]; then
	bootstrap_if_possible
fi

if [ "${STEP_CA_RENEW}" = true ]; then
    create_and_renew_certificate
fi

exec "${@}"
