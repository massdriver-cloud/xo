package massdriver

import (
  "github.com/rs/zerolog/log"
)

// Get a Connection object from Massdriver
func GetConnection(connId string) (string, error) {
  //json, err = client.get(fmt.Sprintf("connection/path?%s", connId))
  var err error
  err = nil
  if err != nil {
    log.Error().Err(err).Msg("")
  }
  json := `{"field1": "value1", "field2": "value2"}`
	return json, err
}

// List Connections from an Organization
func ListConnections(orgId string) (string, error) {
  //json, err = client.get(fmt.Sprintf("connection/org?%s", orgId))
  var err error
  err = nil
  if err != nil {
    log.Error().Err(err).Msg("")
  }
  json := `["connId1", "connId2": "connId3"]`
	return json, err
}
