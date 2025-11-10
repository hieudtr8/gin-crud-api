package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// Employee holds the schema definition for the Employee entity.
type Employee struct {
	ent.Schema
}

// Fields of the Employee.
func (Employee) Fields() []ent.Field {
	return []ent.Field{
		// Primary key - UUID type
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable().
			Comment("Unique identifier for the employee"),

		// Employee name - required field
		field.String("name").
			NotEmpty().
			Comment("Name of the employee"),

		// Employee email - required and unique
		field.String("email").
			NotEmpty().
			Unique().
			Comment("Email address of the employee (must be unique)"),

		// Foreign key to Department
		field.UUID("department_id", uuid.UUID{}).
			Comment("ID of the department this employee belongs to"),

		// Timestamps
		field.Time("created_at").
			Default(time.Now).
			Immutable().
			Comment("Timestamp when employee was created"),

		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			Comment("Timestamp when employee was last updated"),
	}
}

// Edges of the Employee (relationships with other entities).
func (Employee) Edges() []ent.Edge {
	return []ent.Edge{
		// Many employees belong to one department
		// This is the "from" side of the relationship
		edge.From("department", Department.Type).
			Ref("employees").           // References the "employees" edge in Department
			Field("department_id").     // Uses department_id field as foreign key
			Unique().                   // Each employee has exactly one department
			Required().                 // Department is required (cannot be null)
			Comment("The department this employee belongs to"),

		// Many-to-many relationship with Project
		// One employee can work on many projects
		// This is the reverse edge (from Employee to Project)
		edge.From("projects", Project.Type).
			Ref("team_members").        // References the "team_members" edge in Project
			Comment("Projects that this employee is working on"),
	}
}

// Indexes of the Employee.
func (Employee) Indexes() []ent.Index {
	return []ent.Index{
		// Index on email for faster lookups
		index.Fields("email").Unique(),
		// Index on department_id for faster joins
		index.Fields("department_id"),
	}
}
