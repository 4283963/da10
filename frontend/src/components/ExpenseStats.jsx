function ExpenseStats({ expenses }) {
  if (!expenses) return null

  let totalAmount = 0
  let windowAmount = 0
  let expressAmount = 0
  let otherAmount = 0

  expenses.forEach((e) => {
    totalAmount += e.amount
    if (e.expense_type === 'window') {
      windowAmount += e.amount
    } else if (e.expense_type === 'express') {
      expressAmount += e.amount
    } else {
      otherAmount += e.amount
    }
  })

  const getExpenseTypeText = (type) => {
    const map = {
      window: '窗口费',
      express: '快递费',
      other: '其他',
    }
    return map[type] || type
  }

  const formatDate = (dateStr) => {
    if (!dateStr) return ''
    return new Date(dateStr).toLocaleDateString('zh-CN')
  }

  return (
    <div className="card">
      <div className="section-title">
        <span>💰</span> 费用对账统计
      </div>

      <div className="stats-grid">
        <div className="stat-card total">
          <div className="stat-amount">¥{totalAmount.toFixed(2)}</div>
          <div className="stat-label">总费用</div>
        </div>
        <div className="stat-card window">
          <div className="stat-amount">¥{windowAmount.toFixed(2)}</div>
          <div className="stat-label">窗口费</div>
        </div>
        <div className="stat-card express">
          <div className="stat-amount">¥{expressAmount.toFixed(2)}</div>
          <div className="stat-label">快递费</div>
        </div>
        <div className="stat-card other">
          <div className="stat-amount">¥{otherAmount.toFixed(2)}</div>
          <div className="stat-label">其他费用</div>
        </div>
      </div>

      {expenses.length > 0 ? (
        <ul className="expense-list">
          {expenses.map((expense) => (
            <li key={expense.id} className="expense-item">
              <div className="expense-left">
                <div className="expense-type">
                  <span className="expense-name">{expense.type_name}</span>
                  <span className={`expense-type-tag ${expense.expense_type}`}>
                    {getExpenseTypeText(expense.expense_type)}
                  </span>
                </div>
                {expense.description && (
                  <div className="expense-desc">{expense.description}</div>
                )}
                <div className="expense-date">
                  {expense.pay_date && `支付日期: ${formatDate(expense.pay_date)}`}
                  {expense.receipt_no && ` · 票据号: ${expense.receipt_no}`}
                </div>
              </div>
              <div className="expense-amount">-¥{expense.amount.toFixed(2)}</div>
            </li>
          ))}
        </ul>
      ) : (
        <div className="empty-state">
          <div className="empty-state-icon">📋</div>
          <div className="empty-state-text">暂无费用记录</div>
        </div>
      )}
    </div>
  )
}

export default ExpenseStats
