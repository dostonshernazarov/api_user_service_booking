package v1

import (
	"context"
	"github.com/google/uuid"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"

	models "api_user_service_booking/api/handlers/models"
	pbu "api_user_service_booking/genproto/user_proto"
	l "api_user_service_booking/pkg/logger"
)

// CreateUser ...
// @Summary CreateUser
// @Description Api for creating a new user
// @Tags user
// @Accept json
// @Produce json
// @Param User body models.User true "createUserModel"
// @Success 200 {object} models.User
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/users/ [post]
func (h *handlerV1) CreateUser(c *gin.Context) {
	var (
		body        models.User
		jspbMarshal protojson.MarshalOptions
	)
	jspbMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to bind json", l.Error(err))
		return
	}

	if body.Id == "" {
		id := uuid.New()
		body.Id = id.String()
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.UserService().Create(ctx, &pbu.User{
		Id:        body.Id,
		FirstName: body.FirtsName,
		LastName:  body.LastName,
		Email:     body.Email,
		Password:  body.Password,
		Birthday:  body.Birthday,
		ImageUrl:  body.ImageUrl,
		CardNum:   body.Card_num,
		Phone:     body.Phone,
		Role:      body.Role,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to create user", l.Error(err))
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetUser gets user by id
// @Summary GetUser
// @Description Api for getting user by id
// @Tags user
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Success 200 {object} models.User
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/users/{id} [get]
func (h *handlerV1) GetUser(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	id := c.Query("id")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.UserService().GetUserByID(
		ctx, &pbu.IdRequest{
			Id: id,
		})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get user", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

//// ListUsers returns list of users
//// @Summary ListUser
//// @Description Api returns list of users
//// @Tags user
//// @Accept json
//// @Produce json
//// @Param page path int64 true "Page"
//// @Param limit path int64 true "Limit"
//// @Succes 200 {object} models.Users
//// @Failure 400 {object} models.StandardErrorModel
//// @Failure 500 {object} models.StandardErrorModel
//// @Router /v1/users/ [get]
//func (h *handlerV1) ListUsers(c *gin.Context) {
//	queryParams := c.Request.URL.Query()
//
//	params, errStr := utils.ParseQueryParams(queryParams)
//	if errStr != nil {
//		c.JSON(http.StatusBadRequest, gin.H{
//			"error": errStr[0],
//		})
//		h.log.Error("failed to parse query params json" + errStr[0])
//		return
//	}
//
//	var jspbMarshal protojson.MarshalOptions
//	jspbMarshal.UseProtoNames = true
//
//	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
//	defer cancel()
//
//	response, err := h.serviceManager.UserService().GetAllUsers(
//		ctx, &pbu.GetAllUsersRequest{
//			Limit: params.Limit,
//			Page:  params.Page,
//		})
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{
//			"error": err.Error(),
//		})
//		h.log.Error("failed to list users", l.Error(err))
//		return
//	}
//
//	c.JSON(http.StatusOK, response)
//}

// UpdateUser updates user by id
// @Summary UpdateUser
// @Description Api returns updates user
// @Tags user
// @Accept json
// @Produce json
// @Succes 200 {Object} models.User
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/users/{id} [put]
func (h *handlerV1) UpdateUser(c *gin.Context) {
	var (
		body        pbu.User
		jspbMarshal protojson.MarshalOptions
	)
	jspbMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to bind json", l.Error(err))
		return
	}
	body.Id = c.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.UserService().UpdateUser(ctx, &body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to update user", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// DeleteUser deletes user by id
// @Summary DeleteUser
// @Description Api deletes user
// @Tags user
// @Accept json
// @Produce json
// @Succes 200 {Object} models.Delete
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/users/{id} [delete]
func (h *handlerV1) DeleteUser(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	guid := c.Param("id")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.UserService().DeleteUserByID(
		ctx, &pbu.IdRequest{
			Id: guid,
		})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to delete user", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

//// GetWithColumnItem returns list of users
//// @Summary GetWithColumnItem
//// @Description Api returns list of users
//// @Tags user
//// @Accept json
//// @Produce json
//// @Param page path int64 true "Page"
//// @Param limit path int64 true "Limit"
//// @Param column path string true "Column"
//// @Param item path string true "Item"
//// @Succes 200 {object} models.Users
//// @Failure 400 {object} models.StandardErrorModel
//// @Failure 500 {object} models.StandardErrorModel
//// @Router /v1/users/columns [get]
//func (h *handlerV1) GetWithColumnItem(c *gin.Context) {
//	queryParams := c.Request.URL.Query()
//
//	columnName := c.Query("column")
//	item := c.Query("item")
//
//	params, errStr := utils.ParseQueryParams(queryParams)
//	if errStr != nil {
//		c.JSON(http.StatusBadRequest, gin.H{
//			"error": errStr[0],
//		})
//		h.log.Error("failed to parse query params json" + errStr[0])
//		return
//	}
//
//	var jspbMarshal protojson.MarshalOptions
//	jspbMarshal.UseProtoNames = true
//
//	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
//	defer cancel()
//
//	response, err := h.serviceManager.UserService().GetWithColumnAndItem(ctx, &pbu.GetWithColumnAndItemReq{
//		Column: columnName,
//		Item:   item,
//		Page:   params.Page,
//		Limit:  params.Limit,
//	})
//
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{
//			"error": err.Error(),
//		})
//		h.log.Error("failed to list users", l.Error(err))
//		return
//	}
//
//	c.JSON(http.StatusOK, response)
//}
