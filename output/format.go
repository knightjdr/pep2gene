package output

import (
	"math"
	"path/filepath"
	"sort"

	"github.com/gingraslab/pep2gene/helpers"
	"github.com/gingraslab/pep2gene/types"
)

// Data is a map of genes to their output summary.
type Data struct {
	Database           string           `json:"database"`
	Enzyme             string           `json:"enzyme,omitempty"`
	FDR                float64          `json:"fdr,omitempty"`
	File               string           `json:"file"`
	Genes              map[string]*Gene `json:"genes"`
	MissedCleavages    int              `json:"missedCleavages,omitempty"`
	PeptideProbability float64          `json:"peptideProbability,omitempty"`
	Pipeline           string           `json:"pipeline"`
}

// Gene is summary of a gene for output.
type Gene struct {
	Name          string             `json:"name"`
	Peptides      map[string]Peptide `json:"peptides"`
	SharedIDs     []string           `json:"sharedIDs"`
	SharedNames   []string           `json:"sharedNames"`
	SpectralCount float64            `json:"spectralCount"`
	Subsumed      []string           `json:"subsumed"`
	Unique        int                `json:"unique"`
	UniqueShared  int                `json:"uniqueShared"`
}

// Peptide is summary of a gene's peptides for output.
type Peptide struct {
	AllottedSpectralCount float64 `json:"allottedSpectralCount"`
	TotalSpectralCount    float64 `json:"totalSpectralCount"`
	Unique                bool    `json:"unique"`
	UniqueShared          bool    `json:"uniqueShared"`
}

func summarizePeptides(
	genes []string,
	peptideCount map[string]float64,
	peptides types.Peptides,
	peptideMap map[string]string,
) map[string]Peptide {
	summary := make(map[string]Peptide, 0)
	for peptide, spectralCount := range peptideCount {
		peptideStats := peptides[peptideMap[peptide]]
		unique := false
		uniqueShared := false
		if len(genes) == 1 {
			unique = peptideStats.Unique
		} else if len(genes) > 1 && helpers.SliceEqual(genes, peptideStats.Genes) {
			uniqueShared = true
		}
		newPeptide := Peptide{
			AllottedSpectralCount: spectralCount,
			TotalSpectralCount:    peptideStats.Modified[peptide],
			Unique:                unique,
			UniqueShared:          uniqueShared,
		}
		summary[peptide] = newPeptide
	}

	return summary
}

// Format data for output.
func Format(
	options types.Parameters,
	genes types.Genes,
	geneIDtoName map[string]string,
	peptides types.Peptides,
	peptideMap map[string]string,
) Data {
	summary := Data{
		Database:           filepath.Base(options.Database),
		Enzyme:             options.Enzyme,
		FDR:                options.FDR,
		File:               filepath.Base(options.File),
		Genes:              make(map[string]*Gene, 0),
		MissedCleavages:    options.MissedCleavages,
		PeptideProbability: options.PeptideProbability,
		Pipeline:           options.Pipeline,
	}

	for geneID, details := range genes {
		// For each gene get a list of IDs and names for any genes with shared peptides.
		sharedIDs := make([]string, 0)
		sharedNames := make([]string, len(details.Shared))

		var unique int
		var uniqueShared int
		if len(details.Shared) > 0 {
			// Get gene names for shared gene IDs, sort and add to gene name list.
			nameMap := make(map[string]string, 0)
			for i, gene := range details.Shared {
				sharedNames[i] = geneIDtoName[gene]
				nameMap[geneIDtoName[gene]] = gene
			}
			sort.Strings(sharedNames)

			// Get gene IDs for sorted shared genes
			for _, sharedName := range sharedNames {
				sharedIDs = append(sharedIDs, nameMap[sharedName])
			}
			uniqueShared = details.UniqueShared
		} else {
			unique = int(math.Round(details.Unique))
		}

		peptideDetails := summarizePeptides(append(details.Shared, geneID), details.PeptideCount, peptides, peptideMap)
		sort.Strings(details.Subsumed)

		summary.Genes[geneID] = &Gene{
			Name:          geneIDtoName[geneID],
			Peptides:      peptideDetails,
			SharedIDs:     sharedIDs,
			SharedNames:   sharedNames,
			SpectralCount: details.Count,
			Subsumed:      details.Subsumed,
			Unique:        unique,
			UniqueShared:  uniqueShared,
		}
	}

	return summary
}
