package services

import (
	"errors"
	"time"
	"transfer-tracker/database"
	"transfer-tracker/models"
)

type TransferNodeService struct{}

func NewTransferNodeService() *TransferNodeService {
	return &TransferNodeService{}
}

func (s *TransferNodeService) UpdateNode(id uint, status models.NodeStatus, description, proofImage, operator, location string) (*models.TransferNode, error) {
	var node models.TransferNode
	if err := database.DB.First(&node, id).Error; err != nil {
		return nil, errors.New("节点不存在")
	}

	now := time.Now()
	if status == models.NodeDoing && node.StartedAt == nil {
		node.StartedAt = &now
	}
	if status == models.NodeDone && node.CompletedAt == nil {
		node.CompletedAt = &now
	}

	node.Status = status
	if description != "" {
		node.Description = description
	}
	if proofImage != "" {
		node.ProofImage = proofImage
	}
	if operator != "" {
		node.Operator = operator
	}
	if location != "" {
		node.Location = location
	}

	if err := database.DB.Save(&node).Error; err != nil {
		return nil, err
	}

	return &node, nil
}

func (s *TransferNodeService) GetByVehicleID(vehicleID uint) ([]models.TransferNode, error) {
	var nodes []models.TransferNode
	if err := database.DB.Where("vehicle_id = ?", vehicleID).Order("id ASC").Find(&nodes).Error; err != nil {
		return nil, err
	}
	return nodes, nil
}

func (s *TransferNodeService) GetProgress(vehicleID uint) (map[string]interface{}, error) {
	nodes, err := s.GetByVehicleID(vehicleID)
	if err != nil {
		return nil, err
	}

	total := len(nodes)
	done := 0
	doing := 0
	for _, n := range nodes {
		if n.Status == models.NodeDone {
			done++
		} else if n.Status == models.NodeDoing {
			doing++
		}
	}

	var progress float64
	if total > 0 {
		progress = float64(done) / float64(total) * 100
	}

	return map[string]interface{}{
		"total":    total,
		"done":     done,
		"doing":    doing,
		"progress": progress,
		"nodes":    nodes,
	}, nil
}
