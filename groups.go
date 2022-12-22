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
		return fmt.Errorf("Cannot Create SubGroup Error: %s %s", *Group.Name, err)
	}

	if len(*Group.ClientRoles) > 0 {

		err := c.SetAssignedClientRoles(ctx, token, ID, target, Group)
		if err != nil {
			return err
		}
	}
	if len(*Group.RealmRoles) > 0 {
		err := c.SetAssignedRealmRoles(ctx, token, ID, target, Group)
		if err != nil {
			return err
		}
	}

	if len(*Group.SubGroups) > 0 {
		for _, group := range *Group.SubGroups {
			c.AddChildGroup(ctx, token, target, ID, &group)
		}
	}

	return nil
}

func (c *Instance) SetAssignedClientRoles(ctx context.Context, token string, GroupID string, target *gocloak.GoCloak, Group *gocloak.Group) error {

	for client, myroles := range *Group.ClientRoles {
		roles := []gocloak.Role{}

		ClientID, err := c.GetClientID(client, ctx, token, c.Kc_target.Realm, target)
		if err != nil {
			return err
		}

		if len(myroles) > 0 {
			for _, r := range myroles {
				role, err := target.GetClientRole(ctx, token, c.Kc_target.Realm, ClientID, r)
				if err != nil {
					return fmt.Errorf("Cannot Get Role!: %s %s %s", r, client, err)
				}
				roles = append(roles, *role)
			}
		}

		log.Println("Assigning Roles for Client: ", client)
		err2 := target.AddClientRolesToGroup(ctx, token, c.Kc_target.Realm, ClientID, GroupID, roles)
		if err2 != nil {
			return fmt.Errorf("Cannot Add Roles To Group Error: %s", err2.Error())
		}

	}
	return nil
}

func (c *Instance) SetAssignedRealmRoles(ctx context.Context, token string, GroupID string, target *gocloak.GoCloak, Group *gocloak.Group) error {

	roles := []gocloak.Role{}
	for _, myrole := range *Group.RealmRoles {

		r, err := c.GetRealmRole(myrole, ctx, token, c.Kc_target.Realm, target)
		if err != nil {
			return err
		}

		roles = append(roles, *r)
	}

	err := target.AddRealmRoleToGroup(ctx, token, c.Kc_target.Realm, GroupID, roles)
	if err != nil {
		return fmt.Errorf("Cannot Add Roles To Group Error: %s", err.Error())
	}

	return nil
}

func (c *Instance) AddGroup(ctx context.Context, token string, target *gocloak.GoCloak, Group *gocloak.Group) error {

	ID, err := target.CreateGroup(ctx, token, c.Kc_target.Realm, *Group)

	if err != nil {
		return fmt.Errorf("Cannot Create Group Error: %s", err.Error())
	}

	if len(*Group.ClientRoles) > 0 {
		err := c.SetAssignedClientRoles(ctx, token, ID, target, Group)
		if err != nil {
			return err
		}
	}
	if len(*Group.ClientRoles) > 0 {
		err := c.SetAssignedRealmRoles(ctx, token, ID, target, Group)
		if err != nil {
			return err
		}
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
