package repository

type News struct {
	ID      int    `json:"id"`
	Author  string `json:"author"`
	Body    string `json:"body"`
	Created string `json:"created"`
}

type NewsSorter []*News

func (a NewsSorter) Len() int           { return len(a) }
func (a NewsSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a NewsSorter) Less(i, j int) bool { return a[i].Created > a[j].Created }
