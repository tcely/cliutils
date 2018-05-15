#!/bin/bash
host="${1:-127.0.0.1}"
port="${2:-443}"

ciphers='ALL:!eNULL'
printf 'Using openssl at: '
command -v openssl
openssl version -a
printf '\nCiphers selected by server at %s using TCP port %s:\n' "$host" "$port"
while : ; do
    nextCipher="$(openssl s_client </dev/null 2>/dev/null -connect "${host}:${port}" -cipher "$ciphers" | awk '$1 == "Cipher" && $2 == ":" {print $3; if (tmpkey) {print tmpkey}} /^Server Temp Key:/ {tmpkey=$0}')"
    [ -z "$nextCipher" ] && break || echo "${nextCipher/$'\n'/ # }"
    ciphers="${ciphers}:-${nextCipher%%$'\n'*}";
done