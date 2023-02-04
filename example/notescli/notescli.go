package notescli

import (
	"context"
	"errors"

	"github.com/soroushj/grpcmock/example/notes"
)

var (
	ErrBadResponse = errors.New("notescli: bad response")
)

type NotesCLI struct {
	client notes.NotesClient
}

func New(client notes.NotesClient) *NotesCLI {
	return &NotesCLI{client}
}

func (nc *NotesCLI) GetNoteText(ctx context.Context, id string) (string, error) {
	resp, err := nc.client.GetNote(ctx, &notes.GetNoteRequest{Id: id})
	if err != nil {
		return "", err
	}
	if resp.Note == nil {
		return "", ErrBadResponse
	}
	return resp.Note.Text, nil
}
