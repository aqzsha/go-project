package api

import "github.com/gin-gonic/gin"

// Movies proxy
// @Summary Movies service proxy
// @Description Proxy all Movies API requests to movies microservice
// @Tags Movies
// @Accept json
// @Produce json
// @Param path path string true "Any movies path"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/film/{path} [get]
// @Router /api/v1/film/{path} [post]
// @Router /api/v1/film/{path} [put]
// @Router /api/v1/film/{path} [delete]
// @Router /api/v1/film/{path} [patch]
// @Router /api/v1/cinema/{path} [get]
// @Router /api/v1/genre/{path} [get]
// @Router /api/v1/hall/{path} [get]
func MoviesSwaggerProxy(_ *gin.Context) {}
