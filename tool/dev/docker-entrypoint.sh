#!/usr/bin/env sh

set -eu

main() {
    set +e
    /app "${@}"
    EXIT_CODE=$?
    set -e

    if [ "${EXIT_CODE}" -eq 0 ]; then
        echo "Application exited successfully."
    else
        echo "Application exited with code ${EXIT_CODE}."
    fi

    echo "Now sleeping."
    sleep infinity
}

main "${@}"
