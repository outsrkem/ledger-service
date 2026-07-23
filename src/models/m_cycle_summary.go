package models

import "github.com/shopspring/decimal"

type CycleSummary struct {
	Income  decimal.Decimal `gorm:"column:income"`  // 总收入（正数）
	Expense decimal.Decimal `gorm:"column:expense"` // 总支出（正数）
}

// CycleSummary 周期汇总
func (d *DB) CycleSummary(from, to int64) (*CycleSummary, error) {
	var res CycleSummary

	db := d.WithInstance().Table(TableNameTransaction)
	db = db.Where("occ_at >= ? AND occ_at < ?", from, to)
	err := db.Select(`
		IFNULL(SUM(CASE WHEN amount > 0 THEN amount ELSE 0 END), 0) AS income,
		IFNULL(ABS(SUM(CASE WHEN amount < 0 THEN amount ELSE 0 END)), 0) AS expense
	`).Scan(&res).Error

	return &res, err
}

// CategoryMainItem 仅大类汇总（饼图主数据）
type CategoryMainItem struct {
	CategoryID   int     `json:"categoryId"`
	CategoryName string  `json:"categoryName"`
	Direction    int     `json:"direction"`
	Amount       float64 `json:"amount"`
}
type CategoryMainResult struct {
	IncomeList  []CategoryMainItem `json:"incomeList"`
	ExpenseList []CategoryMainItem `json:"expenseList"`
}

// CategorySubItem 小类明细
type CategorySubItem struct {
	SubID   int     `json:"subId"`
	SubName string  `json:"subName"`
	Amount  float64 `json:"amount"`
}

// CategoryParentWithSub 大类+下属小类双层结构
type CategoryParentWithSub struct {
	ParentID    int               `json:"parentId"`
	ParentName  string            `json:"parentName"`
	Direction   int               `json:"direction"`
	TotalAmount float64           `json:"totalAmount"`
	ChildList   []CategorySubItem `json:"childList"`
}
type CategoryDetailResult struct {
	IncomeParentList  []CategoryParentWithSub `json:"incomeParentList"`
	ExpenseParentList []CategoryParentWithSub `json:"expenseParentList"`
}

