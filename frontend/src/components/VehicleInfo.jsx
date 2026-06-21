function VehicleInfo({ vehicle }) {
  if (!vehicle) return null

  return (
    <div className="vehicle-info">
      <div className="vehicle-info-item">
        <span className="vehicle-info-label">车架号 (VIN)</span>
        <span className="vehicle-info-value">{vehicle.vin}</span>
      </div>
      {vehicle.plate_number && (
        <div className="vehicle-info-item">
          <span className="vehicle-info-label">原车牌号</span>
          <span className="vehicle-info-value">{vehicle.plate_number}</span>
        </div>
      )}
      {vehicle.owner_name && (
        <div className="vehicle-info-item">
          <span className="vehicle-info-label">车主姓名</span>
          <span className="vehicle-info-value">{vehicle.owner_name}</span>
        </div>
      )}
      {vehicle.target_city && (
        <div className="vehicle-info-item">
          <span className="vehicle-info-label">转入城市</span>
          <span className="vehicle-info-value">{vehicle.target_city}</span>
        </div>
      )}
      {vehicle.remark && (
        <div className="vehicle-info-item" style={{ gridColumn: '1 / -1' }}>
          <span className="vehicle-info-label">备注</span>
          <span className="vehicle-info-value">{vehicle.remark}</span>
        </div>
      )}
    </div>
  )
}

export default VehicleInfo
