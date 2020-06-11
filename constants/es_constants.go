package constants

const (
	IndexName    = "news_index"
	DocType      = "_doc"
	IndexMapping = `
	{
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
)
