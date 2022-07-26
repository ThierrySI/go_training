# go_training
few tasks/projects written in Go

===============
  web-service
===============
1) return all items
$> curl 'http://localhost:8080/albums'

2) add a new item
$> curl -X POST -d '{"id": "4", "title": "Miracle", "artist": "Celine Dion", "price": 10.99}' 'http://localhost:8080/albums'

3) return a specific item
$> curl 'http://localhost:8080/albums/2'

4) delete item 
$> curl -X DELETE 'http://localhost:8080/albums/2'

5) update item (=delete + insert)
$> curl -X PUT -d '{"title": "Live in Vegas", "artist": "Elvis", "price": 12.29}' 'http://localhost:8080/albums/3'
