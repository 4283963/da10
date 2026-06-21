import { useState } from 'react'

function SearchBox({ onSearch, loading }) {
  const [vin, setVin] = useState('')

  const handleSubmit = (e) => {
    e.preventDefault()
    if (vin.trim()) {
      onSearch(vin.trim().toUpperCase())
    }
  }

  return (
    <form className="search-box" onSubmit={handleSubmit}>
      <input
        type="text"
        value={vin}
        onChange={(e) => setVin(e.target.value)}
        placeholder="请输入车架号（VIN）查询过户进度"
        maxLength={50}
      />
      <button type="submit" className="btn" disabled={loading || !vin.trim()}>
        {loading ? '查询中...' : '查询进度'}
      </button>
    </form>
  )
}

export default SearchBox
