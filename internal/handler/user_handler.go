
// Withdraw 회원 탈퇴
func (h *UserHandler) Withdraw(c *gin.Context) {
    var req dto.WithdrawRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    err := h.userService.Withdraw(c.Request.Context(), &req)
    if err != nil {
        h.handleWithdrawError(c, err)
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "message": "회원 탈퇴가 완료되었습니다. 그동안 이용해 주셔서 감사합니다.",
    })
}

func (h *UserHandler) handleWithdrawError(c *gin.Context, err error) {
    switch {
    case errors.Is(err, service.ErrUnauthorized):
        c.JSON(http.StatusUnauthorized, gin.H{"error": "인증이 필요합니다"})
    case errors.Is(err, service.ErrWrongPassword):
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "비밀번호가 일치하지 않습니다",
            "code":  "WRONG_PASSWORD",
        })
    case errors.Is(err, service.ErrCannotWithdrawAdmin):
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "관리자는 탈퇴할 수 없습니다. 다른 관리자에게 문의하세요.",
            "code":  "ADMIN_CANNOT_WITHDRAW",
        })
    default:
        c.JSON(http.StatusInternalServerError, gin.H{"error": "서버 오류"})
    }
}
