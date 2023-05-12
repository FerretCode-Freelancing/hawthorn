package containers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ferretcode-freelancing/hawthorn/orchestrator"
)

type Response struct {
	Jobs []orchestrator.Job `json:"job"`
}

func List(w http.ResponseWriter, r *http.Request, o orchestrator.Orchestrator) error {
	session, _ := store.Get(r, "hawthorn")

	if session == nil {
		http.Error(w, "you are not authenticated", http.StatusForbidden)

		return errors.New("you are not authenticated")
	}

	jobs := o.List()

	response := Response{Jobs: jobs}

	stringified, err := json.Marshal(response)

	if err != nil {
		http.Error(w, "there was an error listing containers", http.StatusInternalServerError)

		return err
	}

	w.WriteHeader(200)
	w.Write(stringified)

	return nil
}
