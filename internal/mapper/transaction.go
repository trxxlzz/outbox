package mapper

import (
	"outbox/internal/model"
	"strconv"
	"time"
)

func MapToTransaction(values map[string]interface{}) model.Transaction {
	amount, _ := strconv.ParseFloat(values["amount"].(string), 64)
	timestampInt, _ := strconv.ParseInt(values["timestamp"].(string), 10, 64)

	return model.Transaction{
		ID:        values["id"].(string),
		UserID:    values["user_id"].(string),
		Amount:    amount,
		Currency:  values["currency"].(string),
		Status:    values["status"].(string),
		Timestamp: time.Unix(timestampInt, 0),
	}
}
