package services

import (
	"fluxton/errs"
	"fluxton/models"
	"fluxton/policies"
	"fluxton/repositories"
	"fluxton/types"
	"github.com/samber/do"
	"time"
)

type DatabaseStatsService interface {
	GetTotalDatabaseSize(databaseName string, authUser models.AuthUser) (string, error)
	GetTotalIndexSize(databaseName string, authUser models.AuthUser) (string, error)
	GetUnusedIndexes(databaseName string, authUser models.AuthUser) ([]types.UnusedIndex, error)
	GetSlowQueries(databaseName string, authUser models.AuthUser) ([]types.SlowQuery, error)
	GetIndexScansPerTable(databaseName string, authUser models.AuthUser) ([]types.IndexScan, error)
	GetSizePerTable(databaseName string, authUser models.AuthUser) ([]types.TableSize, error)
	GetRowCountPerTable(databaseName string, authUser models.AuthUser) ([]types.TableRowCount, error)
	GetAll(databaseName string, authUser models.AuthUser) (models.DatabaseStat, error)
}

type DatabaseStatsServiceImpl struct {
	connectionService ConnectionService
	adminPolicy       *policies.AdminPolicy
	databaseRepo      *repositories.DatabaseRepository
}

func NewDatabaseStatsService(injector *do.Injector) (DatabaseStatsService, error) {
	connectionService := do.MustInvoke[ConnectionService](injector)
	policy := policies.NewAdminPolicy()
	databaseRepo := do.MustInvoke[*repositories.DatabaseRepository](injector)

	return &DatabaseStatsServiceImpl{
		connectionService: connectionService,
		adminPolicy:       policy,
		databaseRepo:      databaseRepo,
	}, nil
}

func (s *DatabaseStatsServiceImpl) GetTotalDatabaseSize(databaseName string, authUser models.AuthUser) (string, error) {
	if !s.adminPolicy.CanAccess(authUser) {
		return "", errs.NewForbiddenError("database_stats.error.forbidden")
	}

	dbStatsRepo, _, err := s.connectionService.GetDatabaseStatsRepo(databaseName, nil)
	if err != nil {
		return "", err
	}

	return dbStatsRepo.GetTotalDatabaseSize()
}

func (s *DatabaseStatsServiceImpl) GetTotalIndexSize(databaseName string, authUser models.AuthUser) (string, error) {
	if !s.adminPolicy.CanAccess(authUser) {
		return "", errs.NewForbiddenError("database_stats.error.forbidden")
	}

	dbStatsRepo, _, err := s.connectionService.GetDatabaseStatsRepo(databaseName, nil)
	if err != nil {
		return "", err
	}

	return dbStatsRepo.GetTotalIndexSize()
}

func (s *DatabaseStatsServiceImpl) GetUnusedIndexes(databaseName string, authUser models.AuthUser) ([]types.UnusedIndex, error) {
	if !s.adminPolicy.CanAccess(authUser) {
		return []types.UnusedIndex{}, errs.NewForbiddenError("database_stats.error.forbidden")
	}

	dbStatsRepo, _, err := s.connectionService.GetDatabaseStatsRepo(databaseName, nil)
	if err != nil {
		return []types.UnusedIndex{}, err
	}

	return dbStatsRepo.GetUnusedIndexes()
}

func (s *DatabaseStatsServiceImpl) GetSlowQueries(databaseName string, authUser models.AuthUser) ([]types.SlowQuery, error) {
	if !s.adminPolicy.CanAccess(authUser) {
		return []types.SlowQuery{}, errs.NewForbiddenError("database_stats.error.forbidden")
	}

	dbStatsRepo, _, err := s.connectionService.GetDatabaseStatsRepo(databaseName, nil)
	if err != nil {
		return []types.SlowQuery{}, err
	}

	return dbStatsRepo.GetSlowQueries()
}

func (s *DatabaseStatsServiceImpl) GetIndexScansPerTable(databaseName string, authUser models.AuthUser) ([]types.IndexScan, error) {
	if !s.adminPolicy.CanAccess(authUser) {
		return []types.IndexScan{}, errs.NewForbiddenError("database_stats.error.forbidden")
	}

	dbStatsRepo, _, err := s.connectionService.GetDatabaseStatsRepo(databaseName, nil)
	if err != nil {
		return []types.IndexScan{}, err
	}

	return dbStatsRepo.GetIndexScansPerTable()
}

func (s *DatabaseStatsServiceImpl) GetSizePerTable(databaseName string, authUser models.AuthUser) ([]types.TableSize, error) {
	if !s.adminPolicy.CanAccess(authUser) {
		return []types.TableSize{}, errs.NewForbiddenError("database_stats.error.forbidden")
	}

	dbStatsRepo, _, err := s.connectionService.GetDatabaseStatsRepo(databaseName, nil)
	if err != nil {
		return []types.TableSize{}, err
	}

	return dbStatsRepo.GetSizePerTable()
}

func (s *DatabaseStatsServiceImpl) GetRowCountPerTable(databaseName string, authUser models.AuthUser) ([]types.TableRowCount, error) {
	if !s.adminPolicy.CanAccess(authUser) {
		return []types.TableRowCount{}, errs.NewForbiddenError("database_stats.error.forbidden")
	}

	dbStatsRepo, _, err := s.connectionService.GetDatabaseStatsRepo(databaseName, nil)
	if err != nil {
		return []types.TableRowCount{}, err
	}

	return dbStatsRepo.GetRowCountPerTable()
}

func (s *DatabaseStatsServiceImpl) GetAll(databaseName string, authUser models.AuthUser) (models.DatabaseStat, error) {
	if !s.adminPolicy.CanAccess(authUser) {
		return models.DatabaseStat{}, errs.NewForbiddenError("database_stats.error.forbidden")
	}

	totalDatabaseSize, err := s.GetTotalDatabaseSize(databaseName, authUser)
	if err != nil {
		return models.DatabaseStat{}, err
	}

	totalIndexSize, err := s.GetTotalIndexSize(databaseName, authUser)
	if err != nil {
		return models.DatabaseStat{}, err
	}

	unusedIndexes, err := s.GetUnusedIndexes(databaseName, authUser)
	if err != nil {
		return models.DatabaseStat{}, err
	}

	tableCounts, err := s.GetRowCountPerTable(databaseName, authUser)
	if err != nil {
		return models.DatabaseStat{}, err
	}

	tableSizes, err := s.GetSizePerTable(databaseName, authUser)
	if err != nil {
		return models.DatabaseStat{}, err
	}

	return models.DatabaseStat{
		DatabaseName: databaseName,
		TotalSize:    totalDatabaseSize,
		IndexSize:    totalIndexSize,
		UnusedIndex:  unusedIndexes,
		TableCount:   tableCounts,
		TableSize:    tableSizes,
		CreatedAt:    time.Now(),
	}, nil
}
