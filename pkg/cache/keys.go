package cache

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/dimasbaguspm/fluxis/pkg/transformer"
	"github.com/jackc/pgx/v5/pgtype"
)

func derive(secret []byte, parts ...string) string {
	prefix := strings.Join(parts, ":")
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(prefix))
	sig := fmt.Sprintf("%x", h.Sum(nil))
	return fmt.Sprintf("%s:%s", prefix, sig)
}

func paramsToString(params interface{}) string {
	if params == nil {
		return "all"
	}
	switch v := params.(type) {
	case string:
		return v
	case pgtype.UUID:
		return transformer.UUIDString(v)
	default:
		b, _ := json.Marshal(params)
		return string(b)
	}
}

func KeySingleBoard(hmacKey string, boardID pgtype.UUID) string {
	return derive([]byte(hmacKey), "board", "single", transformer.UUIDString(boardID))
}

func KeyPagedBoards(hmacKey string, params interface{}) string {
	return derive([]byte(hmacKey), "board", "paged", paramsToString(params))
}

func KeyPagedBoardColumns(hmacKey string, params interface{}) string {
	return derive([]byte(hmacKey), "board", "columns", paramsToString(params))
}

func KeySingleActiveSprint(hmacKey string, projectID pgtype.UUID) string {
	return derive([]byte(hmacKey), "sprint", "active", transformer.UUIDString(projectID))
}

func KeySingleSprint(hmacKey string, sprintID pgtype.UUID) string {
	return derive([]byte(hmacKey), "sprint", "single", transformer.UUIDString(sprintID))
}

func KeySingleProject(hmacKey string, projectID pgtype.UUID) string {
	return derive([]byte(hmacKey), "project", "single", transformer.UUIDString(projectID))
}

func KeySingleProjectByKey(hmacKey string, orgID pgtype.UUID, key string) string {
	return derive([]byte(hmacKey), "project", "by-key", transformer.UUIDString(orgID), key)
}

func KeySingleUser(hmacKey string, userID pgtype.UUID) string {
	return derive([]byte(hmacKey), "user", "single", transformer.UUIDString(userID))
}

func KeyPagedProjects(hmacKey string, params interface{}) string {
	return derive([]byte(hmacKey), "project", "paged", paramsToString(params))
}

func KeyPagedSprints(hmacKey string, params interface{}) string {
	return derive([]byte(hmacKey), "sprint", "paged", paramsToString(params))
}

func KeyPagedBoardTickets(hmacKey string, params interface{}) string {
	return derive([]byte(hmacKey), "ticket", "board", paramsToString(params))
}

func KeyPagedSprintTickets(hmacKey string, params interface{}) string {
	return derive([]byte(hmacKey), "ticket", "sprint", paramsToString(params))
}

func KeyPagedProjectBacklog(hmacKey string, params interface{}) string {
	return derive([]byte(hmacKey), "ticket", "backlog", paramsToString(params))
}

func KeySingleTicket(hmacKey string, ticketID pgtype.UUID) string {
	return derive([]byte(hmacKey), "ticket", "single", transformer.UUIDString(ticketID))
}

func KeySingleBoardColumn(hmacKey string, boardColumnID pgtype.UUID) string {
	return derive([]byte(hmacKey), "board-column", "single", transformer.UUIDString(boardColumnID))
}

func KeySingleOrg(hmacKey string, orgID pgtype.UUID) string {
	return derive([]byte(hmacKey), "org", "single", transformer.UUIDString(orgID))
}

func KeyPagedOrganizations(hmacKey string, params interface{}) string {
	return derive([]byte(hmacKey), "org", "paged", paramsToString(params))
}
