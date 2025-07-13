#!/bin/sh

if [ -z "$1" ];then
  set -- /usr/local/bin/ledger -c /etc/ledger/ledger.yaml
fi

exec "$@"
