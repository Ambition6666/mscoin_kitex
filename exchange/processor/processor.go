package processor

import (
	"common/database"
	"common/tools"
	"context"
	"encoding/json"
	"exchange/domain"
	"exchange/model"
	"exchange/rpc"
	"exchange/utils"
	"github.com/cloudwego/kitex/pkg/klog"
	"grpc_common/kitex_gen/market"
	"sort"
	"sync"
)

// CoinTradeFactory 货币交易所
// CoinTradeFactory 用于管理多个 CoinTrade 实例的工厂类
type CoinTradeFactory struct {
	tradeMap map[string]*CoinTrade // 存储交易对符号与 CoinTrade 实例的映射
	mux      sync.RWMutex          // 读写锁，确保并发安全
}

// InitCoinTradeFactory 初始化一个 CoinTradeFactory 实例
func InitCoinTradeFactory() *CoinTradeFactory {
	return &CoinTradeFactory{
		tradeMap: make(map[string]*CoinTrade), // 初始化交易映射表
	}
}

// GetCoinTrade 根据交易对符号获取对应的 CoinTrade 实例
func (f *CoinTradeFactory) GetCoinTrade(symbol string) *CoinTrade {
	f.mux.RLock()         // 读锁，防止读取时被写入操作影响
	defer f.mux.RUnlock() // 方法结束时解锁

	return f.tradeMap[symbol] // 返回指定交易对的 CoinTrade 实例
}

// AddCoinTrade 添加一个新的 CoinTrade 实例到工厂
func (f *CoinTradeFactory) AddCoinTrade(symbol string, trade *CoinTrade) {
	f.mux.Lock()         // 写锁，确保写操作的并发安全
	defer f.mux.Unlock() // 方法结束时解锁

	_, ok := f.tradeMap[symbol]
	if !ok { // 如果交易对不存在，则添加新的 CoinTrade 实例
		f.tradeMap[symbol] = trade
	}
}

func (f *CoinTradeFactory) Init() {
	ctx := context.Background()
	exchangeCoinRes, err := rpc.GetMarketClient().FindVisibleExchangeCoins(ctx, &market.MarketReq{})
	if err != nil {
		klog.Error(err)
		return
	}

	for _, v := range exchangeCoinRes.List {
		f.AddCoinTrade(v.Symbol, NewCoinTrade(v.Symbol, utils.GetRocketMQProducer()))
	}
}

// LimitPriceQueue 表示限价队列，支持并发安全的操作
type LimitPriceQueue struct {
	mux  sync.RWMutex // 读写锁，确保对队列的并发访问安全
	list TradeQueue   // 存储按价格排序的交易队列
}

// LimitPriceMap 表示每个价格点的订单列表
type LimitPriceMap struct {
	price float64                // 该限价订单的价格
	list  []*model.ExchangeOrder // 该价格点下的所有订单
}

// TradeQueue 定义了一个交易队列，按价格降序排序
type TradeQueue []*LimitPriceMap

// Len 返回队列的长度
func (t TradeQueue) Len() int {
	return len(t)
}

// Less 实现降序排列，比较两个价格点，返回是否 i 大于 j
func (t TradeQueue) Less(i, j int) bool {
	// 降序排列，根据价格比较
	return t[i].price > t[j].price
}

// Swap 交换队列中的两个元素
func (t TradeQueue) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

// TradeTimeQueue 定义了一个交易队列，按时间升序排序
type TradeTimeQueue []*model.ExchangeOrder

// Len 返回队列的长度
func (t TradeTimeQueue) Len() int {
	return len(t)
}

// Less 实现升序排列，比较两个订单的时间，返回是否 i 的时间早于 j
func (t TradeTimeQueue) Less(i, j int) bool {
	// 升序排列，根据订单时间比较
	return t[i].Time < t[j].Time
}

// Swap 交换队列中的两个订单
func (t TradeTimeQueue) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

