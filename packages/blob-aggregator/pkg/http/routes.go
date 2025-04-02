package http

func (srv *Server) configureRoutes() {
	srv.echo.GET("/healthz", srv.Health)
	srv.echo.GET("/", srv.Health)
	srv.echo.POST("/queueProposal", srv.queue_proposal)
}
