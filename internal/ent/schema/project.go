//go:generate ent generate .
package schema

import (
	"ariga.io/atlas/sql/schema"
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Project holds the schema definition for the Project entity.
type Project struct {
	ent.Schema
}

// Fields of the Project.
func (Project) Fields() []ent.Field {
	return []ent.Field{
		field.String("project_id").
			NotEmpty().
			Comment("Unique identifier of the Project"),
		field.String("parent_proj_id").
			Default("").
			NotEmpty().
			Comment("Identifier of the parent Project"),
		field.String("desc").
			Default("").
			Comment("Optional description of the Project"),
		field.String("location").
			Default("Unknown location").
			Comment("Geographical location description (province, city, district)"),
		field.Other("coordinate", &Point{}).
			SchemaType(map[string]string{
				"mysql": "POINT", // Specify the MySQL type
			}).
			Comment("Geographical coordinates (latitude, longitude)"),
		field.Time("create_time").
			DefaultNow().
			Immutable().
			Comment("Timestamp when the Project was created"),
		field.Time("last_update").
			DefaultNow().
			UpdateDefaultNow().
			Comment("Timestamp when the Project was last updated"),
	}
}

// Edges of the Project.
func (Project) Edges() []ent.Edge {
	return nil
}

// Indexes of the Project.
func (Project) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("parent_proj_id").
			Comment("Index to optimize Project tree queries"),
	}
}
func (Project) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),
		entsql.Annotation{
			Table:     "sys_users",
			Charset:   "utf8mb4",
			Collation: "utf8mb4_unicode_ci",
			Options:   "ENGINE = InnoDB",
		},
		schema.Comment("Project information sheet"),
	}
}
