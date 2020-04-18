package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const (
	baseurl = `https://dynalist.io/api/v1/`
)

type Client struct {
	token  string
	client *http.Client
}

func New(token string) (*Client, error) {
	if token == "" {
		return nil, errors.New("Invalid dynalist token provided")
	}
	return &Client{
		token:  token,
		client: &http.Client{},
	}, nil
}

func (api *Client) FileList() (*Response, error) {
	var res Response
	param := struct {
		Token string `json:"token"`
	}{
		Token: api.token,
	}
	if err := api.post("file/list", param, &res); err != nil {
		return nil, err
	}
	if err := res.GetError(); err != nil {
		return nil, err
	}
	return &res, nil
}

func (api *Client) FileEdit(changes []*Change) (*Response, error) {
	var res Response
	param := struct {
		Token   string    `json:"token"`
		Changes []*Change `json:"changes"`
	}{
		Token:   api.token,
		Changes: changes,
	}
	if err := api.post("file/edit", &param, &res); err != nil {
		return nil, err
	}
	if err := res.GetError(); err != nil {
		return nil, err
	}
	return &res, nil
}

func (api *Client) DocRead(fileID string) (*Response, error) {
	var res Response
	param := struct {
		Token  string `json:"token"`
		FileID string `json:"file_id"`
	}{
		Token:  api.token,
		FileID: fileID,
	}
	if err := api.post("doc/read", &param, &res); err != nil {
		return nil, err
	}
	if err := res.GetError(); err != nil {
		return nil, err
	}
	return &res, nil
}

func (api *Client) DocEdit(fileID string, changes []*Change) (*Response, error) {
	var res Response
	param := struct {
		Token   string    `json:"token"`
		FileID  string    `json:"file_id"`
		Changes []*Change `json:"changes"`
	}{
		Token:   api.token,
		FileID:  fileID,
		Changes: changes,
	}
	if err := api.post("doc/edit", &param, &res); err != nil {
		return nil, err
	}
	if err := res.GetError(); err != nil {
		return nil, err
	}
	return &res, nil
}

func (api *Client) InboxAdd(change *Change) (*Response, error) {
	var res Response
	param := struct {
		Change
		Token string `json:"token"`
	}{
		Token:  api.token,
		Change: *change,
	}
	if err := api.post("inbox/add", &param, &res); err != nil {
		return nil, err
	}
	if err := res.GetError(); err != nil {
		return nil, err
	}
	return &res, nil
}

func (api *Client) post(url string, in, out interface{}) error {
	b := &bytes.Buffer{}
	err := json.NewEncoder(b).Encode(in)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", baseurl+url, b)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	res, err := api.client.Do(req)
	if err != nil {
		return err
	}
	if res.Body == nil {
		return errors.New("no body in the response")
	}
	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(out)
	if err != nil {
		return err
	}

	return nil
}

type Change struct {
	Action   Action `json:"action"`
	Index    int    `json:"index,omitempty"`
	NodeID   string `json:"node_id,omitempty"`
	ParentID string `json:"parent_id,omitempty"`
	Content  string `json:"content,omitempty"`
	Type     Type   `json:"type,omitempty"`
	FileID   string `json:"file_id,omitempty"`
	Title    string `json:"title,omitempty"`
	Note     string `json:"note,omitempty"`
	Checked  bool   `json:"checked,omitempty"`
}

func NewChange(action Action) *Change {
	return &Change{Action: action}
}

type Action string

const (
	ActionInsert Action = "action"
	ActionEdit   Action = "edit"
	ActionMove   Action = "move"
	ActionDelete Action = "delete"
)

type Response struct {
	Code       Code   `json:"_code"`
	Msg        string `json:"_msg"`
	RootFileID string `json:"root_file_id,omitempty"`
	Files      []File `json:"files,omitempty"`
	Results    []bool `json:"results,omitempty"`
	Title      string `json:"title,omitempty"`
	Nodes      []Node `json:"nodes,omitempty"`
}

func (r Response) GetError() error {
	if r.Code == CodeOK {
		return nil
	}

	return fmt.Errorf("Dynalist API error %s: %s", r.Code, r.Msg)
}

type Code string

const (
	CodeOK Code = "Ok"
	//Your request is not valid JSON.
	CodeInvalid Code = "Invalid"
	//You've hit the limit on how many requests you can send.
	CodeTooManyRequests Code = "TooManyRequests"
	//Your secret token is invalid.
	CodeInvalidToken Code = "InvalidToken"
	//Server unable to handle the request.
	CodeLockFail Code = "LockFail"
	//You don't have permission to access this document.
	CodeUnauthorized Code = "Unauthorized"
	//The document you're requesting is not found.
	CodeNotFound Code = "NotFound"
	//The node (item) you're requesting is not found.
	CodeNodeNotFound Code = "NodeNotFound"
	//Inbox location is not configured, or invalid.
	CodeNoInbox Code = "NoInbox"
)

type File struct {
	ID         string     `json:"id"`
	Title      string     `json:"title"`
	Type       Type       `json:"type"`
	Permission Permission `json:"permission"`
	Collapsed  bool       `json:"collapsed,omitempty"`
	Children   []string   `json:"children,omitempty"`
}

type Type string

const (
	TypeDocument Type = "document"
	TypeFolder   Type = "folder"
)

type Permission int

const (
	PermissionNoAccess Permission = iota
	PermissionReadOnly
	PermissionEditRights
	PermissionManage
	PermissionOwner
)

type Node struct {
	ID        string   `json:"id"`
	Content   string   `json:"content"`
	Note      string   `json:"note,omitempty"`
	Checked   bool     `json:"checked,omitempty"`
	Collapsed bool     `json:"collapsed,omitempty"`
	Parent    string   `json:"parent,omitempty"`
	Children  []string `json:"children,omitempty"`
}
