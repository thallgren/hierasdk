package routes

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/lyraproj/hierasdk/hiera"
	"github.com/lyraproj/hierasdk/register"
	"github.com/lyraproj/hierasdk/vf"
)

func TestNoHandler(t *testing.T) {
	register.Clean()
	err := catch(func() error {
		Register()
		return nil
	})
	if err == nil {
		t.Error(`expected panic did not occur`)
	}
}

func TestDataDigHandler(t *testing.T) {
	register.Clean()
	register.DataDig(`my_dd_func`, func(ctx hiera.ProviderContext, key vf.Slice) vf.Data {
		if key.Equal(vf.Slice{vf.String(`config`), vf.String(`path`)}) {
			return vf.String(`/a/b`)
		}
		return nil
	})
	testRequestResponse(t, "/data_dig/my_dd_func", url.Values{`key`: {`["config", "path"]`}}, http.StatusOK, `"/a/b"`)
	testRequestResponse(t, "/data_dig/my_dd_func", url.Values{`key`: {`"config"`}}, http.StatusNotFound, `404 value not found`)
	testRequestResponse(t, "/data_dig/my_rd_func", nil, http.StatusNotFound, `404 page not found`)
}

func TestLookupKeyHandler(t *testing.T) {
	register.Clean()
	register.LookupKey(`my_lk_func`, func(ctx hiera.ProviderContext, key string) vf.Data {
		if key == `host` {
			return ctx.ToData(`example.com`)
		}
		return nil
	})
	testRequestResponse(t, "/lookup_key/my_lk_func", url.Values{`key`: {`host`}}, http.StatusOK, `"example.com"`)
	testRequestResponse(t, "/lookup_key/my_lk_func", url.Values{`key`: {``}}, http.StatusNotFound, `404 value not found`)
	testRequestResponse(t, "/lookup_key/my_lk_func", url.Values{`key`: {`port`}}, http.StatusNotFound, `404 value not found`)
	testRequestResponse(t, "/lookup_key/my_rk_func", url.Values{`key`: {`host`}}, http.StatusNotFound, `404 page not found`)
}

func TestDataHashHandler(t *testing.T) {
	register.Clean()
	register.DataHash(`my_dh_func`, func(ctx hiera.ProviderContext) vf.Data {
		return ctx.ToData(map[string]string{`host`: `example.com`})
	})
	testRequestResponse(t, "/data_hash/my_dh_func", nil, http.StatusOK, `{"host":"example.com"}`)
	testRequestResponse(t, "/data_hash/my_rh_func", nil, http.StatusNotFound, `404 page not found`)
}

func TestDataHashHandler_options(t *testing.T) {
	register.Clean()
	register.DataHash(`my_dh_func`, func(ctx hiera.ProviderContext) vf.Data {
		return ctx.Option(`map_to_deliver`)
	})
	testRequestResponse(t, "/data_hash/my_dh_func",
		url.Values{`options`: {`{"map_to_deliver": {"host": "example.com"}}`}}, http.StatusOK, `{"host":"example.com"}`)
	testRequestResponse(t, "/data_hash/my_dh_func",
		url.Values{`options`: {`{"no_map_to_deliver": {"host": "example.com"}}`}}, http.StatusNotFound, `404 value not found`)
	testRequestResponse(t, "/data_hash/my_dh_func", nil, http.StatusNotFound, `404 value not found`)
}

func TestDataHashHandler_panic(t *testing.T) {
	register.Clean()
	register.DataHash(`my_dh_string_panic`, func(ctx hiera.ProviderContext) vf.Data {
		panic(`goodbye`)
	})
	register.DataHash(`my_dh_error_panic`, func(ctx hiera.ProviderContext) vf.Data {
		panic(errors.New(`goodbye error`))
	})
	register.DataHash(`my_dh_int_panic`, func(ctx hiera.ProviderContext) vf.Data {
		panic(44)
	})
	testRequestResponse(t, "/data_hash/my_dh_string_panic", nil, http.StatusInternalServerError, `goodbye`)
	testRequestResponse(t, "/data_hash/my_dh_error_panic", nil, http.StatusInternalServerError, `goodbye error`)
	testRequestResponse(t, "/data_hash/my_dh_int_panic", nil, http.StatusInternalServerError, `error 44`)
}

func TestDataHashHandler_post(t *testing.T) {
	register.Clean()
	register.DataHash(`my_dh_func`, func(ctx hiera.ProviderContext) vf.Data {
		return ctx.ToData(map[string]string{`host`: `example.com`})
	})
	r, err := http.NewRequest("POST", "/data_hash/my_dh_func", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := Register()
	handler.ServeHTTP(rr, r)
	status := rr.Code
	expectedStatus := http.StatusMethodNotAllowed
	if status != expectedStatus {
		t.Errorf("handler returned wrong status code: got %v want %v", status, expectedStatus)
	}
}

func testRequestResponse(t *testing.T, path string, query url.Values, expectedStatus int, expectedBody string) {
	t.Helper()
	r, err := http.NewRequest("GET", path, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(query) > 0 {
		r.URL.RawQuery = query.Encode()
	}

	rr := httptest.NewRecorder()
	handler := Register()
	handler.ServeHTTP(rr, r)

	status := rr.Code
	if status != expectedStatus {
		t.Errorf("handler returned wrong status code: got %v want %v", status, expectedStatus)
	}

	if status == http.StatusOK {
		expectedType := `application/json`
		actualType := rr.Header().Get(`Content-Type`)
		if expectedType != actualType {
			t.Errorf("handler returned unexpected content path: got %q want %q", actualType, expectedType)
		}
	}

	// Check the response body is what we expect.
	body := strings.TrimSpace(rr.Body.String())
	if body != expectedBody {
		t.Errorf("handler returned unexpected body: got %s want %s", body, expectedBody)
	}
}
