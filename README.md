# 常用的go组件

## cache

支持Get,Set,GetMulti,Increment,Decrement,IsExist,FlushAll 操作
 * mecache不支持decrement小于0的,进行该操作时,将会返回error