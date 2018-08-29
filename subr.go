// subr.go

/*
Package ssllabs contains SSLLabs-related functions.
*/
package ssllabs

import (
	"encoding/json"
	"log"

	"github.com/pkg/errors"
	"net/http"
)

func myRedirect(req *http.Request, via []*http.Request) error {
	return nil
}

// Display for one report
func (rep *LabsReport) String() {
	host := rep.Host
	if len(rep.Endpoints) != 0 {
		grade := rep.Endpoints[0].Grade
		//details := rep.Endpoints[0].Details
		log.Printf("Looking at %s — grade %s", host, grade)
	}
}

// ParseResults unmarshals the json payload
func ParseResults(content []byte) (r []LabsReport, err error) {
	var data []LabsReport

	err = json.Unmarshal(content, &data)
	return data, errors.Wrap(err, "unmarshal")
}
