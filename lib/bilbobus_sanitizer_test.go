package transit

import (
	"testing"
	"io/ioutil"
	"reflect"
	"log"
	"os"
)

var parseAgencyLinesFileTestCases = []struct {
	source string // input
	expectedLines map[string]AgencyLine // expected result
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
		map[string]AgencyLine{
			"03": AgencyLine {"03", "MOON - PLUTO URANO"},
			"46": AgencyLine {"46", "EARTH - SUNSUNSUN"},
			"A3": AgencyLine {"A3", "SUN VENUS - SATURN/NEPTUNE"},
			"G1": AgencyLine {"G1", "MERCURY - SATURN MARS"},
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
		if !linesAreEqual(lines,tc.expectedLines) || err != nil != tc.expectedError {
			t.Errorf("ParseAgencyLinesFile(#%v): expected (%v), actual (%v)", i, tc.expectedLines, lines)
		}
		os.Remove(path)
	}
	log.Printf("------------------------------------------------ ")
}

// Checks equality for maps of lines
func linesAreEqual (actual, expected map[string]AgencyLine) bool {
	if actual == nil && expected ==nil {
		return true
	}
	return reflect.DeepEqual(actual, expected)
}
