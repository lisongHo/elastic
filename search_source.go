// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
	"fmt"
)

// SearchSource enables users to build the search source.
// It resembles the SearchSourceBuilder in Elasticsearch.
type SearchSource struct {
	query                    Query                  // query
	postQuery                Query                  // post_filter
	sliceQuery               Query                  // slice
	from                     int                    // from
	size                     int                    // size
	explain                  *bool                  // explain
	version                  *bool                  // version
	seqNoAndPrimaryTerm      *bool                  // seq_no_primary_term
	sorters                  []Sorter               // sort
	trackScores              *bool                  // track_scores
	trackTotalHits           interface{}            // track_total_hits
	searchAfterSortValues    []interface{}          // search_after
	minScore                 *float64               // min_score
	timeout                  string                 // timeout
	terminateAfter           *int                   // terminate_after
	storedFieldNames         []string               // stored_fields
	docvalueFields           DocvalueFields         // docvalue_fields
	scriptFields             []*ScriptField         // script_fields
	fetchSourceContext       *FetchSourceContext    // _source
	aggregations             map[string]Aggregation // aggregations / aggs
	highlight                *Highlight             // highlight
	globalSuggestText        string
	suggesters               []Suggester // suggest
	rescores                 []*Rescore  // rescore
	defaultRescoreWindowSize *int
	indexBoosts              IndexBoosts // indices_boost
	stats                    []string    // stats
	innerHits                map[string]*InnerHit
	collapse                 *CollapseBuilder // collapse
	profile                  bool             // profile
	// TODO extBuilders []SearchExtBuilder // ext
	pointInTime     *PointInTime // pit
	runtimeMappings RuntimeMappings
}

// NewSearchSource initializes a new SearchSource.
func NewSearchSource() *SearchSource {
	return &SearchSource{
		from:         -1,
		size:         -1,
		aggregations: make(map[string]Aggregation),
		innerHits:    make(map[string]*InnerHit),
	}
}

// Query sets the query to use with this search source.
func (s *SearchSource) Query(query Query) *SearchSource {
	s.query = query
	return s
}

// Profile specifies that this search source should activate the
// Profile API for queries made on it.
func (s *SearchSource) Profile(profile bool) *SearchSource {
	s.profile = profile
	return s
}

// PostFilter will be executed after the query has been executed and
// only affects the search hits, not the aggregations.
// This filter is always executed as the last filtering mechanism.
func (s *SearchSource) PostFilter(postFilter Query) *SearchSource {
	s.postQuery = postFilter
	return s
}

// Slice allows partitioning the documents in multiple slices.
// It is e.g. used to slice a scroll operation, supported in
// Elasticsearch 5.0 or later.
// See https://www.elastic.co/guide/en/elasticsearch/reference/7.0/search-request-scroll.html#sliced-scroll
// for details.
func (s *SearchSource) Slice(sliceQuery Query) *SearchSource {
	s.sliceQuery = sliceQuery
	return s
}

// From index to start the search from. Defaults to 0.
func (s *SearchSource) From(from int) *SearchSource {
	s.from = from
	return s
}

// Size is the number of search hits to return. Defaults to 10.
func (s *SearchSource) Size(size int) *SearchSource {
	s.size = size
	return s
}

// MinScore sets the minimum score below which docs will be filtered out.
func (s *SearchSource) MinScore(minScore float64) *SearchSource {
	s.minScore = &minScore
	return s
}

// Explain indicates whether each search hit should be returned with
// an explanation of the hit (ranking).
func (s *SearchSource) Explain(explain bool) *SearchSource {
	s.explain = &explain
	return s
}

// Version indicates whether each search hit should be returned with
// a version associated to it.
func (s *SearchSource) Version(version bool) *SearchSource {
	s.version = &version
	return s
}

// SeqNoAndPrimaryTerm indicates whether SearchHits should be returned with the
// sequence number and primary term of the last modification of the document.
func (s *SearchSource) SeqNoAndPrimaryTerm(enabled bool) *SearchSource {
	s.seqNoAndPrimaryTerm = &enabled
	return s
}

// Timeout controls how long a search is allowed to take, e.g. "1s" or "500ms".
func (s *SearchSource) Timeout(timeout string) *SearchSource {
	s.timeout = timeout
	return s
}