// CoinTrade 交易处理器
type CoinTrade struct {
	buyMarketQueue   TradeTimeQueue              // 买方市场订单队列
	sellMarketQueue  TradeTimeQueue              // 卖方市场订单队列
	buyLimitQueue    *LimitPriceQueue            // 买方限价订单队列（价格从高到低排序）
	sellLimitQueue   *LimitPriceQueue            // 卖方限价订单队列（价格从低到高排序）
	symbol           string                      // 交易对符号（如 BTC/USDT）
	buyTradePlate    *TradePlate                 // 买盘（用于存储买方订单信息）
	sellTradePlate   *TradePlate                 // 卖盘（用于存储卖方订单信息）
	rocketMQProducer *database.RocketMQProducer  // rocketmq 客户端，用于处理消息队列
	orderDomain      *domain.ExchangeOrderDomain // 订单领域逻辑的接口，用于处理订单相关操作
}

// NewCoinTrade CoinTrade 构造器
func NewCoinTrade(symbol string, client *database.RocketMQProducer) *CoinTrade {
	t := &CoinTrade{
		symbol:           symbol,
		rocketMQProducer: client,
		orderDomain:      domain.NewExchangeOrderDomain(),
	}
	t.init()
	return t
}

// 初始化交易处理器，创建买盘和卖盘，初始化队列
func (t *CoinTrade) init() {
	t.buyTradePlate = NewTradePlate(t.symbol, model.BUY)   // 创建买盘
	t.sellTradePlate = NewTradePlate(t.symbol, model.SELL) // 创建卖盘
	t.buyLimitQueue = &LimitPriceQueue{}                   // 初始化买方限价订单队列
	t.sellLimitQueue = &LimitPriceQueue{}                  // 初始化卖方限价订单队列
	t.initQueue()                                          // 初始化订单队列
}

// Trade 处理传入的交易订单，根据订单类型和方向执行相应的交易匹配和订单管理操作
func (t *CoinTrade) Trade(order *model.ExchangeOrder) {
	// 根据订单方向选择相应的限价单队列和市价单队列
	var limitPriceList *LimitPriceQueue
	var marketPriceList TradeTimeQueue

	if order.Direction == model.BUY {
		limitPriceList = t.sellLimitQueue
		marketPriceList = t.sellMarketQueue
	} else {
		limitPriceList = t.buyLimitQueue
		marketPriceList = t.buyMarketQueue
	}

	// 根据订单类型执行不同的匹配策略
	if order.Type == model.MarketPrice {
		// 市价单与限价单匹配
		t.matchLimitPriceWithMP(marketPriceList, order)
	} else if order.Type == model.LimitPrice {
		// 限价单先与限价单匹配，然后与市价单匹配
		t.matchLimitPriceWithLP(limitPriceList, order)
		if order.Status == model.Trading {
			// 如果限价单未完成，再与市价单匹配
			t.matchLimitPriceWithMP(marketPriceList, order)
		}
		if order.Status == model.Trading {
			// 如果订单仍在交易中，添加到限价单队列
			t.addLimitPriceOrder(order)
			if order.Direction == model.BUY {
				t.sendTradePlateMsg(t.buyTradePlate)
			} else {
				t.sendTradePlateMsg(t.sellTradePlate)
			}
		}
	} else {
		// 市价单与限价单匹配
		t.matchMarketPriceWithLP(limitPriceList, order)
	}
}

// GetTradePlate 根据方向返回对应的交易盘
func (t *CoinTrade) GetTradePlate(direction int) *TradePlate {
	if direction == model.BUY {
		return t.buyTradePlate
	}
	return t.sellTradePlate
}

// sendTradePlateMsg 发送交易盘消息到 RocketMQ
func (t *CoinTrade) sendTradePlateMsg(plate *TradePlate) {
	bytes, _ := json.Marshal(plate.Result(24))
	data := database.RocketMQData{
		Topic: "exchange_order_trade_plate",
		Key:   []string{plate.Symbol},
		Data:  bytes,
	}
	t.rocketMQProducer.Send(data)
}

