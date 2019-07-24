package main

import (
	"errors"
	"os"
	"testing"

	"github.com/knightjdr/gene-peptide/types"
	"github.com/stretchr/testify/assert"
)

func TestParseFlags(t *testing.T) {
	// Argument unmocking.
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	// TEST1: returns parameters with correct flags specied.
	os.Args = []string{
		"cmd",
		"-db", "database.fasta",
		"-enzyme", "Trypsin",
		"-fdr", "0.05",
		"-file", "peptide.txt",
		"-missedcleavages", "2",
		"-output", "dsv",
		"-pepprob", "0.9",
		"-pipeline", "MSPLIT_DDA",
	}
	wantArgs := types.Parameters{
		Database:           "database.fasta",
		Enzyme:             "trypsin",
		FDR:                0.05,
		File:               "peptide.txt",
		MissedCleavages:    2,
		OutFormat:          "dsv",
		PeptideProbability: 0.9,
		Pipeline:           "MSPLIT_DDA",
	}
	args, err := parseFlags()
	assert.Nil(t, err, "Should not return an error when all required command line arguments are present")
	assert.Equal(t, wantArgs, args, "Should return arguments as parameters")

	// TEST2: returns default parameters when missing.
	os.Args = []string{
		"cmd",
		"-enzyme", "Unknown",
	}
	wantArgs = types.Parameters{
		Database:           "",
		Enzyme:             "",
		FDR:                0.01,
		File:               "",
		MissedCleavages:    0,
		OutFormat:          "tsv",
		PeptideProbability: 0.85,
		Pipeline:           "TPP",
	}
	wantErr := errors.New("missing FASTA database; missing search result peptide file")
	args, err = parseFlags()
	assert.NotNil(t, err, "Should return error when missing arguments")
	assert.Equal(t, wantArgs, args, "Should return default parameters")
	assert.Equal(t, wantErr, err, "Should return correct error message")

	// TEST3: returns default parameters when specifed args are not recognized.
	os.Args = []string{
		"cmd",
		"-db", "database.fasta",
		"-enzyme", "Unknown",
		"-fdr", "0.05",
		"-file", "peptide.txt",
		"-missedcleavages", "2",
		"-output", "Unknown",
		"-pepprob", "0.9",
		"-pipeline", "Unknown",
	}
	args, _ = parseFlags()
	assert.Equal(t, "", args.Enzyme, "Should clear enzyme when arg not recognized")
	assert.Equal(t, "tsv", args.OutFormat, "Should set output format to tsv when arg not recognized")
	assert.Equal(t, "TPP", args.Pipeline, "Should set pipeline to TPP when arg not recognized")
}
