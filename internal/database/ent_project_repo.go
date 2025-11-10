package database

import (
	"context"
	"fmt"
	"time"

	"gin-crud-api/internal/ent"
	"gin-crud-api/internal/ent/employee"
	"gin-crud-api/internal/ent/project"
	"gin-crud-api/internal/graph/model"
	"gin-crud-api/internal/logger"

	"github.com/google/uuid"
)

// EntProjectRepo implements ProjectRepository using EntGo
type EntProjectRepo struct {
	client *ent.Client
}

// NewEntProjectRepo creates a new project repository using EntGo
func NewEntProjectRepo(client *ent.Client) ProjectRepository {
	return &EntProjectRepo{client: client}
}

// Save creates a new project in the database
func (r *EntProjectRepo) Save(proj *model.Project) error {
	ctx := context.Background()
	log := logger.WithComponent("ProjectRepo")

	log.Debug().
		Str("project_id", proj.ID).
		Str("name", proj.Name).
		Msg("Saving project to database")

	// Parse the UUID string to UUID type
	id, err := uuid.Parse(proj.ID)
	if err != nil {
		log.Error().
			Err(err).
			Str("project_id", proj.ID).
			Msg("Invalid project ID format")
		return fmt.Errorf("invalid project ID: %w", err)
	}

	// Parse start and end dates
	startDate, err := time.Parse("2006-01-02", proj.StartDate)
	if err != nil {
		log.Error().
			Err(err).
			Str("start_date", proj.StartDate).
			Msg("Invalid start date format")
		return fmt.Errorf("invalid start date format: %w", err)
	}

	endDate, err := time.Parse("2006-01-02", proj.EndDate)
	if err != nil {
		log.Error().
			Err(err).
			Str("end_date", proj.EndDate).
			Msg("Invalid end date format")
		return fmt.Errorf("invalid end date format: %w", err)
	}

	// Create project using EntGo's type-safe builder
	create := r.client.Project.
		Create().
		SetID(id).
		SetName(proj.Name).
		SetStatus(project.Status(proj.Status)).
		SetPriority(project.Priority(proj.Priority)).
		SetStartDate(startDate).
		SetEndDate(endDate).
		SetBudget(proj.Budget)

	// Set optional description
	if proj.Description != nil {
		create = create.SetDescription(*proj.Description)
	}

	// Add team members if provided
	if len(proj.TeamMembers) > 0 {
		teamMemberIDs := make([]uuid.UUID, len(proj.TeamMembers))
		for i, member := range proj.TeamMembers {
			memberID, err := uuid.Parse(member.ID)
			if err != nil {
				return fmt.Errorf("invalid team member ID %s: %w", member.ID, err)
			}
			teamMemberIDs[i] = memberID
		}
		create = create.AddTeamMemberIDs(teamMemberIDs...)
	}

	_, err = create.Save(ctx)
	if err != nil {
		log.Error().
			Err(err).
			Str("project_id", proj.ID).
			Msg("Failed to save project to database")
		return fmt.Errorf("failed to save project: %w", err)
	}

	log.Debug().
		Str("project_id", proj.ID).
		Str("name", proj.Name).
		Msg("Project saved successfully")

	return nil
}

// FindByID retrieves a project by its ID with team members
func (r *EntProjectRepo) FindByID(id string) (*model.Project, error) {
	ctx := context.Background()
	log := logger.WithComponent("ProjectRepo")

	log.Debug().
		Str("project_id", id).
		Msg("Finding project by ID")

	// Parse the UUID string
	uid, err := uuid.Parse(id)
	if err != nil {
		log.Error().
			Err(err).
			Str("project_id", id).
			Msg("Invalid project ID format")
		return nil, fmt.Errorf("invalid project ID: %w", err)
	}

	// Query project with team members using EntGo
	entProj, err := r.client.Project.
		Query().
		Where(project.ID(uid)).
		WithTeamMembers().
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			log.Debug().
				Str("project_id", id).
				Msg("Project not found in database")
			return nil, ErrNotFound
		}
		log.Error().
			Err(err).
			Str("project_id", id).
			Msg("Database error while finding project")
		return nil, fmt.Errorf("failed to find project: %w", err)
	}

	log.Debug().
		Str("project_id", entProj.ID.String()).
		Str("name", entProj.Name).
		Msg("Project found successfully")

	// Convert EntGo entity to GraphQL model
	return entProjectToModel(entProj), nil
}

// FindAll retrieves all projects from the database with team members
func (r *EntProjectRepo) FindAll() ([]*model.Project, error) {
	ctx := context.Background()
	log := logger.WithComponent("ProjectRepo")

	log.Debug().Msg("Finding all projects")

	// Query all projects with team members
	entProjs, err := r.client.Project.
		Query().
		WithTeamMembers().
		All(ctx)

	if err != nil {
		log.Error().
			Err(err).
			Msg("Database error while finding all projects")
		return nil, fmt.Errorf("failed to find all projects: %w", err)
	}

	// Convert EntGo entities to GraphQL models
	projects := make([]*model.Project, len(entProjs))
	for i, entProj := range entProjs {
		projects[i] = entProjectToModel(entProj)
	}

	log.Debug().
		Int("count", len(projects)).
		Msg("All projects found successfully")

	return projects, nil
}

