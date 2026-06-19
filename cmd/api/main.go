
func main() {
    // Hub 시작
    hub := ws.NewHub()
    go hub.Run()

    // 서비스 생성
    notifService := service.NewNotificationService(hub, notifRepo)

    // 라우터 설정
    r := router.SetupRouter(hub, notifService)

    r.Run(":8080")
}