// initQueue 初始化交易队列
func (t *CoinTrade) initQueue() {
	ctx := context.Background()
	list, err := t.orderDomain.FindTradingOrders(ctx)
	if err != nil {
		klog.Error(err)
		return
	}
	for _, v := range list {
		if v.Direction == model.BUY {
			// 处理买单
			if v.Type == model.MarketPrice {
				// 市价买单
				t.buyMarketQueue = append(t.buyMarketQueue, v)
			} else {
				// 限价买单
				isPut := false
				for _, bv := range t.buyLimitQueue.list {
					if bv.price == v.Price {
						bv.list = append(bv.list, v)
						isPut = true
						break
					}
				}
				if !isPut {
					plm := &LimitPriceMap{
						price: v.Price,
					}
					plm.list = append(plm.list, v)
					t.buyLimitQueue.list = append(t.buyLimitQueue.list, plm)
				}
				t.buyTradePlate.Add(v)
			}
		} else {
			// 处理卖单
			if v.Type == model.MarketPrice {
				// 市价卖单
				t.sellMarketQueue = append(t.sellMarketQueue, v)
			} else {
				// 限价卖单
				isPut := false
				for _, bv := range t.sellLimitQueue.list {
					if bv.price == v.Price {
						bv.list = append(bv.list, v)
						isPut = true
						break
					}
				}
				if !isPut {
					plm := &LimitPriceMap{
						price: v.Price,
					}
					plm.list = append(plm.list, v)
					t.sellLimitQueue.list = append(t.sellLimitQueue.list, plm)
				}
				t.sellTradePlate.Add(v)
			}
		}
	}
	// 对队列进行排序
	sort.Sort(t.sellMarketQueue)
	sort.Sort(t.buyMarketQueue)
	sort.Sort(t.buyLimitQueue.list)
	sort.Sort(sort.Reverse(t.sellLimitQueue.list))
}

// matchMarketPriceWithLP 处理市价单与限价单的匹配
func (t *CoinTrade) matchMarketPriceWithLP(lpList *LimitPriceQueue, focusedOrder *model.ExchangeOrder) {
	lpList.mux.Lock()
	defer lpList.mux.Unlock()
	buyNotify := false
	sellNotify := false
	for _, v := range lpList.list {
		var delOrders []string
		for _, matchOrder := range v.list {
			if matchOrder.MemberId == focusedOrder.MemberId {
				// 自己不与自己交易
				continue
			}
			// 计算可交易数量
			price := matchOrder.Price
			matchAmount := tools.SubFloor(matchOrder.Amount, matchOrder.TradedAmount, 8)
			focuseAmount := tools.SubFloor(focusedOrder.Amount, focusedOrder.TradedAmount, 8)
			if focusedOrder.Direction == model.BUY {
				// 买单计算以市价单为基准
				focuseAmount = tools.DivFloor(tools.SubFloor(focusedOrder.Amount, focusedOrder.Turnover, 8), price, 8)
			}
			if matchAmount >= focuseAmount {
				// 匹配成功
				matchOrder.TradedAmount = tools.AddFloor(matchOrder.TradedAmount, focuseAmount, 8)
				focusedOrder.TradedAmount = tools.AddFloor(focusedOrder.TradedAmount, focuseAmount, 8)
				to := tools.MulFloor(price, focuseAmount, 8)
				focusedOrder.Turnover = tools.AddFloor(focusedOrder.Turnover, to, 8)
				matchOrder.Turnover = tools.AddFloor(matchOrder.Turnover, to, 8)
				focusedOrder.Status = model.Completed
				if tools.SubFloor(matchOrder.Amount, matchOrder.TradedAmount, 8) <= 0 {
					// 如果匹配订单完成，从队列中删除
					matchOrder.Status = model.Completed
					delOrders = append(delOrders, matchOrder.OrderId)
				}
				if matchOrder.Direction == model.BUY {
					t.buyTradePlate.Remove(matchOrder, focuseAmount)
					buyNotify = true
				} else {
					t.sellTradePlate.Remove(matchOrder, focuseAmount)
					sellNotify = true
				}
				break
			} else {
				// 部分匹配
				to := tools.MulFloor(price, matchAmount, 8)
				matchOrder.TradedAmount = tools.AddFloor(matchOrder.TradedAmount, matchAmount, 8)
				matchOrder.Turnover = tools.AddFloor(matchOrder.Turnover, to, 8)
				matchOrder.Status = model.Completed
				delOrders = append(delOrders, matchOrder.OrderId)
				focusedOrder.TradedAmount = tools.AddFloor(focusedOrder.TradedAmount, matchAmount, 8)
				focusedOrder.Turnover = tools.AddFloor(focusedOrder.Turnover, to, 8)
				// 继续匹配
				if matchOrder.Direction == model.BUY {
					t.buyTradePlate.Remove(matchOrder, matchAmount)
					buyNotify = true
				} else {
					t.sellTradePlate.Remove(matchOrder, matchAmount)
					sellNotify = true
				}
				continue
			}
		}
		for _, orderId := range delOrders {
			for index, order := range v.list {
				if order.OrderId == orderId {
					v.list = append(v.list[:index], v.list[index+1:]...)
					break
				}
			}
		}
	}
	// 判断订单是否完成，未完成则放入队列
	if focusedOrder.Status == model.Trading {
		t.addMarketPriceOrder(focusedOrder)
	}
	// 通知买卖盘更新
	if buyNotify {
		t.sendTradePlateMsg(t.buyTradePlate)
	}
	if sellNotify {
		t.sendTradePlateMsg(t.sellTradePlate)
	}
}

