package resources

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/dimasbaguspm/fluxis/internal/models"
	"github.com/dimasbaguspm/fluxis/internal/services"
)

type ProjectResource struct {
	projectSrv services.ProjectService
}

func NewProjectResource(projectSrv services.ProjectService) ProjectResource {
	return ProjectResource{projectSrv}
}

func (pr ProjectResource) Routes(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "project-get-paginated",
		Method:      http.MethodGet,
		Path:        "/projects",
		Summary:     "Get Projects",
		Tags:        []string{"Project"},
		Security: []map[string][]string{
			{"bearer": {}},
		},
	}, pr.getPaginated)
	huma.Register(api, huma.Operation{
		OperationID: "project-get-detail",
		Method:      http.MethodGet,
		Path:        "/projects/{projectId}",
		Summary:     "Get Project detail",
		Tags:        []string{"Project"},
		Security: []map[string][]string{
			{"bearer": {}},
		},
	}, pr.getDetail)
	huma.Register(api, huma.Operation{
		OperationID: "project-create",
		Method:      http.MethodPost,
		Path:        "/projects/{projectId}",
		Summary:     "Create single project",
		Tags:        []string{"Project"},
		Security: []map[string][]string{
			{"bearer": {}},
		},
	}, pr.create)
	huma.Register(api, huma.Operation{
		OperationID: "project-update",
		Method:      http.MethodPatch,
		Path:        "/projects/{projectId}",
		Summary:     "Update single project",
		Tags:        []string{"Project"},
		Security: []map[string][]string{
			{"bearer": {}},
		},
	}, pr.update)
	huma.Register(api, huma.Operation{
		OperationID: "project-delete",
		Method:      http.MethodDelete,
		Path:        "/projects/{projectId}",
		Summary:     "Delete single project",
		Tags:        []string{"Project"},
		Security: []map[string][]string{
			{"bearer": {}},
		},
	}, pr.delete)
}

func (pr ProjectResource) getPaginated(ctx context.Context, input *models.ProjectSearchModel) (*struct{ Body models.ProjectPaginatedModel }, error) {
	respSrv, err := pr.projectSrv.GetPaginated(ctx, *input)
	if err != nil {
		return nil, err
	}

	return &struct{ Body models.ProjectPaginatedModel }{
		Body: respSrv,
	}, nil
}

func (pr ProjectResource) getDetail(ctx context.Context, input *struct {
	Path string `path:"projectId"`
}) (*struct{ Body models.ProjectModel }, error) {
	respSrv, err := pr.projectSrv.GetDetail(ctx, input.Path)
	if err != nil {
		return nil, err
	}

	return &struct{ Body models.ProjectModel }{
		Body: respSrv,
	}, nil
}

func (pr ProjectResource) create(ctx context.Context, input *struct {
	Body models.ProjectCreateModel
}) (*struct{ Body models.ProjectModel }, error) {
	respSrc, err := pr.projectSrv.Create(ctx, input.Body)
	if err != nil {
		return nil, err
	}

	return &struct{ Body models.ProjectModel }{
		Body: respSrc,
	}, nil
}

func (pr ProjectResource) update(ctx context.Context, input *struct {
	Path string `path:"projectId"`
	Body models.ProjectUpdateModel
}) (*struct{ Body models.ProjectModel }, error) {
	respSrc, err := pr.projectSrv.Update(ctx, input.Path, input.Body)
	if err != nil {
		return nil, err
	}

	return &struct{ Body models.ProjectModel }{
		Body: respSrc,
	}, nil
}

func (pr ProjectResource) delete(ctx context.Context, input *struct {
	Path string `path:"projectId"`
}) (*struct{}, error) {
	err := pr.projectSrv.Delete(ctx, input.Path)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