// QueryCategoryMain 仅按大类汇总（现有表无冗余字段）
// QueryCategoryMain 纯原生表结构，大类汇总查询
func (d *DB) QueryCategoryMain(from, to int64) (*CategoryMainResult, error) {
	type row struct {
		Pid       int     `gorm:"column:pid"`
		PName     string  `gorm:"column:pname"`
		Direction int     `gorm:"column:direction"`
		TotalAmt  float64 `gorm:"column:total_amt"`
	}
	var rows []row

	// 手动指定表别名t，消除instance_id模糊字段报错
	err := d.db.
		Table("ledger_transaction t").
		Where("t.instance_id = ?", d.instanceId).
		Where("t.occ_at >= ? AND t.occ_at <= ?", from, to).
		Joins("LEFT JOIN ledger_category self_cat ON t.category_id = self_cat.kid").
		Joins("LEFT JOIN ledger_caterela rel ON rel.child_id = self_cat.kid").
		Joins("LEFT JOIN ledger_category parent_cat ON " +
			"(self_cat.layer=2 AND rel.parent_id = parent_cat.kid) OR " +
			"(self_cat.layer=1 AND self_cat.kid = parent_cat.kid)").
		Select(`
			parent_cat.kid pid,
			parent_cat.name pname,
			parent_cat.direction,
			SUM(ABS(t.amount)) total_amt
		`).
		Group("pid, pname, parent_cat.direction").
		Order("total_amt DESC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	// 初始化空切片，避免json输出null
	res := &CategoryMainResult{
		IncomeList:  make([]CategoryMainItem, 0),
		ExpenseList: make([]CategoryMainItem, 0),
	}

	for _, r := range rows {
		item := CategoryMainItem{
			CategoryID:   r.Pid,
			CategoryName: r.PName,
			Direction:    r.Direction,
			Amount:       r.TotalAmt,
		}
		if r.Direction == 1 {
			res.IncomeList = append(res.IncomeList, item)
		} else {
			res.ExpenseList = append(res.ExpenseList, item)
		}
	}
	return res, nil
}

// QueryCategoryDetail 大类下带出全部小类明细，现有原生表结构
func (d *DB) QueryCategoryDetail(from, to int64) (*CategoryDetailResult, error) {
	type row struct {
		Pid       int     `gorm:"column:pid"`
		PName     string  `gorm:"column:pname"`
		Direction int     `gorm:"column:direction"`
		Sid       int     `gorm:"column:sid"`
		SName     string  `gorm:"column:sname"`
		SubAmt    float64 `gorm:"column:sub_amt"`
	}
	var rows []row

	// 关键修复：WHERE 使用 t.instance_id 明确指定交易表别名，消除模糊列
	err := d.db.
		Table("ledger_transaction t").
		Where("t.instance_id = ?", d.instanceId).
		Where("t.occ_at >= ? AND t.occ_at <= ?", from, to).
		Joins("LEFT JOIN ledger_category self_cat ON t.category_id = self_cat.kid").
		Joins("LEFT JOIN ledger_caterela rel ON rel.child_id = self_cat.kid").
		Joins("LEFT JOIN ledger_category parent_cat ON " +
			"(self_cat.layer=2 AND rel.parent_id = parent_cat.kid) OR " +
			"(self_cat.layer=1 AND self_cat.kid = parent_cat.kid)").
		Select(`
			parent_cat.kid pid,
			parent_cat.name pname,
			parent_cat.direction,
			self_cat.kid sid,
			self_cat.name sname,
			SUM(ABS(t.amount)) sub_amt
		`).
		Group("pid, pname, parent_cat.direction, sid, sname").
		Order("pid, sub_amt DESC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	// 初始化空切片，避免JSON返回null
	res := &CategoryDetailResult{
		IncomeParentList:  make([]CategoryParentWithSub, 0),
		ExpenseParentList: make([]CategoryParentWithSub, 0),
	}
	parentMap := make(map[int]*CategoryParentWithSub)

	for _, r := range rows {
		// 不存在则新建大类
		if _, ok := parentMap[r.Pid]; !ok {
			parentMap[r.Pid] = &CategoryParentWithSub{
				ParentID:    r.Pid,
				ParentName:  r.PName,
				Direction:   r.Direction,
				TotalAmount: 0,
				ChildList:   make([]CategorySubItem, 0),
			}
		}
		p := parentMap[r.Pid]
		amt := r.SubAmt
		p.TotalAmount += amt
		p.ChildList = append(p.ChildList, CategorySubItem{
			SubID:   r.Sid,
			SubName: r.SName,
			Amount:  amt,
		})
	}

	// 拆分收入/支出
	for _, item := range parentMap {
		if item.Direction == 1 {
			res.IncomeParentList = append(res.IncomeParentList, *item)
		} else {
			res.ExpenseParentList = append(res.ExpenseParentList, *item)
		}
	}

	return res, nil
}

// QueryCategoryStat 分类构成统一入口
// groupLayer 1=仅大类 2=大类+小类明细
func (d *DB) QueryCategoryStat(from, to int64, groupLayer int) (any, error) {
	if groupLayer == 1 {
		return d.QueryCategoryMain(from, to)
	}
	return d.QueryCategoryDetail(from, to)
}

// ---------------------收支排行

// ExpenseRankItem 支出分类排行项（支持大类）
type ExpenseRankItem struct {
	CategoryId   int     `json:"categoryId"`
	CategoryName string  `json:"categoryName"`
	TotalAmount  float64 `json:"totalAmount"`
}

// ExpenseRankResult 支出排行返回
type ExpenseRankResult struct {
	List []ExpenseRankItem `json:"list"`
}

// QueryExpenseRank 收支分类排行，direction 1收入 2支出
func (d *DB) QueryExpenseRank(from, to int64, topNum int, direction int) (*ExpenseRankResult, error) {
	type row struct {
		Pid      int     `gorm:"column:pid"`
		PName    string  `gorm:"column:pname"`
		TotalAmt float64 `gorm:"column:total_amt"`
	}
	var rows []row

	err := d.db.
		Table("ledger_transaction t").
		Where("t.instance_id = ?", d.instanceId).
		Where("t.occ_at >= ? AND t.occ_at <= ?", from, to).
		Joins("LEFT JOIN ledger_category self_cat ON t.category_id = self_cat.kid").
		Joins("LEFT JOIN ledger_caterela rel ON rel.child_id = self_cat.kid").
		Joins("LEFT JOIN ledger_category parent_cat ON "+
			"(self_cat.layer=2 AND rel.parent_id = parent_cat.kid) OR "+
			"(self_cat.layer=1 AND self_cat.kid = parent_cat.kid)").
		// 动态传入收支类型，不再写死2
		Where("parent_cat.direction = ?", direction).
		Select(`
			parent_cat.kid pid,
			parent_cat.name pname,
			SUM(ABS(t.amount)) total_amt
		`).
		Group("pid, pname").
		Order("total_amt DESC").
		Limit(topNum).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	res := &ExpenseRankResult{
		List: make([]ExpenseRankItem, 0),
	}
	for _, r := range rows {
		res.List = append(res.List, ExpenseRankItem{
			CategoryId:   r.Pid,
			CategoryName: r.PName,
			TotalAmount:  r.TotalAmt,
		})
	}
	return res, nil
}

type SouZhiResult struct {
	Income  decimal.Decimal `gorm:"column:income" json:"income"`   // 总收入
	Expense decimal.Decimal `gorm:"column:expense" json:"expense"` // 总支出（正数）
}

// CountSouZhi 一个函数实现 按年/按月/自定义日期 统计收支
func (d *DB) CountSouZhi(timeType string, from, to string) (*SouZhiResult, error) {
	var res SouZhiResult

	// 基础DB
	db := d.WithInstance().Table("ledger_transaction")

	// ====================================================
	// 把结束日期 +1 天，并用 < 代替 <=，实现日期闭区间全覆盖
	// 例如：2026-03-31 → 变成 2026-04-01 00:00:00
	// 这样能查到 2026-03-31 全天所有数据
	// ====================================================
	db = db.Where(`
		STR_TO_DATE(occ_time, '%Y-%m-%dT%H:%i:%s+0800') >= ? 
		AND 
		STR_TO_DATE(occ_time, '%Y-%m-%dT%H:%i:%s+0800') < ? + INTERVAL 1 DAY
	`, from, to)

	// 统计收入、支出
	err := db.Select(`
		IFNULL(SUM(CASE WHEN amount > 0 THEN amount ELSE 0 END), 0) AS income,
		IFNULL(ABS(SUM(CASE WHEN amount < 0 THEN amount ELSE 0 END)), 0) AS expense
	`).Scan(&res).Error

	return &res, err
}
