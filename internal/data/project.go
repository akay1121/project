package data

import (
	"context"
	"github.com/jinzhu/copier"
	"math"
	v1 "project/api/project/v1"
	"project/internal/biz"
	"project/internal/ent"
)

// projectRepo implements the interface [biz.ProjectRepository] described in the package [project/internal/biz].
type projectRepo struct {
	db    *Data
	cache *Cache
}

// NewProjectRepository creates a new project repository implementation instance and returns the interface value.
func NewProjectRepository(database *Data, cache *Cache) biz.ProjectRepository {
	return &projectRepo{db: database, cache: cache}
}

// Helper function to convert from ent.Project to biz.Project
func convertToBizProject(p *ent.Project) (proj *biz.Project, err error) {
	proj = &biz.Project{}
	if err = copier.Copy(proj, p); err != nil {
		return
	}
	// Remove fields that shouldn't be exposed
	proj.CreateTime = p.CreateTime.Format("2006-01-02 15:04:05")
	proj.LastUpdate = p.LastUpdate.Format("2006-01-02 15:04:05")
	return
}

// Add creates a new project and saves it to the database.
func (r *projectRepo) Add(ctx context.Context, p *biz.Project) (err error) {
	// Check if parent project exists if parentProjId is provided
	if p.ParentProjID != "" {
		var parentExists bool
		if parentExists, err = r.db.Client.Project.Query().
			Where(project.ProjectID(p.ParentProjID)).
			Exist(ctx); err != nil {
			return err
		}
		if !parentExists {
			return v1.ErrorProjectNotFound("Parent project not found")
		}
	}

	// Create new project in the database
	_, err = r.db.Client.Project.Create().
		SetProjectID(p.ProjectID).
		SetParentProjID(p.ParentProjID).
		SetDesc(p.Desc).
		SetLocation(p.Location).
		SetCoordinate(p.Coordinate).
		Save(ctx)
	if err == nil {
		// Cache the project name to avoid duplicate entries
		r.cache.Client.BFAdd(ctx, "project:names", p.ProjectID)
	}
	return
}

// Remove sets the deleted flag for a project (soft delete).
func (r *projectRepo) Remove(ctx context.Context, p *biz.Project) (err error) {
	var proj *ent.Project
	if proj, err = r.db.Client.Project.Query().Where(project.ProjectID(p.ProjectID)).First(ctx); err != nil {
		return err
	}

	// Set the deleted flag to true (soft delete)
	return r.db.Client.Project.Update().Where(project.ProjectID(proj.ProjectID)).
		SetDeleted(true).Exec(ctx)
}

// Update modifies an existing project in the database.
func (r *projectRepo) Update(ctx context.Context, p *biz.Project) error {
	return r.db.Client.Project.Update().
		Where(project.ProjectID(p.ProjectID)).
		SetDesc(p.Desc).
		SetLocation(p.Location).
		SetCoordinate(p.Coordinate).
		Exec(ctx)
}

// FindById retrieves a project by its ID.
func (r *projectRepo) FindById(ctx context.Context, id string) (proj *biz.Project, err error) {
	var p *ent.Project
	if p, err = r.db.Client.Project.Query().Where(project.ProjectID(id)).First(ctx); err != nil {
		return nil, err
	}
	return convertToBizProject(p)
}

// FindByName retrieves a project by its name.
func (r *projectRepo) FindByName(ctx context.Context, name string) (proj *biz.Project, err error) {
	var p *ent.Project
	if p, err = r.db.Client.Project.Query().Where(project.ProjectID(name)).First(ctx); err != nil {
		return nil, err
	}
	return convertToBizProject(p)
}

// RecoverById recovers a soft-deleted project by its ID.
func (r *projectRepo) RecoverById(ctx context.Context, id string) (err error) {
	var proj *ent.Project
	if proj, err = r.db.Client.Project.Query().Where(project.ProjectID(id)).First(ctx); err != nil {
		return err
	}
	return r.db.Client.Project.Update().Where(project.ProjectID(proj.ProjectID)).
		SetDeleted(false).Exec(ctx)
}

// IsProjectIDExist checks if a project ID exists in the database.
func (r *projectRepo) IsProjectIDExist(ctx context.Context, projectID string) (bool, error) {
	// Check if project ID exists using Bloom filter
	if r.cache.Client.BFExists(ctx, "project:names", projectID).Val() { // If already in cache, check DB
		exists, err := r.db.Client.Project.Query().Where(project.ProjectID(projectID)).Exist(ctx)
		return exists, err
	}
	return false, nil
}
func (r *projectRepo) GetProjectPath(ctx context.Context, projectId string) ([]*biz.Project, error) {
	var projectPath []*biz.Project
	currentProject, err := r.db.Client.Project.Query().Where(project.IDEQ(projectId)).First(ctx)
	if err != nil {
		return nil, err
	}

	// Add the current project to the path
	projectPath = append(projectPath, currentProject)

	// Recursively add parent projects to the path
	for currentProject.ParentProjID != "" {
		parentProject, err := r.db.Client.Project.Query().Where(project.IDEQ(currentProject.ParentProjID)).First(ctx)
		if err != nil {
			return nil, err
		}
		projectPath = append(projectPath, parentProject)
		currentProject = parentProject
	}

	// Reverse the project path to return from root to current project
	for i := 0; i < len(projectPath)/2; i++ {
		projectPath[i], projectPath[len(projectPath)-1-i] = projectPath[len(projectPath)-1-i], projectPath[i]
	}

	var projectList []*v1.Project
	for _, p := range projectPath {
		proj, err := convertToBizProject(p)
		if err != nil {
			return nil, err
		}
		projectList = append(projectList, proj)
	}

	return projectList, nil
}
func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const EarthRadius = 6371 // Earth radius in kilometers

	// Convert degrees to radians
	lat1 = lat1 * math.Pi / 180
	lon1 = lon1 * math.Pi / 180
	lat2 = lat2 * math.Pi / 180
	lon2 = lon2 * math.Pi / 180

	// Haversine formula
	deltaLat := lat2 - lat1
	deltaLon := lon2 - lon1
	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1)*math.Cos(lat2)*math.Sin(deltaLon/2)*math.Sin(deltaLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return EarthRadius * c // Returns distance in kilometers
}
func (r *projectRepo) FindNearbyProjects(ctx context.Context, location *v1.GeoPoint, radius float64) ([]*biz.Project, error) {
	var nearbyProjects []*v1.Project
	// Retrieve all projects
	projects, err := r.db.Client.Project.Query().All(ctx)
	if err != nil {
		return nil, err
	}

	// Filter projects within the radius
	for _, p := range projects {
		// Skip projects that do not have a valid coordinate
		if p.Coordinate == nil {
			continue
		}

		// Calculate the distance between the current project and the target location
		distance := haversine(location.Latitude, location.Longitude, p.Coordinate.Latitude, p.Coordinate.Longitude)

		// If the distance is within the radius, add the project to the result
		if distance <= radius {
			proj, err := convertToBizProject(p)
			if err != nil {
				return nil, err
			}
			nearbyProjects = append(nearbyProjects, proj)
		}
	}

	return nearbyProjects, nil
}
