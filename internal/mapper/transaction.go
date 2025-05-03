package mapper

import (
	"log"
	"outbox/internal/model"
	"strconv"
	"time"
)

func MapToTransaction(values map[string]interface{}) model.Transaction {
	getString := func(key string) string {
		if v, ok := values[key]; ok && v != nil {
			return v.(string)
		}
		log.Printf("missing or nil field: %s", key)
		return ""
	}

	getFloat := func(key string) float64 {
		valStr := getString(key)
		f, err := strconv.ParseFloat(valStr, 64)
		if err != nil {
			log.Printf("invalid float in field %s: %v", key, err)
			return 0
		}
		return f
	}

	getInt64 := func(key string) int64 {
		valStr := getString(key)
		i, err := strconv.ParseInt(valStr, 10, 64)
		if err != nil {
			log.Printf("invalid int64 in field %s: %v", key, err)
			return 0
		}
		return i
	}

	return model.Transaction{
		ID:        getString("id"),
		UserID:    getString("user_id"),
		Amount:    getFloat("amount"),
		Currency:  getString("currency"),
		Status:    getString("status"),
		Timestamp: time.Unix(getInt64("timestamp"), 0),
	}
}
