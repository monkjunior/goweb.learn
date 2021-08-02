# Teaching myself web development with Go

```bash
#Run the project
go get github.com/pilu/fresh
fresh
```

- How fast can i init this project ? 8m:20

- Web app usually does not need to be go getable :D -> Just name repo as the domain you want to host.

- When you do a GET request you dont have your body of request.

- Serve mux acts like a router of HTTP reuqests.

- Dynamic reloading for local development: 

```bash
go get github.com/pilu/fresh
cd path/go/myapp
fresh
```
- Should explicit set the Content-Type header if you can.

- HTTP status codes
[net/http header write](https://pkg.go.dev/net/http#Header.Write)

> If WriteHeader has not yet been called, Write calls
> WriteHeader(http.StatusOK) before writing the data.

- ServeMux: pattern matching and why we should not use this mux?

    - Exactly match, care of your traling slask */*

    - Longest pattern match

    - "/"  matches everything

    - No dependencies

    - TOO simple

    - Useful for simple testing or set up a middleware for different router, domain-based handling.

    - Missing features: regex matching pattern, named url params, dynamic url.

- [julienschmidt/httprouter](https://github.com/julienschmidt/httprouter)

    - Simple, fast, Mem & CPU friendly

    - If your application needs to serve HTTP files, this is the best choice

    - If your application spends most of the time dealing with DB, read things from disk and so on, this router should not be your primary concern

    - Named parms -> Dynamic URL

    - Support HTTP method for route matching
