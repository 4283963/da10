package models

import "time"

type NodeStatus string

const (
	NodePending NodeStatus = "pending"
	NodeDoing   NodeStatus = "doing"
	NodeDone    NodeStatus = "done"
)

type NodeType string

const (
	NodeBeijingArchive    NodeType = "beijing_archive"
	NodeBeijingInspection NodeType = "beijing_inspection"
	NodeBeijingSubmit     NodeType = "beijing_submit"
	NodeArchiveDelivery   NodeType = "archive_delivery"
	NodeLocalRegister     NodeType = "local_register"
	NodeLocalInspection   NodeType = "local_inspection"
	NodePlateMaking       NodeType = "plate_making"
	NodePlateDelivery     NodeType = "plate_delivery"
	NodeCompleted         NodeType = "completed"
)

type TransferNode struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	VehicleID   uint       `gorm:"index;not null" json:"vehicle_id"`
	NodeType    NodeType   `gorm:"size:50;not null" json:"node_type"`
	NodeName    string     `gorm:"size:100;not null" json:"node_name"`
	Status      NodeStatus `gorm:"size:20;default:pending" json:"status"`
	Description string     `gorm:"size:500" json:"description"`
	Operator    string     `gorm:"size:50" json:"operator"`
	Location    string     `gorm:"size:100" json:"location"`
	StartedAt   *time.Time `json:"started_at"`
	CompletedAt *time.Time `json:"completed_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func GetDefaultNodes(vehicleID uint) []TransferNode {
	nodeConfigs := []struct {
		NodeType NodeType
		NodeName string
	}{
		{NodeBeijingArchive, "北京提档"},
		{NodeBeijingInspection, "北京验车"},
		{NodeBeijingSubmit, "北京提交材料"},
		{NodeArchiveDelivery, "档案快递"},
		{NodeLocalRegister, "外地落户登记"},
		{NodeLocalInspection, "外地验车"},
		{NodePlateMaking, "新车牌制作"},
		{NodePlateDelivery, "车牌快递"},
		{NodeCompleted, "全部完成"},
	}

	nodes := make([]TransferNode, len(nodeConfigs))
	for i, cfg := range nodeConfigs {
		nodes[i] = TransferNode{
			VehicleID: vehicleID,
			NodeType:  cfg.NodeType,
			NodeName:  cfg.NodeName,
			Status:    NodePending,
		}
	}
	return nodes
}