// addMarketPriceOrder 将市价单添加到对应的队列中
func (t *CoinTrade) addMarketPriceOrder(order *model.ExchangeOrder) {
	if order.Type != model.MarketPrice {
		return
	}
	if order.Direction == model.BUY {
		t.buyMarketQueue = append(t.buyMarketQueue, order)
		sort.Sort(t.buyMarketQueue)
	} else {
		t.sellMarketQueue = append(t.sellMarketQueue, order)
		sort.Sort(t.sellMarketQueue)
	}
}

// addLimitPriceOrder 将限价单添加到对应的队列中
func (t *CoinTrade) addLimitPriceOrder(order *model.ExchangeOrder) {
	if order.Type != model.LimitPrice {
		return
	}
	if order.Direction == model.BUY {
		isPut := false
		for _, v := range t.buyLimitQueue.list {
			if v.price == order.Price {
				v.list = append(v.list, order)
				isPut = true
				break
			}
		}
		if !isPut {
			plm := &LimitPriceMap{
				price: order.Price,
			}
			plm.list = append(plm.list, order)
			t.buyLimitQueue.list = append(t.buyLimitQueue.list, plm)
			sort.Sort(t.buyLimitQueue.list)
		}
	} else {
		isPut := false
		for _, v := range t.sellLimitQueue.list {
			if v.price == order.Price {
				v.list = append(v.list, order)
				isPut = true
				break
			}
		}
		if !isPut {
			plm := &LimitPriceMap{
				price: order.Price,
			}
			plm.list = append(plm.list, order)
			t.sellLimitQueue.list = append(t.sellLimitQueue.list, plm)
			sort.Sort(sort.Reverse(t.sellLimitQueue.list))
		}
	}
}

func (t *CoinTrade) matchLimitPriceWithMP(mpList TradeTimeQueue, focusedOrder *model.ExchangeOrder) {
	//市价单时间是 从旧到新 先去匹配之前的单
	var delOrders []string
	for _, matchOrder := range mpList {
		if matchOrder.MemberId == focusedOrder.MemberId {
			//自己不与自己交易
			continue
		}
		price := focusedOrder.Price
		//可交易的数量
		matchAmount := tools.SubFloor(matchOrder.Amount, matchOrder.TradedAmount, 8)
		focusedAmount := tools.SubFloor(focusedOrder.Amount, focusedOrder.TradedAmount, 8)
		if matchAmount >= focusedAmount {
			//能够进行匹配，直接完成即可
			matchOrder.TradedAmount = tools.AddFloor(matchOrder.TradedAmount, focusedAmount, 8)
			focusedOrder.TradedAmount = tools.AddFloor(focusedOrder.TradedAmount, focusedAmount, 8)
			to := tools.MulFloor(price, focusedAmount, 8)
			focusedOrder.Turnover = tools.AddFloor(focusedOrder.Turnover, to, 8)
			matchOrder.Turnover = tools.AddFloor(matchOrder.Turnover, to, 8)
			focusedOrder.Status = model.Completed
			if tools.SubFloor(matchOrder.Amount, matchOrder.TradedAmount, 8) <= 0 {
				matchOrder.Status = model.Completed
				delOrders = append(delOrders, matchOrder.OrderId)
			}
			break
		} else {
			to := tools.MulFloor(price, matchAmount, 8)
			matchOrder.TradedAmount = tools.AddFloor(matchOrder.TradedAmount, matchAmount, 8)
			matchOrder.Turnover = tools.AddFloor(matchOrder.Turnover, to, 8)
			matchOrder.Status = model.Completed
			delOrders = append(delOrders, matchOrder.OrderId)
			focusedOrder.TradedAmount = tools.AddFloor(focusedOrder.TradedAmount, matchAmount, 8)
			focusedOrder.Turnover = tools.AddFloor(focusedOrder.Turnover, to, 8)
			//还得继续下一轮匹配
			continue
		}
	}
	for _, orderId := range delOrders {
		for index, order := range mpList {
			if order.OrderId == orderId {
				mpList = append(mpList[:index], mpList[index+1:]...)
				break
			}
		}
	}
}

