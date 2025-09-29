package simpleBlockchain

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"log"
	"strconv"
)

type TXOutput struct {
	Value        int
	ScriptPubKey string
}

type TXInput struct {
	Txid      []byte
	Vout      int
	ScriptSig string
} 

type Transaction struct {
	ID   []byte
	Vin  []TXInput
	Vout []TXOutput
}

func (tx *Transaction) SetId() {
	data := make([][]byte, 0)
	data = append(data, tx.ID)
	for _, in := range tx.Vin {
		data = append(data, in.Txid, Int64ToBytes(int64(in.Vout)), []byte(in.ScriptSig))
	}
	for _, out := range tx.Vout {
		data = append(data, Int64ToBytes(int64(out.Value)), []byte(out.ScriptPubKey))
	}
	pereparedData := bytes.Join(data, []byte{})
	hash := sha256.Sum256(pereparedData)
	fmt.Printf("%xbbbbbbb\n", hash)
	tx.ID = hash[:]
}

func  NewCoinbaseTransaction(adress string) *Transaction {
	fmt.Printf("Награда: %s\n", adress)
	out := TXOutput {
		Value: 10, 
		ScriptPubKey: adress,
	}
	in := TXInput {
		Txid: nil,
		Vout: -1,
		ScriptSig: "_",
	}
	tx := &Transaction{
		ID: nil, 
		Vin: []TXInput{in},
		Vout: []TXOutput{out},
	}
	tx.SetId()
	fmt.Println("Создана coinbase транзакция")
	return tx
}

func NewTransaction(adrFrom, adrTo string, amount int, bc *Blockchain) *Transaction {
	sum, outputs, txns, idx := bc.GetFreeOutputs(adrFrom)
	fmt.Println("Баланс кошелька: " + strconv.Itoa(sum) + " " + adrFrom)
	if sum < amount {
		log.Fatal("Ноу мани")
	}
	sum = 0
	inputsInTx := make([]TXInput, 0)
	for  i := 0; sum < amount; i++ {
		sum += outputs[i].Value
		inputsInTx = append(inputsInTx, TXInput{txns[i], idx[i], adrFrom})
	}
	outputsInTx := make([]TXOutput, 0)
	if sum != amount {
		outputsInTx = append(outputsInTx, TXOutput{sum - amount, adrFrom})
	}
	outputsInTx = append(outputsInTx, TXOutput{amount, adrTo})
	fmt.Println("Транзакция создана: " + adrFrom + " " + adrTo + " " + strconv.Itoa(amount))
	tx := &Transaction{
		nil, inputsInTx, outputsInTx}
	tx.SetId()
	return tx
}

