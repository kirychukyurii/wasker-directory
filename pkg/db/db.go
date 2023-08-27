package db

import (
	"context"
	"github.com/kirychukyurii/wasker-directory/pkg/utils"

	sq "github.com/Masterminds/squirrel"
	pgxzero "github.com/jackc/pgx-zerolog"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/rs/zerolog"

	"github.com/kirychukyurii/wasker-directory/internal/config"
	"github.com/kirychukyurii/wasker-directory/pkg/logger"
	"github.com/kirychukyurii/wasker-directory/pkg/server/rpc/interceptor"
)

type Database struct {
	Pool *pgxpool.Pool
}

func New(config config.Config, log logger.Logger) Database {
	ctx := context.Background()

	contextFunc := func(ctx context.Context, zeroCtx zerolog.Context) zerolog.Context {
		reqId := utils.FromContext(ctx, interceptor.XRequestIDCtxKey{})
		if reqId != nil {
			return zeroCtx.Str(interceptor.XRequestIDKey, reqId.(string))
		}

		return zeroCtx
	}

	// TODO: log messages with level trace (now it info)
	dblogger := pgxzero.NewLogger(log.Logger, pgxzero.WithoutPGXModule(), pgxzero.WithContextFunc(contextFunc))
	dblevel, err := tracelog.LogLevelFromString(log.GetLevel().String())
	if err != nil {
		log.Err(err).Msg("setup a pgx tracing level to default: info")
		dblevel = tracelog.LogLevelInfo
	}

	cfg, err := pgxpool.ParseConfig(config.Database.DSN())
	if err != nil {
		log.Fatal().Err(err).Msg("parsing database connection string")
	}

	cfg.ConnConfig.Tracer = &tracelog.TraceLog{
		Logger:   dblogger,
		LogLevel: dblevel,
	}

	dbpool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("create connection pool")
	}

	conn, err := dbpool.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		log.Fatal().Err(err).Msg("acquire connection from pool")
	}

	if err = conn.Ping(ctx); err != nil {
		log.Fatal().Err(err).Msg("pinging connection pool")
	}

	return Database{
		Pool: dbpool,
	}
}

func (a *Database) Dialect() sq.StatementBuilderType {
	return sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
}
