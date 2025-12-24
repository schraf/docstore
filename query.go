package docstore

// QueryFilter represents a function to limit results in a Query operation.
type QueryFilter[T DocData] func(doc *Document[T]) bool

// QuerySort represents a function to compare documents for sorting in a Query operation.
type QuerySort[T DocData] func(doc1 *Document[T], doc2 *Document[T]) bool

// Query represents a query for documents.
type Query[T DocData] struct {
	Filters []QueryFilter[T]
	SortBy  *QuerySort[T]
	Limit   int
}

// QueryResult represents the result of a query operation.
type QueryResult[T DocData] struct {
	Documents []*Document[T]
	Total     int
}
