local redis = require("redis")
local time  = require("time")
local re    = require("re")
local match = re.match
local var   = fbi.var
local log   = fbi.log
local ERROR = fbi.ERROR
local rdkey = "test-fbi-"..var.addr

--------------
-- 自定义一些方法
-- 保证在某个时间段
local function calcTime(now, interval)

    return now - (now - time.zero)%interval
end

--------------
-- 只做数据增加
local function scanRecord(pipeline, tag, interval)
    if var.status > 499 or var.status < 400 then
        return
    end
    pipeline.incr(rdkey, tag.."-"..tostring(interval))
end

local function recorder()
    local pipeline = redis.pipeline
    pipeline.new()
    scanRecord(pipeline, "40x", 300)
    pipeline.exec()
    pipeline.close()
end

---------------
-- 做判断，根据条件修改值
local function scanCheck(pipeline, tag, interval, max, lock)

    -- 5分钟内10个以上40x请求,否则封锁10分钟
    -- status = "40x", interval = 300, max = 10, lock = 600

    -- 如果不是40x，直接结束函数
    if var.status > 499 or var.status < 400 then
        return
    end

    local now = time.unix()
    local l, err = redis.hmget(rdkey, "last-lock")
    if l ~= nil then
        if now - l < lock then
            log(ERROR, "old bad ip : ", var.addr)
            pipeline.hmset(rdkey, "last-lock", now)
            return
        end
    end
    -- 如果还没过锁定期，不做下列操作

    local n, err = redis.hmget(rdkey, "last-"..tostring(interval))
    if n == nil then
        n = 0
    end
    if err == nil then
        if now - n < interval then
            
            local m1, err = redis.hmget(rdkey, tag.."-"..tostring(interval))
            if tonumber(m1) > max then
                log(ERROR, "new bad ip : ", var.addr)
                pipeline.hmset(rdkey, "islock", 1)
                pipeline.hmset(rdkey, "last-lock", now)
            end
        else
            -- 更新时间key,更新这个时间段的值
            redis.hmset(rdkey, "last-"..tostring(interval), calcTime(now , interval) ) 
            redis.hmset(rdkey, tag.."-"..tostring(interval), 1)
        end
    end
end

local function checker()
    local pipeline = redis.pipeline
    pipeline.new()
    scanCheck(pipeline, "40x", 300, 100, 600)
    pipeline.exec()
    pipeline.close()
end


---------------
-- 运行
recorder()
checker()

if var.uri == "/test/fbi/uri/2" then
    print(var.addr)
end

