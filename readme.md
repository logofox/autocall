INFO:

// 预期效果 ： 自动外呼。根据线路的并发量，一通接一通的往外呼电话。
// 实现思路 ： 设置全局chan数组，数组长度为线路的并发量。开启n个线程， 每线程一通，chan 写入值，每挂断一通电话，chan 清空一个，如此不断的进行，直到需要外拨的号码拨完.