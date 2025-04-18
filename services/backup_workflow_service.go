package services

import (
	"fluxton/configs"
	"fluxton/models"
	"fluxton/repositories"
	"fluxton/utils"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"github.com/samber/do"
	"os"
	"time"
)

// BackupWorkflowService is responsible for executing the backup workflow and runs as go routine
type BackupWorkflowService interface {
	Create(databaseName string, backupUUID uuid.UUID)
	Delete(databaseName string, backupUUID uuid.UUID)
}

type BackupWorkflowServiceImpl struct {
	s3Service  S3Service
	backupRepo *repositories.BackupRepository
}

func NewBackupWorkflowService(injector *do.Injector) (BackupWorkflowService, error) {
	s3Service, err := NewS3Service()
	if err != nil {
		return nil, err
	}

	backupRepo := do.MustInvoke[*repositories.BackupRepository](injector)

	return &BackupWorkflowServiceImpl{
		s3Service:  s3Service,
		backupRepo: backupRepo,
	}, nil
}

// Create pg_dump, copy file, ensure bucket exists, and upload to S3
func (s *BackupWorkflowServiceImpl) Create(databaseName string, backupUUID uuid.UUID) {
	backupFilePath := fmt.Sprintf("/tmp/%s.sql", backupUUID)

	// 1. Execute pg_dump
	if err := s.executePgDump(databaseName, backupFilePath); err != nil {
		s.handleBackupFailure(backupUUID, models.BackupStatusCreatingFailed, err.Error())

		return
	}

	// 2. Copy backup file to app container
	if err := s.copyBackupToAppContainer(backupFilePath, backupUUID); err != nil {
		s.handleBackupFailure(backupUUID, models.BackupStatusCreatingFailed, err.Error())

		return
	}

	// 3. Ensure backup bucket exists
	if err := s.ensureBackupBucketExists(); err != nil {
		s.handleBackupFailure(backupUUID, models.BackupStatusCreatingFailed, err.Error())

		return
	}

	fileBytes, err := s.readBackupFile(backupUUID)
	if err != nil {
		s.handleBackupFailure(backupUUID, models.BackupStatusCreatingFailed, err.Error())

		return
	}

	// 4. Upload backup to S3
	err = s.uploadBackupToS3(databaseName, backupUUID, fileBytes)
	if err != nil {
		s.handleBackupFailure(backupUUID, models.BackupStatusCreatingFailed, err.Error())

		return
	}

	// 5. Update backup status to completed
	err = s.backupRepo.UpdateStatus(backupUUID, models.BackupStatusCreated, "", time.Now())
	if err != nil {
		s.handleBackupFailure(backupUUID, models.BackupStatusCreatingFailed, err.Error())
	}

	// 6. Remove backup file from app container
	err = os.Remove(backupFilePath)
	if err != nil {
		log.Errorf("failed to remove backup file: %s", err)
	}
}

// Delete removes the backup file from s3
func (s *BackupWorkflowServiceImpl) Delete(databaseName string, backupUUID uuid.UUID) {
	filePath := fmt.Sprintf("%s/%s.sql", databaseName, backupUUID)
	err := s.s3Service.DeleteFile(configs.BackupBucketName, filePath)
	if err != nil {
		s.handleBackupFailure(backupUUID, models.BackupStatusDeletingFailed, err.Error())
	}

	_, err = s.backupRepo.Delete(backupUUID)
	if err != nil {
		log.Errorf("failed to delete backup: %s", err)
	}
}

// executePgDump runs the pg_dump command in fluxton_db container to create a database backup
func (s *BackupWorkflowServiceImpl) executePgDump(databaseName, backupFilePath string) error {
	command := []string{
		"docker",
		"exec",
		os.Getenv("DATABASE_CONTAINER_NAME"),
		"pg_dump",
		"-U",
		os.Getenv("DATABASE_USER"),
		"-d",
		databaseName,
		"-f", backupFilePath,
	}

	err := utils.ExecuteCommand(command)
	if err != nil {
		log.Errorf("pg_dump failed: %s", err)
	}

	return err
}

// copyBackupToAppContainer copies the backup file from the fluxton_db container to the fluxton_app container
func (s *BackupWorkflowServiceImpl) copyBackupToAppContainer(backupFilePath string, backupUUID uuid.UUID) error {
	dockerCpCommand := []string{
		"docker",
		"cp",
		fmt.Sprintf("%s:%s", os.Getenv("DATABASE_CONTAINER_NAME"), backupFilePath), // Source
		fmt.Sprintf("/tmp/%s.sql", backupUUID),                                     // Destination inside app container
	}

	err := utils.ExecuteCommand(dockerCpCommand)
	if err != nil {
		log.Errorf("failed to copy backup file: %s", err)
	}

	return err
}

// ensureBackupBucketExists checks if the backup bucket exists and creates it if necessary
func (s *BackupWorkflowServiceImpl) ensureBackupBucketExists() error {
	bucketExists := s.s3Service.BucketExists(configs.BackupBucketName)
	if !bucketExists {
		_, err := s.s3Service.CreateBucket(configs.BackupBucketName)
		if err != nil {
			log.Errorf("failed to create backup bucket: %s", err)

			return err
		}
	}

	return nil
}

// readBackupFile reads the backup file from the app container
func (s *BackupWorkflowServiceImpl) readBackupFile(backupUUID uuid.UUID) ([]byte, error) {
	fileBytes, err := os.ReadFile(fmt.Sprintf("/tmp/%s.sql", backupUUID))
	if err != nil {
		log.Errorf("failed to read backup file: %s", err)
	}

	return fileBytes, err
}

// uploadBackupToS3 uploads the backup file to S3
func (s *BackupWorkflowServiceImpl) uploadBackupToS3(databaseName string, backupUUID uuid.UUID, fileBytes []byte) error {
	filePath := fmt.Sprintf("%s/%s.sql", databaseName, backupUUID)
	err := s.s3Service.UploadFile(configs.BackupBucketName, filePath, fileBytes)
	if err != nil {
		log.Errorf("failed to upload backup to S3: %s", err)
	}

	return err
}

// handleBackupFailure updates the backup status to appropriate state and logs the error
func (s *BackupWorkflowServiceImpl) handleBackupFailure(backupUUID uuid.UUID, status, errorMessage string) {
	err := s.backupRepo.UpdateStatus(backupUUID, status, errorMessage, time.Now())
	if err != nil {
		log.Errorf("failed to update backup status: %s", err)
	}

	log.Errorf("backup failed: %s", errorMessage)
}
