package integration_test

import (
	"net/http"
	"testing"

	"github.com/VetKA-org/vetka/pkg/entity"
	"github.com/VetKA-org/vetka/pkg/schema"
	"github.com/levigross/grequests"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func getSpecies(t *testing.T, token, title string) uuid.UUID {
	t.Helper()

	headers := map[string]string{"Authorization": token}
	opts := grequests.RequestOptions{Headers: headers}

	resp := doGetReq(t, "api/v1/species?title="+title, &opts)

	require.Equal(t, http.StatusOK, resp.StatusCode)

	body := schema.ListSpeciesResponse{}
	err := resp.JSON(&body)
	require.NoError(t, err)
	require.Equal(t, 1, len(body.Data))

	return body.Data[0].ID
}

func newRegisterPatientReq(t *testing.T, token string) schema.RegisterPatientRequest {
	t.Helper()

	speciesID := getSpecies(t, token, "Cat")

	return schema.RegisterPatientRequest{
		Name:      "murka",
		SpeciesID: speciesID,
		Gender:    entity.Female,
		Breed:     "Munchkin",
		Birth:     *timeFromString(t, "2018-01-10T00:00:00Z"),
	}
}

func listPatients(t *testing.T, token string) schema.ListPatientsResponse {
	t.Helper()

	headers := map[string]string{"Authorization": token}
	opts := grequests.RequestOptions{Headers: headers}

	resp := doGetReq(t, "api/v1/patients", &opts)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	patients := schema.ListPatientsResponse{}
	require.NoError(t, resp.JSON(&patients))

	return patients
}

type PatientsTestSuite struct {
	suite.Suite

	token string
}

func (ts *PatientsTestSuite) SetupTest() {
	ts.token = doLoginAsHead(ts.T())
}

func (ts *PatientsTestSuite) TestRegisterPatientWithoutRequiredFields() {
	require := ts.Require()

	headers := map[string]string{"Authorization": ts.token}
	req := schema.RegisterPatientRequest{}
	opts := grequests.RequestOptions{Headers: headers, JSON: req}

	resp := doPostReq(ts.T(), "api/v1/patients", &opts)
	require.Equal(http.StatusBadRequest, resp.StatusCode)

	failures := extractErrors(ts.T(), resp)
	expected := []string{"name", "speciesid", "gender", "breed", "birth"}
	require.Equal(expected, failures.Fields())
}

func (ts *PatientsTestSuite) TestRegisterPatientWithBadName() {
	t := ts.T()
	headers := map[string]string{"Authorization": ts.token}

	req := newRegisterPatientReq(t, ts.token)
	req.Name = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

	opts := grequests.RequestOptions{Headers: headers, JSON: req}

	resp := doPostReq(t, "api/v1/patients", &opts)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	failures := extractErrors(t, resp)
	require.Len(t, failures.Errors, 1)
	require.Equal(t, "name", failures.Errors[0].Field)
}

func (ts *PatientsTestSuite) TestRegisterPatientWithBadGender() {
	t := ts.T()
	headers := map[string]string{"Authorization": ts.token}

	req := newRegisterPatientReq(t, ts.token)
	req.Gender = "?"

	opts := grequests.RequestOptions{Headers: headers, JSON: req}

	resp := doPostReq(t, "api/v1/patients", &opts)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	failures := extractErrors(t, resp)
	require.Len(t, failures.Errors, 1)
	require.Equal(t, "gender", failures.Errors[0].Field)
}

func (ts *PatientsTestSuite) TestRegisterPatientWithBadBreed() {
	t := ts.T()
	headers := map[string]string{"Authorization": ts.token}

	req := newRegisterPatientReq(t, ts.token)
	req.Breed = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

	opts := grequests.RequestOptions{Headers: headers, JSON: req}

	resp := doPostReq(t, "api/v1/patients", &opts)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	failures := extractErrors(t, resp)
	require.Len(t, failures.Errors, 1)
	require.Equal(t, "breed", failures.Errors[0].Field)
}

func (ts *PatientsTestSuite) TestRegisterPatientWithBadVaccinationDate() {
	t := ts.T()
	headers := map[string]string{"Authorization": ts.token}

	req := newRegisterPatientReq(t, ts.token)
	req.Birth = *timeFromString(t, "2020-01-10T00:00:00Z")
	req.VaccinatedAt = timeFromString(t, "2019-01-10T00:00:00Z")

	opts := grequests.RequestOptions{Headers: headers, JSON: req}

	resp := doPostReq(t, "api/v1/patients", &opts)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	failures := extractErrors(t, resp)
	require.Len(t, failures.Errors, 1)
	require.Equal(t, "vaccinatedat", failures.Errors[0].Field)
}

func (ts *PatientsTestSuite) TestRegisterPatientWithBadSterilizationDate() {
	t := ts.T()
	headers := map[string]string{"Authorization": ts.token}

	req := newRegisterPatientReq(t, ts.token)
	req.Birth = *timeFromString(t, "2020-01-10T00:00:00Z")
	req.SterilizedAt = timeFromString(t, "2019-01-10T00:00:00Z")

	opts := grequests.RequestOptions{Headers: headers, JSON: req}

	resp := doPostReq(t, "api/v1/patients", &opts)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	failures := extractErrors(t, resp)
	require.Len(t, failures.Errors, 1)
	require.Equal(t, "sterilizedat", failures.Errors[0].Field)
}

func (ts *PatientsTestSuite) TestRegisterPatientWithUnknownSpecies() {
	t := ts.T()
	headers := map[string]string{"Authorization": ts.token}

	req := newRegisterPatientReq(t, ts.token)
	req.SpeciesID = uuid.NewV4()
	opts := grequests.RequestOptions{Headers: headers, JSON: req}

	resp := doPostReq(t, "api/v1/patients", &opts)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func (ts *PatientsTestSuite) TestRegisterPatientWithMinimalInfo() {
	t := ts.T()
	require := ts.Require()
	headers := map[string]string{"Authorization": ts.token}

	req := newRegisterPatientReq(t, ts.token)
	req.Name = "Minimal info"
	opts := grequests.RequestOptions{Headers: headers, JSON: req}

	resp := doPostReq(t, "api/v1/patients", &opts)
	require.Equal(http.StatusOK, resp.StatusCode)

	var created *entity.Patient

	patients := listPatients(t, ts.token)
	for i := range patients.Data {
		if patients.Data[i].Name == req.Name {
			created = &patients.Data[i]

			break
		}
	}

	require.NotNil(created)
	require.NotEqual(uuid.UUID{}, created.ID)
	require.Equal(req.Name, created.Name)
	require.Equal(req.SpeciesID, created.SpeciesID)
	require.Equal(req.Gender, created.Gender)
	require.Equal(req.Breed, created.Breed)
	require.False(req.Aggressive)
	require.Nil(req.VaccinatedAt)
	require.Nil(req.SterilizedAt)
}

func (ts *PatientsTestSuite) TestRegisterPatientWithFullInfo() {
	t := ts.T()
	require := ts.Require()
	headers := map[string]string{"Authorization": ts.token}

	req := newRegisterPatientReq(t, ts.token)
	req.Name = "Full info"
	req.Aggressive = true
	req.VaccinatedAt = timeFromString(t, "2019-01-10T00:00:00Z")
	req.SterilizedAt = timeFromString(t, "2019-01-10T00:00:00Z")

	opts := grequests.RequestOptions{Headers: headers, JSON: req}

	resp := doPostReq(t, "api/v1/patients", &opts)
	require.Equal(http.StatusOK, resp.StatusCode)

	var created *entity.Patient

	patients := listPatients(t, ts.token)
	for i := range patients.Data {
		if patients.Data[i].Name == req.Name {
			created = &patients.Data[i]

			break
		}
	}

	require.NotNil(created)
	require.NotEqual(uuid.UUID{}, created.ID)
	require.Equal(req.Name, created.Name)
	require.Equal(req.SpeciesID, created.SpeciesID)
	require.Equal(req.Gender, created.Gender)
	require.Equal(req.Breed, created.Breed)
	require.Equal(req.Aggressive, created.Aggressive)
	require.Equal(req.VaccinatedAt, created.VaccinatedAt)
	require.Equal(req.SterilizedAt, created.SterilizedAt)
}
