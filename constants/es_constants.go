package constants

const (
	IndexName    = "news_index"
	DocType      = "_doc"
	IndexMapping = `
	{
		"settings":{
			"number_of_shards":1,
			"number_of_replicas":0
		},
		"mappings":{
		   "properties":{
			  "id":{
				 "type":"integer"
			  },
			  "created":{
				 "type":"date"
			  }
		   }
		}
	}
	`
	PageSize = 10
)
