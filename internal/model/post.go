package model

type Post struct {
    ID       uint   `gorm:"primaryKey"`
    Title    string
    UserID   uint
    User     User   `gorm:"foreignKey:UserID"`  // 관계 정의
}

type User struct {
    ID   uint   `gorm:"primaryKey"`
    Name string
}
