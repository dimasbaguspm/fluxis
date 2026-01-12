package resources

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/dimasbaguspm/fluxis/internal/models"
	"github.com/dimasbaguspm/fluxis/internal/services"
)

type StatusResource struct {
	statusSrv services.StatusService
}

func NewStatusResource(statusSrv services.StatusService) StatusResource {
	return StatusResource{statusSrv}
}

func (sr StatusResource) Routes(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "status-get-by-project",
		Method:      http.MethodGet,
		Path:        "/projects/{projectId}/statuses",
		Summary:     "Get statuses for a project",
		Tags:        []string{"Project", "Status"},
		Security:    []map[string][]string{{"bearer": {}}},
	}, sr.getByProject)

	huma.Register(api, huma.Operation{
		OperationID: "status-get",
		Method:      http.MethodGet,
		Path:        "/statuses/{statusId}",
		Summary:     "Get a status by id",
		Tags:        []string{"Project", "Status"},
		Security:    []map[string][]string{{"bearer": {}}},
	}, sr.get)

	huma.Register(api, huma.Operation{
		OperationID: "status-create",
		Method:      http.MethodPost,
		Path:        "/projects/{projectId}/statuses",
		Summary:     "Create a status for a project",
		Tags:        []string{"Project", "Status"},
		Security:    []map[string][]string{{"bearer": {}}},
	}, sr.create)

	huma.Register(api, huma.Operation{
		OperationID: "status-update",
		Method:      http.MethodPatch,
		Path:        "/statuses/{statusId}",
		Summary:     "Update a status",
		Tags:        []string{"Project", "Status"},
		Security:    []map[string][]string{{"bearer": {}}},
	}, sr.update)

	huma.Register(api, huma.Operation{
		OperationID: "status-delete",
		Method:      http.MethodDelete,
		Path:        "/statuses/{statusId}",
		Summary:     "Delete a status",
		Tags:        []string{"Project", "Status"},
		Security:    []map[string][]string{{"bearer": {}}},
	}, sr.delete)

	huma.Register(api, huma.Operation{
		OperationID: "status-reorder",
		Method:      http.MethodPost,
		Path:        "/projects/{projectId}/statuses/reorder",
		Summary:     "Reorder statuses for a project",
		Tags:        []string{"Project", "Status"},
		Security:    []map[string][]string{{"bearer": {}}},
	}, sr.reorder)
}

func (sr StatusResource) getByProject(ctx context.Context, input *struct {
	Path string `path:"projectId"`
}) (*struct{ Body []models.StatusModel }, error) {
	resp, err := sr.statusSrv.GetByProject(ctx, input.Path)
	if err != nil {
		return nil, err
	}
	return &struct{ Body []models.StatusModel }{Body: resp}, nil
}

func (sr StatusResource) get(ctx context.Context, input *struct {
	Path string `path:"statusId"`
}) (*struct{ Body models.StatusModel }, error) {
	resp, err := sr.statusSrv.GetDetail(ctx, input.Path)
	if err != nil {
		return nil, err
	}
	return &struct{ Body models.StatusModel }{Body: resp}, nil
}

func (sr StatusResource) create(ctx context.Context, input *struct {
	Path string `path:"projectId"`
	Body models.StatusCreateModel
}) (*struct{ Body models.StatusModel }, error) {
	resp, err := sr.statusSrv.Create(ctx, input.Path, input.Body)
	if err != nil {
		return nil, err
	}
	return &struct{ Body models.StatusModel }{Body: resp}, nil
}

func (sr StatusResource) update(ctx context.Context, input *struct {
	Path string `path:"statusId"`
	Body models.StatusUpdateModel
}) (*struct{ Body models.StatusModel }, error) {
	resp, err := sr.statusSrv.Update(ctx, input.Path, input.Body)
	if err != nil {
		return nil, err
	}
	return &struct{ Body models.StatusModel }{Body: resp}, nil
}

func (sr StatusResource) delete(ctx context.Context, input *struct {
	Path string `path:"statusId"`
}) (*struct{}, error) {
	err := sr.statusSrv.Delete(ctx, input.Path)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (sr StatusResource) reorder(ctx context.Context, input *struct {
	Path string `path:"projectId"`
	Body models.StatusReorderModel
}) (*struct{ Body []models.StatusModel }, error) {
	resp, err := sr.statusSrv.Reorder(ctx, input.Path, input.Body.IDs)
	if err != nil {
		return nil, err
	}
	return &struct{ Body []models.StatusModel }{Body: resp}, nil
}
