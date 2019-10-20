package firstAssignment

// speciesStorage represents a unified way of accessing species data.
type speciesStorage interface {
	Init()
	Add(s Species) error
	Count() int
	Get(key uint64) (Species, bool)
	GetAll() []Species
}


type ResultList struct {
	AllSpecies[] 					Species `json:"results"`
}


type Species struct {
	Key             			uint64 `json:"key"`
	Kingdom      				  string `json:"kingdom"`
  Phylum       					string `json:"phylum"`
	Order       					string `json:"order"`
	Family      					string `json:"family"`
	Genus     					  string `json:"genus"`
	ScientificName     	  string `json:"scientificName"`
	CanonicalName      		string `json:"canonicalName"`
	Year     						  string `json:"year"`
}

//	Species struct crated in memory
type SpeciesDB struct {
	species map[uint64]Species
}

// Initialided for use
func (db *SpeciesDB) Init() {
	db.species = make(map[uint64]Species)
}

// Adds a new species to the map, into the slot with id = speciesKey (unique value)
func (db *SpeciesDB) Add(s Species) error {
	db.species[s.Key] = s
	return nil
}

// Returns length of the entire map, meaning amount of species stored
func (db *SpeciesDB) Count() int {
	return len(db.species)
}

// Fetches a specifically requested species
func (db *SpeciesDB) Get(keyID uint64) (Species, bool) {
	s, ok := db.species[keyID]
	return s, ok
}

// Fetched all species in an array
func (db *SpeciesDB) GetAll() []Species {
	all := make([]Species, 0, db.Count())
	for _, s := range db.species {
		all = append(all, s)
	}
	return all
}
