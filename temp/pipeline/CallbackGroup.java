package com.yqg.recall.common.util.pipeline;

import lombok.extern.slf4j.Slf4j;
import org.apache.commons.lang3.tuple.Pair;

import java.util.LinkedList;
import java.util.List;
import java.util.concurrent.Callable;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Future;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.function.BiConsumer;
import java.util.function.Consumer;
import java.util.stream.Collectors;

@Slf4j
public class CallbackGroup {
  final AtomicInteger count = new AtomicInteger();
  ExecutorService service;
  List<Pair<String, Future>> futureInfos;
  List<Pair<String, Callable>> runnables = new LinkedList<>();
  Consumer<List<Pair<String, Future>>> callback = null;

  public CallbackGroup(ExecutorService service) {
    this.service = service;
    this.futureInfos = new LinkedList<>();
  }

  public void add(String logInfo, Runnable runnable) {
    runnables.add(Pair.of(logInfo, () -> {
      try {
        runnable.run();
      } finally {
        if (count.addAndGet(-1) == 0 && callback != null) {
          service.submit(() -> callback.accept(futureInfos));
        }
      }
      return 1;
    }));
  }

  public void add(String logInfo, Callable runnable) {
    runnables.add(Pair.of(logInfo, () -> {
      try {
        return runnable.call();
      } finally {
        if (count.addAndGet(-1) == 0 && callback != null) {
          service.submit(() -> callback.accept(futureInfos));
        }
      }
    }));
  }

  public void execute(Consumer<List<Pair<String, Future>>> callback) {
    this.callback = callback;
    run();
  }

  private void run() {
    count.set(runnables.size());
    futureInfos = runnables.stream().map(it -> Pair.of(it.getLeft(), (Future) service.submit(it.getRight()))).collect(Collectors.toList());
  }

  public <Value> void execute(BiConsumer<String, Exception> onError, Consumer<List<Value>> callback) {
    this.callback = list -> callback.accept(list.stream().map(pair -> {
      try {
        return (Value) pair.getRight().get();
      } catch (Exception e) {
        onError.accept(pair.getLeft(), e);
        return null;
      }
    }).collect(Collectors.toList()));

    run();
  }
}
