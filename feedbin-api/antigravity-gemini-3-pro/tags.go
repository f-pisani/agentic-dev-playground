package feedbin

import "context"

type tagRenameRequest struct {
	OldName string `json:"old_name"`
	NewName string `json:"new_name"`
}

// RenameTag renames a tag.
func (c *Client) RenameTag(ctx context.Context, oldName, newName string) ([]Tagging, *Response, error) {
	reqBody := &tagRenameRequest{OldName: oldName, NewName: newName}
	req, err := c.NewRequest("POST", "tags.json", reqBody)
	if err != nil {
		return nil, nil, err
	}

	var taggings []Tagging
	resp, err := c.Do(ctx, req, &taggings)
	if err != nil {
		return nil, resp, err
	}

	return taggings, resp, nil
}

type tagDeleteRequest struct {
	Name string `json:"name"`
}

// DeleteTag deletes a tag.
func (c *Client) DeleteTag(ctx context.Context, name string) ([]Tagging, *Response, error) {
	reqBody := &tagDeleteRequest{Name: name}
	req, err := c.NewRequest("DELETE", "tags.json", reqBody)
	if err != nil {
		return nil, nil, err
	}

	var taggings []Tagging
	resp, err := c.Do(ctx, req, &taggings)
	if err != nil {
		return nil, resp, err
	}

	return taggings, resp, nil
}
