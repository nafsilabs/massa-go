package client

import (
	"context"
	"fmt"
	"strconv"
	"strings"
)

type DatastoreEntryResponse struct {
	CandidateValue []byte `json:"candidate_value"`
	FinalValue     []byte `json:"final_value"`
}

type DatastoreEntryData struct {
	Address string        `json:"address"`
	Key     JSONableSlice `json:"key"`
}

type JSONableSlice []byte

func (u JSONableSlice) MarshalJSON() ([]byte, error) {
	if u == nil {
		return []byte("null"), nil
	}

	// Build a JSON array of numeric bytes: [5,0,0,0,...]
	var b strings.Builder
	b.WriteByte('[')
	for i, by := range u {
		if i > 0 {
			b.WriteByte(',')
		}
		b.WriteString(strconv.Itoa(int(by)))
	}
	b.WriteByte(']')
	return []byte(b.String()), nil
}

func NewDatastoreEntry(address string, key []byte) DatastoreEntryData {
	return DatastoreEntryData{
		Address: address,
		Key:     key,
	}
}

func FetchDatastoreEntry(client *Client, address string, key []byte) (*DatastoreEntryResponse, error) {
	entries := make([]DatastoreEntryData, 1)
	entries[0] = NewDatastoreEntry(address, key)

	response, err := FetchDatastoreEntries(client, entries)
	if err != nil {
		return nil, err
	}

	return &response[0], nil
}

func ContractDatastoreEntries(client *Client, address string, keys [][]byte) ([]DatastoreEntryResponse, error) {
	entries := make([]DatastoreEntryData, len(keys))

	for i, key := range keys {
		entries[i] = NewDatastoreEntry(address, key)
	}

	response, err := FetchDatastoreEntries(client, entries)
	if err != nil {
		return nil, fmt.Errorf("calling get_datastore_entries '%+v': %w", entries, err)
	}

	return response, nil
}

func FetchDatastoreEntries(client *Client, entries []DatastoreEntryData) ([]DatastoreEntryResponse, error) {
	data := make([][]DatastoreEntryData, 1)
	data[0] = entries

	response, err := client.RPCClient.Call(
		context.Background(),
		"get_datastore_entries",
		data,
	)
	if err != nil {
		return nil, fmt.Errorf("calling get_datastore_entries '%+v': %w", entries, err)
	}

	if response.Error != nil {
		return nil, response.Error
	}

	var entry []DatastoreEntryResponse

	err = response.GetObject(&entry)
	if err != nil {
		return nil, fmt.Errorf("parsing get_datastore_entries jsonrpc response '%+v': %w", response, err)
	}

	return entry, nil
}
