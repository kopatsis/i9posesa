package assets

import "go.mongodb.org/mongo-driver/bson/primitive"

// Deprecated
// type ImageSet struct {
// 	ID       primitive.ObjectID `bson:"_id,omitempty"`
// 	Low      []string           `bson:"low"`
// 	Mid      []string           `bson:"mid"`
// 	High     []string           `bson:"high"`
// 	Original []string           `bson:"original"`
// 	Name     string             `bson:"name"`
// }

type ExerPosition struct {
	ImageSetID    string  `bson:"imageset"`
	Hardcoded     bool    `bson:"hardcoded"`
	HardcodedSecs float32 `bson:"hardcodedsecs"`
	MaxSecs       float32 `bson:"maxsecs"` // ?
	PercentSecs   float32 `bson:"percentsecs"`
}

// Exists in DB as part of other entry
type StrPosition struct {
	ImageSetID  string  `bson:"imageset"`
	PercentSecs float32 `bson:"percentsecs"`
}

// Exists in DB as actual entry
type Exercise struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	BackendID      string             `bson:"backendID"`
	Name           string             `bson:"name"`
	Parent         string             `bson:"parent"`
	MaxSecs        float32            `bson:"maxsecs"`
	MinSecs        float32            `bson:"minsecs"`
	ImageSetID0    string             `bson:"imageset0"`
	PositionSlice1 []ExerPosition     `bson:"positions1"`
	PositionSlice2 []ExerPosition     `bson:"positions2"`
	SampleID       string             `bson:"sampleid"`
}

// Exists in DB as actual entry
type DynamicStr struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	BackendID      string             `bson:"backendID"`
	Name           string             `bson:"name"`
	Secs           float32            `bson:"secs"`
	SeparateSets   bool               `bson:"separate"`
	PositionSlice1 []StrPosition      `bson:"positions1"`
	PositionSlice2 []StrPosition      `bson:"positions2"`
	SampleID       string             `bson:"sampleid"`
}

// Exists in DB as actual entry
type StaticStr struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	BackendID   string             `bson:"backendID"`
	Name        string             `bson:"name"`
	ImageSetID1 string             `bson:"imageset1"`
	ImageSetID2 string             `bson:"imageset2"`
	SampleID    string             `bson:"sampleid"`
}

// Exists in DB as part of other entry
type Rep struct {
	Positions []string  `bson:"positions"`
	Times     []float32 `bson:"times"`
	FullTime  float32   `bson:"fulltime"`
}

// Exists in DB as actual entry
type Sample struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name"`
	Description string             `bson:"description"`
	Reps        Rep                `bson:"reps"`
	Type        string             `bson:"type"`
	ExOrStID    string             `bson:"exorstid"`
}

type TransitionRep struct {
	ImageSetIDs []string  `bson:"imagesetids"`
	Times       []float32 `bson:"times"`
	FullTime    float32   `bson:"fulltime"`
}

type TransitionMatrix struct {
	ID            primitive.ObjectID    `bson:"_id,omitempty"`
	FastMatrix    [11][11]TransitionRep `bson:"fastmatrix"`
	RegularMatrix [11][11]TransitionRep `bson:"regularmatrix"`
	SlowMatrix    [11][11]TransitionRep `bson:"slowmatrix"`
}
