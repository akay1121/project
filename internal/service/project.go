package service

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	v1 "project/api/project/v1"
	"project/internal/biz"
)

// ProjectService is the service interface for other services or users to call.
//
// It implements the gRPC service interface by embedding the UnimplementedProjectManagementServer struct,
// and overriding the RPC methods, so that the servers would automatically provide the interfaces.
type ProjectService struct {
	v1.UnimplementedProjectManagementServer
	// mgr is the business layer operation collection, which implements the service interface
	mgr *biz.ProjectManager
}

func NewProjectService(mgr *biz.ProjectManager) *ProjectService {
	return &ProjectService{mgr: mgr}
}

func (s *ProjectService) AddProject(ctx context.Context, proj *v1.Project) (empty *emptypb.Empty, err error) {
	// Validate the project input before adding
	//if valid := proj.Validate(); valid != nil {
	//	return nil, v1.ErrorMalformedInput("Malformed project information: %v", valid)
	//}
	err = s.mgr.Add(ctx, proj)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *ProjectService) UpdateProject(ctx context.Context, proj *v1.Project) (empty *emptypb.Empty, err error) {
	// Validate the project input before updating
	//if valid := proj.Validate(); valid != nil {
	//	return nil, v1.ErrorMalformedInput("Malformed project information: %v", valid)
	//}
	err = s.mgr.Update(ctx, proj)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *ProjectService) FindProjectByName(ctx context.Context, name *v1.ProjectName) (proj *v1.Project, err error) {
	// Validate the project name input before searching
	//if valid := name.Validate(); valid != nil {
	//	return nil, v1.ErrorMalformedInput("Malformed project name: %v", valid)
	//}
	return s.mgr.GetByName(ctx, name.Name)
}

func (s *ProjectService) FindProjectById(ctx context.Context, pid *v1.ProjectId) (proj *v1.Project, err error) {
	// Validate the project ID input before searching
	//if valid := pid.Validate(); valid != nil {
	//	return nil, v1.ErrorMalformedInput("Malformed project ID: %v", valid)
	//}
	return s.mgr.GetById(ctx, pid.Id)
}

func (s *ProjectService) RemoveProjectById(ctx context.Context, pid *v1.ProjectId) (empty *emptypb.Empty, err error) {
	err = s.mgr.RemoveById(ctx, pid.Id)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *ProjectService) SearchBranchProjects(ctx context.Context, id *v1.ProjectId) ([]*v1.Project, error) {
	projects, err := s.mgr.GetProjectPath(ctx, id.Id)
	if err != nil {
		return nil, err
	}

	return projects, nil

}
func (s *ProjectService) FindNearbyProjects(ctx context.Context, req *v1.FindNearbyProjectsRequest) (*v1.FindNearbyProjectsResponse, error) {
	projects, err := s.mgr.FindNearbyProjects(ctx, req.CurrentLocation, req.Radius)
	if err != nil {
		return nil, err
	}

	return &v1.FindNearbyProjectsResponse{
		NearbyProjects: projects,
	}, nil
}
