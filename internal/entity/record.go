package entity

type Record struct {
	URL  string
	Size int
	Ttl  int64
}

type RecordDto struct {
	URL  string
	Size int
}
