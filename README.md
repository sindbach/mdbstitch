# mdbstitch
MongoDB Stitch API (Experiment) 

### Example

Import package 

```
import (
        "github.com/sindbach/mdbstitch/mdbstitch"
)
```

Initialise StitchClient 

```
  sc := mdbstitch.StitchClient{}
  sc.AppId = "stitch-application-name"
  auth, err := sc.APIKeyAuth("SomeAPIKeyLongStringFromStitch")
  sc.AccessToken = auth.AccessToken
  query := bson.M{
          "field": "value",
  }
  projection := bson.M{}
  limit := 5
  result, err := sc.Query("find", "dbName", "collectionName", &query, &projection, limit)
  
```
