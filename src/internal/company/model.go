package company

// Соответсвует таблице Company
// Используется для внутренней передачи данных
type Company struct {
	INN          string  `json:"inn"`
	KPP          *string `json:"kpp,omitempty"`
	OGRN         *string `json:"ogrn,omitempty"`
	OrgName      *string `json:"org_name,omitempty"`
	OrgShortName *string `json:"org_short_name,omitempty"`
}

// CreateCompanyRequest содержит данные для создания компании через POST /company
type CreateCompanyRequest struct {
	Company
}

// CompanyResponse содержит данные для ответа на запросы, связанные с компаниями
type CompanyResponse struct {
	INN          string  `json:"inn"`
	KPP          *string `json:"kpp,omitempty"`
	OGRN         *string `json:"ogrn,omitempty"`
	OrgName      *string `json:"org_name,omitempty"`
	OrgShortName *string `json:"org_short_name,omitempty"`
}