func (t *CoinTrade) matchLimitPriceWithLP(lpList *LimitPriceQueue, focusedOrder *model.ExchangeOrder) {
	lpList.mux.Lock()
	defer lpList.mux.Unlock()
	buyNotify := false
	sellNotify := false
	var completeOrders []*model.ExchangeOrder
	for _, v := range lpList.list {
		var delOrders []string
		for _, matchOrder := range v.list {
			if matchOrder.MemberId == focusedOrder.MemberId {
				//自己不与自己交易
				continue
			}
			if focusedOrder.Direction == model.BUY {
				//买单 matchOrder为限价卖单 价格从低到高
				if matchOrder.Price > focusedOrder.Price {
					//最低卖价 比 买入价高 直接退出
					break
				}
			}
			if focusedOrder.Direction == model.SELL {
				if matchOrder.Price < focusedOrder.Price {
					//最高买价 比 卖价 低 直接退出
					break
				}
			}
			price := matchOrder.Price
			//可交易的数量
			matchAmount := tools.SubFloor(matchOrder.Amount, matchOrder.TradedAmount, 8)
			focuseAmount := tools.SubFloor(focusedOrder.Amount, focusedOrder.TradedAmount, 8)
			if matchAmount <= 0 {
				//证明已经交易完成
				matchOrder.Status = model.Completed
				delOrders = append(delOrders, matchOrder.OrderId)
				completeOrders = append(completeOrders, matchOrder)
				continue
			}
			if matchAmount >= focuseAmount {
				//能够进行匹配，直接完成即可
				matchOrder.TradedAmount = tools.AddFloor(matchOrder.TradedAmount, focuseAmount, 8)
				focusedOrder.TradedAmount = tools.AddFloor(focusedOrder.TradedAmount, focuseAmount, 8)
				to := tools.MulFloor(price, focuseAmount, 8)
				focusedOrder.Turnover = tools.AddFloor(focusedOrder.Turnover, to, 8)
				matchOrder.Turnover = tools.AddFloor(matchOrder.Turnover, to, 8)
				focusedOrder.Status = model.Completed

				if tools.SubFloor(matchOrder.Amount, matchOrder.TradedAmount, 8) <= 0 {
					//matchorder也完成了 需要从匹配列表中删除
					matchOrder.Status = model.Completed
					delOrders = append(delOrders, matchOrder.OrderId)
					completeOrders = append(completeOrders, matchOrder)
				}
				if matchOrder.Direction == model.BUY {
					t.buyTradePlate.Remove(matchOrder, focuseAmount)
					buyNotify = true
				} else {
					t.sellTradePlate.Remove(matchOrder, focuseAmount)
					sellNotify = true
				}
				break
			} else {
				to := tools.MulFloor(price, matchAmount, 8)
				matchOrder.TradedAmount = tools.AddFloor(matchOrder.TradedAmount, matchAmount, 8)
				matchOrder.Turnover = tools.AddFloor(matchOrder.Turnover, to, 8)
				matchOrder.Status = model.Completed
				delOrders = append(delOrders, matchOrder.OrderId)
				completeOrders = append(completeOrders, matchOrder)
				focusedOrder.TradedAmount = tools.AddFloor(focusedOrder.TradedAmount, matchAmount, 8)
				focusedOrder.Turnover = tools.AddFloor(focusedOrder.Turnover, to, 8)
				//还得继续下一轮匹配
				if matchOrder.Direction == model.BUY {
					t.buyTradePlate.Remove(matchOrder, matchAmount)
					buyNotify = true
				} else {
					t.sellTradePlate.Remove(matchOrder, matchAmount)
					sellNotify = true
				}
				continue
			}
		}
		for _, orderId := range delOrders {
			for index, order := range v.list {
				if order.OrderId == orderId {
					v.list = append(v.list[:index], v.list[index+1:]...)
					break
				}
			}
		}
	}
	//通知买卖盘更新
	if buyNotify {
		t.sendTradePlateMsg(t.buyTradePlate)
	}
	if sellNotify {
		t.sendTradePlateMsg(t.sellTradePlate)
	}
	t.onCompleteHandle(completeOrders)
}

