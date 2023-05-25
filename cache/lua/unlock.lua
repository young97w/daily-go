--1.检查是不是自己的锁
--2.删除
if redis.call('get',KEYS[1]) == ARGV[1] then
    --是自己的锁
    return redis.call('del',KEYS[1])
else
    --不是自己的锁
    return 0
end
