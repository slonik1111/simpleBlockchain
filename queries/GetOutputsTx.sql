SELECT amount, recipient
FROM outputs
WHERE txid = $1
ORDER BY idx;