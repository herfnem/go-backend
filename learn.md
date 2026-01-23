As your application grows, dumping every route into `main()` becomes a nightmare. To keep things clean, you’ll want to **split** your logic into different modules and use **dynamic routing** to handle variables like user IDs or slugs.

Since **Go 1.22**, the standard library's `ServeMux` received a massive upgrade, allowing you to handle dynamic parameters and HTTP methods directly without external libraries.

---

## 1. Dynamic Routing (The Modern Way)

You can now define wildcards in your paths using `{name}`. You retrieve these values using `r.PathValue("name")`.

```go
// Example of a dynamic route for a user profile
mux.HandleFunc("GET /users/{id}", func(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")
    fmt.Fprintf(w, "Viewing profile for User ID: %s", id)
})

```

---

## 2. Route Splitting (Project Structure)

To keep your code "Gopher-friendly," it’s common practice to group related routes into separate files or packages.

### The "Routes" File (`routes.go`)

Instead of cluttering `main.go`, create a function that registers all your routes and returns a handler.

```go
package main

import "net/http"

func registerRoutes() http.Handler {
    mux := http.NewServeMux()

    // Static routes
    mux.HandleFunc("GET /", handleHome)

    // Dynamic routes
    mux.HandleFunc("GET /posts/{slug}", handlePostView)
    mux.HandleFunc("POST /posts/create", handlePostCreate)

    return mux
}

```

### The "Handlers" File (`handlers.go`)

This is where the actual business logic lives.

```go
package main

import (
    "fmt"
    "net/http"
)

func handleHome(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Welcome to the Home Page"))
}

func handlePostView(w http.ResponseWriter, r *http.Request) {
    slug := r.PathValue("slug")
    fmt.Fprintf(w, "Displaying post: %s", slug)
}

func handlePostCreate(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Post created successfully!"))
}

```

### The "Entry Point" (`main.go`)

Now your main file is incredibly lean and only focuses on starting the engine.

```go
package main

import (
    "fmt"
    "net/http"
)

func main() {
    router := registerRoutes()

    server := &http.Server{
        Addr:    ":8080",
        Handler: router,
    }

    fmt.Println("Server is running on http://localhost:8080")
    server.ListenAndServe()
}

```

---

## Why split them this way?

1. **Readability:** You don't have to scroll through 500 lines of code to find one endpoint.
2. **Method Specificity:** By using `"GET /path"`, Go automatically returns a `405 Method Not Allowed` if someone tries to `POST` to it—no extra `if` statements required.
3. **Scalability:** If you need to add Middleware (like logging or authentication), you can wrap the entire `mux` in `registerRoutes()`.

Would you like to see how to add **Middleware** to this setup to log every incoming request?
