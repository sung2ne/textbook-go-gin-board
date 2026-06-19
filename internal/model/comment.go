package model

type Comment struct {
    ID     uint   `gorm:"primaryKey"`
    PostID uint
    UserID uint
    Body   string
    User   User `gorm:"foreignKey:UserID"`
}

type Post struct {
    ID       uint      `gorm:"primaryKey"`
    Title    string
    UserID   uint
    User     User      `gorm:"foreignKey:UserID"`
    Comments []Comment `gorm:"foreignKey:PostID"`
}
