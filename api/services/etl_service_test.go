package services

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseJikanDate_Nil(t *testing.T) {
	result := parseJikanDate(nil)
	assert.Equal(t, time.Time{}, result)
}

func TestParseJikanDate_Empty(t *testing.T) {
	s := ""
	result := parseJikanDate(&s)
	assert.Equal(t, time.Time{}, result)
}

func TestParseJikanDate_Invalid(t *testing.T) {
	s := "not-a-date"
	result := parseJikanDate(&s)
	assert.Equal(t, time.Time{}, result)
}

func TestParseJikanDate_Valid(t *testing.T) {
	s := "2024-01-15T00:00:00+00:00"
	result := parseJikanDate(&s)
	assert.False(t, result.IsZero())
	assert.Equal(t, 2024, result.Year())
	assert.Equal(t, time.January, result.Month())
	assert.Equal(t, 15, result.Day())
}

func TestNormalizeStatus_FinishedAiring(t *testing.T) {
	assert.Equal(t, "Completed", normalizeStatus("Finished Airing"))
}

func TestNormalizeStatus_CurrentlyAiring(t *testing.T) {
	assert.Equal(t, "Airing", normalizeStatus("Currently Airing"))
}

func TestNormalizeStatus_NotYetAired(t *testing.T) {
	assert.Equal(t, "Upcoming", normalizeStatus("Not yet aired"))
}

func TestNormalizeStatus_Unknown(t *testing.T) {
	assert.Equal(t, "SomeOtherStatus", normalizeStatus("SomeOtherStatus"))
}

func TestNewETLService(t *testing.T) {
	mc := &mockCache{}
	svc := NewETLService(nil, mc)
	assert.NotNil(t, svc)
}
