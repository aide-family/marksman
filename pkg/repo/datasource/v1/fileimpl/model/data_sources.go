// Package model defines file-backed datasource models.
package model

type DataSourceModel struct {
	ID           uint32            `yaml:"id"`
	UID          int64             `yaml:"uid"`
	NamespaceUID int64             `yaml:"namespace_uid"`
	Type         string            `yaml:"type"`
	Engine       string            `yaml:"engine"`
	Name         string            `yaml:"name"`
	Status       int8              `yaml:"status"`
	Endpoint     string            `yaml:"endpoint"`
	Description  string            `yaml:"description"`
	Config       map[string]string `yaml:"config"`
	Metadata     map[string]string `yaml:"metadata"`
	CreatedAt    int64             `yaml:"created_at"`
	UpdatedAt    int64             `yaml:"updated_at"`
	DeletedAt    int64             `yaml:"deleted_at"`
	Creator      int64             `yaml:"creator"`
}
