package com.yqg.recall.common.util.pipeline;

public class NonReentrantLockByWait {
  //是否被锁
  private volatile boolean locked = false;

  //加锁
  public synchronized void lock() {
    //当某个线程获取锁成功，其他线程进入等待状态
    while (locked) {
      try {
        wait();
      } catch (InterruptedException e) {
        e.printStackTrace();
      }
    }
    //加锁成功，locked 设置为 true
    locked = true;
  }

  //释放锁
  public synchronized void unlock() {
    locked = false;
    notify();
  }
}
