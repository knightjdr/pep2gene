// Package types contains type definitions used throughout pep2gene
package types

import (
	"github.com/gingraslab/pep2gene/helpers"
)

// Gene contains based information: 1) is the gene subsumed by another, 2) The
// (modified) peptides that match to it and their spectral counts, 3) the peptides that
// match to it to, 4) any genes that it completely shares its peptides with, 5) the
// spectral count, 6) any genes it subsumes, 7) the number of unique peptides it
// has (used for distributing spectral counts, so this can be a decimal for genes
// that share peptides), 8) the number of unique peptides shared between genes that
// perfectly share their peptides
type Gene struct {
	Count        float64
	IsSubsumed   bool
	PeptideCount map[string]float64
	Peptides     []string
	Shared       []string
	Subsumed     []string
	Unique       float64
	UniqueShared int
}

// Copy will copy a Gene to a new pointer
func (g Gene) Copy() *Gene {
	copyGeneState := &Gene{
		Count:        g.Count,
		IsSubsumed:   g.IsSubsumed,
		Unique:       g.Unique,
		UniqueShared: g.UniqueShared,
	}
	if g.PeptideCount != nil {
		copyGeneState.PeptideCount = helpers.CopyStringFloatMap(g.PeptideCount)
	}
	if g.Peptides != nil {
		peptides := make([]string, len(g.Peptides))
		copy(peptides, g.Peptides)
		copyGeneState.Peptides = peptides
	}
	if g.Shared != nil {
		shared := make([]string, len(g.Shared))
		copy(shared, g.Shared)
		copyGeneState.Shared = shared
	}
	if g.Subsumed != nil {
		subsumed := make([]string, len(g.Subsumed))
		copy(subsumed, g.Subsumed)
		copyGeneState.Subsumed = subsumed
	}
	return copyGeneState
}

// Genes is a map of gene IDS to their peptide and gene info
type Genes map[string]*Gene

// Parameters for command line arguments.
type Parameters struct {
	Database                    string
	Enzyme                      string
	FDR                         float64
	File                        string
	IgnoreDecoys                bool
	IgnoreInvalid               bool
	InferEnzyme                 bool
	MapFile                     string
	MissedCleavages             int
	Mscore                      float64
	MscorePeptideExperimentWide float64
	OutFormat                   string
	PeakGroupRank               int
	PeptideProbability          float64
	Pipeline                    string
}

// Peptide contains the amino acid "Sequence" for a peptide and the "Modified" version, and
// the intensity when using intensity based data (e.g. OpenSwath)
type Peptide struct {
	Modified  string
	Intensity float64
	Sequence  string
}

// PeptideStat contains the spectral count for a peptide and the individual counts
// for its modified forms
type PeptideStat struct {
	Count    float64
	Genes    []string
	Modified map[string]float64
	Unique   bool
}

// Copy will copy a PeptideStat to a new pointer
func (p PeptideStat) Copy() *PeptideStat {
	copyPeptideState := &PeptideStat{
		Count:  p.Count,
		Unique: p.Unique,
	}
	if p.Genes != nil {
		genes := make([]string, len(p.Genes))
		copy(genes, p.Genes)
		copyPeptideState.Genes = genes
	}
	if p.Modified != nil {
		copyPeptideState.Modified = helpers.CopyStringFloatMap(p.Modified)
	}
	return copyPeptideState
}

// Protein contains the protein name, gene ID and sequence for a protein
type Protein struct {
	GeneID          string
	GeneName        string
	MatchedPeptides []string
	Sequence        string
	Valid           bool
}

// Peptides is a map of peptides to their spectral counts and modified forms
type Peptides map[string]*PeptideStat
