
// 알림 (인증 필요)
protected.GET("/notifications", notificationHandler.GetNotifications)
protected.GET("/notifications/unread-count", notificationHandler.GetUnreadCount)
protected.PUT("/notifications/:id/read", notificationHandler.MarkAsRead)
protected.PUT("/notifications/read-all", notificationHandler.MarkAllAsRead)
