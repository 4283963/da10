import { useState, useRef } from 'react'
import { uploadAPI, nodeAPI } from '../services/api.js'

function TransferTimeline({ nodes, vin, onNodesUpdated }) {
  const [previewNode, setPreviewNode] = useState(null)
  const [editingNodeId, setEditingNodeId] = useState(null)
  const [uploadingId, setUploadingId] = useState(null)
  const [updatingId, setUpdatingId] = useState(null)
  const [formData, setFormData] = useState({
    status: '',
    description: '',
    proof_image: '',
    operator: '',
    location: '',
  })
  const fileInputRef = useRef(null)
  const uploadTargetNodeId = useRef(null)

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

  const handleDotClick = (node) => {
    setPreviewNode(node)
  }

  const closeModal = () => {
    setPreviewNode(null)
  }

  const startEdit = (node) => {
    setEditingNodeId(node.id)
    setFormData({
      status: node.status,
      description: node.description || '',
      proof_image: node.proof_image || '',
      operator: node.operator || '',
      location: node.location || '',
    })
  }

  const cancelEdit = () => {
    setEditingNodeId(null)
    setFormData({ status: '', description: '', proof_image: '', operator: '', location: '' })
  }

  const handleFormChange = (e) => {
    const { name, value } = e.target
    setFormData((prev) => ({ ...prev, [name]: value }))
  }

  const triggerUpload = (nodeId) => {
    uploadTargetNodeId.current = nodeId
    if (fileInputRef.current) {
      fileInputRef.current.value = ''
      fileInputRef.current.click()
    }
  }

  const handleFileChange = async (e) => {
    const file = e.target.files?.[0]
    if (!file) return

    const nodeId = uploadTargetNodeId.current
    if (!nodeId) return

    setUploadingId(nodeId)
    try {
      const result = await uploadAPI.uploadImage(file)
      const imageUrl = result.data.url

      const node = nodes.find((n) => n.id === nodeId)
      if (!node) return

      await nodeAPI.update(nodeId, {
        status: node.status,
        proof_image: imageUrl,
      })

      if (onNodesUpdated) {
        await onNodesUpdated()
      }
    } catch (err) {
      alert(typeof err === 'string' ? err : '上传失败，请重试')
    } finally {
      setUploadingId(null)
      uploadTargetNodeId.current = null
    }
  }

  const submitUpdate = async () => {
    if (!editingNodeId || !formData.status) return

    setUpdatingId(editingNodeId)
    try {
      await nodeAPI.update(editingNodeId, formData)
      if (onNodesUpdated) {
        await onNodesUpdated()
      }
      cancelEdit()
    } catch (err) {
      alert(typeof err === 'string' ? err : '更新失败，请重试')
    } finally {
      setUpdatingId(null)
    }
  }

  return (
    <div className="card">
      <input
        type="file"
        ref={fileInputRef}
        accept="image/*"
        style={{ display: 'none' }}
        onChange={handleFileChange}
      />

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
            <div className="timeline-dot-wrapper">
              <div
                className={`timeline-dot ${node.status} clickable ${node.proof_image ? 'has-proof' : ''}`}
                onClick={() => handleDotClick(node)}
                title={node.proof_image ? '点击查看凭证照片' : '点击查看详情'}
              >
                {node.status === 'done' ? '✓' : index + 1}
              </div>
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

              {node.proof_image && editingNodeId !== node.id && (
                <img
                  src={node.proof_image}
                  alt={node.node_name}
                  className="node-proof-thumb"
                  onClick={() => handleDotClick(node)}
                />
              )}

              {editingNodeId !== node.id && (
                <div className="node-actions">
                  <button
                    className="edit-toggle-btn"
                    onClick={() => startEdit(node)}
                  >
                    ✏️ 更新节点
                  </button>
                  <button
                    className={`upload-btn ${uploadingId === node.id ? 'uploading' : ''}`}
                    onClick={() => triggerUpload(node.id)}
                    disabled={uploadingId === node.id}
                  >
                    {uploadingId === node.id ? '⬆️ 上传中...' : '📎 上传凭证'}
                  </button>
                  {node.proof_image && (
                    <button
                      className="view-proof-btn"
                      onClick={() => handleDotClick(node)}
                    >
                      🔍 查看凭证
                    </button>
                  )}
                </div>
              )}

              {editingNodeId === node.id && (
                <div className="update-form">
                  <div className="form-group">
                    <label>节点状态</label>
                    <select name="status" value={formData.status} onChange={handleFormChange}>
                      <option value="pending">待办理</option>
                      <option value="doing">进行中</option>
                      <option value="done">已完成</option>
                    </select>
                  </div>
                  <div className="form-group">
                    <label>办理说明</label>
                    <textarea
                      name="description"
                      value={formData.description}
                      onChange={handleFormChange}
                      placeholder="填写办理情况说明..."
                    />
                  </div>
                  <div className="form-group">
                    <label>凭证照片URL（也可点下方按钮上传）</label>
                    <input
                      type="text"
                      name="proof_image"
                      value={formData.proof_image}
                      onChange={handleFormChange}
                      placeholder="https://..."
                    />
                  </div>
                  <div className="form-group">
                    <label>经办人</label>
                    <input
                      type="text"
                      name="operator"
                      value={formData.operator}
                      onChange={handleFormChange}
                      placeholder="办事员姓名..."
                    />
                  </div>
                  <div className="form-group">
                    <label>办理地点</label>
                    <input
                      type="text"
                      name="location"
                      value={formData.location}
                      onChange={handleFormChange}
                      placeholder="车管所/快递站等..."
                    />
                  </div>

                  <div className="node-actions">
                    <button
                      className="upload-btn"
                      onClick={() => triggerUpload(node.id)}
                      disabled={uploadingId === node.id}
                    >
                      {uploadingId === node.id ? '⬆️ 上传中...' : '📎 上传新凭证'}
                    </button>
                  </div>

                  {formData.proof_image && (
                    <img
                      src={formData.proof_image}
                      alt="预览"
                      className="node-proof-thumb"
                      onClick={() =>
                        setPreviewNode({ ...node, proof_image: formData.proof_image })
                      }
                    />
                  )}

                  <div className="form-actions">
                    <button className="btn btn-secondary btn-small" onClick={cancelEdit}>
                      取消
                    </button>
                    <button
                      className="btn btn-small"
                      onClick={submitUpdate}
                      disabled={updatingId === node.id}
                    >
                      {updatingId === node.id ? '保存中...' : '保存修改'}
                    </button>
                  </div>
                </div>
              )}
            </div>
          </li>
        ))}
      </ul>

      {previewNode && (
        <div className="modal-overlay" onClick={closeModal}>
          <div className="modal-content" onClick={(e) => e.stopPropagation()}>
            <div className="modal-header">
              <div className="modal-title">
                {previewNode.node_name}
                <span
                  className={`status-badge ${previewNode.status}`}
                  style={{ marginLeft: 10 }}
                >
                  {getStatusText(previewNode.status)}
                </span>
              </div>
              <button className="modal-close" onClick={closeModal}>
                ✕
              </button>
            </div>
            <div className="modal-body">
              {(previewNode.location || previewNode.operator) && (
                <div style={{ marginBottom: 16, color: '#6b7280', fontSize: 14 }}>
                  {previewNode.location && <span>📍 {previewNode.location} &nbsp;&nbsp;</span>}
                  {previewNode.operator && <span>👤 经办人: {previewNode.operator}</span>}
                  {previewNode.started_at && (
                    <span style={{ marginLeft: 12 }}>
                      开始: {formatTime(previewNode.started_at)}
                    </span>
                  )}
                  {previewNode.completed_at && (
                    <span style={{ marginLeft: 12 }}>
                      完成: {formatTime(previewNode.completed_at)}
                    </span>
                  )}
                </div>
              )}
              {previewNode.description && (
                <div
                  style={{
                    padding: '14px 16px',
                    background: '#f9fafb',
                    borderRadius: 10,
                    marginBottom: 16,
                    color: '#374151',
                    fontSize: 14,
                    lineHeight: 1.6,
                  }}
                >
                  {previewNode.description}
                </div>
              )}
              {previewNode.proof_image ? (
                <img
                  src={previewNode.proof_image}
                  alt={previewNode.node_name}
                  className="proof-image-preview"
                />
              ) : (
                <div className="proof-image-placeholder">
                  <div style={{ fontSize: 48, marginBottom: 12 }}>🖼️</div>
                  <div style={{ fontSize: 15, marginBottom: 4 }}>该节点暂无凭证照片</div>
                  <div style={{ fontSize: 13 }}>办事员可在下方点击「上传凭证」添加</div>
                </div>
              )}
            </div>
          </div>
        </div>
      )}
    </div>
  )
}

export default TransferTimeline
