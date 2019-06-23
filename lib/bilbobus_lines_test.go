package transit

import (
	"testing"
	"io/ioutil"
	"reflect"
	"log"
	"os"
)

var isNightly = true
var isNotNightly = false

var parseAgencyLinesFileTestCases = []struct {
	source string // input
	expectedLines []Line // expected result
	expectedError bool // shall return an error?
}{
	{`<a href="/wiki/Robotic_spacecraft" title="Robotic spacecraft">robotic spacecraft</a> missions such as <i><a href="/wiki/New_Horizons" title="New Horizons">New Horizons</a></i>;<sup id="cite_ref-NYT-20150828_13-0" class="reference"><a href="#cite_note-NYT-20150828-13">&#91;12&#93;</a></sup> and researching <a href="/wiki/Astrophysics" title="Astrophysics">astrophysics</a> topics, such as the <a href="/wiki/Big_Bang" title="Big Bang">Big Bang</a>, through the <a href="/wiki/Great_Observatories_program" title="Great Observatories program">Great Observatories</a> and associated programs.<sup id="cite_ref-14" class="reference"><a href="#cite_note-14">&#91;13&#93;</a></sup>
</p>
<div id="toc" class="toc"><input type="checkbox" role="button" id="toctogglecheckbox" class="toctogglecheckbox" style="display:none" /><div class="toctitle" lang="en" dir="ltr"><h2>Contents</h2><span class="toctogglespan"><label class="toctogglelabel" for="toctogglecheckbox"></label></span></div>
<ul>
<li class="toclevel-1 tocsection-1"><a href="#Creation"><span class="tocnumber">1</span> <span class="toctext">Creation</span></a></li>
    <div class="select selectlong">
        <select name="urano" id="jupiter" class="select">
            <option value="0000" selected="selected">Seleccione una l&iacute;nea</option>
            <optgroup label="Moon&uacute;aa">
            <option value="03">03 - MOON - PLUTO URANO</option>
            <option value="46">46 - EARTH - SUNSUNSUN</option>
            <option value="A3">A3 - SUN VENUS - SATURN/NEPTUNE</option></optgroup>
            <optgroup label="Planets">
            <option value="G1">G1 - MERCURY - SATURN MARS</option>
        </select>
        <i></i>
    </div>
</section>
<p class="aut"><b>Class:</b></p>`,
		[]Line {
			Line {"I03", "03", 3,"MOON - PLUTO URANO", "FORWARD", nil, nil, &isNotNightly },
			Line {"V03", "03", 3,"PLUTO URANO - MOON", "BACKWARD", nil, nil, &isNotNightly },
			Line {"I46", "46",46,"EARTH - SUNSUNSUN","FORWARD", nil, nil, &isNotNightly },
			Line {"V46", "46",46,"SUNSUNSUN - EARTH","BACKWARD", nil, nil, &isNotNightly },
			Line {"IA3", "A3", 9001, "SUN VENUS - SATURN/NEPTUNE", "FORWARD",nil, nil, &isNotNightly },
			Line {"VA3", "A3", 9001, "SATURN/NEPTUNE - SUN VENUS", "BACKWARD",nil, nil, &isNotNightly },
			Line {"IG1", "G1", 9002, "MERCURY - SATURN MARS", "FORWARD",nil, nil, &isNightly },
			Line {"VG1", "G1", 9002, "SATURN MARS - MERCURY", "BACKWARD",nil, nil, &isNightly },
		},
		false,
	},
	{`<a href="/wiki/Robotic_spacecraft" title="Ro</p>`,
		nil,
		true,
	},
}


func TestParseAgencyLinesFile(t *testing.T) {
	log.Printf("---------- TestParseAgencyLinesFile ------------ ")
	path := "TestParseAgencyLinesFile_source.html"
	for i, tc := range parseAgencyLinesFileTestCases {
		err := ioutil.WriteFile(path, []byte(tc.source), 0644)
		if err!=nil {
			t.Errorf("Error creating file: %v", err)
		}
		lines, err := ParseAgencyLinesFile(path)
		if !linesAreEqual(*lines,tc.expectedLines) || err != nil != tc.expectedError {
			t.Errorf("ParseAgencyLinesFile(#%v): expected (%v), actual (%v)", i, tc.expectedLines, lines)
		}
		os.Remove(path)
	}
	log.Printf("------------------------------------------------ ")
}

// Checks equality for maps of lines
func linesAreEqual (actual, expected []Line) bool {
	if actual == nil && expected ==nil {
		return true
	}
	return reflect.DeepEqual(actual, expected)
}
