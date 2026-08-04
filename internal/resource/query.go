package resource

type Query struct {
	Predicate string `json:"predicate" doc:"A PostgreSQL-compatible SQL-standard SQL/JSON Path expression that can be used to filter **Place**s by their OpenStreetMap [tags](https://wiki.openstreetmap.org/wiki/Tags) while **query**ing them." example:"$.historic == \"yes\" && $.tourism == \"attraction\""`
}
