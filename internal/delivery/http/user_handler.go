package http

import (
	"github.com/gin-gonic/gin"
	_ "user-simple-crud/internal/delivery/http/response"
	"user-simple-crud/internal/entity"
	"user-simple-crud/internal/model"
	service "user-simple-crud/internal/services"
)

type UserHTTPHandler struct {
	Handler
	UserService service.UserService
}

func NewUserHTTPHandler(user service.UserService) *UserHTTPHandler {
	return &UserHTTPHandler{
		UserService: user,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Registers a new user with the provided username and password
// @Tags Users
// @Accept json
// @Produce json
// @Param register body entity.UserLogin true "Registration Request"
// @Success 200 {object} response.DataResponse{data=entity.UserLogin} "success"
// @Failure 400 {object} response.DataResponse "error"
// @Router /auth/register [post]
func (h UserHTTPHandler) Register(ctx *gin.Context) {
	request := entity.UserLogin{}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		h.BadRequestJSON(ctx, err.Error())
		return
	}
	if errException := h.UserService.Create(ctx, &request); errException != nil {
		h.ExceptionJSON(ctx, errException)
		return
	}

	h.DataJSON(ctx, request)
}

// Login godoc
// @Summary User login
// @Description Authenticates the user and returns an access token
// @Tags Users
// @Accept json
// @Produce json
// @Param login body entity.UserLogin true "Login Request"
// @Success 200 {object} response.DataResponse{data=service.UserLoginResponse} "success"
// @Failure 400 {object} response.DataResponse "error"
// @Router /auth/login [post]
func (h UserHTTPHandler) Login(ctx *gin.Context) {
	request := entity.UserLogin{}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		h.BadRequestJSON(ctx, err.Error())
		return
	}
	result, errException := h.UserService.Login(ctx, &request)
	if errException != nil {
		h.ExceptionJSON(ctx, errException)
		return
	}

	h.DataJSON(ctx, result)
}

// Create godoc
// @Summary Create a new book
// @Description Creates a new book with the provided details
// @Tags Users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "format: Bearer <JWT TOKEN>"
// @Param notification-list body entity.UserLogin true "User Request"
// @Success 200 {object} response.DataResponse{data=entity.UserLogin} "success"
// @Failure 400 {object} response.DataResponse "error"
// @Router /users [post]
func (h UserHTTPHandler) Create(ctx *gin.Context) {
	request := entity.UserLogin{}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		h.BadRequestJSON(ctx, err.Error())
		return
	}
	if errException := h.UserService.Create(ctx, &request); errException != nil {
		h.ExceptionJSON(ctx, errException)
		return
	}

	h.DataJSON(ctx, request)
}

// List godoc
// @Summary List users
// @Description Retrieves a paginated list of users with optional ordering and filtering
// @Tags Users
// @Accept json
// @Produce json
// @Param Authorization header string true "format: Bearer <JWT TOKEN>"
// @Param pageSize query string false "Number of items per page"
// @Param page query string false "Page number"
// @Param filter query string false "Filter rules<br><br>### Rules Filter<br>rule:<br>  * {Name of Field}:{value}:{Symbol}<br><br>Symbols:<br>  * eq (=)<br>  * lt (<)<br>  * gt (>)<br>  * lte (<=)<br>  * gte (>=)<br>  * in (in)<br>  * like (like)<br><br>Field list:<br>  * id<br>  * username<br>  * email"
// @Param sort query string false "Sort rules:<br><br>### Rules Sort<br>rule:<br>  * {Name of Field}:{Symbol}<br><br>Symbols:<br>  * asc<br>  * desc<br><br>Field list:<br>  * id<br>  * title<br>  * isbn<br>  * author_id"
// @Success 200 {object} response.PaginationResponse{data=[]entity.User,pagination=model.Pagination} "success"
// @Failure 400 {object} response.DataResponse "error"
// @Router /users [get]
func (h UserHTTPHandler) List(ctx *gin.Context) {
	var req model.ListReq
	var err error
	req.Page, req.Order, req.Filter, err = h.ParsePaginationParams(ctx)
	if err != nil {
		h.BadRequestJSON(ctx, err.Error())
		return
	}
	result, errException := h.UserService.List(ctx, req)
	if errException != nil {
		h.ExceptionJSON(ctx, errException)
		return
	}

	h.DataJSON(ctx, result)
}

// FindOne godoc
// @Summary Get details of a book
// @Description Retrieves the details of a specific book by ID
// @Tags Users
// @Accept json
// @Produce json
// @Param Authorization header string true "format: Bearer <JWT TOKEN>"
// @Param id path string true "User ID (UUID format)"
// @Success 200 {object} response.DataResponse{data=entity.User} "success"
// @Failure 400 {object} response.DataResponse "error"
// @Router /users/{id} [get]
func (h UserHTTPHandler) FindOne(ctx *gin.Context) {
	idParam := ctx.Param("id")
	result, errException := h.UserService.FindOne(ctx, idParam)
	if errException != nil {
		h.ExceptionJSON(ctx, errException)
		return
	}

	h.DataJSON(ctx, result)
}

// Update godoc
// @Summary Update an existing book
// @Description Updates an existing book with the provided details
// @Tags Users
// @Accept json
// @Produce json
// @Param Authorization header string true "format: Bearer <JWT TOKEN>"
// @Param id path string true "User ID (UUID format)"
// @Param book body entity.UserLogin true "Updated User details"
// @Success 200 {object} response.DataResponse{data=entity.UserLogin} "success"
// @Failure 400 {object} response.DataResponse "error"
// @Router /users/{id} [put]
func (h UserHTTPHandler) Update(ctx *gin.Context) {
	// Get Info
	idParam := ctx.Param("id")
	request := entity.UserLogin{}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		h.BadRequestJSON(ctx, err.Error())
		return
	}

	if errException := h.UserService.Update(ctx, idParam, &request); errException != nil {
		h.ExceptionJSON(ctx, errException)
		return
	}

	h.DataJSON(ctx, request)
}

// Delete godoc
// @Summary Delete an existing book
// @Description Deletes an existing book by ID
// @Tags Users
// @Accept json
// @Produce json
// @Param Authorization header string true "format: Bearer <JWT TOKEN>"
// @Param id path string true "User ID (UUID format)"
// @Success 200 {object} response.SuccessResponse "success"
// @Failure 400 {object} response.SuccessResponse "error"
// @Router /users/{id} [delete]
func (h UserHTTPHandler) Delete(ctx *gin.Context) {
	idParam := ctx.Param("id")
	if errException := h.UserService.Delete(ctx, idParam); errException != nil {
		h.ExceptionJSON(ctx, errException)
		return
	}

	h.SuccessMessageJSON(ctx, idParam+" has been deleted")
}
