/*
Copyright © 2024 Lance haig <lnhaig@gmail.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package internal

import "time"

// ListResponse struct for Returning a list of Dids
type ListResponse struct {
	Cursor string `json:"cursor"`
	Repos  []struct {
		Did  string `json:"did"`
		Head string `json:"head"`
		Rev  string `json:"rev"`
	} `json:"repos"`
}

// RepoDetailResponse struct fro details views
type RepoDetailResponse struct {
	Did       string    `json:"did"`
	Handle    string    `json:"handle"`
	Email     string    `json:"email"`
	IndexedAt time.Time `json:"indexedAt"`
	InvitedBy struct {
		Code       string    `json:"code"`
		Available  int       `json:"available"`
		Disabled   bool      `json:"disabled"`
		ForAccount string    `json:"forAccount"`
		CreatedBy  string    `json:"createdBy"`
		CreatedAt  time.Time `json:"createdAt"`
		Uses       []struct {
			UsedBy string    `json:"usedBy"`
			UsedAt time.Time `json:"usedAt"`
		} `json:"uses"`
	} `json:"invitedBy"`
	Invites         []any `json:"invites"`
	InvitesDisabled bool  `json:"invitesDisabled"`
}

// CreateAccount struct used to create a new account
type CreateAccount struct {
	Email      string `json:"email"`
	Handle     string `json:"handle"`
	Password   string `json:"password"`
	InviteCode string `json:"inviteCode"`
}

// InviteCodeResponse struct
type InviteCodeResponse struct {
	Code string
}

// CreateResponseData struct
type CreateResponseData struct {
	Email    string
	Handle   string
	Password string
	Did      string
}

// DeleteDid struct
type DeleteDid struct {
	Did string `json:"did"`
}

// TakeDown struct
type TakeDown struct {
	Subject struct {
		Type string `json:"$type"`
		Did  string `json:"did"`
	} `json:"subject"`
	Takedown struct {
		Applied bool   `json:"applied"`
		Ref     string `json:"ref"`
	} `json:"takedown"`
}

// UpdatePassword struct
type UpdatePassword struct {
	Did      string `json:"did"`
	Password string `json:"password"`
}

// CrawlHost struct
type CrawlHost struct {
	HostName string `json:"hostname"`
}
