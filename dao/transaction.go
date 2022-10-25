package dao

import "bancario/model"

type TransactionDb model.TransactionDb

var transactions = []TransactionDb{
	{"5bbdadf782ebac06a695a8e7", "1", 100.10, "24/10/2022", "2"},
	{"5bbdadf782ebac06a695a8e8", "2", 50.00, "25/10/2022", "1"},
	{"5bbdadf782ebac06a695a8e9", "1", 50.55, "25/10/2022", "3"},
}

func (d TransactionDb) GetTransactionsBySourceId(sourceId string) []TransactionDb {
	var transactionsDb []TransactionDb
	for _, transaction := range transactions {
		if transaction.SourceId == sourceId {
			transactionsDb = append(transactionsDb, transaction)
		}
	}

	return transactionsDb
}

func (d TransactionDb) GetTransactionsByDestinationId(destinationId string) []TransactionDb {
	var transactionsDb []TransactionDb
	for _, transaction := range transactions {
		if transaction.DestinationId == destinationId {
			transactionsDb = append(transactionsDb, transaction)
		}
	}

	return transactionsDb
}
