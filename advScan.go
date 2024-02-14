package main

//func parseJSONToModel(src interface{}, dest interface{}) error {
//	var data []byte
//
//	if b, ok := src.([]byte); ok {
//		data = b
//	} else if s, ok := src.(string); ok {
//		data = []byte(s)
//	} else if src == nil {
//		return nil
//	}
//
//	return json.Unmarshal(data, dest)
//}
//type Registries []string
//func (r *Registries) Scan(src interface{}) error {
//	return parseJSONToModel(src, r)
//}
