`KC_Ops` is `Golang` CLI for manipulating `Keycloak` assets, mainly migrating particular Clients, Roles and Groups from one KC instance to another , for cases where migrating the whole Realms is not feasible.

Usage of kc_ops:

  `--clients`

        Whether to migrate Clients

  `--groups`

        Whether to migrate Groups

  `--path` string

        Path of the config file (default "config/config.yml")

  `--roles`

        Whether to migrate Roles