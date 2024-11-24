package biz

import (
	"context"
	v1 "project/api/project/v1"
	"project/internal/ent"
)

type Project = v1.Project

// ProjectRepository represents the interface for operating project entities stored in the database.
type ProjectRepository interface {
	Add(ctx context.Context, project *Project) error
	Remove(ctx context.Context, project *Project) error
	Update(ctx context.Context, project *Project) error
	FindById(ctx context.Context, id string) (*Project, error)
	FindByName(ctx context.Context, name string) (*Project, error)
	RecoverById(ctx context.Context, id string) error
	IsProjectIDExist(ctx context.Context, projectID string) (bool, error)
	GetProjectPath(ctx context.Context, projectId string) ([]*Project, error)
	FindNearbyProjects(ctx context.Context, location *v1.GeoPoint, radius float64) ([]*v1.Project, error)
}

type ProjectManager struct {
	repo ProjectRepository
}

func NewProjectManager(repo ProjectRepository) *ProjectManager {
	return &ProjectManager{repo: repo}
}

func (m *ProjectManager) Add(ctx context.Context, project *Project) error {
	// may be you need to validate the project before adding it,then you can add it
	return m.repo.Add(ctx, project)
}

func (m *ProjectManager) RemoveById(ctx context.Context, id string) (err error) {
	var proj *Project
	if proj, err = m.repo.FindById(ctx, id); err != nil {
		if ent.IsNotFound(err) {
			return v1.ErrorProjectNotFound("Cannot find the specified project with id %v", id)
		}
		return err
	}
	return m.repo.Remove(ctx, proj)
}

func (m *ProjectManager) Update(ctx context.Context, project *Project) (err error) {
	return m.repo.Update(ctx, project)
}

func (m *ProjectManager) GetById(ctx context.Context, id string) (proj *Project, err error) {
	return m.repo.FindById(ctx, id)
}

func (m *ProjectManager) GetByName(ctx context.Context, name string) (proj *Project, err error) {
	return m.repo.FindByName(ctx, name)
}

func (m *ProjectManager) RecoverById(ctx context.Context, id string) (err error) {
	return m.repo.RecoverById(ctx, id)
}

func (m *ProjectManager) IsProjectIDExist(ctx context.Context, projectID string) (bool, error) {
	return m.repo.IsProjectIDExist(ctx, projectID)
}
func (m *ProjectManager) GetProjectPath(ctx context.Context, projectId string) ([]*v1.Project, error) {
	return m.repo.GetProjectPath(ctx, projectId)
}
func (m *ProjectManager) FindNearbyProjects(ctx context.Context, location *v1.GeoPoint, radius float64) ([]*v1.Project, error) {
	return m.repo.FindNearbyProjects(ctx, location, radius)
}