// Update updates an existing project
func (r *EntProjectRepo) Update(proj *model.Project) error {
	ctx := context.Background()
	log := logger.WithComponent("ProjectRepo")

	log.Debug().
		Str("project_id", proj.ID).
		Msg("Updating project")

	// Parse the UUID string
	id, err := uuid.Parse(proj.ID)
	if err != nil {
		log.Error().
			Err(err).
			Str("project_id", proj.ID).
			Msg("Invalid project ID format")
		return fmt.Errorf("invalid project ID: %w", err)
	}

	// Parse dates
	startDate, err := time.Parse("2006-01-02", proj.StartDate)
	if err != nil {
		return fmt.Errorf("invalid start date format: %w", err)
	}

	endDate, err := time.Parse("2006-01-02", proj.EndDate)
	if err != nil {
		return fmt.Errorf("invalid end date format: %w", err)
	}

	// Update project using EntGo
	update := r.client.Project.
		UpdateOneID(id).
		SetName(proj.Name).
		SetStatus(project.Status(proj.Status)).
		SetPriority(project.Priority(proj.Priority)).
		SetStartDate(startDate).
		SetEndDate(endDate).
		SetBudget(proj.Budget)

	// Set optional description
	if proj.Description != nil {
		update = update.SetDescription(*proj.Description)
	} else {
		update = update.ClearDescription()
	}

	// Update team members if provided
	if proj.TeamMembers != nil {
		// Clear existing team members first
		update = update.ClearTeamMembers()

		// Add new team members
		if len(proj.TeamMembers) > 0 {
			teamMemberIDs := make([]uuid.UUID, len(proj.TeamMembers))
			for i, member := range proj.TeamMembers {
				memberID, err := uuid.Parse(member.ID)
				if err != nil {
					return fmt.Errorf("invalid team member ID %s: %w", member.ID, err)
				}
				teamMemberIDs[i] = memberID
			}
			update = update.AddTeamMemberIDs(teamMemberIDs...)
		}
	}

	_, err = update.Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			log.Debug().
				Str("project_id", proj.ID).
				Msg("Project not found for update")
			return ErrNotFound
		}
		log.Error().
			Err(err).
			Str("project_id", proj.ID).
			Msg("Failed to update project")
		return fmt.Errorf("failed to update project: %w", err)
	}

	log.Debug().
		Str("project_id", proj.ID).
		Msg("Project updated successfully")

	return nil
}

// Delete deletes a project by its ID
func (r *EntProjectRepo) Delete(id string) error {
	ctx := context.Background()
	log := logger.WithComponent("ProjectRepo")

	log.Debug().
		Str("project_id", id).
		Msg("Deleting project")

	// Parse the UUID string
	uid, err := uuid.Parse(id)
	if err != nil {
		log.Error().
			Err(err).
			Str("project_id", id).
			Msg("Invalid project ID format")
		return fmt.Errorf("invalid project ID: %w", err)
	}

	// Delete project using EntGo
	err = r.client.Project.DeleteOneID(uid).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			log.Debug().
				Str("project_id", id).
				Msg("Project not found for deletion")
			return ErrNotFound
		}
		log.Error().
			Err(err).
			Str("project_id", id).
			Msg("Failed to delete project")
		return fmt.Errorf("failed to delete project: %w", err)
	}

	log.Debug().
		Str("project_id", id).
		Msg("Project deleted successfully")

	return nil
}

// FindByStatus retrieves all projects with a specific status
func (r *EntProjectRepo) FindByStatus(status model.ProjectStatus) ([]*model.Project, error) {
	ctx := context.Background()
	log := logger.WithComponent("ProjectRepo")

	log.Debug().
		Str("status", string(status)).
		Msg("Finding projects by status")

	// Query projects by status
	entProjs, err := r.client.Project.
		Query().
		Where(project.StatusEQ(project.Status(status))).
		WithTeamMembers().
		All(ctx)

	if err != nil {
		log.Error().
			Err(err).
			Str("status", string(status)).
			Msg("Database error while finding projects by status")
		return nil, fmt.Errorf("failed to find projects by status: %w", err)
	}

	// Convert EntGo entities to GraphQL models
	projects := make([]*model.Project, len(entProjs))
	for i, entProj := range entProjs {
		projects[i] = entProjectToModel(entProj)
	}

	log.Debug().
		Str("status", string(status)).
		Int("count", len(projects)).
		Msg("Projects found by status successfully")

	return projects, nil
}

