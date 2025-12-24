package docstore

import "sort"

func sortDocuments[T DocData](docs []*Document[T], sortFunc QuerySort[T]) []*Document[T] {
	sorter := docSortHelper[T]{
		docs:     docs,
		lessFunc: sortFunc,
	}

	sorter.Sort()

	return sorter.docs
}

type docSortHelper[T DocData] struct {
	docs     []*Document[T]
	lessFunc func(a, b *Document[T]) bool
}

func (d *docSortHelper[T]) Sort() {
	sort.Sort(d)
}

func (d *docSortHelper[T]) Len() int {
	return len(d.docs)
}

func (d *docSortHelper[T]) Less(i, j int) bool {
	return d.lessFunc(d.docs[i], d.docs[j])
}

func (d *docSortHelper[T]) Swap(i, j int) {
	d.docs[i], d.docs[j] = d.docs[j], d.docs[i]
}
