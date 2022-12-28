package main

import (
	"context"
	"fmt"

	"github.com/Nerzal/gocloak/v12"
)

func (c *Instance) GetClientID(Name string, ctx context.Context, token string, realm string, source *gocloak.GoCloak) (ID string, err error) {

	clients, err := source.GetClients(
		ctx,
		token,
		realm,
		gocloak.GetClientsParams{
			ClientID: &Name,
		},
	)
	if err != nil || len(clients) == 0 {
		fmt.Errorf("Cannot get Client ID Error: %s", err)
		return "", err
	}
	return *clients[0].ID, nil

}

func (c *Instance) GetClient(Name string, ctx context.Context, token string, realm string, source *gocloak.GoCloak) (Client *gocloak.Client, err error) {

	clients, err := source.GetClients(
		ctx,
		token,
		realm,
		gocloak.GetClientsParams{
			ClientID: &Name,
		},
	)
	if err != nil {
		fmt.Errorf("Cannot get Client ID Error: %s", err)
		return nil, err
	}

	if len(clients) > 0 {
		Client = clients[0]
	}
	return Client, nil

}

func (c *Instance) GetClientRoles(Name string, ctx context.Context, token string, source *gocloak.GoCloak) (Roles []*gocloak.Role, err error) {

	ID, err := c.GetClientID(Name, ctx, token, c.Kc_source.Realm, source)
	if err != nil {
		return nil, err
	}
	realm := c.Kc_source.Realm
	Roles = nil

	MyRoles, err := source.GetClientRoles(ctx, token, realm, ID, gocloak.GetRoleParams{})
	if err != nil {
		return nil, err
	}
	for _, role := range MyRoles {

		MyRole, err := source.GetClientRole(ctx, token, realm, ID, *role.Name)
		if err != nil {
			return nil, err
		}
		Roles = append(Roles, MyRole)
	}
	return Roles, nil
}

func (c *Instance) AddClient(ctx context.Context, token string, target *gocloak.GoCloak, Client *gocloak.Client, Roles []*gocloak.Role) error {

	_, err := target.CreateClient(ctx, token, c.Kc_target.Realm, *Client)
	if err != nil {
		return err
	}
	for _, Role := range Roles {
		target.CreateClientRole(ctx, token, c.Kc_target.Realm, *Client.ID, *Role)
	}
	return nil
}

func (c *Instance) UpdateClient(ctx context.Context, token string, target *gocloak.GoCloak, Client *gocloak.Client, Roles []*gocloak.Role) error {
	err := target.UpdateClient(ctx, token, c.Kc_target.Realm, *Client)
	if err != nil {
		return err
	}
	for _, Role := range Roles {
		target.CreateClientRole(ctx, token, c.Kc_target.Realm, *Client.ID, *Role)
		err := target.UpdateRole(ctx, token, c.Kc_target.Realm, *Client.ID, *Role)
		if err != nil {
			return err
		}

	}

	return nil
}
