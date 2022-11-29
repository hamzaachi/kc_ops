package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Nerzal/gocloak/v12"
	"scm.eadn.dz/DevOps/kc_ops/config"
)

func Kc_AddRealmRoles(ctx context.Context, instance Instance, source, target *gocloak.GoCloak) error {
	for _, role := range instance.Kc_source.Roles {

		token, err := source.LoginAdmin(ctx, instance.Kc_source.Username, instance.Kc_source.Password, "master")
		if err != nil {
			fmt.Errorf("Something wrong with the credentials on source" + err.Error())
			return err
		}

		log.Println("Getting Role Info...", role)
		MyRole, err := instance.GetRealmRole(role, ctx, token.AccessToken, source)
		if err != nil {
			return err
		}
		token2, err2 := target.LoginAdmin(ctx, instance.Kc_target.Username, instance.Kc_target.Password, "master")
		if err2 != nil {
			fmt.Errorf("Something wrong with the credentials on target" + err2.Error())
			return err2
		}
		log.Println("Adding Role ...", role)
		err2 = instance.AddRealmRole(ctx, token2.AccessToken, target, MyRole)
		if err2 != nil {
			return err2
		}
	}
	return nil
}

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

func Kc_AddGroups(ctx context.Context, instance Instance, source, target *gocloak.GoCloak) error {
	for _, group := range instance.Kc_source.Groups {

		token, err := source.LoginAdmin(ctx, instance.Kc_source.Username, instance.Kc_source.Password, "master")
		if err != nil {
			fmt.Errorf("Something wrong with the credentials on source" + err.Error())
			return err
		}

		log.Println("Getting Group Info...", group)
		MyGroup, err := instance.GetGroup(group, ctx, token.AccessToken, source)
		if err != nil {
			return err
		}
		token2, err2 := target.LoginAdmin(ctx, instance.Kc_target.Username, instance.Kc_target.Password, "master")
		if err2 != nil {
			fmt.Errorf("Something wrong with the credentials on target" + err2.Error())
			return err2
		}
		log.Println("Adding Group ...", group)

		err2 = instance.AddGroup(ctx, token2.AccessToken, target, MyGroup)
		if err2 != nil {
			return err2
		}
	}
	return nil
}

func main() {
	ConfigFile := flag.String("path", "config/config.yml", "Path of the config file")
	clients := flag.Bool("clients", false, "Whether to migrate Clients")
	roles := flag.Bool("roles", false, "Whether to migrate Roles")
	groups := flag.Bool("groups", false, "Whether to migrate Groups")

	flag.Parse()
	if _, err := os.Stat(*ConfigFile); err != nil {
		log.Fatalf("Error: File %s does not exist", *ConfigFile)
	}

	var conf = &config.Config{}
	conf, err := config.New(*ConfigFile)
	if err != nil {
		panic("cannot load config file")
	}
	instance := &Instance{conf}

	source := gocloak.NewClient(instance.Kc_source.Url)
	ctx := context.Background()

	target := gocloak.NewClient(instance.Kc_target.Url)
	if *clients {
		err = KC_AddClients(ctx, *instance, source, target)
		if err != nil {
			panic("Something is wrong, cannot add clients: " + err.Error())
		}
	}

	if *roles {
		err = Kc_AddRealmRoles(ctx, *instance, source, target)
		if err != nil {
			panic("Something is wrong, cannot add roles: " + err.Error())
		}
	}

	if *groups {

		err := Kc_AddGroups(ctx, *instance, source, target)
		if err != nil {
			panic(err.Error())
		}

	}
}
