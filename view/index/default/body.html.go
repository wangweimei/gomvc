// Code generated by hero.
// DO NOT EDIT!
package template

import (
	"bytes"

	"github.com/shiyanhui/hero"
)

func Body(r string, buffer *bytes.Buffer) {
	buffer.WriteString(`<html>

<head>
    <title>gomvc</title>
</head>

<body>
    `)
	buffer.WriteString(`
    <p>`)
	hero.EscapeHTML(r, buffer)
	buffer.WriteString(`</p>
`)

	buffer.WriteString(`
</body>

</html>`)

}