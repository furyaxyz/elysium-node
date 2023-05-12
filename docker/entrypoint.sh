#!/bin/bash

set -e

if [ "$1" = 'elysium' ]; then
    echo "Initializing Elysium Node with command:"

    if [[ -n "$NODE_STORE" ]]; then
        echo "elysium "${NODE_TYPE}" init --p2p.network "${P2P_NETWORK}" --node.store "${NODE_STORE}""
        elysium "${NODE_TYPE}" init --p2p.network "${P2P_NETWORK}" --node.store "${NODE_STORE}"
    else
        echo "elysium "${NODE_TYPE}" init --p2p.network "${P2P_NETWORK}""
        elysium "${NODE_TYPE}" init --p2p.network "${P2P_NETWORK}"
    fi

    echo ""
    echo ""
fi

echo "Starting Elysium Node with command:"
echo "$@"
echo ""

exec "$@"
