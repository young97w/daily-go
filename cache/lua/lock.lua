---
--- Generated by EmmyLua(https://github.com/EmmyLua)
--- Created by yangyuan.
--- DateTime: 2023/5/25 22:15
---
val = redis.call('get',KEYS[1])
if val == false then
    --说明锁没被持有
    return redis.call('set',KEYS[1],ARGV[1],'EX',ARGV[2])
elseif val == ARGV[1] then
    -- 自己的锁，续期就行
    redis.call('expire',KEYS[1],ARGV[2])
    return 'OK'
else
    --别人的锁
    return ''
end