package app

import (
	"fmt"
	"log"
	"net"

	v1 "github.com/RenterRus/sausage-profile/docs/proto/v1"
	protoServe "github.com/RenterRus/sausage-profile/internal/controller/grpc"
	"github.com/RenterRus/sausage-profile/internal/repo/psql"
	"github.com/RenterRus/sausage-profile/internal/usecase"
	"github.com/RenterRus/sausage-profile/internal/usecase/hashing"
	"github.com/RenterRus/sausage-profile/internal/usecase/jwt"
	"github.com/RenterRus/sausage-profile/internal/usecase/otp"
	"github.com/sourcegraph/conc/pool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const secretSize = 32

type App struct {
	conf *Config
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
	p := pool.New().WithErrors()

	// grpc
	p.Go(func() error {
		lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", a.conf.GRPC.Host, a.conf.GRPC.Port))
		if err != nil {
			return fmt.Errorf("failed to listen: %v", err)
		}

		s := grpc.NewServer()
		reflection.Register(s)

		defer func() {
			s.Stop()
		}()

		users, err := psql.NewDBManager(fmt.Sprintf("%s://%s:%s@%s:%d/%s",
			a.conf.PSQL.Provider, a.conf.PSQL.Username, a.conf.PSQL.Password,
			a.conf.PSQL.Host, a.conf.PSQL.Port, a.conf.PSQL.DBName))
		if err != nil {
			return fmt.Errorf("Run.NewDBManager: %w", err)
		}

		v1.RegisterAuthServiceServer(s, protoServe.NewManager(usecase.NewProfileManager(usecase.ProfileConf{
			OtpRepo:    otp.NewOTPManager(hashing.NewHashingManager([]byte(a.conf.OTPHashSecretKey)[:secretSize]), a.conf.Issuer),
			Hash:       hashing.NewHashingManager([]byte(a.conf.JWTHashSecretKey)[:secretSize]),
			JwtManager: jwt.NewJWTManager([]byte(a.conf.JWTSecretKey)[:secretSize]),
			UsersRepo:  users,
		})))

		log.Printf("gRPC server listening on %s", fmt.Sprintf("%s:%d", a.conf.GRPC.Host, a.conf.GRPC.Port))

		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}

		return nil
	})

	if err := p.Wait(); err != nil {
		return fmt.Errorf("Run: %w", err)
	}

	return nil
}
