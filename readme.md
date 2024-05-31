# Saving Simulator

Compound interest calculation for saving money 统一按日计息

## Core Objects

* 账户
    * 余额
    * 冻结金额
* 交易流水
    * 时间
    * 金额
    * 类型

## Saving Mode

如果设置了 complex 参数，则程序自动匹配最有理财产品进行购买

### Simple saving mode
按周期进行储蓄，直接买入活期+，不进行只能理财操作
### Complex saving mode
当余额满足智能理财金额时，将主动赎回活期+，并购买最优理财产品

## Saving Strategy

* 按日计息
* 存款金额
    * 周 (500)
    * 月 (2210)
    * 季度 (6630)
    * 半年 (13260)
    * 年 (26520)

* 理财产品 

  手动录入，`products.json`

* 买入

  trx: 时间，金额，类型 -- debit
  
  -> 购买记录 金额，起息日，到期日，利率

* 赎回

  trx: 时间，金额，类型-- credit

  -> 赎回记录 到账日期，金额，赎回日期（可有可无）

* 利息计算

  基于金额按日进行利息计算，按照结算日进行结算

* 结算
  
  汇总每日交易trx，更新账户余额
  
  周五，周六，节假日 只计息，不结算





