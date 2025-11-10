package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// Project holds the schema definition for the Project entity.
type Project struct {
	ent.Schema
}

// Fields of the Project.
func (Project) Fields() []ent.Field {
	return []ent.Field{
		// Primary key - UUID type
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable().
			Comment("Unique identifier for the project"),

		// Project name - required field
		field.String("name").
			NotEmpty().
			Comment("Name of the project"),

		// Project description - optional
		field.String("description").
			Optional().
			Comment("Detailed description of the project"),

		// Project status - enum field
		field.Enum("status").
			Values("ACTIVE", "COMPLETED", "ON_HOLD").
			Default("ACTIVE").
			Comment("Current status of the project"),

		// Project priority - enum field
		field.Enum("priority").
			Values("HIGH", "MEDIUM", "LOW").
			Default("MEDIUM").
			Comment("Priority level of the project"),

		// Start date
		field.Time("start_date").
			Comment("Project start date"),

		// End date
		field.Time("end_date").
			Comment("Project end date (deadline)"),

		// Budget amount
		field.Float("budget").
			Positive().
			Comment("Project budget amount"),

		// Timestamps for tracking creation and updates
		field.Time("created_at").
			Default(time.Now).
			Immutable().
			Comment("Timestamp when project was created"),

		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			Comment("Timestamp when project was last updated"),
	}
}

// Edges of the Project.
func (Project) Edges() []ent.Edge {
	return []ent.Edge{
		// Many-to-many relationship with Employee
		// One project can have many employees (team members)
		// One employee can work on many projects
		edge.To("team_members", Employee.Type).
			Comment("Employees assigned to this project (team members)"),
	}
}

// Indexes of the Project.
func (Project) Indexes() []ent.Index {
	return []ent.Index{
		// Index on status for faster filtering by status
		index.Fields("status"),
		// Index on priority for faster filtering by priority
		index.Fields("priority"),
		// Composite index on dates for range queries
		index.Fields("start_date", "end_date"),
	}
}
