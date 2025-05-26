package option

import "encoding/json"

func (o Option[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.ptr)
}

func (o *Option[T]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		o.ptr = nil
		return nil
	}
	if err := json.Unmarshal(data, &o.ptr); err != nil {
		return err
	}
	return nil
}

func (o Option[T]) IsZero() bool {
	return o.IsNone()
}
func (o Option[T]) IsEmpty() bool {
	return o.IsNone()
}
