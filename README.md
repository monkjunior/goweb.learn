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