// TimeoutInMillis controls how many milliseconds a search is allowed
// to take before it is canceled.
func (s *SearchSource) TimeoutInMillis(timeoutInMillis int) *SearchSource {
	s.timeout = fmt.Sprintf("%dms", timeoutInMillis)
	return s
}

// TerminateAfter specifies the maximum number of documents to collect for
// each shard, upon reaching which the query execution will terminate early.
func (s *SearchSource) TerminateAfter(terminateAfter int) *SearchSource {
	s.terminateAfter = &terminateAfter
	return s
}

// Sort adds a sort order.
func (s *SearchSource) Sort(field string, ascending bool) *SearchSource {
	s.sorters = append(s.sorters, SortInfo{Field: field, Ascending: ascending})
	return s
}

// SortWithInfo adds a sort order.
func (s *SearchSource) SortWithInfo(info SortInfo) *SearchSource {
	s.sorters = append(s.sorters, info)
	return s
}

// SortBy	adds a sort order.
func (s *SearchSource) SortBy(sorter ...Sorter) *SearchSource {
	s.sorters = append(s.sorters, sorter...)
	return s
}

func (s *SearchSource) hasSort() bool {
	return len(s.sorters) > 0
}

// TrackScores is applied when sorting and controls if scores will be
// tracked as well. Defaults to false.
func (s *SearchSource) TrackScores(trackScores bool) *SearchSource {
	s.trackScores = &trackScores
	return s
}

// TrackTotalHits controls how the total number of hits should be tracked.
// Defaults to 10000 which will count the total hit accurately up to 10,000 hits.
//
// See https://www.elastic.co/guide/en/elasticsearch/reference/7.0/search-request-track-total-hits.html
// for details.
func (s *SearchSource) TrackTotalHits(trackTotalHits interface{}) *SearchSource {
	s.trackTotalHits = trackTotalHits
	return s
}

// SearchAfter allows a different form of pagination by using a live cursor,
// using the results of the previous page to help the retrieval of the next.
//
// See https://www.elastic.co/guide/en/elasticsearch/reference/7.0/search-request-search-after.html
func (s *SearchSource) SearchAfter(sortValues ...interface{}) *SearchSource {
	for _, v := range sortValues {
		// To search after a null date , we need to convert it to the java Long.MIN_VALUE
		if num, ok := v.(float64); ok && num == -9223372036854776000 {
			v = -9223372036854775808
		}
		s.searchAfterSortValues = append(s.searchAfterSortValues, v)
	}
	return s
}

// Aggregation adds an aggreation to perform as part of the search.
func (s *SearchSource) Aggregation(name string, aggregation Aggregation) *SearchSource {
	s.aggregations[name] = aggregation
	return s
}

// DefaultRescoreWindowSize sets the rescore window size for rescores
// that don't specify their window.
func (s *SearchSource) DefaultRescoreWindowSize(defaultRescoreWindowSize int) *SearchSource {
	s.defaultRescoreWindowSize = &defaultRescoreWindowSize
	return s
}

// Highlight adds highlighting to the search.
func (s *SearchSource) Highlight(highlight *Highlight) *SearchSource {
	s.highlight = highlight
	return s
}

// Highlighter returns the highlighter.
func (s *SearchSource) Highlighter() *Highlight {
	if s.highlight == nil {
		s.highlight = NewHighlight()
	}
	return s.highlight
}

// GlobalSuggestText defines the global text to use with all suggesters.
// This avoids repetition.
func (s *SearchSource) GlobalSuggestText(text string) *SearchSource {
	s.globalSuggestText = text
	return s
}

// Suggester adds a suggester to the search.
func (s *SearchSource) Suggester(suggester Suggester) *SearchSource {
	s.suggesters = append(s.suggesters, suggester)
	return s
}

// Rescorer adds a rescorer to the search.
func (s *SearchSource) Rescorer(rescore *Rescore) *SearchSource {
	s.rescores = append(s.rescores, rescore)
	return s
}

// ClearRescorers removes all rescorers from the search.
func (s *SearchSource) ClearRescorers() *SearchSource {
	s.rescores = make([]*Rescore, 0)
	return s
}

