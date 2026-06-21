function TransferTimeline({ nodes }) {
  if (!nodes || nodes.length === 0) return null

  const doneCount = nodes.filter((n) => n.status === 'done').length
  const doingCount = nodes.filter((n) => n.status === 'doing').length
  const total = nodes.length
  const progress = total > 0 ? Math.round((doneCount / total) * 100) : 0

  const getStatusText = (status) => {
    const map = {
      done: '已完成',
      doing: '进行中',
      pending: '待办理',
    }
    return map[status] || status
  }

  const formatTime = (timeStr) => {
    if (!timeStr) return ''
    const t = new Date(timeStr)
    return t.toLocaleString('zh-CN', {
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
    })
  }

  return (
    <div className="card">
      <div className="section-title">
        <span>🚗</span> 过户节点追踪
      </div>

      <div className="progress-bar-container">
        <div className="progress-bar-header">
          <span>
            已完成 {doneCount} / {total} 项（进行中 {doingCount} 项）
          </span>
          <span>{progress}%</span>
        </div>
        <div className="progress-bar">
          <div className="progress-bar-fill" style={{ width: `${progress}%` }} />
        </div>
      </div>

      <ul className="timeline">
        {nodes.map((node, index) => (
          <li key={node.id} className="timeline-item">
            <div className={`timeline-dot ${node.status}`}>
              {node.status === 'done' ? '✓' : index + 1}
            </div>
            <div className="timeline-content">
              <div className={`timeline-title ${node.status}`}>
                {node.node_name}
                <span className={`status-badge ${node.status}`} style={{ marginLeft: 8 }}>
                  {getStatusText(node.status)}
                </span>
              </div>
              <div className="timeline-meta">
                {node.location && <span>📍 {node.location} &nbsp;&nbsp;</span>}
                {node.operator && <span>👤 {node.operator} &nbsp;&nbsp;</span>}
                {node.started_at && <span>开始: {formatTime(node.started_at)} &nbsp;&nbsp;</span>}
                {node.completed_at && <span>完成: {formatTime(node.completed_at)}</span>}
              </div>
              {node.description && (
                <div className="timeline-desc">{node.description}</div>
              )}
            </div>
          </li>
        ))}
      </ul>
    </div>
  )
}

export default TransferTimeline
