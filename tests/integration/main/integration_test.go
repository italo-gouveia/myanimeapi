package main

/*import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestIntegration(t *testing.T) {
	// Setup the test database
	db, err := sql.Open("postgres", "user=youruser password=yourpassword dbname=yourdbname host=localhost port=5432 sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Setup router
	router := mux.NewRouter()
	RegisterRoutes(router)

	// Create a request
	req := httptest.NewRequest(http.MethodGet, "/anime/1", nil)
	rr := httptest.NewRecorder()

	// Call the handler
	router.ServeHTTP(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)
}
*/
