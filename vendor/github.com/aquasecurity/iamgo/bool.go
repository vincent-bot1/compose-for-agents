package iamgo

import "github.com/aquasecurity/jfather"

type Bool struct {
	inner bool
	r     Range
}

func (b *Bool) UnmarshalJSONWithMetadata(node jfather.Node) error {
	b.r.StartLine = node.Range().Start.Line
	b.r.EndLine = node.Range().End.Line
	return node.Decode(&b.inner)
}
