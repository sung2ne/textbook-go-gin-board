package migration

import "gorm.io/gorm"

func AddIndexes(db *gorm.DB) error {
    // 인덱스 추가
    if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_posts_created ON posts(created_at DESC)").Error; err != nil {
        return err
    }

    // 복합 인덱스
    if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_posts_user_status ON posts(user_id, status)").Error; err != nil {
        return err
    }

    return nil
}
