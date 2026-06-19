
func (s *userService) RestoreAccount(ctx context.Context, email, password string) error {
    // 소프트 삭제된 사용자 조회
    var user domain.User
    err := s.db.Unscoped().
        Where("email = ? AND deleted_at IS NOT NULL", email).
        First(&user).Error
    if err != nil {
        return ErrUserNotFound
    }

    // 삭제된 지 30일 이내인지 확인
    if time.Since(user.DeletedAt.Time) > 30*24*time.Hour {
        return errors.New("복구 기간이 만료되었습니다")
    }

    // 비밀번호 확인
    if !s.passwordHasher.Compare(user.Password, password) {
        return ErrWrongPassword
    }

    // 복구 처리
    return s.db.Unscoped().Model(&user).Update("deleted_at", nil).Error
}
