# go_training
few tasks/projects written in Go

# web-service in server side

- start HTTP Server
```
~/ProjectsGO/go_training/web-service > go run .
                                     
[GIN] 2022/07/26 - 16:35:36 | 200 |     168.041µs |       127.0.0.1 | GET      "/albums"
[GIN] 2022/07/26 - 16:35:53 | 201 |     185.583µs |       127.0.0.1 | POST     "/albums"
[GIN] 2022/07/26 - 16:36:02 | 200 |      38.792µs |       127.0.0.1 | GET      "/albums/2"
[GIN] 2022/07/26 - 16:36:12 | 200 |      45.083µs |       127.0.0.1 | DELETE   "/albums/2"
[GIN] 2022/07/26 - 16:36:24 | 200 |      58.708µs |       127.0.0.1 | PUT      "/albums/3"

```

# web-service in client side

- return all items 
```
curl 'http://localhost:8080/albums'
```
- add a new item
```
curl -X POST -d '{"id": "4", "title": "Miracle", "artist": "Celine Dion", "price": 10.99}' 'http://localhost:8080/albums'
```
- return a specific item
```
curl 'http://localhost:8080/albums/2'
```
- delete item 
```
curl -X DELETE 'http://localhost:8080/albums/2'
```
- update item (=delete + insert)
```
curl -X PUT -d '{"title": "Live in Vegas", "artist": "Elvis", "price": 12.29}' 'http://localhost:8080/albums/3'
```