// FetchSource indicates whether the response should contain the stored
// _source for every hit.
func (s *SearchSource) FetchSource(fetchSource bool) *SearchSource {
	if s.fetchSourceContext == nil {
		s.fetchSourceContext = NewFetchSourceContext(fetchSource)
	} else {
		s.fetchSourceContext.SetFetchSource(fetchSource)
	}
	return s
}

// FetchSourceContext indicates how the _source should be fetched.
func (s *SearchSource) FetchSourceContext(fetchSourceContext *FetchSourceContext) *SearchSource {
	s.fetchSourceContext = fetchSourceContext
	return s
}

// FetchSourceIncludeExclude specifies that _source should be returned
// with each hit, where "include" and "exclude" serve as a simple wildcard
// matcher that gets applied to its fields
// (e.g. include := []string{"obj1.*","obj2.*"}, exclude := []string{"description.*"}).
func (s *SearchSource) FetchSourceIncludeExclude(include, exclude []string) *SearchSource {
	s.fetchSourceContext = NewFetchSourceContext(true).
		Include(include...).
		Exclude(exclude...)
	return s
}

// NoStoredFields indicates that no fields should be loaded, resulting in only
// id and type to be returned per field.
func (s *SearchSource) NoStoredFields() *SearchSource {
	s.storedFieldNames = []string{}
	return s
}

// StoredField adds a single field to load and return (note, must be stored) as
// part of the search request. If none are specified, the source of the
// document will be returned.
func (s *SearchSource) StoredField(storedFieldName string) *SearchSource {
	s.storedFieldNames = append(s.storedFieldNames, storedFieldName)
	return s
}

// StoredFields	sets the fields to load and return as part of the search request.
// If none are specified, the source of the document will be returned.
func (s *SearchSource) StoredFields(storedFieldNames ...string) *SearchSource {
	s.storedFieldNames = append(s.storedFieldNames, storedFieldNames...)
	return s
}

// DocvalueField adds a single field to load from the field data cache
// and return as part of the search request.
func (s *SearchSource) DocvalueField(fieldDataField string) *SearchSource {
	s.docvalueFields = append(s.docvalueFields, DocvalueField{Field: fieldDataField})
	return s
}

// DocvalueField adds a single docvalue field to load from the field data cache
// and return as part of the search request.
func (s *SearchSource) DocvalueFieldWithFormat(fieldDataFieldWithFormat DocvalueField) *SearchSource {
	s.docvalueFields = append(s.docvalueFields, fieldDataFieldWithFormat)
	return s
}

// DocvalueFields adds one or more fields to load from the field data cache
// and return as part of the search request.
func (s *SearchSource) DocvalueFields(docvalueFields ...string) *SearchSource {
	for _, f := range docvalueFields {
		s.docvalueFields = append(s.docvalueFields, DocvalueField{Field: f})
	}
	return s
}

// DocvalueFields adds one or more docvalue fields to load from the field data cache
// and return as part of the search request.
func (s *SearchSource) DocvalueFieldsWithFormat(docvalueFields ...DocvalueField) *SearchSource {
	s.docvalueFields = append(s.docvalueFields, docvalueFields...)
	return s
}

// ScriptField adds a single script field with the provided script.
func (s *SearchSource) ScriptField(scriptField *ScriptField) *SearchSource {
	s.scriptFields = append(s.scriptFields, scriptField)
	return s
}

// ScriptFields adds one or more script fields with the provided scripts.
func (s *SearchSource) ScriptFields(scriptFields ...*ScriptField) *SearchSource {
	s.scriptFields = append(s.scriptFields, scriptFields...)
	return s
}

// IndexBoost sets the boost that a specific index will receive when the
// query is executed against it.
func (s *SearchSource) IndexBoost(index string, boost float64) *SearchSource {
	s.indexBoosts = append(s.indexBoosts, IndexBoost{Index: index, Boost: boost})
	return s
}

// IndexBoosts sets the boosts for specific indices.
func (s *SearchSource) IndexBoosts(boosts ...IndexBoost) *SearchSource {
	s.indexBoosts = append(s.indexBoosts, boosts...)
	return s
}

// Stats group this request will be aggregated under.
func (s *SearchSource) Stats(statsGroup ...string) *SearchSource {
	s.stats = append(s.stats, statsGroup...)
	return s
}

