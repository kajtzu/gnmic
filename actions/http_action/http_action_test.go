package http_action

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/karimra/gnmic/actions"
	"github.com/karimra/gnmic/formatters"
)

type item struct {
	input  *formatters.EventMsg
	output map[string]interface{}
}

var testset = map[string]struct {
	actionType string
	action     map[string]interface{}
	tests      []item
}{
	"default_values": {
		actionType: actionType,
		action: map[string]interface{}{
			"type":  "http",
			"URL":   "http://localhost:8080",
			"debug": true,
		},
		tests: []item{
			{
				input: &formatters.EventMsg{
					Name: "sub1",
					Tags: map[string]string{
						"tag1": "1",
					},
				},
				output: map[string]interface{}{
					"name": "sub1",
					"tags": map[string]interface{}{
						"tag1": "1",
					},
				},
			},
		},
	},
}

func TestEventAddTag(t *testing.T) {
	for name, ts := range testset {
		if ai, ok := actions.Actions[ts.actionType]; ok {
			t.Log("found action")
			a := ai()
			err := a.Init(ts.action, nil)
			if err != nil {
				t.Errorf("failed to initialize action: %v", err)
				return
			}
			t.Logf("action: %+v", a)
			mux := http.NewServeMux()
			mux.Handle("/", echo())
			ah, ok := a.(*httpAction)
			if !ok {
				t.Errorf("failed to assert action type: %T", a)
				t.Fail()
				return
			}
			// start http server
			urlAddr, err := url.Parse(ah.URL)
			if err != nil {
				t.Errorf("failed to parse URL: %v", err)
				t.Fail()
				return
			}
			s := &http.Server{
				Addr:    urlAddr.Host,
				Handler: mux,
			}
			go func() {
				if err := s.ListenAndServe(); err != nil {
					if !errors.Is(err, http.ErrServerClosed) {
						t.Logf("failed to start http server: %v", err)
					}
				}
			}()
			// wait for server
			time.Sleep(time.Second)
			//
			for i, item := range ts.tests {
				t.Run(name, func(t *testing.T) {
					t.Logf("running test item %d", i)
					res, err := a.Run(item.input)
					if err != nil {
						t.Errorf("failed at %s item %d, %v", name, i, err)
						t.Fail()
						return
					}
					var result interface{}
					err = json.Unmarshal(res.([]byte), &result)
					if err != nil {
						t.Errorf("failed at %s item %d, %v", name, i, err)
					}
					if !reflect.DeepEqual(result, item.output) {
						t.Errorf("failed at %s item %d, expected %+v, got: %+v", name, i, item.output, result)
					}
				})
			}
			ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
			s.Shutdown(ctx)
			cancel()
		} else {
			t.Errorf("action %s not found", ts.actionType)
		}
	}
}

func echo() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "%v", err)
			return
		}
		//log.Println(string(b))
		fmt.Fprint(w, string(b))
	})
}