// FindByEmployeeID retrieves all projects that an employee is working on
func (r *EntProjectRepo) FindByEmployeeID(employeeID string) ([]*model.Project, error) {
	ctx := context.Background()
	log := logger.WithComponent("ProjectRepo")

	log.Debug().
		Str("employee_id", employeeID).
		Msg("Finding projects by employee ID")

	// Parse the UUID string
	empID, err := uuid.Parse(employeeID)
	if err != nil {
		log.Error().
			Err(err).
			Str("employee_id", employeeID).
			Msg("Invalid employee ID format")
		return nil, fmt.Errorf("invalid employee ID: %w", err)
	}

	// Query projects by employee ID
	entProjs, err := r.client.Project.
		Query().
		Where(project.HasTeamMembersWith(
			employee.ID(empID),
		)).
		WithTeamMembers().
		All(ctx)

	if err != nil {
		log.Error().
			Err(err).
			Str("employee_id", employeeID).
			Msg("Database error while finding projects by employee ID")
		return nil, fmt.Errorf("failed to find projects by employee ID: %w", err)
	}

	// Convert EntGo entities to GraphQL models
	projects := make([]*model.Project, len(entProjs))
	for i, entProj := range entProjs {
		projects[i] = entProjectToModel(entProj)
	}

	log.Debug().
		Str("employee_id", employeeID).
		Int("count", len(projects)).
		Msg("Projects found by employee ID successfully")

	return projects, nil
}

// AddTeamMember adds an employee to a project's team
func (r *EntProjectRepo) AddTeamMember(projectID string, employeeID string) error {
	ctx := context.Background()
	log := logger.WithComponent("ProjectRepo")

	log.Debug().
		Str("project_id", projectID).
		Str("employee_id", employeeID).
		Msg("Adding team member to project")

	// Parse UUIDs
	projID, err := uuid.Parse(projectID)
	if err != nil {
		return fmt.Errorf("invalid project ID: %w", err)
	}

	empID, err := uuid.Parse(employeeID)
	if err != nil {
		return fmt.Errorf("invalid employee ID: %w", err)
	}

	// Add team member using EntGo
	_, err = r.client.Project.
		UpdateOneID(projID).
		AddTeamMemberIDs(empID).
		Save(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return ErrNotFound
		}
		log.Error().
			Err(err).
			Str("project_id", projectID).
			Str("employee_id", employeeID).
			Msg("Failed to add team member to project")
		return fmt.Errorf("failed to add team member: %w", err)
	}

	log.Debug().
		Str("project_id", projectID).
		Str("employee_id", employeeID).
		Msg("Team member added successfully")

	return nil
}

// RemoveTeamMember removes an employee from a project's team
func (r *EntProjectRepo) RemoveTeamMember(projectID string, employeeID string) error {
	ctx := context.Background()
	log := logger.WithComponent("ProjectRepo")

	log.Debug().
		Str("project_id", projectID).
		Str("employee_id", employeeID).
		Msg("Removing team member from project")

	// Parse UUIDs
	projID, err := uuid.Parse(projectID)
	if err != nil {
		return fmt.Errorf("invalid project ID: %w", err)
	}

	empID, err := uuid.Parse(employeeID)
	if err != nil {
		return fmt.Errorf("invalid employee ID: %w", err)
	}

	// Remove team member using EntGo
	_, err = r.client.Project.
		UpdateOneID(projID).
		RemoveTeamMemberIDs(empID).
		Save(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return ErrNotFound
		}
		log.Error().
			Err(err).
			Str("project_id", projectID).
			Str("employee_id", employeeID).
			Msg("Failed to remove team member from project")
		return fmt.Errorf("failed to remove team member: %w", err)
	}

	log.Debug().
		Str("project_id", projectID).
		Str("employee_id", employeeID).
		Msg("Team member removed successfully")

	return nil
}

// entProjectToModel converts an EntGo project entity to a GraphQL model
func entProjectToModel(entProj *ent.Project) *model.Project {
	proj := &model.Project{
		ID:        entProj.ID.String(),
		Name:      entProj.Name,
		Status:    model.ProjectStatus(entProj.Status),
		Priority:  model.ProjectPriority(entProj.Priority),
		StartDate: entProj.StartDate.Format("2006-01-02"),
		EndDate:   entProj.EndDate.Format("2006-01-02"),
		Budget:    entProj.Budget,
	}

	// Set optional description
	if entProj.Description != "" {
		proj.Description = &entProj.Description
	}

	// Convert team members if loaded
	if entProj.Edges.TeamMembers != nil {
		proj.TeamMembers = make([]*model.Employee, len(entProj.Edges.TeamMembers))
		for i, entEmp := range entProj.Edges.TeamMembers {
			proj.TeamMembers[i] = &model.Employee{
				ID:           entEmp.ID.String(),
				Name:         entEmp.Name,
				Email:        entEmp.Email,
				DepartmentID: entEmp.DepartmentID.String(),
			}
		}
	}

	return proj
}
