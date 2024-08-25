package gqmarshal

import (
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func MarshalTimestamp(t *timestamppb.Timestamp) graphql.Marshaler {
	if t == nil || !t.IsValid() {
		return graphql.Null
	}

	return graphql.WriterFunc(func(w io.Writer) {
		tt := t.AsTime()
		_, err := io.WriteString(w, strconv.Quote(tt.Format(time.RFC3339)))
		if err != nil {
			return
		}
	})
}

func UnmarshalTimestamp(v interface{}) (*timestamppb.Timestamp, error) {
	str, ok := v.(string)
	if !ok {
		return nil, fmt.Errorf("timestamp must be a string")
	}
	t, err := time.Parse(time.RFC3339, str)
	if err != nil {
		return nil, err
	}
	return timestamppb.New(t), nil
}
