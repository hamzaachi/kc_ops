package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Nerzal/gocloak/v12"
	"scm.eadn.dz/DevOps/kc_ops/config"
)

func KC_AddClients(ctx context.Context, instance Instance, source, target *gocloak.GoCloak) error {
	for _, client := range instance.Kc_source.Clients {

		token, err := source.LoginAdmin(ctx, instance.Kc_source.Username, instance.Kc_source.Password, "master")

		log.Println("Getting Client Info...", client)

		Roles, err := instance.GetClientRoles(client, ctx, token.AccessToken, source)
		if err != nil {
			fmt.Errorf("Something wrong with getting roles: %s", err)
			return err
		}

		kc_client, err := instance.GetClient(client, ctx, token.AccessToken, source)
		if err != nil {
			fmt.Errorf("Something wrong with getting client value: %s", err)
			return err
		}

		token2, err2 := target.LoginAdmin(ctx, instance.Kc_target.Username, instance.Kc_target.Password, "master")
		if err2 != nil {
			fmt.Errorf("Something wrong with the credentials on target" + err2.Error())
			return err2
		}

		log.Println("Adding Client...", client)

		err2 = instance.AddClient(ctx, token2.AccessToken, target, kc_client, Roles)
		if err2 != nil {
			log.Println("Something wrong, cannot add client" + err2.Error())
		}

	}
	return nil
}

func main() {
	var conf = &config.Config{}

	conf, err := config.New("config/config.yml")
	if err != nil {
		panic("cannot load config file")
	}
	instance := &Instance{conf}

	source := gocloak.NewClient(instance.Kc_source.Url)
	ctx := context.Background()

	target := gocloak.NewClient(instance.Kc_target.Url)

	err = KC_AddClients(ctx, *instance, source, target)
	if err != nil {
		panic("Something wrong, cannot add client" + err.Error())
	}
}
