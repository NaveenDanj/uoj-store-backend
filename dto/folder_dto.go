package dto

type CreateFolderDTO struct {
	Name     string `json:"name" binding:"required"`
	ParentId uint   `json:"parent_id" binding:"required"`
}

type RenameFolderDTO struct {
	Name     string `json:"name" binding:"required"`
	FolderId uint   `json:"folder_id" binding:"required"`
}

type DeleteFolderDTO struct {
	FolderId uint `json:"folder_id" binding:"required"`
}

type MoveFolderDTO struct {
	FolderId            uint `json:"folder_id" binding:"required"`
	DestinationFolderID uint `json:"destination_folder_id" binding:"required"`
}

type MoveToTrashDTO struct {
	FolderId uint `json:"folder_id" binding:"required"`
}
