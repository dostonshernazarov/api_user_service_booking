package v1

import (
	"api_user_service_booking/api/handlers/models"
	"encoding/json"
	"go.uber.org/zap"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Security      ApiKeyAuth
// @Summary       Get list of policeis
// @Description   This API gets list of policies
// @Tags          rbac
// @Accept        json
// @Produce       json
// @Param         role query string true "Role"
// @Succes        200 {object} models.ListPolePolicyResponse
// @Failure       404 {object} models.Error
// @Failure       500 {object} models.Error
// @Router        /v1/rbac/policy [GET]
func (h *handlerV1) ListAllPolicies(ctx *gin.Context) {
	role := ctx.Query("role")

	var resp models.ListPolicyResponse

	for _, p := range h.enforcer.GetFilteredPolicy(0, role) {
		resp.Policies = append(resp.Policies, &models.Policy{
			Role:     p[0],
			Endpoint: p[1],
			Method:   p[2],
		})
	}

	ctx.JSON(http.StatusOK, resp)
}

// @Security      ApiKeyAuth
// @Summary       Get list of roles
// @Description   This API get list of roles
// @Tags          rbac
// @Accept        json
// @Produce       json
// @Param         limit query int false "limit"
// @Param         offset query int false "offset"
// @Succes        200 {object} []string
// @Failure       404 {object} models.Error
// @Failure       500 {object} models.Error
// @Router        /v1/rbac/roles [GET]
func (h *handlerV1) ListAllRoles(ctx *gin.Context) {
	resp := h.enforcer.GetAllRoles()
	ctx.JSON(http.StatusOK, resp)
}

// @Security      ApiKeyAuth
// @Summary       Create new role for user
// @Description   This API create new role for user
// @Tags          rbac
// @Accept        json
// @Produce       json
// @Succes        200 {object} models.CreateNewRoleRequest
// @Failure       404 {object} models.Error
// @Failure       500 {object} models.Error
// @Router        /v1/rbac/create [POST]
func (h *handlerV1) CreateNewRole(ctx *gin.Context) {
	var bodyReq models.CreateNewRoleRequest

	if err := json.NewDecoder(ctx.Request.Body).Decode(&bodyReq); err != nil {
		h.log.Error("error rbac create new role ", zap.Error(err))
		return
	}

	if _, err := h.enforcer.AddRoleForUser(bodyReq.UserId, bodyReq.Path); err != nil {
		h.log.Error("error while add role for user", zap.Error(err))
		return
	}
	h.enforcer.SavePolicy()
	ctx.JSON(http.StatusOK, bodyReq)
}
