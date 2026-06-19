
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
    var req ForgotPasswordRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user, err := h.userService.FindByEmail(c.Request.Context(), req.Email)
    if err != nil {
        // 보안상 존재 여부를 알리지 않음
        c.JSON(http.StatusOK, gin.H{"message": "이메일을 확인해주세요"})
        return
    }

    token, _ := h.authService.GenerateResetToken(user.ID)
    resetLink := fmt.Sprintf("%s/reset-password?token=%s", h.config.BaseURL, token)

    // 우선순위 높음 (비밀번호 재설정은 긴급)
    h.emailService.SendPasswordReset(c.Request.Context(), user.Email, user.Username, resetLink)

    c.JSON(http.StatusOK, gin.H{"message": "이메일을 확인해주세요"})
}
