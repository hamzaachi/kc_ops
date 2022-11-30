package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Nerzal/gocloak/v12"
	"scm.eadn.dz/DevOps/kc_ops/config"
)

type Instance struct {
	*config.Config
}

func (c *Instance) GetRealmRole(Name string, ctx context.Context, token string, source *gocloak.GoCloak) (*gocloak.Role, error) {
	Role, err := source.GetRealmRole(ctx, token, c.Kc_source.Realm, Name)
	if err != nil {
		return nil, fmt.Errorf("Cannot get Role Info Error: %s", err)
	}
	return Role, nil
}

func (c *Instance) AddRealmRole(ctx context.Context, token string, target *gocloak.GoCloak, Role *gocloak.Role) error {
	_, err := target.CreateRealmRole(ctx, token, c.Kc_target.Realm, *Role)
	if err != nil {
		return fmt.Errorf("Cannot Create Realm Role Error: %s", err)
	}
	return nil
}

func (c *Instance) GetClientID(Name string, ctx context.Context, token string, source *gocloak.GoCloak) (ID string, err error) {

	clients, err := source.GetClients(
		ctx,
		token,
		c.Kc_source.Realm,
		gocloak.GetClientsParams{
			ClientID: &Name,
		},
	)
	if err != nil {
		fmt.Errorf("Cannot get Client ID Error: %s", err)
		return "", err
	}
	return *clients[0].ID, nil

}

func (c *Instance) GetClient(Name string, ctx context.Context, token string, source *gocloak.GoCloak) (Client *gocloak.Client, err error) {

	clients, err := source.GetClients(
		ctx,
		token,
		c.Kc_source.Realm,
		gocloak.GetClientsParams{
			ClientID: &Name,
		},
	)
	if err != nil {
		fmt.Errorf("Cannot get Client ID Error: %s", err)
		return nil, err
	}
	Client = clients[0]
	return Client, nil

}

func (c *Instance) GetClientRoles(Name string, ctx context.Context, token string, source *gocloak.GoCloak) (Roles []*gocloak.Role, err error) {

	ID, err := c.GetClientID(Name, ctx, token, source)
	if err != nil {
		return nil, err
	}
	realm := c.Kc_source.Realm
	Roles = nil

	MyRoles, err := source.GetClientRoles(ctx, token, realm, ID, gocloak.GetRoleParams{})
	if err != nil {
		fmt.Errorf("Cannot get roles error: %", err)
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
		fmt.Println("Oh, no. Cannot create Client on Target: %s", err.Error())
		return err
	}
	for _, Role := range Roles {
		target.CreateClientRole(ctx, token, c.Kc_target.Realm, *Client.ID, *Role)
	}
	return nil
}

func (c *Instance) GetGroup(Name string, ctx context.Context, token string, source *gocloak.GoCloak) (*gocloak.Group, error) {
	groups, err := source.GetGroups(
		ctx,
		token,
		c.Kc_source.Realm,
		gocloak.GetGroupsParams{
			Search: &Name,
		},
	)
	if err != nil {
		fmt.Errorf("Cannot get Group ID Error: %s", err.Error())
		return nil, err
	}

	Group, err := source.GetGroup(ctx, token, c.Kc_source.Realm, *groups[0].ID)
	if err != nil {
		return nil, fmt.Errorf("Cannot get Group Info Error: %s", err.Error())
	}
	return Group, nil
}
func (c *Instance) AddChildGroup(ctx context.Context, token string, target *gocloak.GoCloak, parent string, Group *gocloak.Group) error {

	log.Println("Adding subgroup: ", *Group.Name)
	Group.ID = nil
	ID, err := target.CreateChildGroup(ctx, token, c.Kc_target.Realm, parent, *Group)
	if err != nil {
		return fmt.Errorf("Cannot Create SubGroup Error: %s", err.Error())
	}
	if len(*Group.SubGroups) > 0 {
		for _, group := range *Group.SubGroups {
			c.AddChildGroup(ctx, token, target, ID, &group)
		}
	}
	return nil
}
func (c *Instance) AddGroup(ctx context.Context, token string, target *gocloak.GoCloak, Group *gocloak.Group) error {

	ID, err := target.CreateGroup(ctx, token, c.Kc_target.Realm, *Group)
	if err != nil {
		return fmt.Errorf("Cannot Create Group Error: %s", err.Error())
	}
	if len(*Group.SubGroups) > 0 {
		for _, group := range *Group.SubGroups {
			err := c.AddChildGroup(ctx, token, target, ID, &group)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
