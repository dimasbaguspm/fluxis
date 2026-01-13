package resources

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/dimasbaguspm/fluxis/internal/models"
	"github.com/dimasbaguspm/fluxis/internal/services"
)

type TaskResource struct {
	taskSrv services.TaskService
}

func NewTaskResource(taskSrv services.TaskService) TaskResource {
	return TaskResource{taskSrv}
}

func (tr TaskResource) Routes(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "task-get-paginated",
		Method:      http.MethodGet,
		Path:        "/tasks",
		Summary:     "Get tasks",
		Tags:        []string{"Task"},
		Security:    []map[string][]string{{"bearer": {}}},
	}, tr.search)

	huma.Register(api, huma.Operation{
		OperationID: "task-get",
		Method:      http.MethodGet,
		Path:        "/tasks/{taskId}",
		Summary:     "Get task detail",
		Tags:        []string{"Task"},
		Security:    []map[string][]string{{"bearer": {}}},
	}, tr.get)

	huma.Register(api, huma.Operation{
		OperationID: "task-create",
		Method:      http.MethodPost,
		Path:        "/tasks",
		Summary:     "Create a task",
		Tags:        []string{"Task"},
		Security:    []map[string][]string{{"bearer": {}}},
	}, tr.create)

	huma.Register(api, huma.Operation{
		OperationID: "task-update",
		Method:      http.MethodPatch,
		Path:        "/tasks/{taskId}",
		Summary:     "Update a task",
		Tags:        []string{"Task"},
		Security:    []map[string][]string{{"bearer": {}}},
	}, tr.update)

	huma.Register(api, huma.Operation{
		OperationID: "task-delete",
		Method:      http.MethodDelete,
		Path:        "/tasks/{taskId}",
		Summary:     "Delete a task",
		Tags:        []string{"Task"},
		Security:    []map[string][]string{{"bearer": {}}},
	}, tr.delete)

	huma.Register(api, huma.Operation{
		OperationID: "task-get-logs",
		Method:      http.MethodGet,
		Path:        "/tasks/{taskId}/logs",
		Summary:     "Get logs for a task",
		Tags:        []string{"Task"},
		Security:    []map[string][]string{{"bearer": {}}},
	}, tr.getLogs)
}

func (tr TaskResource) search(ctx context.Context, input *models.TaskSearchModel) (*struct{ Body models.TaskPaginatedModel }, error) {
	resp, err := tr.taskSrv.GetPaginated(ctx, *input)
	if err != nil {
		return nil, err
	}
	return &struct{ Body models.TaskPaginatedModel }{Body: resp}, nil
}

func (tr TaskResource) get(ctx context.Context, input *struct {
	Path string `path:"taskId"`
}) (*struct{ Body models.TaskModel }, error) {
	resp, err := tr.taskSrv.GetDetail(ctx, input.Path)
	if err != nil {
		return nil, err
	}
	return &struct{ Body models.TaskModel }{Body: resp}, nil
}

func (tr TaskResource) create(ctx context.Context, input *struct{ Body models.TaskCreateModel }) (*struct{ Body models.TaskModel }, error) {
	resp, err := tr.taskSrv.Create(ctx, input.Body)
	if err != nil {
		return nil, err
	}
	return &struct{ Body models.TaskModel }{Body: resp}, nil
}

func (tr TaskResource) update(ctx context.Context, input *struct {
	Path string `path:"taskId" format:"uuid"`
	Body models.TaskUpdateModel
}) (*struct{ Body models.TaskModel }, error) {
	resp, err := tr.taskSrv.Update(ctx, input.Path, input.Body)
	if err != nil {
		return nil, err
	}
	return &struct{ Body models.TaskModel }{Body: resp}, nil
}

func (tr TaskResource) delete(ctx context.Context, input *struct {
	Path string `path:"taskId" format:"uuid"`
}) (*struct{}, error) {
	err := tr.taskSrv.Delete(ctx, input.Path)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (tr TaskResource) getLogs(ctx context.Context, input *struct {
	Path string `path:"taskId" format:"uuid"`
	models.LogSearchModel
}) (*struct{ Body models.LogPaginatedModel }, error) {
	resp, err := tr.taskSrv.GetLogs(ctx, input.Path, input.LogSearchModel)
	if err != nil {
		return nil, err
	}
	return &struct{ Body models.LogPaginatedModel }{Body: resp}, nil
}
