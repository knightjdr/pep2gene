package read

import (
	"testing"

	"github.com/gingraslab/pep2gene/fs"
	"github.com/gingraslab/pep2gene/types"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

var msplitDIAText = `#File	Scan#	Mz	z	Peptide	Mz	z	cosine	Name	#Peak(Query)	#Peaks(match)	#shared	fraction-matched	relative-alpha	IonCount	stat0	stat1	stat2	stat3	stat4	stat5	stat6	stat7	stat8	stat9	stat10	stat11	stat12	stat13	stat14	stat15	stat16	stat17	stat18	stat19
file.mzXML	6621	504.171630859375	-1	ABC	500.2589	2	0.954460742970582	1/P23588	250	18	15.0	0.976607265/913863	0.6892903514128037	234504.0	15235.50000150902	649.426	-35.6	0.954437268004796	2.085352659383829	0.9544372680040669	0.9826628787216884	0.8278467816035263	0.7562122650863585	0.7418724241945006	0.8519031016742648	19.0	26.0	14.0	0.0	0.8667127915068145	0.7217749203777925	1.576969258056507	0.00103182792285102fdr	0.0027279812938425566
file.mzXML	6691	504.171630859375	-1	DEF	500.2589	2	0.942778746992883	1/P23588	254	18	13.0	0.9475501296104784	0.3541936460591167	219253.0	15235.50000150902	656.293	-35.6	0.9262740821850953	1.9975947884349832	0.9569986599883877	0.9748245393591848	0.717960653174932	0.6412735691053147	0.5863169362137685	0.7744284006975275	22.0	28.0	12.0	0.0	0.7969855606277485	0.6045790919607626	1.2810047396058777	0.002774825835399693fdr	0.0074917590650284685
file.mzXML	6792	404.9359130859375	-1	GHI	414.7141	2	0.8941625947259864	1/P61978	216	19	11.0	0.8290137307288692	0.5830895432823082	232297.0	4857.700001282378	666.202	-40.8	0.894928911805412	2.4019024816486736	0.8949289118048592	0.9029416031275561	0.8146240447613425	0.6201268787305588	0.4828139240373466	0.7479778926516389	19.0	23.0	10.0	0.0	0.7825492695226175	0.5544942589450436	1.4894842889599753	0.0014542281683996218fdr	0.0032585083272990588
file.mzXML	6838	679.2222900390625	-1	JK[147]L	685.787	2	0.9044174425625956	1/P23588	232	27	16.0	0.885799647019989	1.0	133992.0	83008.50000245772	670.713	-43.29999999999999	0.8893152490266515	2.6160115934413177	0.897489901168478	0.9202925568468254	0.7212759221057556	0.6667879134951098	0.47689314919254777	0.7908414687408863	16.0	38.0	15.0	0.0	0.0	0.6030546194548965	1.7443249120697537	6.566604127579737E-4fdr	0.0018058690744920992
`

func TestMsplitDIARawSequence(t *testing.T) {
	peptide := "JK[147]L"
	expected := "JKL"
	assert.Equal(t, expected, msplitDIARawSequence(peptide), "Should strip modifications from peptide")
}

func TestMsplitDIA(t *testing.T) {
	// Mock fs.
	oldFs := fs.Instance
	defer func() { fs.Instance = oldFs }()
	fs.Instance = afero.NewMemMapFs()

	// Create test directory and files.
	fs.Instance.MkdirAll("test", 0755)
	afero.WriteFile(
		fs.Instance,
		"test/testfile.txt",
		[]byte(msplitDIAText),
		0444,
	)

	file, _ := fs.Instance.Open("test/testfile.txt")
	actualPeptides, actualPeptideMap := msplitDIA(file)

	expectedPeptideMap := map[string]string{
		"ABC":      "ABC",
		"DEF":      "DEF",
		"GHI":      "GHI",
		"JK[147]L": "JKL",
	}
	expectedPeptides := []types.Peptide{
		{Modified: "ABC", Sequence: "ABC"},
		{Modified: "DEF", Sequence: "DEF"},
		{Modified: "GHI", Sequence: "GHI"},
		{Modified: "JK[147]L", Sequence: "JKL"},
	}
	assert.Equal(t, expectedPeptides, actualPeptides, "Should parse correct peptides from file")
	assert.Equal(t, expectedPeptideMap, actualPeptideMap, "Should create a map of modified peptides to raw sequence")
}
