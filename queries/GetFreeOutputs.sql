SELECT o.txid, o.idx, o.amount
FROM outputs o
LEFT JOIN inputs i 
    ON o.txid = i.prev_txid 
   AND o.idx = i.prev_idx
WHERE i.txid IS NULL 
  AND o.recipient = $1;
