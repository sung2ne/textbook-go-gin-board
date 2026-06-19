
func SetupRouter(hub *ws.Hub) *gin.Engine {
    r := gin.Default()

    wsHandler := handler.NewWSHandler(hub)

    // WebSocket 엔드포인트
    r.GET("/ws", authMiddleware, wsHandler.HandleWebSocket)

    return r
}
