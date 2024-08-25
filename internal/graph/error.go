package graph

import (
	"encoding/json"
	"fmt"

	"github.com/vektah/gqlparser/v2/gqlerror"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func handleError(service string, err error) error {
	result := gqlerror.Error{
		Message: err.Error(),
		Extensions: map[string]interface{}{
			"service": service,
		},
	}

	st, ok := status.FromError(err)
	if !ok {
		return &result
	}

	result.Extensions["code"] = st.Code().String()

	for _, detail := range st.Details() {
		switch info := detail.(type) {
		case proto.Message:
			jsonBytes, marshalErr := protojson.Marshal(info)
			if marshalErr != nil {
				return fmt.Errorf("protojson.Marshal: %w", marshalErr)
			}

			var extensions map[string]interface{}
			if unmarshalErr := json.Unmarshal(jsonBytes, &extensions); unmarshalErr != nil {
				return fmt.Errorf("json.Unmarshal: %w", unmarshalErr)
			}
			extensions["code"] = st.Code().String()
			extensions["service"] = service

			return &gqlerror.Error{
				Message:    err.Error(),
				Extensions: extensions,
			}
		}
	}

	return &result
}
