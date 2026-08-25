package app

import (
	"context"
	"fmt"
	"log"
	"net"

	v1 "github.com/RenterRus/sausage-auth/docs/proto/v1"
	protoServe "github.com/RenterRus/sausage-auth/internal/controller/grpc"
	"github.com/RenterRus/sausage-auth/internal/repo/inmem"
	"github.com/RenterRus/sausage-auth/internal/repo/psql"
	"github.com/RenterRus/sausage-auth/internal/usecase"
	"github.com/RenterRus/sausage-auth/internal/usecase/hashing"
	"github.com/RenterRus/sausage-auth/internal/usecase/jwt"
	"github.com/RenterRus/sausage-auth/internal/usecase/otp"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const SECRET_SIZE = 32

type App struct {
	conf        *Config
	cacheClient *memcache.Client
	pgx         *pgxpool.Pool
	s           *grpc.Server
}

func NewApp(configPath string) (*App, error) {
	lastSlash := 0
	for i, v := range configPath {
		if v == '/' {
			lastSlash = i
		}
	}

	conf, err := ReadConfig(configPath[:lastSlash], configPath[lastSlash+1:])
	if err != nil {
		return nil, fmt.Errorf("ReadConfig: %w", err)
	}

	return &App{
		conf: conf,
	}, nil
}

func (a *App) Run() error {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", a.conf.GRPC.Host, a.conf.GRPC.Port))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	a.s = grpc.NewServer()
	reflection.Register(a.s)

	defer func() {
		a.s.Stop()
	}()

	a.pgx, err = pgxpool.New(context.Background(), fmt.Sprintf("%s://%s:%s@%s:%d/%s",
		a.conf.PSQL.Provider, a.conf.PSQL.Username, a.conf.PSQL.Password,
		a.conf.PSQL.Host, a.conf.PSQL.Port, a.conf.PSQL.DBName))
	if err != nil {
		return fmt.Errorf("Run.NewDBManager: %w", err)
	}
	defer func() {
		a.pgx.Close()
	}()

	a.cacheClient = memcache.New(fmt.Sprintf("%s:%d", a.conf.Memcache.Host, a.conf.Memcache.Port))
	defer func() {
		a.cacheClient.Close()
	}()

	v1.RegisterAuthServiceServer(a.s, protoServe.NewManager(usecase.NewProfileManager(usecase.ProfileConf{
		OtpRepo: otp.NewOTPManager(hashing.NewHashingManager([]byte(a.conf.OTPHashSecretKey)[:SECRET_SIZE]), a.conf.Issuer),
		Hash:    hashing.NewHashingManager([]byte(a.conf.JWTHashSecretKey)[:SECRET_SIZE]),
		JwtManager: jwt.NewJWTManager(jwt.JwtManagerConf{
			Key:        []byte(a.conf.JWTSecretKey)[:SECRET_SIZE],
			AccessExp:  a.conf.AccessExp,
			RefreshExp: a.conf.RefreshExp,
		}),
		UsersRepo: psql.NewDBManager(a.pgx),
		AccCache: inmem.NewAccessCache(inmem.AccessCacheConf{
			AccessExp: a.conf.AccessExp,
			Client:    a.cacheClient,
		}),
	})))

	log.Printf("gRPC server listening on %s", fmt.Sprintf("%s:%d", a.conf.GRPC.Host, a.conf.GRPC.Port))

	if err := a.s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	return nil
}

func (a *App) Close() {
	a.cacheClient.Close()

	a.pgx.Close()

	a.s.Stop()
}
