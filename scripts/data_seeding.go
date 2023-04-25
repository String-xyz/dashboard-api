package scripts

import (
	"context"
	"os"

	"github.com/String-xyz/platform-admin-api/api"
	"github.com/String-xyz/platform-admin-api/env"
	"github.com/String-xyz/platform-admin-api/pkg/model"
	"github.com/String-xyz/platform-admin-api/pkg/repository"
	"github.com/String-xyz/platform-admin-api/pkg/store"
	"github.com/rs/zerolog"
)

// Duplicating this instead of making api.NewRepos public
func newRepos(config api.APIConfig) repository.Repositories {
	return repository.Repositories{
		Platform:         repository.NewPlatform(config.DB),
		PlatformMember:   repository.NewPlatformMember(config.DB),
		MemberToRole:     repository.NewMemberToRole(config.DB),
		MemberToPlatform: repository.NewMemberToPlatform(config.DB),
		MemberRole:       repository.NewMemberRole(config.DB),
		MemberInvite:     repository.NewMemberInvite(config.DB),
		Apikey:           repository.NewApikey(config.DB),
	}
}

func getEnv(key string) string {
	value, err := env.Get(key)
	if err != nil {
		panic(err)
	}
	return value
}

func DataSeeding() {
	// Initialize repos
	env.LoadEnv() // removed the err since in cloud this wont be loaded
	port, err := env.Get("PORT")
	if err != nil {
		panic(err)
	}

	lg := zerolog.New(os.Stdout)

	// Note: This will panic if the env is set to use docker and you run this script from the command line
	config := api.APIConfig{
		DB:     store.MustNewPG(),
		Port:   port,
		Logger: &lg,
	}

	repos := newRepos(config)
	nullCtx := context.Background()

	memberId := getEnv("MEMBER_ROLE_MEMBER_ID")
	_, err := repos.MemberRole.Create(nullCtx, model.MemberRole{ID: memberId, Name: "Member"})
	if err != nil {
		panic(err)
	}

	adminId := getEnv("MEMBER_ROLE_ADMIN_ID")
	_, err = repos.MemberRole.Create(nullCtx, model.MemberRole{ID: adminId, Name: "Admin"})
	if err != nil {
		panic(err)
	}

	ownerId := getEnv("MEMBER_ROLE_OWNER_ID")
	_, err = repos.MemberRole.Create(nullCtx, model.MemberRole{ID: ownerId, Name: "Owner"})
	if err != nil {
		panic(err)
	}
}

func MockSeeding() {
	// Initialize repos
	env.LoadEnv() // removed the err since in cloud this wont be loaded
	port, err := env.Get("PORT")
	if err != nil {
		panic(err)
	}

	lg := zerolog.New(os.Stdout)

	// Note: This will panic if the env is set to use docker and you run this script from the command line
	config := api.APIConfig{
		DB:     store.MustNewPG(),
		Port:   port,
		Logger: &lg,
	}

	repos := newRepos(config)
	nullCtx := context.Background()

	memberId := getEnv("MEMBER_ROLE_MEMBER_ID")
	_, err = repos.MemberRole.Create(nullCtx, model.MemberRole{ID: memberId, Name: "Member"})
	if err != nil {
		panic(err)
	}

	adminId := getEnv("MEMBER_ROLE_ADMIN_ID")
	_, err = repos.MemberRole.Create(nullCtx, model.MemberRole{ID: adminId, Name: "Admin"})
	if err != nil {
		panic(err)
	}

	ownerId := getEnv("MEMBER_ROLE_OWNER_ID")
	_, err = repos.MemberRole.Create(nullCtx, model.MemberRole{ID: ownerId, Name: "Owner"})
	if err != nil {
		panic(err)
	}
}

// func nullString(str string) sql.NullString {
// 	return sql.NullString{String: str, Valid: true}
// }
