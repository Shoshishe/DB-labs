local cur_time = redis.call('TIME')
local key = KEYS[1]
local max_num = 100
local window = 60
local request_count = redis.call('ZCARD', key)

if request_count < max_num then
    redis.call('ZADD', key, cur_time[1], cur_time[1]..cur_time[2])
    redis.call('EXPIRE', key, window)
    return 0
end
return 1