
import (
    customValidator "goboardapi/internal/validator"
    "github.com/gin-gonic/gin/binding"
    "github.com/go-playground/validator/v10"
)

func main() {
    // ...

    // 커스텀 유효성 검사기 등록
    if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
        customValidator.RegisterCustomValidators(v)
    }

    // ...
}
