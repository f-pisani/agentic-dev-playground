package feedbin

import "net/http"

// IconsService handles communication with the icon related
// methods of the Feedbin API.
type IconsService service

// List retrieves the icons (favicons) for all feeds the user is subscribed to.
// Feedbin API docs: https://github.com/feedbin/feedbin-api/blob/master/content/icons.md#get-v2iconsjson
func (s *IconsService) List() ([]*Icon, *http.Response, error) {
	req, err := s.client.newRequest(http.MethodGet, "icons.json", nil)
	if err != nil {
		return nil, nil, err
	}

	var icons []*Icon
	resp, err := s.client.do(req, &icons)
	if err != nil {
		return nil, resp, err
	}

	return icons, resp, nil
}
