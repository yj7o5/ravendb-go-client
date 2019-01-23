package ravendb

import "strings"

type IndexCreator interface {
}

// Note: AbstractIndexCreationTask combines functionality of AbstractMultiMapIndexCreationTask

// AbstractIndexCreationTask is for creating an index
// TODO: rename to IndexCreationTask
type AbstractIndexCreationTask struct {

	// for a single map index, set Map
	// for multiple map index, set Maps
	Map  string
	Maps []string

	Reduce string

	Conventions       *DocumentConventions
	AdditionalSources map[string]string
	Priority          IndexPriority
	LockMode          IndexLockMode

	StoresStrings         map[string]FieldStorage
	IndexesStrings        map[string]FieldIndexing
	AnalyzersStrings      map[string]string
	IndexSuggestions      []string
	TermVectorsStrings    map[string]FieldTermVector
	SpatialOptionsStrings map[string]*SpatialOptions

	OutputReduceToCollection string

	// Note: in Go IndexName must provided explicitly
	// In Java it's dynamically calculated as getClass().getSimpleName()
	IndexName string
}

// NewAbstractIndexCreationTask creates AbstractIndexCreationTask
// Note: in Java we subclass AbstractIndexCreationTask and indexName is derived
// from derived class name. In Go we don't subclass and must provide index name
// manually
func NewAbstractIndexCreationTask(indexName string) *AbstractIndexCreationTask {
	panicIf(indexName == "", "indexName cannot be empty")
	return &AbstractIndexCreationTask{
		StoresStrings:         make(map[string]FieldStorage),
		IndexesStrings:        make(map[string]FieldIndexing),
		AnalyzersStrings:      make(map[string]string),
		TermVectorsStrings:    make(map[string]FieldTermVector),
		SpatialOptionsStrings: make(map[string]*SpatialOptions),

		IndexName: indexName,
	}
}

// CreateIndexDefinition creates IndexDefinition
func (t *AbstractIndexCreationTask) CreateIndexDefinition() *IndexDefinition {
	if t.Conventions == nil {
		t.Conventions = NewDocumentConventions()
	}

	indexDefinitionBuilder := NewIndexDefinitionBuilder(t.GetIndexName())
	indexDefinitionBuilder.indexesStrings = t.IndexesStrings
	indexDefinitionBuilder.analyzersStrings = t.AnalyzersStrings
	indexDefinitionBuilder.setMap(t.Map)
	indexDefinitionBuilder.reduce = t.Reduce
	indexDefinitionBuilder.storesStrings = t.StoresStrings
	indexDefinitionBuilder.suggestionsOptions = t.IndexSuggestions
	indexDefinitionBuilder.termVectorsStrings = t.TermVectorsStrings
	indexDefinitionBuilder.spatialIndexesStrings = t.SpatialOptionsStrings
	indexDefinitionBuilder.outputReduceToCollection = t.OutputReduceToCollection
	indexDefinitionBuilder.additionalSources = t.AdditionalSources

	// validate for single map (Map set), don't validate multiple map (Maps)
	validate := len(t.Maps) == 0

	def := indexDefinitionBuilder.toIndexDefinition(t.Conventions, validate)
	for _, m := range t.Maps {
		def.Maps = append(def.Maps, m)
	}
	return def
}

// IsMapReduce returns true if this is map-reduce index
func (t *AbstractIndexCreationTask) IsMapReduce() bool {
	return t.Reduce != ""
}

// GetIndexName returns index name
func (t *AbstractIndexCreationTask) GetIndexName() string {
	panicIf(t.IndexName == "", "indexName must be set by 'sub-class' to be equivalent of Java's getClass().getSimpleName()")
	return strings.Replace(t.IndexName, "_", "/", -1)
}

// Execute executes index in specified document store
// TODO: remove conventions argument, can set AbstractIndexCreationTask.Conventions
func (t *AbstractIndexCreationTask) Execute(store *DocumentStore, conventions *DocumentConventions, database string) error {
	return t.putIndex(store, conventions, database)
}

func (t *AbstractIndexCreationTask) putIndex(store *DocumentStore, conventions *DocumentConventions, database string) error {
	oldConventions := t.Conventions
	defer func() { t.Conventions = oldConventions }()

	conv := conventions
	if conv == nil {
		conv = t.Conventions
	}
	if conv == nil {
		conv = store.GetConventions()
	}
	t.Conventions = conv

	indexDefinition := t.CreateIndexDefinition()
	indexDefinition.Name = t.GetIndexName()
	indexDefinition.LockMode = t.LockMode
	indexDefinition.Priority = t.Priority

	op := NewPutIndexesOperation(indexDefinition)
	if database == "" {
		database = store.GetDatabase()
	}
	return store.Maintenance().ForDatabase(database).Send(op)
}

// Index registers field to be indexed
func (t *AbstractIndexCreationTask) Index(field string, indexing FieldIndexing) {
	t.IndexesStrings[field] = indexing
}

// Spatial registers field to be spatially indexed
func (t *AbstractIndexCreationTask) Spatial(field string, indexing func() *SpatialOptions) {
	v := indexing()
	t.SpatialOptionsStrings[field] = v
}

// StoreAllFields selects if we're storing all fields or not
func (t *AbstractIndexCreationTask) StoreAllFields(storage FieldStorage) {
	t.StoresStrings[IndexingFieldAllFields] = storage
}

// Store registers field to be stored
func (t *AbstractIndexCreationTask) Store(field string, storage FieldStorage) {
	t.StoresStrings[field] = storage
}

// Analyze registers field to be analyzed
func (t *AbstractIndexCreationTask) Analyze(field string, analyzer string) {
	t.AnalyzersStrings[field] = analyzer
}

// TermVector registers field to have term vectors
func (t *AbstractIndexCreationTask) TermVector(field string, termVector FieldTermVector) {
	t.TermVectorsStrings[field] = termVector
}

// Suggestion registers field to be indexed as suggestions
func (t *AbstractIndexCreationTask) Suggestion(field string) {
	t.IndexSuggestions = append(t.IndexSuggestions, field)
}
