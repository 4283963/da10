import { useState } from 'react'
import SearchBox from './components/SearchBox.jsx'
import VehicleInfo from './components/VehicleInfo.jsx'
import TransferTimeline from './components/TransferTimeline.jsx'
import ExpenseStats from './components/ExpenseStats.jsx'
import { vehicleAPI } from './services/api.js'

function App() {
  const [vehicle, setVehicle] = useState(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [hasSearched, setHasSearched] = useState(false)
  const [lastVin, setLastVin] = useState('')

  const handleSearch = async (vin) => {
    setLastVin(vin)
    setLoading(true)
    setError('')
    setVehicle(null)
    setHasSearched(true)

    try {
      const data = await vehicleAPI.getByVIN(vin)
      setVehicle(data.data)
    } catch (err) {
      setError(typeof err === 'string' ? err : '查询失败，请稍后重试')
    } finally {
      setLoading(false)
    }
  }

  const refreshVehicle = async () => {
    if (!lastVin) return
    try {
      const data = await vehicleAPI.getByVIN(lastVin)
      setVehicle(data.data)
    } catch (err) {
      console.error('刷新数据失败:', err)
    }
  }

  return (
    <div className="container">
      <div className="header">
        <h1>🚗 车辆过户进度追踪</h1>
        <p>北京转外地 · 实时追踪过户节点和费用明细</p>
      </div>

      <div className="card">
        <SearchBox onSearch={handleSearch} loading={loading} />

        {loading && (
          <div className="loading">
            <div style={{ fontSize: 20, marginBottom: 8 }}>🔍</div>
            <div>正在查询车架号信息...</div>
          </div>
        )}

        {error && !loading && <div className="error-message">{error}</div>}

        {!vehicle && !loading && !error && hasSearched && (
          <div className="empty-state">
            <div className="empty-state-icon">😔</div>
            <div className="empty-state-text">未找到该车架号的信息</div>
            <div className="empty-state-hint">请检查车架号是否正确，或联系管理员录入</div>
          </div>
        )}

        {!hasSearched && !loading && (
          <div className="empty-state">
            <div className="empty-state-icon">📝</div>
            <div className="empty-state-text">请输入车架号开始查询</div>
            <div className="empty-state-hint">支持查询过户节点进度和费用统计</div>
          </div>
        )}
      </div>

      {vehicle && !loading && (
        <>
          <div className="card">
            <div className="section-title">
              <span>🚙</span> 车辆基本信息
            </div>
            <VehicleInfo vehicle={vehicle} />
          </div>

          <TransferTimeline
            nodes={vehicle.nodes}
            vin={vehicle.vin}
            onNodesUpdated={refreshVehicle}
          />

          <ExpenseStats expenses={vehicle.expenses} />
        </>
      )}
    </div>
  )
}

export default App
