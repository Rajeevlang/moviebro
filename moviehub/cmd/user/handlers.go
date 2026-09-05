package main

import (
	"moviesapi/internal/user"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	srv user.Service
}

func NewHandler(svc user.Service) *Handler {

	return &Handler{
		srv: svc,
	}
}

func (h *Handler) Register(c *gin.Context) {
	var req user.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input data", "details": err.Error()})
		return
	}

	createdUser, err := h.srv.RegisterUser(c.Request.Context(), &req)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create user"})
		return
	}

	resp := user.UserResponse{
		ID:        createdUser.ID,
		Username:  createdUser.Username,
		Email:     createdUser.Email,
		CreatedAt: createdUser.CreatedAt,
	}

	c.JSON(201, resp)
}

func (h *Handler) Login(c *gin.Context) {
	var req user.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input data", "details": err.Error()})
		return
	}

	token, err := h.srv.LoginUser(c.Request.Context(), &req)
	if err != nil {
		// 401 Unauthorized is the correct status code for failed login
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"token": token})
}

func (h *Handler) LoginGoogle(c *gin.Context) {

}

func (h *Handler) Gloginbefore(c *gin.Context) {
	state, err := user.GenerateState()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"server": "trying again in sometime"})
		return
	}

	expiration := time.Now().Add(10 * time.Minute) // Give the user 10 mins to log in
	cookie := http.Cookie{
		Name:     "openid",
		Value:    state,
		Expires:  expiration,
		Path:     "/",
		HttpOnly: true,                 // JavaScript cannot access this cookie (prevents XSS theft)
		Secure:   true,                 // Only send over HTTPS (essential for production)
		SameSite: http.SameSiteLaxMode, // Required for cross-site redirects to work properly
	}
	http.SetCookie(c.Writer, &cookie)

	authurl := user.Config.AuthCodeURL(state)

	//Redirecte the user to the googleloginHandler
	c.Redirect(http.StatusTemporaryRedirect, authurl)
}

func (h *Handler) Callback(c *gin.Context) {
	usercookie, err := c.Request.Cookie("openid")
	state := c.Query("state")

	if err != nil || state != usercookie.Value {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid state parameter"})
		return
	}

	code := c.Query("code")
	token, err := user.Config.Exchange(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to exchange token"})
		return
	}

	idtoken := token.Extra("id_token").(string)

	req := &user.GoogleLoginRequest{
		IDToken: idtoken,
	}

	jwtToken, err := h.srv.LoginGoogle(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": jwtToken})
}
