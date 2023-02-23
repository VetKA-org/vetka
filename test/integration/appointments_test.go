package integration_test

import (
	"net/http"
	"strings"

	"github.com/VetKA-org/vetka/pkg/entity"
	"github.com/VetKA-org/vetka/pkg/schema"
	"github.com/levigross/grequests"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type AppointmentsTestSuite struct {
	suite.Suite

	token      string
	assigneeID uuid.UUID
}

func (ts *AppointmentsTestSuite) SetupTest() {
	ts.token = doLoginAsHead(ts.T())

	rawToken := strings.TrimPrefix(ts.token, "Bearer ")
	decodedToken, err := entity.DecodeTokenUnverified(rawToken)
	require.NoError(ts.T(), err)

	ts.assigneeID, err = uuid.FromString(decodedToken.Subject)
	require.NoError(ts.T(), err)
}

func (ts *AppointmentsTestSuite) TestCreateAppointmentWithoutRequiredFields() {
	require := ts.Require()

	headers := map[string]string{"Authorization": ts.token}
	req := schema.CreateAppointmentRequest{}
	opts := grequests.RequestOptions{Headers: headers, JSON: req}

	resp := doPostReq(ts.T(), "api/v1/appointments", &opts)
	require.Equal(http.StatusBadRequest, resp.StatusCode)

	failures := extractErrors(ts.T(), resp)
	expected := []string{"patientid", "assigneeid", "scheduledfor", "reason"}
	require.Equal(expected, failures.Fields())
}

func (ts *AppointmentsTestSuite) TestCreateAppointmentWithBadReason() {
	t := ts.T()
	headers := map[string]string{"Authorization": ts.token}

	req := schema.CreateAppointmentRequest{
		PatientID:    uuid.NewV4(),
		AssigneeID:   ts.assigneeID,
		ScheduledFor: *timeFromString(ts.T(), "2023-02-22T10:30:00Z"),
		Reason:       strings.Repeat("a", 256),
	}

	opts := grequests.RequestOptions{Headers: headers, JSON: req}

	resp := doPostReq(t, "api/v1/appointments", &opts)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	failures := extractErrors(t, resp)
	require.Len(t, failures.Errors, 1)
	require.Equal(t, "reason", failures.Errors[0].Field)
}

func (ts *AppointmentsTestSuite) TestCreateAppointmentForUnknownPatient() {
	t := ts.T()
	headers := map[string]string{"Authorization": ts.token}

	req := schema.CreateAppointmentRequest{
		PatientID:    uuid.NewV4(),
		AssigneeID:   ts.assigneeID,
		ScheduledFor: *timeFromString(ts.T(), "2023-02-22T10:30:00Z"),
		Reason:       "Vaccination",
	}

	opts := grequests.RequestOptions{Headers: headers, JSON: req}

	resp := doPostReq(t, "api/v1/appointments", &opts)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)
}
