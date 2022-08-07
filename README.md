# go_training
few tasks/projects written in Go

`Source` https://go.dev/doc/tutorial/web-service-gin

# web-service in server side

- start HTTP Server
```
~/ProjectsGO/go_training/web-service/cmd > go run .

[GIN] 2022/07/26 - 16:35:36 | 200 |     168.041µs |       127.0.0.1 | GET      "/albums"
[GIN] 2022/07/26 - 16:35:53 | 201 |     185.583µs |       127.0.0.1 | POST     "/albums"
[GIN] 2022/07/26 - 16:36:02 | 200 |      38.792µs |       127.0.0.1 | GET      "/albums/2"
[GIN] 2022/07/26 - 16:36:12 | 200 |      45.083µs |       127.0.0.1 | DELETE   "/albums/2"
[GIN] 2022/07/26 - 16:36:24 | 200 |      58.708µs |       127.0.0.1 | PUT      "/albums/3"

```

# web-service in client side

- create table in CH and add data
```
curl 'http://localhost:8080/albums/initdb'
```
```
-- SQL command executed in ClickHouse
DROP TABLE IF EXISTS albums
CREATE TABLE IF NOT EXISTS albums ( code String, title String, artist String, price Float32 ) engine=Memory
INSERT 3 rows/albums
```


- return all items
```
curl 'http://localhost:8080/albums'
```
- add a new item
```
curl -X POST -d '{"code": "4", "title": "Miracle", "artist": "Celine Dion", "price": 10.99}' 'http://localhost:8080/albums'
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

# Testing

- How to run
```
 ~/ProjectsGO/go_training/web-service > cd pkg/album
 ~/ProjectsGO/go_training/web-service/pkg/album > go test -v
```
- Results
```
=== RUN   TestCreateTableAndData
    albumDB_test.go:94: TestCreateTableAndData | LOG | Number of Album | received:3, wanted:3
--- PASS: TestCreateTableAndData (0.06s)
=== RUN   TestSelectAlbum
    albumDB_test.go:94: TestSelectAlbum | LOG | Number of Album | received:1, wanted:1
--- PASS: TestSelectAlbum (0.02s)
=== RUN   TestSelectAlbums
    albumDB_test.go:94: TestSelectAlbums | LOG | Number of Album | received:3, wanted:3
--- PASS: TestSelectAlbums (0.03s)
=== RUN   TestInsertNewAlbum
    albumDB_test.go:94: TestInsertNewAlbum | LOG | Number of Album | received:1, wanted:1
    albumDB_test.go:94: TestInsertNewAlbum | LOG | Number of Album | received:4, wanted:4
--- PASS: TestInsertNewAlbum (0.03s)
=== RUN   TestDeleteAlbum
    albumDB_test.go:94: TestDeleteAlbum | LOG | Number of Album | received:1, wanted:1
    albumDB_test.go:94: TestDeleteAlbum | LOG | Number of Album | received:3, wanted:3
--- PASS: TestDeleteAlbum (0.04s)
PASS
ok      go_training/web-service/pkg/album       0.655s
```

- How to measure test coverage
```
 ~/ProjectsGO/go_training/web-service/pkg/album > go test -cover
 PASS
 coverage: 44.6% of statements
 ok      go_training/web-service/pkg/album       0.610s
```
