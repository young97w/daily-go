package errorPage

import (
	"geektime/web/v6"
	"net/http"
	"testing"
)

func TestErrPgBuilder_Build(t *testing.T) {
	errBd := NewErrPg()
	errBd.AddPage(http.StatusNotFound, []byte(myPage))

	s := web.NewHTTPServer(web.ServerWithMiddleware(errBd.Build()))

	s.Get("/user", func(ctx *web.Context) {
		ctx.RespStatusCode = 404
	})

	err := s.Start(":8081")
	if err != nil {
		panic(err)
	}
}

const myPage = `<!DOCTYPE html>
<html>
<head>
    <title>404 Not Found</title>
    <style>
        body {
            background-color: #f8f9fa;
            font-family: Arial, sans-serif;
            color: #343a40;
        }
        .container {
            max-width: 800px;
            margin: 0 auto;
            padding: 20px;
            text-align: center;
        }
        h1 {
            font-size: 60px;
            font-weight: bold;
            margin-top: 100px;
            margin-bottom: 10px;
        }
        p {
            font-size: 24px;
            margin-top: 0;
            margin-bottom: 40px;
        }
        a {
            color: #007bff;
            text-decoration: none;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>404</h1>
        <p>The page you requested could not be found.</p>
    </div>
</body>
</html>
`
