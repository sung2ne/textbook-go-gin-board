
// WithdrawRequest 회원 탈퇴 요청
type WithdrawRequest struct {
    Password string `json:"password" binding:"required"`
    Reason   string `json:"reason"` // 탈퇴 사유 (선택)
}