// InnerHit adds an inner hit to return with the result.
func (s *SearchSource) InnerHit(name string, innerHit *InnerHit) *SearchSource {
	s.innerHits[name] = innerHit
	return s
}

// Collapse adds field collapsing.
func (s *SearchSource) Collapse(collapse *CollapseBuilder) *SearchSource {
	s.collapse = collapse
	return s
}

// PointInTime specifies an optional PointInTime to be used in the context
// of this search.
func (s *SearchSource) PointInTime(pointInTime *PointInTime) *SearchSource {
	s.pointInTime = pointInTime
	return s
}

// RuntimeMappings specifies optional runtime mappings.
func (s *SearchSource) RuntimeMappings(runtimeMappings RuntimeMappings) *SearchSource {
	s.runtimeMappings = runtimeMappings
	return s
}

// Source returns the serializable JSON for the source builder.
func (s *SearchSource) Source() (interface{}, error) {
	source := make(map[string]interface{})

	if s.from != -1 {
		source["from"] = s.from
	}
	if s.size != -1 {
		source["size"] = s.size
	}
	if s.timeout != "" {
		source["timeout"] = s.timeout
	}
	if s.terminateAfter != nil {
		source["terminate_after"] = *s.terminateAfter
	}
	if s.query != nil {
		src, err := s.query.Source()
		if err != nil {
			return nil, err
		}
		source["query"] = src
	}
	if s.postQuery != nil {
		src, err := s.postQuery.Source()
		if err != nil {
			return nil, err
		}
		source["post_filter"] = src
	}
	if s.minScore != nil {
		source["min_score"] = *s.minScore
	}
	if s.version != nil {
		source["version"] = *s.version
	}
	if s.explain != nil {
		source["explain"] = *s.explain
	}
	if s.profile {
		source["profile"] = s.profile
	}
	if s.fetchSourceContext != nil {
		src, err := s.fetchSourceContext.Source()
		if err != nil {
			return nil, err
		}
		source["_source"] = src
	}
	if s.storedFieldNames != nil {
		switch len(s.storedFieldNames) {
		case 1:
			source["stored_fields"] = s.storedFieldNames[0]
		default:
			source["stored_fields"] = s.storedFieldNames
		}
	}
	if len(s.docvalueFields) > 0 {
		src, err := s.docvalueFields.Source()
		if err != nil {
			return nil, err
		}
		source["docvalue_fields"] = src
	}
	if len(s.scriptFields) > 0 {
		sfmap := make(map[string]interface{})
		for _, scriptField := range s.scriptFields {
			src, err := scriptField.Source()
			if err != nil {
				return nil, err
			}
			sfmap[scriptField.FieldName] = src
		}
		source["script_fields"] = sfmap
	}
	if len(s.sorters) > 0 {
		var sortarr []interface{}
		for _, sorter := range s.sorters {
			src, err := sorter.Source()
			if err != nil {
				return nil, err
			}
			sortarr = append(sortarr, src)
		}
		source["sort"] = sortarr
	}
	if v := s.trackScores; v != nil {
		source["track_scores"] = *v
	}
	if v := s.trackTotalHits; v != nil {
		source["track_total_hits"] = v
	}
	if len(s.searchAfterSortValues) > 0 {
		source["search_after"] = s.searchAfterSortValues
	}
	if s.sliceQuery != nil {
		src, err := s.sliceQuery.Source()
		if err != nil {
			return nil, err
		}
		source["slice"] = src
	}
	if len(s.indexBoosts) > 0 {
		src, err := s.indexBoosts.Source()
		if err != nil {
			return nil, err
		}
		source["indices_boost"] = src
	}
	if len(s.aggregations) > 0 {
		aggsMap := make(map[string]interface{})
		for name, aggregate := range s.aggregations {
			src, err := aggregate.Source()
			if err != nil {
				return nil, err
			}
			aggsMap[name] = src
		}
		source["aggregations"] = aggsMap
	}
	if s.highlight != nil {
		src, err := s.highlight.Source()
		if err != nil {
			return nil, err
		}
		source["highlight"] = src
	}
	if len(s.suggesters) > 0 {
		suggesters := make(map[string]interface{})
		for _, s := range s.suggesters {
			src, err := s.Source(false)
			if err != nil {
				return nil, err
			}
			suggesters[s.Name()] = src
		}
		if s.globalSuggestText != "" {
			suggesters["text"] = s.globalSuggestText
		}
		source["suggest"] = suggesters
	}
	if len(s.rescores) > 0 {
		// Strip empty rescores from request
		var rescores []*Rescore
		for _, r := range s.rescores {
			if !r.IsEmpty() {
				rescores = append(rescores, r)
			}
		}
		if len(rescores) == 1 {
			rescores[0].defaultRescoreWindowSize = s.defaultRescoreWindowSize
			src, err := rescores[0].Source()
			if err != nil {
				return nil, err
			}
			source["rescore"] = src
		} else {
			var slice []interface{}
			for _, r := range rescores {
				r.defaultRescoreWindowSize = s.defaultRescoreWindowSize
				src, err := r.Source()
				if err != nil {
					return nil, err
				}
				slice = append(slice, src)
			}
			source["rescore"] = slice
		}
	}
	if len(s.stats) > 0 {
		source["stats"] = s.stats
	}
	// TODO ext builders

	if s.collapse != nil {
		src, err := s.collapse.Source()
		if err != nil {
			return nil, err
		}
		source["collapse"] = src
	}

	if v := s.seqNoAndPrimaryTerm; v != nil {
		source["seq_no_primary_term"] = *v
	}

	if len(s.innerHits) > 0 {
		// Top-level inner hits
		// See http://www.elastic.co/guide/en/elasticsearch/reference/1.5/search-request-inner-hits.html#top-level-inner-hits
		// "inner_hits": {
		//   "<inner_hits_name>": {
		//     "<path|type>": {
		//       "<path-to-nested-object-field|child-or-parent-type>": {
		//         <inner_hits_body>,
		//         [,"inner_hits" : { [<sub_inner_hits>]+ } ]?
		//       }
		//     }
		//   },
		//   [,"<inner_hits_name_2>" : { ... } ]*
		// }
		m := make(map[string]interface{})
		for name, hit := range s.innerHits {
			if hit.path != "" {
				src, err := hit.Source()
				if err != nil {
					return nil, err
				}
				path := make(map[string]interface{})
				path[hit.path] = src
				m[name] = map[string]interface{}{
					"path": path,
				}
			} else if hit.typ != "" {
				src, err := hit.Source()
				if err != nil {
					return nil, err
				}
				typ := make(map[string]interface{})
				typ[hit.typ] = src
				m[name] = map[string]interface{}{
					"type": typ,
				}
			} else {
				// TODO the Java client throws here, because either path or typ must be specified
				_ = m
			}
		}
		source["inner_hits"] = m
	}

	// Point in Time
	if s.pointInTime != nil {
		src, err := s.pointInTime.Source()
		if err != nil {
			return nil, err
		}
		source["pit"] = src
	}

	if s.runtimeMappings != nil {
		src, err := s.runtimeMappings.Source()
		if err != nil {
			return nil, err
		}
		source["runtime_mappings"] = src
	}

	return source, nil
}

// MarshalJSON enables serializing the type as JSON.
func (q *SearchSource) MarshalJSON() ([]byte, error) {
	if q == nil {
		return nilByte, nil
	}
	src, err := q.Source()
	if err != nil {
		return nil, err
	}
	return json.Marshal(src)
}

// -- IndexBoosts --

// IndexBoost specifies an index by some boost factor.
type IndexBoost struct {
	Index string
	Boost float64
}

// Source generates a JSON-serializable output for IndexBoost.
func (b IndexBoost) Source() (interface{}, error) {
	return map[string]interface{}{
		b.Index: b.Boost,
	}, nil
}

// IndexBoosts is a slice of IndexBoost entities.
type IndexBoosts []IndexBoost

// Source generates a JSON-serializable output for IndexBoosts.
func (b IndexBoosts) Source() (interface{}, error) {
	var boosts []interface{}
	for _, ib := range b {
		src, err := ib.Source()
		if err != nil {
			return nil, err
		}
		boosts = append(boosts, src)
	}
	return boosts, nil
}
