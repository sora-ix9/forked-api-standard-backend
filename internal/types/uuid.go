package types

import (
	"database/sql/driver"
	"errors"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

type UUID [16]byte

func (u UUID) String() string {
	return uuid.UUID(u[:]).String()
}

func (u UUID) Value() (driver.Value, error) {
	return uuid.UUID(u[:]).String(), nil
}

func (u *UUID) Scan(src interface{}) error {
	switch v := src.(type) {
	case []byte:
		if len(v) == 16 {
			*u = UUID(v)
			return nil
		}
		return errors.New("Invalid UUID bytes length")
	case string:
		id, err := uuid.ParseBytes([]byte(v))
		if err != nil {
			return err
		}
		*u = UUID(id[:])
		return nil
	default:
		return errors.New("Invalid UUID source")
	}
}

func (u UUID) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.String())
}

func (u *UUID) UnmarshalJSON(data []byte) error {
	var str string
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}

	id, err := uuid.ParseBytes([]byte(str))
	if err != nil {
		return err
	}

	*u = UUID(id[:])
	return nil
}

func (u UUID) MarshalBSON() ([]byte, error) {
	return bson.Marshal(u.String())
}

func (u *UUID) UnmarshalBSON(data []byte) error {
	var str string
	err := bson.Unmarshal(data, &str)
	if err != nil {
		return err
	}

	id, err := uuid.Parse(str)
	if err != nil {
		return err
	}

	*u = UUID(id)
	return nil
}
