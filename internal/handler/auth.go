package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/CodingFervor/carbon-emission-management/internal/model"
	"github.com/CodingFervor/carbon-emission-management/internal/repository"
	"github.com/CodingFervor/carbon-emission-management/pkg/jwt"
	"github.com/CodingFervor/carbon-emission-management/pkg/password"
	"github.com/CodingFervor/carbon-emission-management/pkg/response"
)

// LoginRequest carries the credentials submitted at login.
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest carries the fields required to create a new account.
type RegisterRequest struct {
	Username       string `json:"username" binding:"required,min=3,max=50"`
	Password       string `json:"password" binding:"required,min=6"`
	Email          string `json:"email" binding:"required,email"`
	OrganizationID int64  `json:"organization_id" binding:"required"`
	Role           string `json:"role"`
}

// Login authenticates a user and returns a signed JWT.
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "invalid request: "+err.Error())
		return
	}
	user, err := h.User.FindByUsername(req.Username)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			response.Unauthorized(c, "invalid credentials")
			return
		}
		response.InternalError(c, "lookup failed")
		return
	}
	if !password.Compare(user.Password, req.Password) {
		response.Unauthorized(c, "invalid credentials")
		return
	}
	if user.Status != "active" {
		response.Forbidden(c, "account disabled")
		return
	}
	token, err := jwt.GenerateToken(user.ID, user.Username, user.Role, user.OrganizationID, 24)
	if err != nil {
		response.InternalError(c, "token generation failed")
		return
	}
	_ = h.Audit.Record(&model.AuditLog{UserID: &user.ID, Action: "login", Entity: "user", EntityID: &user.ID, IPAddress: c.ClientIP()})
	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  gin.H{"id": user.ID, "username": user.Username, "role": user.Role, "organization_id": user.OrganizationID},
	})
}

// Register creates a new user account.
func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "invalid request: "+err.Error())
		return
	}
	hash, err := password.Hash(req.Password)
	if err != nil {
		response.InternalError(c, "password hashing failed")
		return
	}
	role := req.Role
	if role == "" {
		role = "viewer"
	}
	u := &model.User{
		Username: req.Username, Password: hash, Email: req.Email,
		Role: role, OrganizationID: req.OrganizationID, Status: "active",
	}
	if err := h.User.Create(u); err != nil {
		response.Fail(c, 409, "username or email already exists")
		return
	}
	response.Created(c, gin.H{"id": u.ID, "username": u.Username})
}

// Profile returns the authenticated user's details.
func (h *Handler) Profile(c *gin.Context) {
	uid := userIDOf(c)
	user, err := h.User.Get(uid)
	if err != nil {
		response.NotFound(c, "user")
		return
	}
	response.OK(c, user)
}
