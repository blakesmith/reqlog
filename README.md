# Golang net/http request logger.

## Example

Given a http HandlerFunc like this:

```go
func home(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "home.html", nil)
}
```

Wrap your handler with a reqlog handler to have it log basic http request data

```go
import (
       "log"
       "net/http"
       "os"
       "github.com/blakesmith/reqlog"
)

func main() {
	logger := log.New(os.Stderr, "", log.LstdFlags)
	http.Handle("/", reqlog.NewHandler(http.HandlerFunc(home), logger))
	http.ListenAndServe(":5555", nil)
}
```

You'll get log output like this:

```
2012/11/15 13:29:53 [request] GET / [693us]
2012/11/15 13:29:53 [request] GET /favicon.ico [55us]
```

