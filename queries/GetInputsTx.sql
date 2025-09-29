SELECT prev_txid, prev_idx, owner
FROM inputs
WHERE txid = $1
ORDER BY idx;