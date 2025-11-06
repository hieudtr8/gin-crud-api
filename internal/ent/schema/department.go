package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Department holds the schema definition for the Department entity.
type Department struct {
	ent.Schema
}

// Fields of the Department.
func (Department) Fields() []ent.Field {
	return []ent.Field{
		// Primary key - UUID type
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable().
			Comment("Unique identifier for the department"),

		// Department name - required field
		field.String("name").
			NotEmpty().
			Comment("Name of the department"),

		// Timestamps for tracking creation and updates
		field.Time("created_at").
			Default(time.Now).
			Immutable().
			Comment("Timestamp when department was created"),

		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			Comment("Timestamp when department was last updated"),
	}
}

// Edges of the Department (relationships with other entities).
func (Department) Edges() []ent.Edge {
	return []ent.Edge{
		// One department has many employees
		// This creates a one-to-many relationship
		edge.To("employees", Employee.Type).
			Comment("Employees belonging to this department"),
	}
}
