package main

import (
	"context"
	"fmt"

	"github.com/Nerzal/gocloak/v12"
)

func (c *Instance) GetRealmRole(Name string, ctx context.Context, token string, realm string, source *gocloak.GoCloak) (*gocloak.Role, error) {
	Role, err := source.GetRealmRole(ctx, token, realm, Name)
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

func (c *Instance) UdateRealmRole(ctx context.Context, token string, target *gocloak.GoCloak, name string, Role *gocloak.Role) error {
	err := target.UpdateRealmRole(ctx, token, c.Kc_target.Realm, name, *Role)
	if err != nil {
		return fmt.Errorf("Cannot Update Realm Role Error: %s", err)
	}
	return nil
}