func (t *CoinTrade) onCompleteHandle(orders []*model.ExchangeOrder) {
	if len(orders) <= 0 {
		return
	}
	for _, order := range orders {
		marshal, err := json.Marshal(order)
		if err != nil {
			klog.Error("封装已完成数据错误:", err)
			return
		}
		klog.Info("准备开始发已完成的数据")
		rocketmqData := database.RocketMQData{
			Topic: "exchange_order_completed",
			Key:   []string{t.symbol},
			Data:  marshal,
		}
		for {
			//保证一定发成功
			t.rocketMQProducer.Send(rocketmqData)
		}
	}

}

// TradePlate 表示交易盘口信息，维护订单的深度信息
type TradePlate struct {
	Items     []*TradePlateItem `json:"items"` // 盘口中的具体项，包含价格和数量等信息
	Symbol    string            // 交易对的符号，例如 "BTC/USDT"
	direction int               // 交易方向，买入或卖出 (BUY 或 SELL)
	maxDepth  int               // 最大深度，限制盘口显示的层级深度
	mux       sync.RWMutex      // 读写锁，确保并发访问时的安全
}

// TradePlate 构造器
func NewTradePlate(symbol string, direction int) *TradePlate {
	return &TradePlate{
		Symbol:    symbol,
		direction: direction,
		maxDepth:  100,
	}
}

// Add 将限价订单加入 TradePlate (盘口)
// 该函数会检查订单的方向、类型，并根据价格与当前盘口进行匹配，
// 如果价格相同则增加相应数量，否则如果未达到最大深度，则添加新项。
func (p *TradePlate) Add(order *model.ExchangeOrder) {
	// 如果订单的方向与当前盘口方向不符，记录错误日志并退出
	if p.direction != order.Direction {
		klog.Error("买卖盘 direction not match，check code...")
		return
	}

	p.mux.Lock() // 加锁以确保并发访问安全
	defer p.mux.Unlock()

	// 市价单不加入买卖盘，直接返回
	if order.Type == model.MarketPrice {
		klog.Error("市价单 不加入买卖盘")
		return
	}

	size := len(p.Items) // 获取当前盘口中的项数

	// 检查是否有与当前订单价格相同的项，如果有则增加该项的数量
	if size > 0 {
		for _, v := range p.Items {
			// 如果是买单且当前项的价格高于订单价格，或是卖单且当前项的价格低于订单价格，则跳过
			if (order.Direction == model.BUY && v.Price > order.Price) ||
				(order.Direction == model.SELL && v.Price < order.Price) {
				continue
			} else if v.Price == order.Price { // 如果价格相同，则增加相应的数量
				v.Amount = tools.AddN(v.Amount, tools.SubFloor(order.Amount, order.TradedAmount, 5), 5)
				return
			} else {
				break // 如果价格不符且需要插入新项，则跳出循环
			}
		}
	}

	// 如果当前盘口深度未达到最大限制，则添加新的 TradePlateItem
	if size < p.maxDepth {
		tpi := &TradePlateItem{
			Amount: tools.SubFloor(order.Amount, order.TradedAmount, 5), // 计算未成交的剩余数量
			Price:  order.Price,                                         // 设置价格
		}
		p.Items = append(p.Items, tpi) // 将新项加入到盘口中
	}
}

// TradePlateResult 表示交易盘口的结果数据，用于展示当前买卖盘的状态信息
type TradePlateResult struct {
	Direction    string            `json:"direction"`    // 买卖方向，可能为 "buy" 或 "sell"
	MaxAmount    float64           `json:"maxAmount"`    // 盘口中最大交易量
	MinAmount    float64           `json:"minAmount"`    // 盘口中最小交易量
	HighestPrice float64           `json:"highestPrice"` // 盘口中最高价格
	LowestPrice  float64           `json:"lowestPrice"`  // 盘口中最低价格
	Symbol       string            `json:"symbol"`       // 交易对的符号，例如 BTC/USDT
	Items        []*TradePlateItem `json:"items"`        // 盘口中的具体项列表，每项包含价格和数量
}

