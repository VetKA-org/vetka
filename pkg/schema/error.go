package schema

type InvalidField struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type BindErrorsResponse struct {
	Errors []InvalidField `json:"errors"`
}

func (r *BindErrorsResponse) Fields() []string {
	fields := make([]string, len(r.Errors))

	for i, item := range r.Errors {
		fields[i] = item.Field
	}

	return fields
}

type ErrorResponse struct {
	Error string `json:"error"`
}
