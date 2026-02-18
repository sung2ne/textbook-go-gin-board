package email

import (
	"bytes"
	"errors"
	"html/template"
)

// TemplateData 이메일 템플릿 데이터
type TemplateData struct {
	Username string
	Link     string
	AppName  string
}

var templates = map[string]string{
	"welcome": `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
</head>
<body>
    <h1>{{.AppName}}에 오신 것을 환영합니다!</h1>
    <p>안녕하세요, {{.Username}}님!</p>
    <p>회원가입을 완료해 주셔서 감사합니다.</p>
</body>
</html>
`,
	"reset_password": `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
</head>
<body>
    <h1>비밀번호 재설정</h1>
    <p>안녕하세요, {{.Username}}님!</p>
    <p>아래 링크를 클릭하여 비밀번호를 재설정하세요.</p>
    <p><a href="{{.Link}}">비밀번호 재설정</a></p>
    <p>이 링크는 1시간 후 만료됩니다.</p>
</body>
</html>
`,
}

// RenderTemplate 이메일 템플릿 렌더링
func RenderTemplate(name string, data TemplateData) (string, error) {
	tmplStr, ok := templates[name]
	if !ok {
		return "", errors.New("template not found: " + name)
	}

	tmpl, err := template.New(name).Parse(tmplStr)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