// AllResult 返回当前交易盘的所有结果，包括最大最小金额、最高最低价格以及所有项
func (p *TradePlate) AllResult() *TradePlateResult {
	result := &TradePlateResult{}
	direction := model.DirectionMap.Value(p.direction) // 根据方向值获取对应的描述
	result.Direction = direction
	result.MaxAmount = p.getMaxAmount()       // 获取最大金额
	result.MinAmount = p.getMinAmount()       // 获取最小金额
	result.HighestPrice = p.getHighestPrice() // 获取最高价格
	result.LowestPrice = p.getLowestPrice()   // 获取最低价格
	result.Symbol = p.Symbol                  // 交易对符号
	result.Items = p.Items                    // 当前所有的挂单项
	return result
}

// Result 返回当前交易盘的部分结果，最多包含 num 个项
func (p *TradePlate) Result(num int) *TradePlateResult {
	if num > len(p.Items) {
		num = len(p.Items) // 如果请求的数量超过当前项数，则限制为当前项数
	}
	result := &TradePlateResult{}
	direction := model.DirectionMap.Value(p.direction) // 根据方向值获取对应的描述
	result.Direction = direction
	result.MaxAmount = p.getMaxAmount()       // 获取最大金额
	result.MinAmount = p.getMinAmount()       // 获取最小金额
	result.HighestPrice = p.getHighestPrice() // 获取最高价格
	result.LowestPrice = p.getLowestPrice()   // 获取最低价格
	result.Symbol = p.Symbol                  // 交易对符号
	result.Items = p.Items[:num]              // 仅返回前 num 项
	return result
}

// getMaxAmount 获取当前交易盘中挂单项的最大金额
func (p *TradePlate) getMaxAmount() float64 {
	if len(p.Items) <= 0 {
		return 0 // 如果没有挂单项，则返回 0
	}
	var amount float64 = 0
	for _, v := range p.Items {
		if v.Amount > amount {
			amount = v.Amount // 更新最大金额
		}
	}
	return amount
}

// getMinAmount 获取当前交易盘中挂单项的最小金额
func (p *TradePlate) getMinAmount() float64 {
	if len(p.Items) <= 0 {
		return 0 // 如果没有挂单项，则返回 0
	}
	var amount float64 = p.Items[0].Amount
	for _, v := range p.Items {
		if v.Amount < amount {
			amount = v.Amount // 更新最小金额
		}
	}
	return amount
}

// getHighestPrice 获取当前交易盘中挂单项的最高价格
func (p *TradePlate) getHighestPrice() float64 {
	if len(p.Items) <= 0 {
		return 0 // 如果没有挂单项，则返回 0
	}
	var price float64 = 0
	for _, v := range p.Items {
		if v.Price > price {
			price = v.Price // 更新最高价格
		}
	}
	return price
}

// getLowestPrice 获取当前交易盘中挂单项的最低价格
func (p *TradePlate) getLowestPrice() float64 {
	if len(p.Items) <= 0 {
		return 0 // 如果没有挂单项，则返回 0
	}
	var price float64 = p.Items[0].Price
	for _, v := range p.Items {
		if v.Price < price {
			price = v.Price // 更新最低价格
		}
	}
	return price
}

// Remove 从交易盘中移除指定价格的挂单项，减少指定的金额
func (p *TradePlate) Remove(order *model.ExchangeOrder, amount float64) {
	for i, item := range p.Items {
		if item.Price == order.Price {
			item.Amount = tools.SubFloor(item.Amount, amount, 8) // 减少挂单项的金额
			if item.Amount <= 0 {
				// 如果金额减少到 0 或以下，则从列表中移除该项
				p.Items = append(p.Items[:i], p.Items[i+1:]...)
			}
			break
		}
	}
}

// TradePlateItem 代表交易盘中的单个挂单项
type TradePlateItem struct {
	Price  float64 `json:"price"`  // 挂单的价格
	Amount float64 `json:"amount"` // 挂单的数量
}
